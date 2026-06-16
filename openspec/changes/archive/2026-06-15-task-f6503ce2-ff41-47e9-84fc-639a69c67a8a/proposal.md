## Why

The `/tasks` route still renders a generic "Task foundation" placeholder, while the task explicitly asks for the Task management component shown in `designs/components.png`. The next valuable slice is a frontend-only implementation that reproduces the designed task table area with static placeholder data, so the route demonstrates the intended UI now without waiting for backend wiring.

## What Changes

This change replaces the `/tasks` placeholder content with a feature-local Task Management panel built under `frontend/src/tasks/` and composed into `frontend/src/tasks/TasksView.vue`.

The implementation includes:

- A reusable `TaskManagementPanel` component for the section 4 Task Management layout from `designs/components.png`
- A `TaskRow` sub-component for row rendering and isolated row-level markup
- Feature-local typed task data and fixtures that mirror the future backend `TaskRow` contract (`id`, `title`, `priority`, `peopleNeeded`, `room`, `status`, `assignedTo`)
- A centralized display mapping that renders the design label `Unassigned` from empty `assignedTo` values instead of introducing a new canonical status
- Existing shared primitives (`Card`, `Button`, `Badge`) and semantic badge variants (`priorityHigh`, `priorityMedium`, `priorityLow`, `secondary`) plus verified `@lucide/vue` icons for the toolbar and people-needed indicator
- A five-column task list/table presentation with horizontal overflow support rather than speculative mobile card-collapse behavior
- A separate priority legend showing High, Medium, and Low pills
- Updated route-render test coverage for `/tasks` so assertions target the new task-management content instead of the old placeholder

The controls remain frontend-only: Filter and Add Task are visible, clickable, and emit placeholder events, but they do not perform filtering, creation, API calls, or mutation behavior.

## Impact

- **Affected frontend files**: `frontend/src/tasks/TasksView.vue`, new task feature components under `frontend/src/tasks/components/`, feature-local types/fixtures under `frontend/src/tasks/`, and `/tasks` route tests in `frontend/tests/app-routes-render.test.ts`
- **Design-system usage**: Reuses existing shared UI primitives and semantic tokens only; no raw hex colors, inline color styles, or page-specific color hacks
- **Backend/API**: No backend handlers, API clients, generated OpenAPI artifacts, persistence, or route wiring changes
- **Verification**: Implementation must update automated tests and pass `scripts/precommit-run`
- **Out of scope**: Backend integration, live filtering, task creation flows, assignment workflows, status persistence, new showcase/catalog surfaces, and speculative responsive redesign outside the initial table-overflow approach
