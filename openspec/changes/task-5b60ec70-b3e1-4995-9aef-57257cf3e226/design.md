## Context

The Great Migration app is a Vue 3 + TypeScript moving-plan dashboard built with shadcn-vue components, Tailwind CSS v4 (OKLCH design tokens), and Vite. The home page at `src/home/HomeView.vue` renders summary statistic cards, a today's plan section, and move notes inside a `<section>` element with `flex flex-1 flex-col gap-6` spacing.

The original hello world v1 banner was removed in a clean-up task (50d758b7), returning the file to its pre-banner state. This v2 task re-validates that the OpenSpec refinement pipeline generates correct, mergeable output after pipeline improvements.

## Goals

- Render a visible greeting banner at the top of the home page at the `/` route as the first child of the `<section>` element
- Use only existing project patterns (Vue 3 script setup, Tailwind v4 utility classes, no new dependencies)
- Keep the change minimal (~8–12 lines of template markup, no script or style changes)
- Clearly distinguish v2 from v1 (text: "👋 Hello world v2 — swarm pipeline test") so pipeline iteration is traceable
- Pass `scripts/precommit-run` and `npm run build` (vue-tsc -b && vite build)

## Non-Goals

- Add any user-facing functionality beyond the banner (no data fetching, routing, state management, or interaction)
- Modify router configuration, app shell, layout, or any file outside `src/home/HomeView.vue` (and the OpenSpec change directory)
- Create new Vue components, composables, shared UI infrastructure, or shadcn-vue component additions
- Write tests (this is pipeline validation, not feature delivery)
- Revert or modify the scripts/check, scripts/test, or other project infrastructure from the v1 upgrade
- Alter archived OpenSpec artifacts from v1 or the clean-up task

## Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Banner placement | First child of `<section>` in HomeView.vue | The existing `flex flex-1 flex-col gap-6` container provides automatic vertical spacing. Placement above the summary grid ensures visibility at the `/` route entry point. Same placement as v1. |
| Styling approach | Plain `<div>` with `bg-muted/40 rounded-lg border p-4 border-l-4 border-primary text-sm` | The `bg-muted/40 rounded-lg border p-4` pattern is proven in TasksView.vue, PeopleView.vue, CalendarView.vue, and HomeView.vue (upcoming work cards). The `border-l-4 border-primary` adds visual accent consistent with v1 precedent. The `--color-primary` CSS variable is defined in `app/styles.css` (both light and dark themes), making `border-primary` a valid Tailwind v4 token. |
| Text content | "👋 Hello world v2 — swarm pipeline test" | Distinct from v1's "👋 Hello world — swarm pipeline test" for pipeline iteration traceability. The "v2" marker makes the two pipeline validations independently identifiable. |
| Feature flag | None | Not needed — the banner is a pipeline validation marker. |
| New component | None | Adding a new `.vue` file or `shadcn-vue add` contradicts non-goals. Inline markup only. |
| OpenSpec spec creation | Include spec under change directory | The `.openspec.yaml` schema is `spec-driven`, requiring companion spec artifacts. The v1 archive provides an exact template to follow. This fully exercises the spec-driven pipeline. |

## Risks & Mitigations

| Risk | Likelihood | Mitigation |
|---|---|---|
| `bg-muted/40` renders poorly in dark mode | Low | Already proven in dark mode by existing usage across 4 views in the codebase |
| Pre-commit hooks fail | Low | The change is ~5 lines of markup with zero new dependencies; any failure indicates a CI config issue, not a design flaw |
| `vue-tsc -b` type-strict compilation error | None | Template-only change introduces no TypeScript risk |
| `border-l-4 border-primary` style has no current codebase precedent | Low | The CSS variable `--color-primary` is defined in both light and dark themes in `app/styles.css`. The pattern was proven in v1 and accepted by all three refinement participants. |

## Traceability

All decisions are grounded in:
- Task prompt: "Test hello world v2 — create a small hello world banner to test swarm functionality inside this repo. We've done this before, but this is another test"
- Exploration dossier: confirms stack, placement, minimal-scope intent, and v2 differentiation requirement
- Round 1 outputs from swarm-architect, swarm-lead-dev, and swarm-reviewer — all three independently recommend the same approach with converging opinions on spec creation
- Accepted decisions: records the consensus to proceed with inline markup in HomeView.vue and include OpenSpec companion artifacts
- Artifact snapshot: initial proposal.md and tasks.md at `openspec/changes/task-5b60ec70-b3e1-4995-9aef-57257cf3e226/`