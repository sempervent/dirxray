package model

// Severity ranks finding importance.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	default:
		return "unknown"
	}
}

// Finding is an evidence-backed issue or highlight.
type Finding struct {
	ID            string
	Title         string
	Severity      Severity
	Rationale     string
	Evidence      []EvidenceItem
	RelatedPaths  []string
}

// EvidenceItem ties a claim to observable facts.
type EvidenceItem struct {
	Label  string
	Detail string
	Path   string
}
