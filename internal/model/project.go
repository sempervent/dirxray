package model

// ProjectSignal is a weighted hint toward an archetype or concern.
type ProjectSignal struct {
	Key         string  // e.g. "go.mod", "package.json"
	Weight      float64 // 0–1 contribution
	Description string
	Paths       []string
}

// ArchetypeID identifies a heuristic project/directory class.
type ArchetypeID string

const (
	ArchetypeGitRepo       ArchetypeID = "git_repo"
	ArchetypeGo            ArchetypeID = "go"
	ArchetypePython        ArchetypeID = "python"
	ArchetypeNode          ArchetypeID = "node_js_ts"
	ArchetypeDocker        ArchetypeID = "docker"
	ArchetypeKubernetes    ArchetypeID = "kubernetes_iac"
	ArchetypeDocsSite      ArchetypeID = "docs_site"
	ArchetypeDataDir       ArchetypeID = "data_directory"
	ArchetypeMixedJunk     ArchetypeID = "mixed_junk"
	ArchetypeGeneric       ArchetypeID = "generic_directory"
)

// ArchetypeScore is a ranked inference with evidence.
type ArchetypeScore struct {
	ID          ArchetypeID
	Score       float64 // 0–1
	Explanation string
	Evidence    []EvidenceItem
}

// DirectorySummary is the high-level narrative for the overview pane.
type DirectorySummary struct {
	ProbablePurpose   string
	TopEntryPoints    []string
	DominantExts      map[string]int64 // ext -> total bytes
	SignalSummary     []string
	ArchetypeScores   []ArchetypeScore
	PrimaryArchetypes []ArchetypeID
}
