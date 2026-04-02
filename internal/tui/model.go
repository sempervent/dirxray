package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"dirxray/internal/model"
)

type tabID int

const (
	tabOverview tabID = iota
	tabTree
	tabFindings
	tabEvidence
	tabData
)

// Model is the Bubble Tea root model.
type Model struct {
	rootAbs  string
	scan     *model.ScanResult
	analysis *model.AnalysisResult

	width  int
	height int
	tab    tabID

	treeRows     []treeRow
	treeCursor   int
	findCursor   int
	selectedPath string

	mainViewport viewport.Model
	helpOpen     bool
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// New builds a Model.
func New(rootAbs string, scan *model.ScanResult, analysis *model.AnalysisResult) *Model {
	m := &Model{
		rootAbs:  rootAbs,
		scan:     scan,
		analysis: analysis,
		tab:      tabOverview,
		width:    80,
		height:   24,
	}
	m.mainViewport = viewport.New(80, 24)
	m.treeRows = flattenTree(scanRoot(scan))
	m.selectedPath = "."
	m.syncViewport()
	return m
}

func scanRoot(s *model.ScanResult) *model.Node {
	if s == nil {
		return nil
	}
	return s.Root
}
