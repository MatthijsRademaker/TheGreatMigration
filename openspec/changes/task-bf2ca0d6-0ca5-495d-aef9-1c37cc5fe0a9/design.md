## Context

The Great Migration app has backend primitives for connectivity (`/api/hello`), planning-window, and dashboard people-availability, but no task backlog API. The homepage design, C4 architecture model, and `/tasks` route all depend on task data with priority, staffing, room, and assignment state. This change adds the first read-only task backlog endpoint to fill that gap.

## Goals / Non-Goals

### Goals
- Define a first read-only `GET /api/tasks/backlog` endpoint with Huma registration and OpenAPI visibility.
- Return backlog rows with priority, people-needed, room, status, and assignment metadata.
- Expose derived summary counts (total, high-priority, unassigned, understaffed) in a single fetch.
- Provide canonical `priorities` and `statuses` legend arrays with design-system color intents.
- Implement in-memory seed data exercising all vocabulary values across at least 10 tasks.
- Cover with backend tests for shape, consistency, vocabulary, and non-regression.
- Remain additive: no auth, persistence, CORS changes, or existing endpoint modifications.

### Non-Goals
- No task create/edit/delete workflow or persistence layer.
- No authentication, authorization, or tenancy.
- No schedule-board, people-availability, move-notes, or full rooms CRUD expansion.
- No homepage or TasksView UI implementation.
- No OpenAPI snapshot regeneration for the frontend (deferred to follow-up).
- No finalization of the long-term task/room domain model beyond the fields needed for this first read contract.
- No `roomsCompleted` summary (requires a room model that does not exist yet).

## Decisions

### Decision 1: Domain-aligned endpoint path `/api/tasks/backlog`

The C4 architecture model shows Tasks as its own feature area that captures and organizes `taskData`, while Dashboard summarizes it. Placing the canonical task source under `/api/tasks/` respects this boundary and avoids duplication when the `/tasks` route eventually needs richer endpoints. This is not a dashboard BFF path — it is the domain source.

### Decision 2: Single `status` field with `assignedTo` array

A single `status` field (values: `backlog`, `ready`, `assigned`) covers the planning lifecycle visible in TasksView.vue's three status buckets. The `assignedTo` array of person-ID strings provides the assignment dimension needed for unassigned/understaffed counts without duplicating person data already available from `/api/dashboard/people-availability`. Splitting into separate planning-status and assignment-state fields would over-model before calendar scheduling coupling is understood.

### Decision 3: Summary includes `understaffedTasks`

`frontend/src/home/HomeView.vue` already renders a static "Under-staffed: 3" summary card. Without this computed value, the dashboard card remains permanently placeholder. The computation is straightforward: `understaffedTasks` = count of tasks where `len(assignedTo) > 0` and `len(assignedTo) < peopleNeeded`.

### Decision 4: Include `priorities` legend, not just `statuses`

Following the established `dashboard.go` pattern, both `priorities` and `statuses` legend arrays are included with `id`, `label`, and `colorIntent` fields. While priority values are documented in the design system, including them in the response provides frontend rendering metadata without requiring hardcoded mappings. Priority colors: high=destructive, medium=warning, low=success. Status colors: backlog=muted, ready=info, assigned=success.

### Decision 5: Field naming — `room` not `roomArea`

The C4 architecture model consistently uses `room` for the task data field. TasksView.vue also references `room` in its task field expectations. The `room` name is consistent with the architecture source of truth.

### Decision 6: `unassigned` defined by empty `assignedTo`, not by status

A task is unassigned when `assignedTo` is empty — this is orthogonal to status. A ready task with no one assigned is still unassigned. This decouples the planning lifecycle from staffing state and avoids ambiguities in derived counts.

### Decision 7: Seed data — representative variety with at least 10 tasks

The seed includes at least 10 tasks exercising all three priority levels (high, medium, low), all three statuses (backlog, ready, assigned), at least 4 distinct room names, and a mix of assignment states (empty `assignedTo`, partially filled, fully filled). Task IDs follow the existing Person ID pattern: `"task-1"`, `"task-2"`, etc.

## Risks

| Risk | Mitigation |
| --- | --- |
| OpenAPI snapshot drift: The frontend snapshot won't include the new endpoint until regenerated. | Explicitly defer snapshot regeneration to a follow-up task. Document in tasks.md. |
| Single `status` field may need to split later if calendar scheduling needs independent planning vs. assignment tracking. | The `assignedTo` array provides an escape valve — assignment state is independently derivable. Document in spec that status granularity may evolve. |
| `room` naming may need adjustment when a rooms domain model emerges. | Acceptable for first slice; rename is a simple refactor. |

## Traceability

| Source | Evidence |
| --- | --- |
| Dossier goals | "Define a first read-only backend contract for backlog task data aligned to the homepage/task design" |
| Dossier acceptance criteria | A new read-only Huma endpoint, response contract with backlog rows and summary, canonical vocabulary, 8+ tasks, backend tests, additive-only |
| Dossier non-goals | No writes, persistence, auth, rooms CRUD, frontend UI work |
| C4 architecture (`workspace.c4`) | taskData = "Move jobs with priority, status, room, staffing needs, and assignments" |
| TasksView.vue | Three status buckets: Backlog, Ready, Assigned |
| HomeView.vue | Static "Under-staffed: 3" summary card |
| design-system-v2.md | Priority color intents: high=destructive, medium=warning, low=success |
| dashboard.go | Established pattern: seeded additive Huma endpoint with summary + data + legend arrays |
| main_test.go | Established pattern: additive test with shape, consistency, and non-regression checks |
| Architect recommendation | Path `/api/tasks/backlog`, single status + assignedTo, priorities + statuses legends, representative seed |
| Lead-dev recommendation | At least 10 tasks, understaffedTasks included, legend with id/label/colorIntent |
| Reviewer recommendation | `room` not `roomArea`, 3 status values matching TasksView, task-id prefix pattern |