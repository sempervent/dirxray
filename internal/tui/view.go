package tui

import (
	"fmt"
	"strings"

	"dirxray/internal/model"
	"dirxray/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model.
func (m *Model) View() string {
	if m.width <= 0 {
		m.width = 80
	}
	if m.height <= 0 {
		m.height = 24
	}

	tabs := m.renderTabs()
	body := m.renderBody()
	footer := m.renderFooter()

	main := lipgloss.JoinVertical(lipgloss.Left, tabs, body)
	if m.helpOpen {
		help := styles.Muted.Render("\n" + helpText())
		main = lipgloss.JoinVertical(lipgloss.Left, main, help)
	}
	return lipgloss.JoinVertical(lipgloss.Left, main, footer)
}

func (m *Model) renderTabs() string {
	names := []string{"1:Overview", "2:Tree", "3:Findings", "4:Evidence", "5:Data"}
	var parts []string
	for i, n := range names {
		tid := tabID(i)
		if tid == m.tab {
			parts = append(parts, styles.TabActive.Render(n))
		} else {
			parts = append(parts, styles.Tab.Render(n))
		}
	}
	return lipgloss.NewStyle().Width(m.width).Render(strings.Join(parts, ""))
}

func (m *Model) renderBody() string {
	switch m.tab {
	case tabTree:
		return m.renderTreePane()
	case tabFindings:
		return m.renderFindingsPane()
	default:
		return m.mainViewport.View()
	}
}

func (m *Model) renderTreePane() string {
	var b strings.Builder
	for i, row := range m.treeRows {
		line := row.Line
		if i == m.treeCursor {
			line = styles.TabActive.Render(line)
		} else {
			line = styles.Muted.Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	if len(m.treeRows) == 0 {
		b.WriteString(styles.Muted.Render("(empty tree)"))
	}
	return lipgloss.NewStyle().Width(m.width).Height(m.bodyHeight()).Render(strings.TrimRight(b.String(), "\n"))
}

func (m *Model) renderFindingsPane() string {
	if m.analysis == nil || len(m.analysis.Findings) == 0 {
		return lipgloss.NewStyle().Width(m.width).Height(m.bodyHeight()).Render(styles.Muted.Render("No findings."))
	}
	var b strings.Builder
	for i, f := range m.analysis.Findings {
		sev := severityStyle(f.Severity).Render(f.Severity.String())
		line := fmt.Sprintf("%s  %s", sev, f.Title)
		if i == m.findCursor {
			line = lipgloss.NewStyle().Reverse(true).Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return lipgloss.NewStyle().Width(m.width).Height(m.bodyHeight()).Render(strings.TrimRight(b.String(), "\n"))
}

func severityStyle(s model.Severity) lipgloss.Style {
	switch s {
	case model.SeverityHigh:
		return styles.SeverityHigh
	case model.SeverityMedium:
		return styles.SeverityMed
	default:
		return styles.SeverityLow
	}
}

func (m *Model) bodyHeight() int {
	h := m.height - 4
	if m.helpOpen {
		h -= 8
	}
	if h < 6 {
		h = 6
	}
	return h
}

func (m *Model) syncViewport() {
	vw := m.width
	if vw < 20 {
		vw = 80
	}
	vh := m.bodyHeight()
	m.mainViewport.Width = vw
	m.mainViewport.Height = vh
	m.mainViewport.SetContent(m.viewportContent())
}

func (m *Model) viewportContent() string {
	switch m.tab {
	case tabOverview:
		return m.renderOverview()
	case tabEvidence:
		return m.renderEvidence()
	case tabData:
		return m.renderData()
	default:
		return ""
	}
}

func (m *Model) renderOverview() string {
	var b strings.Builder
	b.WriteString(styles.Title.Render("Directory overview"))
	b.WriteString("\n\n")
	b.WriteString(styles.Muted.Render("Root: "))
	b.WriteString(m.rootAbs)
	b.WriteString("\n\n")

	if m.scan != nil {
		st := m.scan.Stats
		b.WriteString(fmt.Sprintf("Files: %d  Dirs: %d  Symlinks: %d  Bytes: %d  Skipped(ignore): %d\n",
			st.Files, st.Dirs, st.Symlinks, st.TotalBytes, st.Skipped))
		if len(m.scan.Notices) > 0 {
			b.WriteString(styles.SeverityMed.Render(fmt.Sprintf("Scan notices: %d\n", len(m.scan.Notices))))
		}
	}
	b.WriteString("\n")

	if m.analysis != nil {
		sum := m.analysis.Summary
		b.WriteString(styles.Title.Render("Probable purpose / archetypes"))
		b.WriteString("\n")
		b.WriteString(sum.ProbablePurpose)
		b.WriteString("\n\n")
		for _, a := range sum.ArchetypeScores {
			if a.Score <= 0 {
				continue
			}
			b.WriteString(fmt.Sprintf("  • %-22s  score=%.2f  %s\n", string(a.ID), a.Score, a.Explanation))
		}
		b.WriteString("\n")
		if len(sum.TopEntryPoints) > 0 {
			b.WriteString(styles.Title.Render("Top entry points"))
			b.WriteString("\n")
			for _, e := range sum.TopEntryPoints {
				b.WriteString("  • " + e + "\n")
			}
			b.WriteString("\n")
		}
		if len(sum.SignalSummary) > 0 {
			b.WriteString(styles.Title.Render("Signals"))
			b.WriteString("\n")
			for _, s := range sum.SignalSummary {
				if s == "" {
					continue
				}
				b.WriteString("  • " + s + "\n")
			}
			b.WriteString("\n")
		}
		if len(m.analysis.Findings) > 0 {
			b.WriteString(styles.Title.Render("Top findings"))
			b.WriteString("\n")
			max := 6
			if len(m.analysis.Findings) < max {
				max = len(m.analysis.Findings)
			}
			for i := 0; i < max; i++ {
				f := m.analysis.Findings[i]
				b.WriteString(fmt.Sprintf("  • [%s] %s\n", f.Severity.String(), f.Title))
			}
		}
	}
	return b.String()
}

func (m *Model) renderEvidence() string {
	var b strings.Builder
	b.WriteString(styles.Title.Render("Evidence"))
	b.WriteString("\n\n")
	path := m.selectedPath
	if path == "" {
		path = "."
	}
	b.WriteString(styles.Muted.Render("Selected path: "))
	b.WriteString(path)
	b.WriteString("\n\n")

	if m.analysis == nil {
		return b.String()
	}
	if ev, ok := m.analysis.EvidenceIndex[path]; ok && len(ev) > 0 {
		for _, e := range ev {
			b.WriteString(styles.Badge.Render(e.Label))
			b.WriteString(" ")
			b.WriteString(e.Detail)
			if e.Path != "" && e.Path != path {
				b.WriteString(styles.Muted.Render("  (" + e.Path + ")"))
			}
			b.WriteString("\n")
		}
	} else {
		b.WriteString(styles.Muted.Render("No indexed evidence for this path. Select a tree row (Enter) or pick a finding.\n"))
	}

	b.WriteString("\n")
	b.WriteString(styles.Title.Render("Archetype evidence (global)"))
	b.WriteString("\n")
	for _, a := range m.analysis.Archetypes {
		if len(a.Evidence) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("[%s]\n", a.ID))
		for _, e := range a.Evidence {
			b.WriteString(fmt.Sprintf("  - %s: %s\n", e.Label, e.Detail))
		}
	}
	return b.String()
}

func (m *Model) renderData() string {
	var b strings.Builder
	b.WriteString(styles.Title.Render("Data"))
	b.WriteString("\n\n")
	if m.analysis == nil || m.analysis.Data == nil {
		b.WriteString(styles.Muted.Render("No data-oriented summary."))
		return b.String()
	}
	d := m.analysis.Data
	b.WriteString(d.DuckDBExplain)
	b.WriteString("\n\n")
	if !d.IsDataHeavy {
		b.WriteString(styles.Muted.Render("No strong data-file signals in this scan."))
		return b.String()
	}
	b.WriteString(fmt.Sprintf("Total data-ish bytes (tracked extensions): %d\n", d.TotalDataBytes))
	for _, note := range d.LayoutNotes {
		b.WriteString("  • " + note + "\n")
	}
	b.WriteString("\n")
	for _, p := range d.Files {
		b.WriteString(fmt.Sprintf("— %s  kind=%s  size=%d\n", p.Path, p.Kind, p.SizeBytes))
		if p.PartitionHint != "" {
			b.WriteString(fmt.Sprintf("    partition hint: %s\n", p.PartitionHint))
		}
		if len(p.ColumnNames) > 0 {
			names := p.ColumnNames
			if len(names) > 12 {
				names = names[:12]
			}
			b.WriteString(fmt.Sprintf("    columns (%d): %s\n", p.ColumnHint, strings.Join(names, ", ")))
		}
		if p.SampleNote != "" {
			b.WriteString(styles.Muted.Render("    " + p.SampleNote + "\n"))
		}
		if p.DuckDBError != "" {
			b.WriteString(styles.SeverityLow.Render("    duckdb: " + p.DuckDBError + "\n"))
		}
	}
	return b.String()
}

func (m *Model) renderFooter() string {
	line := "q quit  ? help  1-5 tabs  j/k move  Enter→Evidence  PgUp/PgDn page (overview/evidence/data)"
	return styles.Footer.Width(m.width).Render(line)
}

func helpText() string {
	return `dirxray — Phase 1 heuristic directory explainer (no LLM).

Views
  Overview — archetypes, signals, scan stats, top findings.
  Tree — annotated inventory; Enter opens Evidence for that path.
  Findings — ranked list; Enter attaches first related path and opens Evidence.
  Evidence — why the engine believes what it believes for the selected path.
  Data — CSV/TSV/Parquet/JSONL hints; DuckDB CLI optional for schema/peek.

This tool is rules-based; scores are relative heuristics, not probabilities.
`
}
