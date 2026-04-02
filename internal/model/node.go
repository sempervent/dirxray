package model

import (
	"time"
)

// NodeKind classifies inventory entries.
type NodeKind int

const (
	NodeFile NodeKind = iota
	NodeDir
	NodeSymlink
)

// RoleHint is a coarse label from the scanner or analyzers.
type RoleHint string

const (
	RoleUnknown     RoleHint = ""
	RoleConfig      RoleHint = "config"
	RoleManifest    RoleHint = "manifest"
	RoleLockfile    RoleHint = "lockfile"
	RoleData        RoleHint = "data"
	RoleDocs        RoleHint = "docs"
	RoleBuild       RoleHint = "build"
	RoleVCS         RoleHint = "vcs"
	RoleBinary      RoleHint = "binary"
	RoleGenerated   RoleHint = "generated_hint"
)

// Node is a file, directory, or symlink in the scanned tree.
type Node struct {
	Path   string // relative to scan root, always slash-separated for cross-platform display
	Name   string
	Kind   NodeKind
	Size   int64
	ModTime time.Time
	Depth  int
	Ext    string // lowercase, includes dot for non-empty

	Children []*Node
	Parent   *Node

	IsHidden  bool
	IsSymlink bool
	// SymlinkTarget is set when IsSymlink; may be empty if read failed.
	SymlinkTarget string

	// ScanErr is a non-fatal error (e.g. permission denied listing dir).
	ScanErr error

	// Analyzer annotations (merged from plugins)
	Badges []string // short labels: archetype signals, suspicion, importance
	Role   RoleHint
}

// IsDir reports whether this node is a directory.
func (n *Node) IsDir() bool {
	return n != nil && n.Kind == NodeDir
}

// FlatFiles returns all file nodes under n (depth-first, excluding symlinks to dirs as files).
func (n *Node) FlatFiles() []*Node {
	if n == nil {
		return nil
	}
	var out []*Node
	var walk func(*Node)
	walk = func(x *Node) {
		if x == nil {
			return
		}
		switch x.Kind {
		case NodeFile, NodeSymlink:
			out = append(out, x)
		case NodeDir:
			for _, c := range x.Children {
				walk(c)
			}
		}
	}
	walk(n)
	return out
}
