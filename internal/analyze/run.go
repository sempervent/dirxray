package analyze

import (
	"dirxray/internal/analyze/plugins"
	"dirxray/internal/model"
)

// Run executes the default registry against a scan.
func Run(rootAbs string, scan *model.ScanResult) (*model.AnalysisResult, error) {
	reg := plugins.NewRegistry()
	ctx := &plugins.Context{RootAbs: rootAbs, Scan: scan}
	var parts []model.PluginResult
	for _, p := range reg.All() {
		pr, err := p.Run(ctx)
		if err != nil {
			return nil, err
		}
		if pr == nil {
			pr = &model.PluginResult{PluginName: p.Name()}
		}
		parts = append(parts, *pr)
	}
	return Merge(scan, parts), nil
}
