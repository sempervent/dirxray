package analyze

import (
	"path/filepath"
	"runtime"
	"testing"

	"dirxray/internal/scan"
)

func testRoot(parts ...string) string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	return filepath.Join(dir, "..", "..", "testdata", filepath.Join(parts...))
}

func TestRunGoArchetype(t *testing.T) {
	root := testRoot("fixture-go")
	sc, err := scan.Scan(scan.Options{Root: root, NoGitignore: true})
	if err != nil {
		t.Fatal(err)
	}
	ar, err := Run(root, sc)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, a := range ar.Archetypes {
		if a.ID == "go" && a.Score > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected go archetype, got %#v", ar.Archetypes)
	}
}

func TestRunDataFixture(t *testing.T) {
	root := testRoot("fixture-data")
	sc, err := scan.Scan(scan.Options{Root: root, NoGitignore: true})
	if err != nil {
		t.Fatal(err)
	}
	ar, err := Run(root, sc)
	if err != nil {
		t.Fatal(err)
	}
	if ar.Data == nil || !ar.Data.IsDataHeavy {
		t.Fatal("expected data-heavy summary")
	}
}

func TestMergeFindingsNestedGit(t *testing.T) {
	root := testRoot("fixture-nested", "outer")
	sc, err := scan.Scan(scan.Options{Root: root, NoGitignore: true})
	if err != nil {
		t.Fatal(err)
	}
	ar, err := Run(root, sc)
	if err != nil {
		t.Fatal(err)
	}
	var nested bool
	for _, f := range ar.Findings {
		if f.Title == "Nested repository marker" {
			nested = true
			break
		}
	}
	if !nested {
		t.Fatalf("expected nested .git finding, got %d findings", len(ar.Findings))
	}
}
