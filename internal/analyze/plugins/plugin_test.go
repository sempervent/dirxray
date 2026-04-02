package plugins

import (
	"path/filepath"
	"runtime"
	"testing"

	"dirxray/internal/scan"
)

func fixture(parts ...string) string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	return filepath.Join(dir, "..", "..", "..", "testdata", filepath.Join(parts...))
}

func TestRegistryRuns(t *testing.T) {
	root := fixture("fixture-docker")
	sc, err := scan.Scan(scan.Options{Root: root, NoGitignore: true})
	if err != nil {
		t.Fatal(err)
	}
	reg := NewRegistry()
	ctx := &Context{RootAbs: root, Scan: sc}
	for _, p := range reg.All() {
		pr, err := p.Run(ctx)
		if err != nil {
			t.Fatalf("%s: %v", p.Name(), err)
		}
		if pr == nil || pr.PluginName == "" {
			t.Fatalf("empty result from %s", p.Name())
		}
	}
}
