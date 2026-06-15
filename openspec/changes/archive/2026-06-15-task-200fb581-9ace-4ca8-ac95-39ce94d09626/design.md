## Context

The Great Migration app is a Vue 3 + TypeScript moving-plan dashboard built with shadcn-vue components, Tailwind CSS v4 (OKLCH design tokens), and Vite. The home page at `src/home/HomeView.vue` renders summary statistic cards, a today's plan section, and move notes inside a flex container with `gap-6` spacing.

This task is a pipeline validation exercise: the banner's content or design is not functionally meaningful — correctness of the swarm pipeline is the actual deliverable.

## Goals

- Render a visible greeting banner at the top of the home page at the `/` route
- Use only existing project patterns (Vue 3 script setup, Tailwind v4 utility classes)
- Keep the change minimal (~8–12 lines of template markup)
- Pass pre-commit hooks and build verification

## Non-Goals

- Add any user-facing functionality beyond the banner
- Modify router configuration, app shell, or layout
- Create new components, composables, or shared UI infrastructure
- Write tests (pipeline validation, not feature delivery)

## Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Banner placement | First child of `<section>` in HomeView.vue | The existing `flex flex-1 flex-col gap-6` container provides automatic vertical spacing. Placement above the summary grid ensures visibility at the `/` route entry point. |
| Styling approach | Plain `<div>` with `bg-muted/40 rounded-lg border p-4 border-l-4 border-primary` | The `bg-muted/40 rounded-lg border p-4` pattern is proven in `tasks/TasksView.vue:21` and `shared/ui/input/Input.vue:27`. The `border-l-4 border-primary` adds visual accent. The shadcn-vue Alert component does not exist in the project and adding it contradicts the minimal-change constraint. |
| Text content | "👋 Hello world — swarm pipeline test" | Provides traceability to the task ID for pipeline validation. The emoji and distinct text make the banner recognizable without being intrusive. |
| Feature flag | None | Not needed — the task is explicitly a permanent banner for pipeline validation. No gating required per accepted decisions. |
| New component | None | Adding a new `.vue` file or `shadcn-vue add alert` contradicts non-goals. Inline markup only. |

## Risks & Mitigations

| Risk | Likelihood | Mitigation |
|---|---|---|
| `bg-muted/40` renders poorly in dark mode | Low | Already proven in dark mode by existing usage in `tasks/TasksView.vue` and `shared/ui/input/Input.vue` |
| Pre-commit hooks fail | Low | The change is ~5 lines of markup with zero new dependencies; any failure indicates a CI config issue, not a design flaw |
| `vue-tsc -b` type-strict compilation error | None | Template-only change introduces no TypeScript risk |

## Traceability

All decisions are grounded in:
- Task prompt: "Test hello world — create a small hello world banner to test swarm functionality"
- Exploration dossier: confirms stack, placement, and minimal-scope intent
- Round 1 outputs from swarm-architect, swarm-lead-dev, and swarm-reviewer — all three independently recommend the same approach with minor styling variations
- Accepted decisions: records the consensus to proceed with inline markup in HomeView.vue
- Artifact snapshot: initial proposal.md and tasks.md at `openspec/changes/task-200fb581-9ace-4ca8-ac95-39ce94d09626/`
