package ignore

import (
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

// Matcher decides whether a path relative to root should be skipped.
type Matcher struct {
	root   string
	ignore *gitignore.GitIgnore
}

// New loads .gitignore at root if present and combines internal defaults.
func New(root string, useGitignore bool) *Matcher {
	root = filepath.Clean(root)
	var lines []string
	lines = append(lines, internalRules()...)
	if useGitignore {
		gi := filepath.Join(root, ".gitignore")
		if b, err := os.ReadFile(gi); err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				lines = append(lines, line)
			}
		}
	}
	compiled := gitignore.CompileIgnoreLines(lines...)
	return &Matcher{root: root, ignore: compiled}
}

func internalRules() []string {
	return []string{
		".git/",
		"node_modules/",
		"vendor/",
		".venv/",
		"venv/",
		"__pycache__/",
		".idea/",
		".vscode/",
		"*.pyc",
		".DS_Store",
		"Thumbs.db",
	}
}

// ShouldSkip returns true if relPath (slash-separated, relative to scan root) matches ignore rules.
func (m *Matcher) ShouldSkip(relPath string, isDir bool) bool {
	if m == nil || m.ignore == nil {
		return false
	}
	p := filepath.ToSlash(relPath)
	if p == "." || p == "" {
		return false
	}
	// go-gitignore expects paths without leading ./
	p = strings.TrimPrefix(p, "./")
	return m.ignore.MatchesPath(p)
}
