package model

// DataFileKind identifies supported data shapes.
type DataFileKind string

const (
	DataUnknown DataFileKind = "unknown"
	DataCSV     DataFileKind = "csv"
	DataTSV     DataFileKind = "tsv"
	DataParquet DataFileKind = "parquet"
	DataJSON    DataFileKind = "json"
	DataJSONL   DataFileKind = "jsonl"
)

// DataFileProfile is lightweight metadata for one file.
type DataFileProfile struct {
	Path         string
	Kind         DataFileKind
	SizeBytes    int64
	RowHint      int // -1 if unknown
	ColumnHint   int
	ColumnNames  []string
	PartitionHint string // e.g. "date=YYYY-MM-DD" from path
	SampleNote   string
	DuckDBUsed   bool
	DuckDBError  string
}

// DataSummary aggregates data-directory intelligence.
type DataSummary struct {
	IsDataHeavy      bool
	LayoutNotes      []string
	Files            []DataFileProfile
	ExtensionMix     map[string]int
	TotalDataBytes   int64
	DuckDBAvailable  bool
	DuckDBExplain    string
}
