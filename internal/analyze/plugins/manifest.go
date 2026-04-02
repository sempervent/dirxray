package plugins

import (
	"path/filepath"
	"strings"

	"dirxray/internal/model"
)

type manifestPlugin struct{}

func NewManifestPlugin() Plugin { return manifestPlugin{} }

func (manifestPlugin) Name() string { return "manifest" }

func (manifestPlugin) Run(ctx *Context) (*model.PluginResult, error) {
	pr := &model.PluginResult{PluginName: "manifest", NodeBadges: map[string][]string{}}
	if ctx.Scan == nil || ctx.Scan.Root == nil {
		return pr, nil
	}
	manifestNames := map[string]string{
		"go.mod": "go", "go.work": "go",
		"package.json": "node", "package-lock.json": "node", "yarn.lock": "node", "pnpm-lock.yaml": "node",
		"pyproject.toml": "python", "setup.py": "python", "requirements.txt": "python", "pipfile": "python",
		"cargo.toml": "rust", "composer.json": "php",
		"dockerfile": "docker", "containerfile": "docker",
		"docker-compose.yml": "compose", "docker-compose.yaml": "compose",
		"makefile": "make", "rakefile": "ruby",
		"mkdocs.yml": "docs", "docusaurus.config.js": "docs",
	}
	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind == model.NodeDir {
			return true
		}
		base := strings.ToLower(filepath.Base(n.Name))
		if base == "dockerfile" || base == "containerfile" {
			rel := n.Path
			if rel == "" {
				rel = "."
			}
			pr.NodeBadges[rel] = append(pr.NodeBadges[rel], "docker")
			pr.Signals = append(pr.Signals, model.ProjectSignal{
				Key: "dockerfile", Weight: 0.35, Description: "Docker build file present", Paths: []string{rel},
			})
			return true
		}
		if stack, ok := manifestNames[base]; ok {
			rel := n.Path
			if rel == "" {
				rel = "."
			}
			pr.NodeBadges[rel] = append(pr.NodeBadges[rel], "manifest:"+stack)
			pr.Signals = append(pr.Signals, model.ProjectSignal{
				Key: base, Weight: 0.4, Description: "manifest " + base, Paths: []string{rel},
			})
		}
		return true
	})
	return pr, nil
}
