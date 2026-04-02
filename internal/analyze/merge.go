package analyze

import (
	"sort"
	"strings"

	"dirxray/internal/model"
)

// Merge combines plugin outputs and builds summary + evidence index.
func Merge(scan *model.ScanResult, parts []model.PluginResult) *model.AnalysisResult {
	ar := &model.AnalysisResult{
		EvidenceIndex: map[string][]model.EvidenceItem{},
	}
	var data *model.DataSummary
	for _, pr := range parts {
		ar.Signals = append(ar.Signals, pr.Signals...)
		ar.Archetypes = append(ar.Archetypes, pr.Archetypes...)
		ar.Findings = append(ar.Findings, pr.Findings...)
		if pr.Data != nil {
			if data == nil {
				data = pr.Data
			} else {
				// prefer heavier data summary
				if len(pr.Data.Files) > len(data.Files) {
					data = pr.Data
				}
			}
		}
		ar.PluginResults = append(ar.PluginResults, pr)
		applyNodeBadges(scan, pr.NodeBadges)
	}
	ar.Data = data

	sort.Slice(ar.Archetypes, func(i, j int) bool {
		return ar.Archetypes[i].Score > ar.Archetypes[j].Score
	})

	sort.Slice(ar.Findings, func(i, j int) bool {
		if ar.Findings[i].Severity != ar.Findings[j].Severity {
			return ar.Findings[i].Severity > ar.Findings[j].Severity
		}
		return ar.Findings[i].Title < ar.Findings[j].Title
	})

	ar.Summary = buildSummary(scan, ar)
	indexEvidence(ar)
	return ar
}

func applyNodeBadges(scan *model.ScanResult, m map[string][]string) {
	if scan == nil || scan.Root == nil || len(m) == 0 {
		return
	}
	model.Walk(scan.Root, func(n *model.Node) bool {
		key := n.Path
		if key == "" {
			key = "."
		}
		if bs, ok := m[key]; ok {
			n.Badges = append(n.Badges, bs...)
		}
		return true
	})
}

func buildSummary(scan *model.ScanResult, ar *model.AnalysisResult) model.DirectorySummary {
	ds := model.DirectorySummary{
		DominantExts: map[string]int64{},
	}
	if scan != nil && scan.Root != nil {
		model.Walk(scan.Root, func(n *model.Node) bool {
			if n.Kind != model.NodeDir && n.Ext != "" {
				ds.DominantExts[n.Ext] += n.Size
			}
			return true
		})
	}
	for _, s := range ar.Signals {
		if s.Description != "" {
			ds.SignalSummary = append(ds.SignalSummary, s.Description)
		}
	}
	ds.ArchetypeScores = ar.Archetypes
	if len(ar.Archetypes) > 0 {
		top := 3
		if len(ar.Archetypes) < top {
			top = len(ar.Archetypes)
		}
		for i := 0; i < top; i++ {
			ds.PrimaryArchetypes = append(ds.PrimaryArchetypes, ar.Archetypes[i].ID)
		}
		ds.ProbablePurpose = string(ar.Archetypes[0].ID)
	} else {
		ds.ProbablePurpose = "unknown"
	}

	// Entry points: shallow manifests / readme
	if scan != nil && scan.Root != nil {
		var eps []string
		model.Walk(scan.Root, func(n *model.Node) bool {
			if n.Kind == model.NodeDir {
				return true
			}
			depth := strings.Count(n.Path, "/")
			if depth > 2 {
				return true
			}
			l := strings.ToLower(n.Name)
			if l == "readme.md" || l == "readme" || l == "go.mod" || l == "package.json" || l == "pyproject.toml" || l == "dockerfile" {
				eps = append(eps, n.Path)
			}
			return true
		})
		sort.Strings(eps)
		if len(eps) > 8 {
			eps = eps[:8]
		}
		ds.TopEntryPoints = eps
	}
	return ds
}

func indexEvidence(ar *model.AnalysisResult) {
	for _, f := range ar.Findings {
		for _, e := range f.Evidence {
			if e.Path != "" {
				ar.EvidenceIndex[e.Path] = append(ar.EvidenceIndex[e.Path], e)
			}
		}
		for _, p := range f.RelatedPaths {
			ar.EvidenceIndex[p] = append(ar.EvidenceIndex[p], model.EvidenceItem{
				Label: "finding", Detail: f.Title, Path: p,
			})
		}
	}
	for _, a := range ar.Archetypes {
		for _, e := range a.Evidence {
			if e.Path != "" {
				ar.EvidenceIndex[e.Path] = append(ar.EvidenceIndex[e.Path], e)
			}
		}
	}
}
