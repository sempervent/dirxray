package detect

import (
	"testing"

	"dirxray/internal/model"
)

func TestKindFromPath(t *testing.T) {
	tests := []struct {
		path string
		want model.DataFileKind
	}{
		{"a/b.csv", model.DataCSV},
		{"x.tsv", model.DataTSV},
		{"p.parquet", model.DataParquet},
		{"l.jsonl", model.DataJSONL},
	}
	for _, tc := range tests {
		if got := KindFromPath(tc.path); got != tc.want {
			t.Fatalf("%s: got %v want %v", tc.path, got, tc.want)
		}
	}
}

func TestPartitionHint(t *testing.T) {
	h := PartitionHintFromPath("data/dt=2024-01-01/file.csv")
	if h != "dt=2024-01-01" {
		t.Fatalf("got %q", h)
	}
}
