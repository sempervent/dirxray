package tui

import (
	"fmt"
	"strings"

	"dirxray/internal/model"
)

type treeRow struct {
	Node  *model.Node
	Depth int
	Line  string
}

func flattenTree(root *model.Node) []treeRow {
	if root == nil {
		return nil
	}
	var rows []treeRow
	var walk func(*model.Node, int)
	walk = func(n *model.Node, d int) {
		if n == nil {
			return
		}
		prefix := strings.Repeat("  ", d)
		icon := "[f]"
		if n.Kind == model.NodeDir {
			icon = "[d]"
		} else if n.Kind == model.NodeSymlink {
			icon = "[l]"
		}
		name := n.Name
		if n.Path == "" {
			name = "."
		}
		suffix := ""
		if len(n.Badges) > 0 {
			suffix = " [" + strings.Join(n.Badges, ", ") + "]"
		}
		if n.ScanErr != nil {
			suffix += " (!)"
		}
		line := fmt.Sprintf("%s%s %s%s", prefix, icon, name, suffix)
		rows = append(rows, treeRow{Node: n, Depth: d, Line: line})
		for _, c := range n.Children {
			walk(c, d+1)
		}
	}
	walk(root, 0)
	return rows
}
