package tui

import (
	"strings"
	"testing"

	"dirxray/internal/model"
)

func TestFlattenTree(t *testing.T) {
	root := &model.Node{Name: ".", Kind: model.NodeDir, Path: ""}
	root.Children = []*model.Node{
		{Name: "a.go", Kind: model.NodeFile, Path: "a.go", Badges: []string{"manifest:go"}},
		{Name: "cmd", Kind: model.NodeDir, Path: "cmd"},
	}
	root.Children[1].Children = []*model.Node{
		{Name: "main.go", Kind: model.NodeFile, Path: "cmd/main.go"},
	}
	rows := flattenTree(root)
	if len(rows) != 4 {
		t.Fatalf("want 4 rows (root + 2 files + cmd), got %d", len(rows))
	}
	if !strings.Contains(rows[1].Line, "a.go") {
		t.Fatalf("expected a.go in row 1, line: %q", rows[1].Line)
	}
}
