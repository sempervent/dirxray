package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			m.helpOpen = !m.helpOpen
			m.syncViewport()
		case "1":
			m.setTab(tabOverview)
		case "2":
			m.setTab(tabTree)
		case "3":
			m.setTab(tabFindings)
		case "4":
			m.setTab(tabEvidence)
		case "5":
			m.setTab(tabData)
		case "up", "k":
			m.move(-1)
		case "down", "j":
			m.move(1)
		case "pgup":
			if m.tab == tabOverview || m.tab == tabEvidence || m.tab == tabData {
				m.mainViewport.ViewUp()
			}
		case "pgdown":
			if m.tab == tabOverview || m.tab == tabEvidence || m.tab == tabData {
				m.mainViewport.ViewDown()
			}
		case "enter":
			if m.tab == tabTree && m.treeCursor >= 0 && m.treeCursor < len(m.treeRows) {
				n := m.treeRows[m.treeCursor].Node
				if n != nil {
					m.selectedPath = n.Path
					if m.selectedPath == "" {
						m.selectedPath = "."
					}
					m.setTab(tabEvidence)
				}
			} else if m.tab == tabFindings && m.analysis != nil && len(m.analysis.Findings) > 0 {
				f := m.analysis.Findings[m.findCursor]
				if len(f.RelatedPaths) > 0 {
					m.selectedPath = f.RelatedPaths[0]
				} else {
					for _, e := range f.Evidence {
						if e.Path != "" {
							m.selectedPath = e.Path
							break
						}
					}
				}
				if m.selectedPath == "" {
					m.selectedPath = "."
				}
				m.setTab(tabEvidence)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncViewport()
	}
	return m, nil
}

func (m *Model) setTab(t tabID) {
	m.tab = t
	m.syncViewport()
}

func (m *Model) move(delta int) {
	switch m.tab {
	case tabTree:
		m.treeCursor += delta
		if m.treeCursor < 0 {
			m.treeCursor = 0
		}
		if m.treeCursor >= len(m.treeRows) && len(m.treeRows) > 0 {
			m.treeCursor = len(m.treeRows) - 1
		}
	case tabFindings:
		nf := 0
		if m.analysis != nil {
			nf = len(m.analysis.Findings)
		}
		m.findCursor += delta
		if m.findCursor < 0 {
			m.findCursor = 0
		}
		if m.findCursor >= nf && nf > 0 {
			m.findCursor = nf - 1
		}
	case tabOverview, tabEvidence, tabData:
		if delta < 0 {
			m.mainViewport.LineUp(-delta)
		} else {
			m.mainViewport.LineDown(delta)
		}
	}
}
