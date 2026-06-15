# task-backlog-api Specification

## Purpose
Defines the first read-only backend API contract for the task backlog, providing seeded task rows and derived summary counts for the homepage dashboard and `/tasks` route.
## Requirements
### Requirement: Backend SHALL expose a Huma-registered read-only endpoint for task backlog data

The backend SHALL register `GET /api/tasks/backlog` via Huma v2, making it appear in the auto-generated `/openapi.json`. The endpoint SHALL be read-only with no query parameters required.

#### Scenario: Endpoint appears in OpenAPI specification
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` includes the `/api/tasks/backlog` endpoint with its response schema and description

#### Scenario: Endpoint returns 200 with no query parameters
- **WHEN** `GET /api/tasks/backlog` is called
- **THEN** the response status is 200 with `Content-Type: application/json`

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
