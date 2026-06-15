## Why

The homepage design shows a Tasks Backlog table and task-derived summary cards (high-priority, under-staffed, unassigned). The C4 architecture model defines `taskData` as move jobs with priority, status, room, staffing needs, and assignments. The `/tasks` route already references these fields and the `/` dashboard has placeholder summary cards that need backend data. However, no backend endpoint exists for task backlog data, leaving both the dashboard and tasks surfaces reliant on static placeholders.

## What Changes

- Add a new read-only `GET /api/tasks/backlog` Huma-registered endpoint serving seeded in-memory task data.
- The endpoint returns a combined payload with `summary` (totalTasks, highPriorityTasks, unassignedTasks, understaffedTasks), `tasks` (backlog rows), `priorities` (canonical priority legend), and `statuses` (canonical task-status legend).
- Each task row includes: `id`, `title`, `priority` (high/medium/low), `peopleNeeded`, `room`, `status` (backlog/ready/assigned), and `assignedTo` (array of person-ID strings).
- Seed data covers at least 10 tasks across all priority and status values, at least 4 distinct rooms, and a mix of assignment states (empty, partially filled, fully filled).
- Backend tests verify response shape, derived summary consistency, canonical vocabulary, OpenAPI publication, and non-regression for existing endpoints.
- Implementation follows the established `backend/dashboard.go` pattern: co-located types, in-memory seed, Huma handler, and registration from `main.go`.

## Impact

- **Affected specs**: New `task-backlog-api` capability.
- **Affected code**: New `backend/tasks.go` (handler, types, seed data), modified `backend/main.go` (registration), modified `backend/main_test.go` (new test function + `newTestAPI` update).
- **Downstream**: `frontend/openapi-snapshot.json` will need regeneration to include the new endpoint; this is explicitly deferred to a follow-up task.
- **No regression**: Existing `/api/hello`, `/api/dashboard/people-availability`, `/api/planning-window`, CORS, and OpenAPI paths remain unchanged.
- **Additive only**: No authentication, persistence, or new dependencies introduced.
