## Why

Task 200FB581 was a pipeline validation exercise that added a "👋 Hello world — swarm pipeline test" banner to the home page and created associated OpenSpec artifacts. That validation is complete — the banner and spec are pipeline-test artifacts that should be removed to keep the project clean.

## What Changes

- **Remove** the greeting banner `<div>` from `src/home/HomeView.vue` (4 lines of template markup)
- **Remove** the active greeting-banner spec at `openspec/specs/greeting-banner/spec.md`
- **No change** to `scripts/check` or `scripts/test` — those were placeholder-to-functional upgrades that are now useful project infrastructure
- **No change** to the archived change artifacts (the `2026-06-15-task-...` directory under `openspec/changes/archive/`) — they are inert and serve as historical record

## Capabilities

### New Capabilities
<!-- None — this change removes code and specs, it doesn't introduce new capabilities -->

### Modified Capabilities
- `greeting-banner`: **Removed entirely.** The greeting-banner spec at `openspec/specs/greeting-banner/spec.md` was created by the archive of the test task and describes pipeline-validation requirements that no longer apply.

## Impact

- `src/home/HomeView.vue` — remove 4 lines of template markup (the banner `<div>`)
- `openspec/specs/greeting-banner/` — remove the entire spec directory
- No functional impact on routing, state, data fetching, or any other component
- The `scripts/check` and `scripts/test` improvements from the test task remain intact
