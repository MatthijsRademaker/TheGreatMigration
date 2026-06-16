## Context

The repository has a rich set of design artifacts and backend contracts that define what the homepage dashboard should look like, but the actual `HomeView.vue` is a placeholder. The key inputs are:

- **Design artifacts**: `designs/home-page.png` (target composition), `designs/components.png` (named regions: KPI Summary Cards, Task Management, People Availability, Schedule Components, Notes Components, Utility Components), `designs/design-system.png` (navigation/layout fragments).
- **Backend contracts**: `openspec/specs/task-backlog-api/spec.md` (summary counts + task rows), `openspec/specs/dashboard-people-availability/spec.md` (people summary + availability matrix), `openspec/specs/dashboard-daily-schedule/spec.md` (per-day schedule columns).
- **Existing UI components**: `AppSidebar.vue` (sidebar with Dashboard/Tasks/Schedule/People/Rooms/Settings navigation, `collapsible="icon"`), `AppShell.vue` (simple page-title header + Planning mode badge, NOT the richer design toolbar), `PeopleAvailability.vue` (reusable availability matrix), `TaskManagementPanel.vue` (reusable task backlog table).
- **Design system**: `docs/design-system-v2.md` defines 12-column grid, calm green palette, Inter typography, spacing/radius/shadow hierarchy.
- **Current HomeView.vue**: Renders Hello world card, Available today (6), Under-staffed (3), Move days, static Today's plan, and Move notes — none matching the target dashboard.

The deliverable is a text-only artifact that makes the target composition legible to a text-only LLM without requiring PNG inspection.

## Goals / Non-Goals

### Goals
- Produce a text-only homepage composition artifact that mirrors the desktop structure shown in `designs/home-page.png`.
- Use labels and section boundaries that map to existing product vocabulary: sidebar, top toolbar, KPI cards, Tasks Backlog, People Availability, Daily Schedule, and Move Notes.
- Make the mockup detailed enough that a text-only LLM can compose a future `HomeView` without needing to inspect PNG assets.
- Document the gap between the current placeholder home route and the intended dashboard composition so follow-up implementation can be scoped cleanly.

### Non-Goals
- Implement the homepage composition in Vue.
- Redesign the information architecture or invent new dashboard sections beyond the existing design artifacts and specs.
- Create pixel-perfect art, responsive breakpoints, or mobile mockups.
- Edit canonical OpenSpec specs or publish final product docs.

## Decisions

### Decision 1: Artifact location and structure
**Decision**: Create `openspec/changes/task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a/home-page-ascii.md` with legend, desktop wireframe, collapsed-sidebar variant, and gap analysis.
**Rationale**: The change-local directory is the correct non-canonical home. Three sections (legend, wireframe, gap analysis) plus a collapsed-sidebar variant address all reviewer requirements without modifying canonical artifacts.
**Sources**: swarm-architect (round 1), swarm-lead-dev (round 1), swarm-reviewer (round 1).

### Decision 2: Exact KPI card enumeration
**Decision**: The four KPI cards SHALL be: (1) Total Tasks — from `task-backlog-api` `summary.totalTasks`, (2) Available People — from `dashboard-people-availability` `summary.availableToday` / `summary.totalPeople`, (3) High Priority Tasks — from `task-backlog-api` `summary.highPriorityTasks`, (4) Move Days — from `planning-window` `planWindowDayCount`. Each card SHALL list its backing API endpoint in the legend.
**Rationale**: The reviewer identified ambiguity between current placeholder cards and the target dashboard KPIs. These four cards are derivable from existing spec contracts and the current `HomeView.vue` Move days card (which already consumes `usePlanningWindow()`).
**Sources**: swarm-reviewer (round 1) non-blocking item 2; `task-backlog-api/spec.md` summary requirement; `dashboard-people-availability/spec.md` summary requirement; `HomeView.vue` `usePlanningWindow()` usage.

### Decision 3: Aspirational component annotations
**Decision**: The Daily Schedule and Move Notes regions SHALL include explicit annotations noting that no Vue components exist yet. Daily Schedule is backed by `dashboard-daily-schedule/spec.md` (API contract only). Move Notes has no spec or component — it is a static design-only panel.
**Rationale**: A text-only LLM reading the ASCII must not attempt to import non-existent components. The legend must distinguish between implemented components (`PeopleAvailability.vue`, `TaskManagementPanel.vue`) and aspirational regions.
**Sources**: swarm-reviewer (round 1) non-blocking item 3; find results confirming no DailySchedule or MoveNotes components exist.

### Decision 4: Toolbar gap annotation
**Decision**: The ASCII wireframe SHALL show the aspirational toolbar from the design (date-range, Today, arrows, alerts, profile chip) but SHALL include an explicit gap annotation noting that `AppShell.vue` currently renders only a page-title header with a Planning mode badge.
**Rationale**: Without this annotation, a future executor would misinterpret the ASCII as representing current state and waste time on mismatched expectations.
**Sources**: swarm-reviewer (round 1) non-blocking item 1; `AppShell.vue` lines showing header composition; `docs/design-system-v2.md` Navigation states section.

### Decision 5: Collapsed sidebar variant
**Decision**: The artifact SHALL include a second narrower ASCII variant showing the layout when the sidebar is collapsed to icon-only, since `AppSidebar` already supports `collapsible="icon"`.
**Rationale**: Future implementers need to see how the content area reflows when the sidebar collapses. Low cost to add, high value for responsive/composition work.
**Sources**: swarm-lead-dev (round 1); `AppSidebar.vue` `<Sidebar collapsible="icon" variant="inset">`.

## Risks

| Risk | Mitigation |
|------|-----------|
| ASCII proportions may not precisely match pixel layout in `designs/home-page.png` since the PNG cannot be directly inspected. | Use 12-column grid from design-system-v2.md as the proportional guide; annotate legend with "design-PNG inferred" where labels derive from specs rather than image inspection. |
| Future executor may interpret ASCII as pixel-accurate spec rather than compositional blueprint. | Artifact SHALL frame itself as a structural reference, not a visual spec, in its preamble. |
| Current AppShell header differs significantly from target toolbar — gap may be misinterpreted as requirement to rewrite AppShell in this same task. | Explicit gap annotation states the toolbar is aspirational and out of scope for this task. |
| KPI card labels and counts in target design PNG cannot be verified without image access. | Derive KPI labels from spec contracts and note the source in the legend. |
| Path confusion: executor could accidentally write to `openspec/specs/` or `docs/`. | Artifact preamble SHALL state the non-canonical constraint explicitly. |

## Traceability

- **Task**: `2bdd3fb5-c4ec-4153-9a49-d6b463894a8a` — "Create ascii art mockup of home page"
- **Dossier**: `2026-06-16T05:06:18.453Z` — exploration dossier with problem framing, goals, acceptance criteria
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Round outputs**: round:1:agent:swarm-architect, round:1:agent:swarm-lead-dev, round:1:agent:swarm-reviewer
- **Non-blocking items**: round-1-swarm-reviewer-blocker-1 (toolbar gap), round-1-swarm-reviewer-blocker-2 (KPI enumeration), round-1-swarm-reviewer-blocker-3 (aspirational components)
- **Artifact snapshot**: `initial`
