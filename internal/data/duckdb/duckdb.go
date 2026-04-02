// Package duckdb runs the DuckDB CLI for local sampling when available.
// Rationale: avoids CGO and cross-compilation friction; install CLI from
// https://duckdb.org/docs/installation/ for Parquet/schema previews.
package duckdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const queryTimeout = 10 * time.Second

// Available returns true if `duckdb` is on PATH.
func Available() bool {
	name := "duckdb"
	if runtime.GOOS == "windows" {
		name = "duckdb.exe"
	}
	_, err := exec.LookPath(name)
	return err == nil
}

// Explain documents integration for operators.
func Explain() string {
	if Available() {
		return "DuckDB CLI on PATH: Parquet DESCRIBE and CSV peek enabled."
	}
	return "DuckDB CLI not on PATH: using byte/header heuristics only. Install: https://duckdb.org/docs/installation/"
}

// DescribeParquet returns column names from DESCRIBE output (best-effort parse).
func DescribeParquet(ctx context.Context, absPath string) ([]string, string, error) {
	if !Available() {
		return nil, "", errors.New("duckdb not available")
	}
	p := filepath.ToSlash(filepath.Clean(absPath))
	p = strings.ReplaceAll(p, "'", "''")
	q := fmt.Sprintf("DESCRIBE SELECT * FROM read_parquet('%s');", p)
	out, err := run(ctx, q)
	if err != nil {
		return nil, "", err
	}
	var cols []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "┌") || strings.HasPrefix(line, "└") || strings.HasPrefix(line, "│") {
			continue
		}
		// column_name | column_type | null | key | default | extra
		parts := strings.Split(line, "|")
		if len(parts) > 0 {
			c := strings.TrimSpace(parts[0])
			if c != "" && c != "column_name" {
				cols = append(cols, c)
			}
		}
	}
	return cols, "DESCRIBE via DuckDB CLI", nil
}

// PeekCSV runs a single-row select to validate CSV readability (optional).
func PeekCSV(ctx context.Context, absPath string) (string, error) {
	if !Available() {
		return "", errors.New("duckdb not available")
	}
	p := filepath.ToSlash(filepath.Clean(absPath))
	p = strings.ReplaceAll(p, "'", "''")
	q := fmt.Sprintf("SELECT COUNT(*) FROM read_csv_auto('%s', header=true, sample_size=-1);", p)
	return run(ctx, q)
}

func run(ctx context.Context, sql string) (string, error) {
	name := "duckdb"
	if runtime.GOOS == "windows" {
		name = "duckdb.exe"
	}
	cctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	cmd := exec.CommandContext(cctx, name, "-c", sql)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(buf.String()))
	}
	return buf.String(), nil
}
