package scan

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dirxray/internal/ignore"
	"dirxray/internal/model"
)

// Scan walks the tree and builds a ScanResult.
func Scan(opts Options) (*model.ScanResult, error) {
	root := filepath.Clean(opts.Root)
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	st, err := os.Lstat(abs)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() && st.Mode()&os.ModeSymlink == 0 {
		return nil, fs.ErrNotExist
	}

	var matcher *ignore.Matcher
	if !opts.NoGitignore {
		matcher = ignore.New(abs, true)
	}

	res := &model.ScanResult{
		RootAbsPath: abs,
		StartedAt:   time.Now(),
	}

	if st.Mode()&os.ModeSymlink != 0 && opts.FollowSymlinks {
		if t, e := filepath.EvalSymlinks(abs); e == nil {
			abs = t
		}
	}

	node, stats, notices := buildTree(abs, abs, matcher, opts, 0, nil)
	res.Root = node
	res.Stats = stats
	res.Notices = notices
	res.FinishedAt = time.Now()
	return res, nil
}

func buildTree(absRoot, absPath string, matcher *ignore.Matcher, opts Options, depth int, parent *model.Node) (*model.Node, model.ScanStats, []model.ScanNotice) {
	var stats model.ScanStats
	var notices []model.ScanNotice

	rel, _ := filepath.Rel(absRoot, absPath)
	relSlash := filepath.ToSlash(rel)
	if relSlash == "." {
		relSlash = ""
	}

	st, err := os.Lstat(absPath)
	if err != nil {
		notices = append(notices, model.ScanNotice{
			Path:    absPath,
			Kind:    "stat_failed",
			Message: err.Error(),
		})
		return nil, stats, notices
	}

	name := filepath.Base(absPath)
	if name == "." {
		name = filepath.Base(absRoot)
	}

	isHidden := isHiddenName(name)
	if !opts.Hidden && isHidden && relSlash != "" {
		return nil, stats, notices
	}

	if matcher != nil && relSlash != "" && matcher.ShouldSkip(relSlash, st.IsDir()) {
		stats.Skipped++
		return nil, stats, notices
	}

	mode := st.Mode()
	isSymlink := mode&os.ModeSymlink != 0
	kind := model.NodeFile
	if mode.IsDir() {
		kind = model.NodeDir
	}
	if isSymlink {
		kind = model.NodeSymlink
	}

	ext := strings.ToLower(filepath.Ext(name))

	n := &model.Node{
		Path:      relSlash,
		Name:      name,
		Kind:      kind,
		Size:      st.Size(),
		ModTime:   st.ModTime(),
		Depth:     depth,
		Ext:       ext,
		Parent:    parent,
		IsHidden:  isHidden,
		IsSymlink: isSymlink,
		Role:      roleHint(name, ext, kind),
	}

	if isSymlink {
		if t, err := os.Readlink(absPath); err == nil {
			n.SymlinkTarget = t
		}
		if !opts.FollowSymlinks {
			n.Kind = model.NodeSymlink
			stats.Symlinks++
			if !st.IsDir() {
				stats.Files++
				stats.TotalBytes += st.Size()
			}
			if depth > stats.MaxDepth {
				stats.MaxDepth = depth
			}
			if isHidden {
				stats.HiddenCount++
			}
			return n, stats, notices
		}
		// follow: resolve for further walk
		if t, err := filepath.EvalSymlinks(absPath); err == nil {
			absPath = t
			st2, err2 := os.Lstat(absPath)
			if err2 != nil {
				notices = append(notices, model.ScanNotice{
					Path:    absPath,
					Kind:    "symlink_target",
					Message: err2.Error(),
				})
				return n, stats, notices
			}
			st = st2
			mode = st.Mode()
			isSymlink = false
			if st.IsDir() {
				kind = model.NodeDir
				n.Kind = model.NodeDir
			} else {
				kind = model.NodeFile
				n.Kind = model.NodeFile
			}
			n.Size = st.Size()
			n.ModTime = st.ModTime()
		}
	}

	if !st.IsDir() {
		stats.Files++
		stats.TotalBytes += st.Size()
		if isHidden {
			stats.HiddenCount++
		}
		if depth > stats.MaxDepth {
			stats.MaxDepth = depth
		}
		return n, stats, notices
	}

	stats.Dirs++
	if isHidden {
		stats.HiddenCount++
	}
	if depth > stats.MaxDepth {
		stats.MaxDepth = depth
	}

	if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
		return n, stats, notices
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		n.ScanErr = err
		if errors.Is(err, fs.ErrPermission) || os.IsPermission(err) {
			notices = append(notices, model.ScanNotice{
				Path:    absPath,
				Kind:    "permission_denied",
				Message: err.Error(),
			})
		} else {
			notices = append(notices, model.ScanNotice{
				Path:    absPath,
				Kind:    "read_dir_failed",
				Message: err.Error(),
			})
		}
		return n, stats, notices
	}

	for _, ent := range entries {
		childAbs := filepath.Join(absPath, ent.Name())
		cRel := filepath.Join(relSlash, ent.Name())
		cRelSlash := filepath.ToSlash(cRel)

		if !opts.Hidden && isHiddenName(ent.Name()) {
			continue
		}
		if matcher != nil && matcher.ShouldSkip(cRelSlash, ent.IsDir()) {
			stats.Skipped++
			continue
		}

		cn, cstats, cnotes := buildTree(absRoot, childAbs, matcher, opts, depth+1, n)
		notices = append(notices, cnotes...)
		stats.Skipped += cstats.Skipped
		if cn == nil {
			continue
		}
		n.Children = append(n.Children, cn)
		stats.Dirs += cstats.Dirs
		stats.Files += cstats.Files
		stats.Symlinks += cstats.Symlinks
		stats.TotalBytes += cstats.TotalBytes
		if cstats.MaxDepth > stats.MaxDepth {
			stats.MaxDepth = cstats.MaxDepth
		}
		stats.HiddenCount += cstats.HiddenCount
	}

	return n, stats, notices
}

func isHiddenName(name string) bool {
	if name == "" {
		return false
	}
	// Unix: dotfiles. Windows: also check attributes in future; dot is still common.
	return strings.HasPrefix(name, ".")
}

func roleHint(name, ext string, kind model.NodeKind) model.RoleHint {
	if kind == model.NodeDir && name == ".git" {
		return model.RoleVCS
	}
	switch strings.ToLower(name) {
	case "go.mod", "cargo.toml", "pyproject.toml", "setup.py", "package.json", "composer.json":
		return model.RoleManifest
	case "go.sum", "package-lock.json", "yarn.lock", "pnpm-lock.yaml", "poetry.lock", "pipfile.lock":
		return model.RoleLockfile
	case "dockerfile", "containerfile":
		return model.RoleConfig
	}
	switch ext {
	case ".yaml", ".yml", ".toml", ".ini", ".cfg", ".env":
		if strings.Contains(strings.ToLower(name), "lock") {
			return model.RoleLockfile
		}
		return model.RoleConfig
	case ".json":
		if strings.Contains(strings.ToLower(name), "lock") {
			return model.RoleLockfile
		}
		return model.RoleConfig
	case ".jsonl", ".ndjson":
		return model.RoleData
	case ".md", ".rst", ".adoc":
		return model.RoleDocs
	case ".csv", ".tsv", ".parquet":
		return model.RoleData
	}
	if ext == ".exe" || ext == ".dll" || ext == ".so" || ext == ".dylib" {
		return model.RoleBinary
	}
	return model.RoleUnknown
}
