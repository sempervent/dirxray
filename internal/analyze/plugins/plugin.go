package plugins

import (
	"dirxray/internal/model"
)

// Context is input for plugins (immutable scan snapshot + root path).
type Context struct {
	RootAbs string
	Scan    *model.ScanResult
}

// Plugin analyzes a scan result and returns partial results.
type Plugin interface {
	Name() string
	Run(ctx *Context) (*model.PluginResult, error)
}

// Registry holds ordered plugins.
type Registry struct {
	plugins []Plugin
}

// NewRegistry returns the default Phase-1 plugin set.
func NewRegistry() *Registry {
	return &Registry{
		plugins: []Plugin{
			NewStructurePlugin(),
			NewManifestPlugin(),
			NewArchetypePlugin(),
			NewDataPlugin(),
			NewSuspiciousPlugin(),
		},
	}
}

// All returns registered plugins (copy).
func (r *Registry) All() []Plugin {
	out := make([]Plugin, len(r.plugins))
	copy(out, r.plugins)
	return out
}

// Register appends a plugin (mainly for tests).
func (r *Registry) Register(p Plugin) {
	r.plugins = append(r.plugins, p)
}
