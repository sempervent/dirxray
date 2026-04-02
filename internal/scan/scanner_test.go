package scan

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdata(parts ...string) string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	return filepath.Join(dir, "..", "..", "testdata", filepath.Join(parts...))
}

func TestScanGoFixture(t *testing.T) {
	root := testdata("fixture-go")
	res, err := Scan(Options{Root: root, NoGitignore: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Root == nil {
		t.Fatal("nil root")
	}
	if res.Stats.Files < 2 {
		t.Fatalf("expected at least 2 files, got %d", res.Stats.Files)
	}
}

func TestScanRespectsMaxDepth(t *testing.T) {
	root := testdata("fixture-python")
	res, err := Scan(Options{Root: root, MaxDepth: 1, NoGitignore: true})
	if err != nil {
		t.Fatal(err)
	}
	// depth 1: root + immediate children only, no src/pkg traversal
	for _, n := range res.Root.Children {
		if n.Name == "src" && len(n.Children) > 0 {
			t.Fatal("max-depth should prevent listing inside src")
		}
	}
}

func TestValidatePathRejectsMissing(t *testing.T) {
	_, err := ValidatePath("/nonexistent/dirxray-path-xyz")
	if err == nil {
		t.Fatal("expected error")
	}
}
