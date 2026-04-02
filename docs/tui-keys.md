# TUI keybindings

| Key | Action |
|-----|--------|
| `1` … `5` | Switch tabs: Overview, Tree, Findings, Evidence, Data |
| `j` / `↓` | Move cursor (tree/findings) or scroll down (viewport tabs) |
| `k` / `↑` | Move cursor or scroll up |
| `PgDn` | Page down (viewport tabs) |
| `PgUp` | Page up (viewport tabs) |
| `Enter` | Tree: select node path and jump to Evidence. Findings: use related path and jump to Evidence |
| `?` | Toggle help overlay |
| `q` / `Ctrl+C` | Quit |

## Tabs

- **Overview** — Stats, dominant archetypes, signals, top findings.
- **Tree** — Flattened pre-order tree with role/manifest badges.
- **Findings** — Severity-sorted list.
- **Evidence** — Indexed evidence for `selectedPath` plus global archetype evidence.
- **Data** — Data profiles, partition hints, DuckDB status.
