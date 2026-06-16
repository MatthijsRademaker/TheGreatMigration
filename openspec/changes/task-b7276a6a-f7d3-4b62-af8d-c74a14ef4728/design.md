## Context

The backend already exposes a daily-schedule read model and the design system already provides the card, badge, button, and avatar primitives needed for the schedule board. The missing piece is the frontend surface: `/calendar` still renders `Schedule board foundation` and a planning-window placeholder grid. This change turns the design images into a text-only implementation brief for a feature-local Vue component that can render immediately from deterministic demo data.

## Goals

- Add `frontend/src/calendar/DailySchedule.vue` as a frontend-only presentational component.
- Make the board implementable from text alone by spelling out the panel, day-column, task-card, priority, assignee, staffing, and placeholder anatomy.
- Reuse existing semantic UI primitives and tokens instead of adding new low-level components.
- Replace the `/calendar` placeholder with the new component.
- Add SSR/frontend tests for the component and the route.

## Non-Goals

- Wiring `GET /api/dashboard/daily-schedule` into the route or component.
- Regenerating OpenAPI/client artifacts or introducing Pinia Colada queries for daily schedule data.
- Changing backend handlers, persistence, migrations, or the existing daily-schedule contract.
- Adding drag/drop, CRUD, assignment editing, date navigation, filters, or alternate board modes beyond static design affordances.
- Replacing the home dashboard `Today’s plan` panel in this slice.

## Decisions

### 1. Integration surface

Render `DailySchedule` from `frontend/src/calendar/CalendarView.vue` and replace the existing placeholder entirely. Remove the current `usePlanningWindow()` import, loading/error branches, and `data-testid="plan-day-column"` success grid so the route becomes a deterministic presentational surface.

### 2. Local data contract and demo fixtures

Define feature-local TypeScript interfaces that mirror the backend field names and nesting for schedule days, task cards, and assigned people. Provide demo defaults in the component itself. The demo slice uses the design-image dates `2 Jul (Tue)` through `5 Jul (Fri)`, the visible helper counts `6 available`, `7 available`, `7 available`, and `5 available`, and the design task-title vocabulary (`Painting hall`, `Steam walls`, `Clean up`, `Sanding`, `Bedroom painting`, `Touch up woodwork`, `Living room finishing`, `2nd floor walls`, `Kitchen painting`, `Final clean`). Those titles are demo-only fixtures; the contract still preserves the backend property names and status values for later wiring.

### 3. Board layout

Use an outer `Card` shell with a header row. The board body should use `overflow-x-auto` around a horizontal `flex` container with consistent gaps. Each day column should remain readable through a fixed minimum width (`min-w-[280px]`) and `shrink-0` rather than collapsing into a tight responsive grid.

### 4. Task-card composition

Each task card is a plain `div`, not a nested `Card`, with `rounded-lg border bg-card shadow-sm p-3` styling and a `border-l-4` priority accent. Map priority to existing badge and border semantics: high -> `priorityHigh`/destructive, medium -> `priorityMedium`/warning, low -> `priorityLow`/success. Show the task title, a right-aligned priority badge, compact assignee metadata using the existing `Avatar` component plus names/initials, and staffing text such as `2 / 2`, `1 / 1`, and at least one under-staffed example through the same typed contract.

### 5. Static affordances

Keep the panel-level `View by: Day` and `Add task` controls plus each column's dashed `Add task` placeholder as visible, non-wired affordances. Use existing `Button` styling surfaces for these controls, but do not introduce event handling, filters, or task-creation behavior in this change.

### 6. Verification

Add a dedicated SSR component test at `frontend/tests/calendar/DailySchedule.test.ts` following the existing `renderToString` pattern, and update `frontend/tests/app-routes-render.test.ts` so `/calendar` asserts the new board content instead of planning-window placeholder text or column counts.

## Risks

- Dense schedule columns can drift from the design if spacing is not kept compact when composing multiple reused primitives inside a single card.
- Demo fixture titles intentionally diverge from the seeded backend schedule names, so the proposal must clearly mark them as design-mode defaults.
- Route tests will fail until the legacy `/calendar` placeholder assertions and `plan-day-column` count checks are removed.

## Conflict Resolution

- **Dates and weekday labels:** Resolved to the exact design-image slice `2 Jul (Tue)` through `5 Jul (Fri)` after verifying `designs/components.png` and `designs/home-page.png`; this removes the earlier weekday mismatch.
- **Demo task names:** Resolved to the design-image titles for visual fidelity, while keeping backend field names and enum values so later API wiring is additive.
- **Header actions:** Resolved to required static controls because `View by: Day` and `Add task` are visible in the Daily Schedule panel design, not optional embellishments.
- **Interactivity:** Resolved to non-interactive affordances only; neither the header controls nor the per-column `Add task` placeholders emit behavior in this slice.
- **Planning-window dependency:** Resolved to full removal from `CalendarView.vue` for this route, with planning-window coverage staying in its own tests rather than driving `/calendar`.
- **Contract coverage vs. design fidelity:** Resolved by keeping the design-image dates and title vocabulary while requiring the component contract and test data to cover all priority values and both staffing statuses.

## Traceability

- Task: `b7276a6a-f7d3-4b62-af8d-c74a14ef4728`
- Dossier: `2026-06-16T05:08:58.721Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Non-blocking review findings: `round-1-swarm-reviewer-blocker-1`, `round-1-swarm-reviewer-blocker-2`, `round-1-swarm-reviewer-blocker-3`, `round-1-swarm-reviewer-blocker-4`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial`
