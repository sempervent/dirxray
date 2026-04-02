package plugins

import (
	"sort"
	"strings"

	"dirxray/internal/model"
)

type structurePlugin struct{}

func NewStructurePlugin() Plugin { return structurePlugin{} }

func (structurePlugin) Name() string { return "structure" }

func (structurePlugin) Run(ctx *Context) (*model.PluginResult, error) {
	if ctx.Scan == nil || ctx.Scan.Root == nil {
		return &model.PluginResult{PluginName: "structure"}, nil
	}
	extBytes := map[string]int64{}
	var topPaths []string
	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind == model.NodeFile || n.Kind == model.NodeSymlink {
			if n.Ext != "" {
				extBytes[n.Ext] += n.Size
			}
			depth := strings.Count(n.Path, "/")
			if depth <= 2 && n.Size > 0 {
				topPaths = append(topPaths, n.Path)
			}
		}
		return true
	})
	sort.Slice(topPaths, func(i, j int) bool {
		return topPaths[i] < topPaths[j]
	})
	if len(topPaths) > 12 {
		topPaths = topPaths[:12]
	}

	var signals []model.ProjectSignal
	if len(extBytes) > 0 {
		type kv struct {
			k string
			v int64
		}
		var pairs []kv
		for k, v := range extBytes {
			pairs = append(pairs, kv{k, v})
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].v > pairs[j].v })
		top := pairs[0].k
		signals = append(signals, model.ProjectSignal{
			Key:         "dominant_extension",
			Weight:      0.2,
			Description: "largest share of bytes by extension: " + top,
		})
	}

	pr := &model.PluginResult{
		PluginName: "structure",
		Signals:    signals,
		NodeBadges: map[string][]string{},
	}
	// Depth and role badges on shallow manifest-like files
	model.Walk(ctx.Scan.Root, func(n *model.Node) bool {
		if n.Kind != model.NodeDir && n.Role != model.RoleUnknown {
			b := string(n.Role)
			rel := n.Path
			if rel == "" {
				rel = "."
			}
			pr.NodeBadges[rel] = append(pr.NodeBadges[rel], b)
		}
		return true
	})
	return pr, nil
}
