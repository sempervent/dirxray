package tui

import (
	"fmt"
	"os"

	"dirxray/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the Bubble Tea program (alt-screen TUI).
func Run(rootAbs string, scan *model.ScanResult, analysis *model.AnalysisResult) error {
	if scan == nil {
		return fmt.Errorf("nil scan result")
	}
	m := New(rootAbs, scan, analysis)
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return err
	}
	if final == nil {
		return nil
	}
	// Allow future: print summary on exit with a flag
	_ = final
	return nil
}

// RunStdoutIfTTY runs the TUI; if stdout is not a terminal, prints a plain summary.
func RunStdoutIfTTY(rootAbs string, scan *model.ScanResult, analysis *model.AnalysisResult) error {
	if f, ok := any(os.Stdout).(*os.File); ok {
		if st, err := f.Stat(); err == nil && st.Mode()&os.ModeCharDevice == 0 {
			// Not a TTY — plain text fallback
			printPlain(rootAbs, scan, analysis)
			return nil
		}
		_ = f
	}
	return Run(rootAbs, scan, analysis)
}

func printPlain(rootAbs string, scan *model.ScanResult, analysis *model.AnalysisResult) {
	fmt.Printf("dirxray (non-TTY output)\nRoot: %s\n", rootAbs)
	if scan != nil {
		s := scan.Stats
		fmt.Printf("files=%d dirs=%d bytes=%d\n", s.Files, s.Dirs, s.TotalBytes)
	}
	if analysis != nil {
		fmt.Printf("top archetype: %s\n", analysis.Summary.ProbablePurpose)
		for _, f := range analysis.Findings {
			fmt.Printf("[%s] %s\n", f.Severity.String(), f.Title)
		}
	}
}
