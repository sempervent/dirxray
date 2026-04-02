# Plugin authoring

Plugins implement `plugins.Plugin`:

```go
type Plugin interface {
    Name() string
    Run(ctx *Context) (*model.PluginResult, error)
}
```

`Context` carries `RootAbs` and the immutable `*model.ScanResult`.

`PluginResult` fields merge into the global analysis:

- `Signals` — weighted `ProjectSignal` hints.
- `Archetypes` — `ArchetypeScore` with evidence.
- `Findings` — evidence-backed `Finding` records.
- `Data` — optional `DataSummary` (only one merged winner by file count today).
- `NodeBadges` — map relative path (`"."` for root files) to short badge strings.

## Registration

`plugins.NewRegistry()` returns the default ordered set. For tests, call `Register` on a custom `Registry`. Dynamic `.so` loading is **not** in scope for Phase 1; compile-time registration keeps the binary hermetic.

## Ordering

Run order matters for badge accumulation (later plugins append badges). Prefer: structure → manifest → archetype → data → suspicious.
