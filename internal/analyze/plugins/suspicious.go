package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dirxray/internal/git"
	"dirxray/internal/model"
)

type suspiciousPlugin struct{}

func NewSuspiciousPlugin() Plugin { return suspiciousPlugin{} }

func (suspiciousPlugin) Name() string { return "suspicious" }

func (suspiciousPlugin) Run(ctx *Context) (*model.PluginResult, error) {
	pr := &model.PluginResult{PluginName: "suspicious"}
	if ctx.Scan == nil || ctx.Scan.Root == nil {
		return pr, nil
	}

	var findings []model.Finding
	seq := 0
	add := func(title string, sev model.Severity, rat string, paths []string, ev ...model.EvidenceItem) {
		seq++
		findings = append(findings, model.Finding{
			ID:            fmt.Sprintf("finding-%d", seq),
			Title:         title,
			Severity:      sev,
			Rationale:     rat,
			Evidence:      ev,
			RelatedPaths:  paths,
		})
	}

	for _, n := range ctx.Scan.Notices {
		if n.Kind == "permission_denied" {
			add("Unreadable path", model.SeverityMedium, "scanner could not read a directory or file; inventory is incomplete.",
				[]string{n.Path},
				model.EvidenceItem{Label: "notice", Detail: n.Message, Path: n.Path})
		}
	}

	if p := findNestedGit(ctx.RootAbs, ctx.Scan.Root); p != "" {
		add("Nested repository marker", model.SeverityLow, "a .git directory exists under a subdirectory of the scan root (ignored during inventory); possible submodule or nested clone.",
			[]string{p},
			model.EvidenceItem{Label: "path", Detail: p, Path: p})
	}

	if git.IsRepoRoot(ctx.RootAbs) {
		model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
			if n.Kind == model.NodeDir {
				return true
			}
			if n.Size > 50*1024*1024 {
				ext := strings.ToLower(n.Ext)
				if ext == ".zip" || ext == ".tar" || ext == ".gz" || ext == "" || ext == ".bin" || ext == ".exe" {
					add("Large binary-like artifact in repo", model.SeverityLow, "very large file in a git repo may be an accidental check-in.",
						[]string{n.Path},
						model.EvidenceItem{Label: "size", Detail: formatSize(n.Size), Path: n.Path})
				}
			}
			return true
		})
	}

	seenLock := map[string]string{}
	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind == model.NodeDir {
			return true
		}
		base := strings.ToLower(filepath.Base(n.Name))
		if base == "package-lock.json" || base == "yarn.lock" || base == "pnpm-lock.yaml" {
			dir := filepath.ToSlash(filepath.Dir(n.Path))
			if dir == "." {
				dir = ""
			}
			if prev, ok := seenLock[dir]; ok && prev != base {
				add("Multiple JS lockfiles", model.SeverityMedium, "more than one package manager lockfile in the same folder suggests mixed tooling.",
					[]string{n.Path},
					model.EvidenceItem{Label: "lockfiles", Detail: prev + " and " + base, Path: n.Path})
			} else {
				seenLock[dir] = base
			}
		}
		return true
	})

	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind != model.NodeDir && git.IsRepoRoot(ctx.RootAbs) {
			base := strings.ToLower(n.Name)
			if base == ".env" || strings.HasPrefix(base, ".env.") {
				add("Environment file in tree", model.SeverityMedium, ".env-style files in a repo tree are often risky; confirm they are not committed with secrets.",
					[]string{n.Path},
					model.EvidenceItem{Label: "file", Detail: base, Path: n.Path})
			}
		}
		return true
	})

	var composePath string
	dockerfilePresent := false
	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind == model.NodeDir {
			return true
		}
		b := strings.ToLower(n.Name)
		if b == "docker-compose.yml" || b == "docker-compose.yaml" {
			composePath = n.Path
		}
		if b == "dockerfile" || b == "containerfile" {
			dockerfilePresent = true
		}
		return true
	})
	if composePath != "" && !dockerfilePresent {
		add("Compose without Dockerfile in scan", model.SeverityInfo, "docker-compose present but no Dockerfile in this tree; may use prebuilt images only.",
			[]string{composePath},
			model.EvidenceItem{Label: "compose", Detail: composePath, Path: composePath})
	}

	pr.Findings = findings
	return pr, nil
}

// findNestedGit returns the first relative path like "sub/.git" under rootAbs, excluding the repo root's own .git.
func findNestedGit(rootAbs string, root *model.Node) string {
	if root == nil {
		return ""
	}
	var hit string
	model.Walk(root, func(n *model.Node) bool {
		if n.Kind != model.NodeDir || hit != "" {
			return true
		}
		rel := n.Path
		if rel == "" || rel == "." {
			return true
		}
		gitPath := filepath.Join(rootAbs, filepath.FromSlash(rel), ".git")
		if st, err := os.Stat(gitPath); err == nil && (st.IsDir() || st.Mode().IsRegular()) {
			hit = filepath.ToSlash(filepath.Join(rel, ".git"))
			return false
		}
		return true
	})
	return hit
}

func formatSize(n int64) string {
	const kb = 1024
	switch {
	case n < kb:
		return fmt.Sprintf("%d B", n)
	case n < kb*kb:
		return fmt.Sprintf("%d KiB", n/kb)
	default:
		return fmt.Sprintf("%d MiB", n/(kb*kb))
	}
}
