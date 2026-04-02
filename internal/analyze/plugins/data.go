package plugins

import (
	"context"

	"dirxray/internal/data"
	"dirxray/internal/model"
)

type dataPlugin struct{}

func NewDataPlugin() Plugin { return dataPlugin{} }

func (dataPlugin) Name() string { return "data" }

func (dataPlugin) Run(ctx *Context) (*model.PluginResult, error) {
	pr := &model.PluginResult{PluginName: "data"}
	if ctx.Scan == nil {
		return pr, nil
	}
	ds := data.Analyze(context.Background(), ctx.RootAbs, ctx.Scan, 48)
	pr.Data = ds
	if ds != nil && ds.IsDataHeavy {
		pr.Signals = append(pr.Signals, model.ProjectSignal{
			Key: "data_files", Weight: 0.35, Description: "directory contains tabular/parquet/jsonl style files",
		})
	}
	return pr, nil
}
