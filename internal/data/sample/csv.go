package sample

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
)

// DelimitedHeader reads the first record with the given field delimiter (',' or '\t').
func DelimitedHeader(absPath string, maxBytes int64, comma rune) ([]string, error) {
	f, err := os.Open(filepath.Clean(absPath))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader = f
	if maxBytes > 0 {
		r = io.LimitReader(f, maxBytes)
	}
	cr := csv.NewReader(r)
	cr.Comma = comma
	cr.ReuseRecord = true
	row, err := cr.Read()
	if err != nil {
		return nil, err
	}
	return row, nil
}
