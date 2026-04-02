package plugins

import (
	"path/filepath"
	"strings"

	"dirxray/internal/git"
	"dirxray/internal/model"
)

type archetypePlugin struct{}

func NewArchetypePlugin() Plugin { return archetypePlugin{} }

func (archetypePlugin) Name() string { return "archetype" }

func (archetypePlugin) Run(ctx *Context) (*model.PluginResult, error) {
	pr := &model.PluginResult{PluginName: "archetype"}
	if ctx.Scan == nil || ctx.Scan.Root == nil {
		return pr, nil
	}

	scores := map[model.ArchetypeID]float64{}
	evidence := map[model.ArchetypeID][]model.EvidenceItem{}
	add := func(id model.ArchetypeID, w float64, label, detail, p string) {
		scores[id] += w
		evidence[id] = append(evidence[id], model.EvidenceItem{Label: label, Detail: detail, Path: p})
	}

	if git.IsRepoRoot(ctx.RootAbs) {
		add(model.ArchetypeGitRepo, 0.5, "git", ".git present at scan root", "")
	}

	var hasGo, hasPy, hasNode, hasDocker, hasDocs, dataScore float64
	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind == model.NodeDir {
			l := strings.ToLower(n.Name)
			if l == "k8s" || l == "kubernetes" || l == "helm" || strings.HasSuffix(l, "charts") {
				add(model.ArchetypeKubernetes, 0.25, "dir", "infra-style directory: "+n.Name, n.Path)
			}
			if l == "docs" || l == "documentation" {
				add(model.ArchetypeDocsSite, 0.15, "dir", "documentation directory", n.Path)
			}
			if l == "data" || l == "datasets" || l == "parquet" || l == "csv" {
				dataScore += 0.1
			}
			return true
		}
		base := strings.ToLower(filepath.Base(n.Name))
		ext := strings.ToLower(n.Ext)
		switch base {
		case "go.mod", "go.sum":
			hasGo += 0.45
			add(model.ArchetypeGo, 0.45, "file", base, n.Path)
		case "package.json", "tsconfig.json":
			hasNode += 0.45
			add(model.ArchetypeNode, 0.45, "file", base, n.Path)
		case "pyproject.toml", "setup.py", "requirements.txt":
			hasPy += 0.45
			add(model.ArchetypePython, 0.45, "file", base, n.Path)
		case "dockerfile", "containerfile", ".dockerignore":
			hasDocker += 0.35
			add(model.ArchetypeDocker, 0.35, "file", base, n.Path)
		case "mkdocs.yml", "docusaurus.config.js", "vitepress.config.ts":
			hasDocs += 0.35
			add(model.ArchetypeDocsSite, 0.35, "file", base, n.Path)
		}
		switch ext {
		case ".yaml", ".yml":
			if strings.Contains(strings.ToLower(n.Path), "kustomization") ||
				strings.Contains(base, "chart") {
				add(model.ArchetypeKubernetes, 0.2, "file", "yaml manifest hint", n.Path)
			}
		case ".csv", ".parquet", ".tsv", ".jsonl":
			dataScore += 0.15
			add(model.ArchetypeDataDir, 0.15, "file", "data extension "+ext, n.Path)
		}
		return true
	})

	if hasGo > 0 {
		scores[model.ArchetypeGo] += hasGo
	}
	if hasPy > 0 {
		scores[model.ArchetypePython] += hasPy
	}
	if hasNode > 0 {
		scores[model.ArchetypeNode] += hasNode
	}
	if hasDocker > 0 {
		scores[model.ArchetypeDocker] += hasDocker
	}
	if hasDocs > 0 {
		scores[model.ArchetypeDocsSite] += hasDocs
	}
	if dataScore > 0 {
		scores[model.ArchetypeDataDir] += dataScore
	}

	// mixed / junk: many extensions, no strong manifest
	strong := scores[model.ArchetypeGo] + scores[model.ArchetypePython] + scores[model.ArchetypeNode]
	extCount := countExtensions(ctx)
	if strong < 0.3 && extCount > 8 && ctx.Scan.Stats.Files > 30 {
		add(model.ArchetypeMixedJunk, 0.35, "heuristic", "many file types, weak manifest signals", "")
	}

	if len(scores) == 0 {
		add(model.ArchetypeGeneric, 0.3, "fallback", "no strong archetype markers", "")
	}

	var ranked []model.ArchetypeScore
	for id, s := range scores {
		if s > 1 {
			s = 1
		}
		ranked = append(ranked, model.ArchetypeScore{
			ID:          id,
			Score:       s,
			Explanation: string(id) + " confidence from path/manifest/layout heuristics",
			Evidence:    evidence[id],
		})
	}
	// sort ranked by score desc in merge; plugin returns unsorted slice
	pr.Archetypes = ranked
	return pr, nil
}

func countExtensions(ctx *Context) int {
	m := map[string]struct{}{}
	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind != model.NodeDir && n.Ext != "" {
			m[n.Ext] = struct{}{}
		}
		return true
	})
	return len(m)
}
