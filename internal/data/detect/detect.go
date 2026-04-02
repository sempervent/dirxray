package detect

import (
	"path/filepath"
	"strings"

	"dirxray/internal/model"
)

// KindFromPath guesses data format from extension and name.
func KindFromPath(relPath string) model.DataFileKind {
	ext := strings.ToLower(filepath.Ext(relPath))
	base := strings.ToLower(filepath.Base(relPath))
	switch ext {
	case ".csv":
		return model.DataCSV
	case ".tsv":
		return model.DataTSV
	case ".parquet":
		return model.DataParquet
	case ".jsonl", ".ndjson":
		return model.DataJSONL
	case ".json":
		if strings.Contains(base, "manifest") || strings.Contains(base, "package") {
			return model.DataUnknown
		}
		return model.DataJSON
	default:
		return model.DataUnknown
	}
}

// PartitionHintFromPath extracts coarse hive-style hints from path segments.
func PartitionHintFromPath(relPath string) string {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	for _, p := range parts {
		if strings.Contains(p, "=") {
			return p
		}
	}
	return ""
}
