## Context

The test hello world task (200FB581) added a pipeline-validation greeting banner to `src/home/HomeView.vue` and created a corresponding OpenSpec spec at `openspec/specs/greeting-banner/`. It also replaced placeholder `scripts/check` and `scripts/test` scripts with working Docker-based implementations. The validation is complete; the banner and spec are test scaffolding that need cleanup.

## Goals / Non-Goals

**Goals:**
- Remove the greeting banner `<div>` from HomeView.vue's template
- Remove the active `openspec/specs/greeting-banner/` spec directory
- Keep `scripts/check` and `scripts/test` as-is (they are now useful project infrastructure)
- Keep archived change artifacts as historical record

**Non-Goals:**
- No revert of scripts/check or scripts/test to placeholders
- No changes to any other source files, routing, or components
- No changes to archived OpenSpec change artifacts

## Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Script retention | Keep scripts/check and scripts/test as-is | These were upgraded from placeholders to working Docker-based scripts. They are now useful project infra. Reverting them would be a step backward — the project gains nothing from returning to `echo "TODO" + exit 1`. |
| Archive retention | Keep archived change artifacts | The archive at `openspec/changes/archive/` is inert by design — it serves as historical record. Removing it offers no benefit and would break the OpenSpec convention of archiving completed changes. |
| Spec removal strategy | Delete the entire `openspec/specs/greeting-banner/` directory | The greeting-banner capability has no reason to exist after the banner is removed. A delta spec marking requirements as REMOVED clarifies intent during this change's lifecycle. |

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| Pre-commit hooks reject empty change | Not a risk — the change touches HomeView.vue and a spec directory. At minimum, the banner removal and spec deletion produce a meaningful diff. |
| Spec directory removal orphans the archive copy | Not a risk — the archive copy at `openspec/changes/archive/2026-06-15-task.../specs/greeting-banner/spec.md` is the authoritative historical record. The live spec is the one being removed. |
