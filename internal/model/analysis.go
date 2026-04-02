package model

// PluginResult is one plugin's contribution before merge.
type PluginResult struct {
	PluginName string
	Signals    []ProjectSignal
	Archetypes []ArchetypeScore
	Findings   []Finding
	Data       *DataSummary
	// NodeBadges maps relative path -> badge strings to merge onto nodes.
	NodeBadges map[string][]string
}

// AnalysisResult is the merged Phase-1 analysis graph.
type AnalysisResult struct {
	Signals          []ProjectSignal
	Archetypes       []ArchetypeScore
	Findings         []Finding
	Data             *DataSummary
	Summary          DirectorySummary
	PluginResults    []PluginResult
	EvidenceIndex    map[string][]EvidenceItem // path -> items mentioning it
}
