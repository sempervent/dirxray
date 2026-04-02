# Architecture

`dirxray` is a single Go binary: CLI (`internal/cli`) parses flags, `internal/scan` builds an in-memory tree, `internal/analyze` runs registered plugins and merges results, and `internal/tui` renders Bubble Tea views.

## Data flow

1. **Scan** — Recursive walk with options (hidden, max depth, symlinks, `.gitignore` + internal ignores). Produces `model.ScanResult` with `model.Node` tree and `ScanNotice` entries for permission errors.
2. **Analyze** — `analyze.Run` executes each `plugins.Plugin` in order, collecting `PluginResult` slices. `analyze.Merge` sorts archetypes/findings, applies per-node badges, builds `DirectorySummary`, and indexes `EvidenceItem` by path.
3. **TUI** — `tui.New` flattens the tree for the Tree tab; Overview/Evidence/Data use a `bubbles/viewport` for scrolling.

## Packages

| Path | Role |
|------|------|
| `cmd/dirxray` | `main` |
| `internal/cli` | Cobra root command |
| `internal/scan` | Filesystem inventory |
| `internal/ignore` | `.gitignore` + built-in deny patterns |
| `internal/git` | `.git` discovery (Phase 1: presence only) |
| `internal/model` | Shared domain types |
| `internal/analyze` | Merge + orchestration |
| `internal/analyze/plugins` | Built-in analyzers |
| `internal/data` | Data detection + CSV/TSV header sampling + DuckDB CLI bridge |
| `internal/tui` | Bubble Tea model, views, styles |

## Cross-platform notes

- Paths in the model use slash-separated **relative** paths for display and evidence keys; `filepath` is used at the OS boundary.
- Hidden files: dot-prefix (Unix convention); Windows-specific hidden attribute is not yet read.
- Symlinks: default is **not** to follow; `--follow-symlinks` resolves targets for traversal.

## DuckDB

Phase 1 uses the **DuckDB CLI** on `PATH` when present (no CGO). See `internal/data/duckdb`. Without the CLI, CSV/TSV headers are still sniffed with Go’s `encoding/csv`.
