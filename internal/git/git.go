package git

import (
	"os"
	"path/filepath"
)

// IsRepoRoot reports whether dir contains a .git entry (file or directory).
func IsRepoRoot(dir string) bool {
	p := filepath.Join(dir, ".git")
	st, err := os.Stat(p)
	return err == nil && (st.IsDir() || st.Mode().IsRegular())
}

// FindRepoRoot walks up from startDir looking for .git; returns empty if none.
func FindRepoRoot(startDir string) string {
	dir := filepath.Clean(startDir)
	for {
		if IsRepoRoot(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
