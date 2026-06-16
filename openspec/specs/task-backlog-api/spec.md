# task-backlog-api Specification

## Purpose
Defines the first read-only backend API contract for the task backlog, providing seeded task rows and derived summary counts for the homepage dashboard and `/tasks` route.
## Requirements
### Requirement: Backend SHALL expose a Huma-registered read-only endpoint for task backlog data

The backend SHALL continue to register `GET /api/tasks/backlog` via Huma and SHALL add `POST /api/tasks`, `PUT /api/tasks/{id}`, and `DELETE /api/tasks/{id}` to `/openapi.json`. `GET /api/tasks/backlog` SHALL remain the canonical read model for task rows, summary counts, priority legend, and status legend. Successful creates, updates, and deletes SHALL be observable through subsequent `GET /api/tasks/backlog` responses rather than a separate task read contract.

#### Scenario: OpenAPI includes the task write surface

- **WHEN** the backend OpenAPI document is fetched
- **THEN** it includes `/api/tasks` and `/api/tasks/{id}` write operations alongside `GET /api/tasks/backlog`

#### Scenario: Subsequent backlog reads reflect successful writes

- **WHEN** a task is created, updated, or deleted successfully
- **THEN** a subsequent `GET /api/tasks/backlog` response reflects the resulting task rows and derived summary counts

### Requirement: Response SHALL contain summary counts, task rows, priority legend, and status legend

The JSON response body SHALL include the following top-level fields:

- `summary`: an object with `totalTasks` (integer), `highPriorityTasks` (integer), `unassignedTasks` (integer), and `understaffedTasks` (integer).
- `tasks`: an array of task row objects.
- `priorities`: an array of canonical priority legend objects.
- `statuses`: an array of canonical task-status legend objects.

#### Scenario: Response shape matches the contract
- **WHEN** `GET /api/tasks/backlog` returns 200
- **THEN** the JSON body contains `summary`, `tasks`, `priorities`, and `statuses` top-level fields

#### Scenario: Summary counts are consistent with task data
- **WHEN** the response is returned
- **THEN** `summary.totalTasks` equals `len(tasks)`
- **AND** `summary.highPriorityTasks` equals the count of tasks where `priority` equals `"high"`
- **AND** `summary.unassignedTasks` equals the count of tasks where `assignedTo` is an empty array
- **AND** `summary.understaffedTasks` equals the count of tasks where `assignedTo` is non-empty and `len(assignedTo) < peopleNeeded`

### Requirement: Each task row SHALL contain canonical fields

Each object in the `tasks` array SHALL include:

- `id` (string): stable task identifier, prefixed `"task-"`.
- `title` (string): human-readable task description.
- `priority` (string): one of `"high"`, `"medium"`, or `"low"`.
- `peopleNeeded` (integer): number of people required for the task, minimum 1.
- `room` (string): the room or area the task belongs to.
- `status` (string): one of `"backlog"`, `"ready"`, or `"assigned"`.
- `assignedTo` (array of strings): person-ID strings for assigned helpers, may be empty.

#### Scenario: Task rows have all required fields with correct types
- **WHEN** the response is returned
- **THEN** every object in `tasks` has non-empty `id`, `title`, `priority`, `room`, and `status` string fields
- **AND** every task has `peopleNeeded` as a positive integer
- **AND** every task has `assignedTo` as an array (possibly empty)

#### Scenario: Task IDs follow the stable prefix pattern
- **WHEN** the response is returned
- **THEN** every task `id` starts with `"task-"` and is unique within the response

### Requirement: Priority vocabulary SHALL be exactly three canonical values

Every `priority` value in `tasks[]` SHALL be one of exactly three values: `"high"`, `"medium"`, or `"low"`. The `priorities` legend SHALL include all three priorities with their design-system metadata:

| id | label | colorIntent |
| --- | --- | --- |
| `high` | `High` | `destructive` |
| `medium` | `Medium` | `warning` |
| `low` | `Low` | `success` |

#### Scenario: All priority values are canonical
- **WHEN** the response is returned
- **THEN** every `priority` value in every task is one of `high`, `medium`, or `low`

#### Scenario: Priority legend includes all three values with metadata
- **WHEN** the response is returned
- **THEN** `priorities` is an array of exactly three objects with `id`, `label`, and `colorIntent` fields matching the canonical design-system vocabulary

### Requirement: Task status vocabulary SHALL be exactly three canonical values

Every `status` value in `tasks[]` SHALL be one of exactly three values: `"backlog"`, `"ready"`, or `"assigned"`. The `statuses` legend SHALL include all three statuses with their metadata:

| id | label | colorIntent |
| --- | --- | --- |
| `backlog` | `Backlog` | `muted` |
| `ready` | `Ready` | `info` |
| `assigned` | `Assigned` | `success` |

#### Scenario: All status values are canonical
- **WHEN** the response is returned
- **THEN** every `status` value in every task is one of `backlog`, `ready`, or `assigned`

#### Scenario: Status legend includes all three values with metadata
- **WHEN** the response is returned
- **THEN** `statuses` is an array of exactly three objects with `id`, `label`, and `colorIntent` fields matching the canonical status vocabulary

### Requirement: Seed data SHALL be in-memory and exercise all vocabulary values

The endpoint handler SHALL read task backlog data from a `Store` interface backed by Postgres tables and sqlc-generated queries. Seed data inserted through migrations SHALL reproduce the current deterministic backlog payload exactly enough to preserve:

- the stable task IDs `task-1` through `task-11`;
- all three priority values (`high`, `medium`, `low`);
- all three backlog status values (`backlog`, `ready`, `assigned`);
- at least four distinct room names; and
- the current mix of empty, partially filled, and fully filled `assignedTo` arrays relative to `peopleNeeded`.

The response summary SHALL continue to be derived from the returned task rows rather than stored as separate persisted counts.

#### Scenario: Seed data preserves the current deterministic backlog rows
- **WHEN** `GET /api/tasks/backlog` is called against seeded Postgres data
- **THEN** `tasks` contains the same stable IDs and seeded assignment variety currently defined by the in-memory backlog seed

#### Scenario: Database-backed rows still exercise all canonical vocabulary values
- **WHEN** the endpoint is called
- **THEN** across all returned tasks, the priorities `high`, `medium`, and `low` each appear at least once
- **AND** the statuses `backlog`, `ready`, and `assigned` each appear at least once
- **AND** at least one task has an empty `assignedTo` array
- **AND** at least one task has a non-empty `assignedTo` with `len(assignedTo) < peopleNeeded`

### Requirement: Existing behavior SHALL remain intact

The task backlog endpoint SHALL remain read-only and additive. `GET /api/hello`, `GET /api/planning-window`, and `GET /api/dashboard/people-availability` SHALL continue to return their expected responses, and `GET /api/tasks/backlog` SHALL preserve its existing JSON contract while switching to Postgres-backed storage.

#### Scenario: Task backlog contract remains stable after persistence wiring
- **WHEN** `GET /api/tasks/backlog` is called after the refactor
- **THEN** the response still contains `summary`, `tasks`, `priorities`, and `statuses`
- **AND** the task rows and derived summary counts still satisfy the existing contract invariants

### Requirement: Backend tests SHALL cover the new endpoint

Backend tests in `backend/main_test.go` SHALL continue to verify the existing task backlog contract, and they SHALL also cover Store-backed success and failure paths for the backlog handler. Integration tests SHALL validate the backlog contract against the seeded Postgres database.

#### Scenario: MockStore-backed task backlog tests pass
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** task backlog tests verify the existing contract invariants using a Store-backed handler registration

#### Scenario: Real-Postgres task backlog integration test passes
- **WHEN** `go test -tags=integration ./...` runs in `backend/`
- **THEN** the integration suite verifies the seeded backlog payload and summary invariants through `GET /api/tasks/backlog`

### Requirement: The Store interface SHALL be injected into the backlog handler

The `registerTasksBacklog` function SHALL accept a `Store` parameter. The handler closure SHALL call `store.GetTaskBacklog(ctx)` and SHALL return a `huma.Error500InternalServerError` if the Store call fails.

#### Scenario: Backlog handler delegates to Store on success
- **WHEN** the handler is called and the Store returns backlog data
- **THEN** the response body matches the Store-backed backlog payload

#### Scenario: Backlog handler returns 500 on Store failure
- **WHEN** the handler is called and the Store returns an error
- **THEN** the response is a 500 Internal Server Error

### Requirement: Task writes SHALL use canonical task fields and validate assignment references

Create and update requests SHALL accept `title`, `priority`, `peopleNeeded`, `room`, `status`, and `assignedTo`. The backend SHALL reject empty titles, priorities outside `high`, `medium`, and `low`, statuses outside `backlog`, `ready`, and `assigned`, `peopleNeeded < 1`, missing room values, and `assignedTo` values that do not reference existing `people.id` records. Updating or deleting an unknown task ID SHALL return `404`.

#### Scenario: Invalid task input is rejected

- **WHEN** a create or update request includes an empty title, non-canonical priority or status, `peopleNeeded < 1`, a missing room, or an unknown assigned person ID
- **THEN** the endpoint returns `400`

#### Scenario: Missing task IDs return not found

- **WHEN** `PUT /api/tasks/{id}` or `DELETE /api/tasks/{id}` addresses a task ID that does not exist
- **THEN** the endpoint returns `404`

### Requirement: Store-backed persistence SHALL create, update, and delete backlog tasks transactionally without weakening the current backlog contract

The Store interface, sqlc queries, and Postgres-backed implementation SHALL support `CreateTask`, `UpdateTask`, and `DeleteTask` over `backlog_tasks` and `backlog_task_assignments`. Task creation SHALL generate a stable server-assigned `task-*` ID and append the task at the next `sort_order`. Task updates SHALL replace the full assignment set transactionally from `assignedTo`. Task deletes SHALL remove both the task row and its assignment rows in the same transaction. The existing summary invariants, canonical legends, and task row field contract returned by `GET /api/tasks/backlog` SHALL remain unchanged.

#### Scenario: Created tasks receive stable backend IDs and appear in backlog reads

- **WHEN** `POST /api/tasks` succeeds
- **THEN** the created task has an `id` prefixed with `task-`
- **AND** a subsequent `GET /api/tasks/backlog` includes that task in the returned rows

#### Scenario: Assignment replacement stays consistent with the read model

- **WHEN** `PUT /api/tasks/{id}` updates `assignedTo`
- **THEN** the persisted assignment rows are replaced transactionally
- **AND** a subsequent `GET /api/tasks/backlog` returns the updated `assignedTo` values with summary counts still derived consistently from the returned tasks

#### Scenario: Deleting a task removes its assignment rows

- **WHEN** `DELETE /api/tasks/{id}` succeeds
- **THEN** subsequent backlog reads no longer include the task
- **AND** no orphaned `backlog_task_assignments` rows remain for that task

### Requirement: The committed frontend API artifacts SHALL reflect the task write surface

When the backend contract gains task write endpoints, the committed OpenAPI snapshot and generated frontend client artifacts SHALL be refreshed so task management can call the write surface type-safely without depending on a live backend during normal verification.

#### Scenario: Snapshot-backed generation includes task write operations

- **WHEN** the frontend API artifacts are refreshed for this change
- **THEN** `frontend/openapi-snapshot.json` includes `/api/tasks` and `/api/tasks/{id}` write paths
- **AND** the committed generated client under `frontend/src/client/` includes typed operations for task create, update, and delete

### Requirement: Backend verification SHALL cover CRUD, validation, assignment persistence, and contract preservation

Backend tests SHALL cover task CRUD success paths, `400` validation failures, `404` missing task IDs, assignment persistence and replacement, delete behavior, OpenAPI path inclusion, and preservation of the existing hello, planning-window, dashboard, and task-backlog read contracts. Integration tests SHALL verify the seeded backlog payload plus successful persisted task writes against Postgres.

#### Scenario: CRUD contract tests pass alongside existing read-path tests

- **WHEN** the backend test suite runs
- **THEN** task write-surface tests pass alongside the existing backlog-read contract tests

#### Scenario: OpenAPI includes the committed task write surface

- **WHEN** backend tests inspect `/openapi.json`
- **THEN** the document includes `/api/tasks` and `/api/tasks/{id}` write paths
