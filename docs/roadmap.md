# Roadmap

## Phase 1 (shipped)

- Recursive scan with ignores, symlinks, depth limits.
- Plugin registry and merge pipeline.
- Archetype and manifest heuristics, data sniffing, suspicion findings.
- Optional DuckDB CLI for Parquet describe / CSV sanity checks.
- Bubble Tea TUI (Overview, Tree, Findings, Evidence, Data).

## Phase 2 (planned)

- Git history, churn, and blame-assisted signals (not yet implemented).
- Diff or snapshot comparison between two scan/analysis runs.
- Findings tied to temporal patterns (stale branches, hot files).

## Phase 3 (planned)

- Light semantic parsing of source (not file-presence only).
- SQL and notebook understanding.
- Richer topology (service graph, import graph) where safe and local.
