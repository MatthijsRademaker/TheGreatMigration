## Context

The task asks for the Task management component shown in `designs/components.png`, but `frontend/src/tasks/TasksView.vue` currently renders only a generic "Task foundation" placeholder. The repository already has the shared primitives and semantic tokens needed for this UI (`Card`, `Button`, `Badge`, Lucide icons, and design-system token classes), and the backend backlog contract already documents the future task row shape. This slice is therefore a static, frontend-only composition task: build the designed task-management panel now, keep the data local, and avoid implementation choices that would fight the later API-wiring change.

## Goals / Non-Goals

### Goals
- Replace the `/tasks` placeholder with a Task Management panel matching section 4 of `designs/components.png`
- Keep all new implementation under the tasks feature area and compose it through `frontend/src/tasks/TasksView.vue`
- Use static placeholder data that includes the design example row and exercises High/Medium/Low priorities, people-needed counts, rooms/areas, and an unassigned visual state
- Reuse existing shared UI primitives and semantic badge/token variants instead of page-specific styling
- Update route-render coverage so `/tasks` is verified through the new static UI

### Non-Goals
- Backend integration, API calls, OpenAPI regeneration, Pinia Colada queries, persistence, or mutation behavior
- Real filtering, task creation, assignment workflows, or status editing
- New shared table primitives, showcase routes, Storybook, or external component documentation tooling
- Redesigning dashboard cards, people views, scheduling UI, or global design tokens
- Speculative responsive card-collapse behavior beyond the initial overflow-based table layout

## Conflict Resolution

### Status vocabulary
Refinement evidence conflicted on whether the component should introduce an `unassigned` status value locally or mirror the backend contract exactly. Resolution: mirror the backend `TaskRow` shape (`status` plus `assignedTo`) and derive the displayed `Unassigned` pill from empty `assignedTo` through a centralized mapping function. This preserves the design label without creating a conflicting canonical status.

### Initial layout behavior
One proposal suggested stacked responsive rows, while multiple reviewers flagged that as speculative scope because the design reference only shows the desktop table. Resolution: implement the five-column table/grid with `overflow-x-auto` support and defer any stacked mobile redesign.

### Control interactivity and forward contract
Evidence left open whether Filter and Add Task should be disabled, decorative, or interactive. Resolution: both controls remain visibly interactive, emit placeholder `filter` and `add-task` events, and perform no data operations.

### Fixture ownership and future wiring
Evidence raised whether the panel should accept a tasks prop immediately or keep local seed data. Resolution: keep feature-local typed fixtures in a dedicated tasks file for this slice, with no backend/client prop wiring yet, while shaping the data to match the future backend contract.

## Decisions

### D1: Use feature-local composition under `frontend/src/tasks/`
Create a `TaskManagementPanel` component and a `TaskRow` sub-component inside the tasks feature area, then render the panel from `TasksView.vue`. Do not promote this UI into shared primitives yet.

### D2: Mirror the backend task row contract in local types
Define a feature-local task row type matching the future backend shape: `id`, `title`, `priority`, `peopleNeeded`, `room`, `status`, and `assignedTo`. Keep fixture data in a dedicated tasks-local file so future API wiring can replace the import instead of rewriting the component.

### D3: Centralize display mapping for the status pill
Render the design label `Unassigned` when a row has no assigned people. Keep that mapping in one typed utility/helper rather than embedding the label directly across templates.

### D4: Reuse existing shared primitives and semantic variants only
Use `Card`, `Button`, and `Badge` for the panel shell and pills. Use `priorityHigh`, `priorityMedium`, and `priorityLow` for row badges and the priority legend, and use the existing `secondary` badge variant for the `Unassigned` pill. Do not add raw colors, inline color styles, or a new shared table primitive.

### D5: Use a five-column table/list presentation with overflow support
Render the columns shown in the design: Task, Priority, People Needed, Room / Area, and Status. Use a table-like grid or semantic table with horizontal overflow support. The People Needed column should render an inline icon-plus-count treatment using verified `@lucide/vue` icons.

### D6: Establish forward toolbar events without backend behavior
The Task Management panel should emit `filter` and `add-task` events from its toolbar controls. The route does not wire those events to data behavior yet.

### D7: Keep testing focused on route-rendered static UI
Update `frontend/tests/app-routes-render.test.ts` so the `/tasks` route asserts the new panel content, including the panel title, representative row content, and toolbar/legend labels. The implementation must still pass `scripts/precommit-run`.

## Risks

| Risk | Mitigation |
| --- | --- |
| Display logic could hardcode `Unassigned` as a canonical status and complicate later API wiring. | Mirror the backend row type and derive the display label from `assignedTo` in one mapping helper. |
| Incorrect Lucide icon names would fail the strict TypeScript build. | Verify the toolbar and people-needed icon exports from `@lucide/vue` before finalizing the template. |
| Responsive layout work could expand beyond the approved slice. | Limit the first implementation to the five-column overflow layout and defer stacked mobile redesign. |
| Ad-hoc styling could bypass semantic tokens. | Restrict the component to existing shared primitives, badge variants, and semantic utility classes only. |

## Traceability

- Task: `f6503ce2-ff41-47e9-84fc-639a69c67a8a`
- Dossier: `2026-06-15T21:18:09.368Z`
- Accepted decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Validated round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial`
