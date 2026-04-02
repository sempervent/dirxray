# dirxray

Terminal-native, **local-first** directory explainer: scan a path, then navigate evidence-backed heuristics in a [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI. No LLM, no web UI—rules and file signals only.

```bash
go run ./cmd/dirxray .
# or
make build && ./bin/dirxray --path /your/project
```

## Features (Phase 1)

- Recursive scan with `.gitignore` + built-in ignores, `--max-depth`, `--hidden`, `--follow-symlinks`, `--no-gitignore`
- **Overview** — archetypes, signals, stats, top findings  
- **Tree** — annotated inventory  
- **Findings** — severity-ranked issues  
- **Evidence** — path-indexed rationale  
- **Data** — CSV/TSV/Parquet/JSONL hints; optional [DuckDB CLI](https://duckdb.org/docs/installation/) for Parquet `DESCRIBE` / CSV checks  

## CLI flags

| Flag | Meaning |
|------|---------|
| `--path` | Directory to scan (default: `.` or first argument) |
| `--hidden` | Include dotfiles |
| `--max-depth` | Limit recursion (0 = unlimited) |
| `--follow-symlinks` | Follow symlinks |
| `--no-gitignore` | Skip `.gitignore` merging |
| `--debug` | Log target path to stderr |

## Docs

- [Architecture](docs/architecture.md)  
- [Plugin authoring](docs/plugins.md)  
- [Heuristics](docs/heuristics.md)  
- [TUI keys](docs/tui-keys.md)  
- [Roadmap / phases](docs/roadmap.md)  

## Development

```bash
make test    # go test ./...
make vet     # go vet ./...
make build   # outputs bin/dirxray
```

Go **1.22+**. Module path: `dirxray`.

## Screenshots

Replace this section with a terminal recording (e.g. [asciinema](https://asciinema.org/)) or screenshots once you capture them.

## License

See [LICENSE](LICENSE).
