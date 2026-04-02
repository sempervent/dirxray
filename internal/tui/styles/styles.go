package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	Muted    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	TabActive = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Padding(0, 1)
	Tab       = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Padding(0, 1)
	Badge     = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Padding(0, 1)
	SeverityHigh = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	SeverityMed  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	SeverityLow  = lipgloss.NewStyle().Foreground(lipgloss.Color("222"))
	Footer   = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
)
