package model

import "time"

// ScanStats summarizes inventory counts.
type ScanStats struct {
	Dirs        int
	Files       int
	Symlinks    int
	TotalBytes  int64
	MaxDepth    int
	HiddenCount int
	Skipped     int // ignored by rules
}

// ScanNotice is a non-fatal scan issue.
type ScanNotice struct {
	Path    string
	Kind    string // permission_denied, symlink_skipped, loop, other
	Message string
}

// ScanResult is the output of the inventory engine.
type ScanResult struct {
	RootAbsPath string
	Root        *Node
	Stats       ScanStats
	Notices     []ScanNotice
	StartedAt   time.Time
	FinishedAt  time.Time
}

// Walk visits every node in pre-order.
func Walk(n *Node, fn func(*Node) bool) {
	if n == nil {
		return
	}
	if !fn(n) {
		return
	}
	for _, c := range n.Children {
		Walk(c, fn)
	}
}
