# Heuristics

All inference is **deterministic and rules-based**. Numeric scores are relative weights for ranking, not calibrated probabilities.

## Archetypes

Signals come from:

- Presence of `.git` at scan root → `git_repo`.
- Manifests: `go.mod`, `package.json`, `pyproject.toml`, Dockerfiles, compose files, k8s-style dirs (`k8s`, `helm`, `charts`), docs configs (`mkdocs.yml`, etc.).
- Data extensions: `.csv`, `.tsv`, `.parquet`, `.jsonl` → `data_directory`.
- Weak manifest + many extensions + file count → `mixed_junk`.

Each archetype score carries `EvidenceItem` paths where applicable.

## Findings (suspicious plugin)

Examples:

- Permission errors from the scanner.
- `.git` directory not at relative path `.git` (nested repo marker).
- Large binary-like files inside a git repo.
- Multiple JS lockfiles in the same directory.
- `.env` files under a git root.
- `docker-compose` without a `Dockerfile` in the scanned tree (informational).

## Role hints (scanner)

Filename/extension heuristics tag nodes with `RoleManifest`, `RoleLockfile`, `RoleConfig`, `RoleData`, etc., for badges and structure signals.

## JSON files

`.json` is treated as **config** by default (lockfile substring → lock). `.jsonl`/`.ndjson` → **data**. Package-style JSON manifests are still caught by filename rules (`package.json`).
