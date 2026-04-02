package data

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	"dirxray/internal/data/detect"
	"dirxray/internal/data/duckdb"
	"dirxray/internal/data/sample"
	"dirxray/internal/model"
)

// Analyze builds a DataSummary from the scan tree (bounded file count).
func Analyze(ctx context.Context, rootAbs string, scan *model.ScanResult, maxFiles int) *model.DataSummary {
	ds := &model.DataSummary{
		ExtensionMix:    map[string]int{},
		DuckDBAvailable: duckdb.Available(),
		DuckDBExplain:   duckdb.Explain(),
	}
	if scan == nil || scan.Root == nil {
		return ds
	}
	if maxFiles <= 0 {
		maxFiles = 40
	}

	model.Walk(scan.Root, func(n *model.Node) bool {
		if n.Kind == model.NodeDir {
			return true
		}
		k := detect.KindFromPath(n.Path)
		if k == model.DataUnknown {
			return true
		}
		ds.ExtensionMix[string(k)]++
		ds.TotalDataBytes += n.Size
		return true
	})

	if len(ds.ExtensionMix) == 0 {
		return ds
	}
	ds.IsDataHeavy = true

	var profiles []model.DataFileProfile
	count := 0
	model.Walk(scan.Root, func(n *model.Node) bool {
		if count >= maxFiles {
			return false
		}
		if n.Kind == model.NodeDir {
			return true
		}
		k := detect.KindFromPath(n.Path)
		if k == model.DataUnknown {
			return true
		}
		abs := filepath.Join(rootAbs, filepath.FromSlash(n.Path))
		ph := detect.PartitionHintFromPath(n.Path)
		prof := model.DataFileProfile{
			Path:          n.Path,
			Kind:          k,
			SizeBytes:     n.Size,
			PartitionHint: ph,
		}
		switch k {
		case model.DataCSV, model.DataTSV:
			comma := rune(',')
			if k == model.DataTSV {
				comma = '\t'
			}
			cols, err := sample.DelimitedHeader(abs, 256*1024, comma)
			if err == nil {
				prof.ColumnNames = cols
				prof.ColumnHint = len(cols)
				prof.SampleNote = "header from first row (Go csv reader)"
			}
			if duckdb.Available() {
				if _, err := duckdb.PeekCSV(ctx, abs); err == nil {
					prof.DuckDBUsed = true
					if prof.SampleNote != "" {
						prof.SampleNote += "; DuckDB read_csv_auto ok"
					} else {
						prof.SampleNote = "DuckDB read_csv_auto ok"
					}
				} else {
					prof.DuckDBError = err.Error()
				}
			}
		case model.DataParquet:
			if duckdb.Available() {
				cols, note, err := duckdb.DescribeParquet(ctx, abs)
				if err == nil {
					prof.ColumnNames = cols
					prof.ColumnHint = len(cols)
					prof.DuckDBUsed = true
					prof.SampleNote = note
				} else {
					prof.DuckDBError = err.Error()
					prof.SampleNote = "parquet present; DuckDB describe failed"
				}
			} else {
				prof.SampleNote = "parquet; install DuckDB CLI for schema"
			}
		case model.DataJSON, model.DataJSONL:
			prof.SampleNote = "structured text; deep parse deferred (Phase 3)"
		}
		profiles = append(profiles, prof)
		count++
		return true
	})
	ds.Files = profiles

	ds.LayoutNotes = layoutNotes(ds)
	return ds
}

func layoutNotes(ds *model.DataSummary) []string {
	var notes []string
	if len(ds.ExtensionMix) > 1 {
		var parts []string
		for k, v := range ds.ExtensionMix {
			parts = append(parts, k+":"+strconv.Itoa(v))
		}
		notes = append(notes, "mixed formats: "+strings.Join(parts, ", "))
	}
	for _, p := range ds.Files {
		if p.PartitionHint != "" {
			notes = append(notes, "partition-like path segment: "+p.PartitionHint)
			break
		}
	}
	return notes
}
