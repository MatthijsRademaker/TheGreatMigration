## MODIFIED Requirements

### Requirement: Backend SHALL expose backlog task read and write endpoints through the existing BFF contract

The backend SHALL continue to register `GET /api/tasks/backlog` via Huma and SHALL add `POST /api/tasks`, `PUT /api/tasks/{id}`, and `DELETE /api/tasks/{id}` to `/openapi.json`. `GET /api/tasks/backlog` SHALL remain the canonical read model for task rows, summary counts, priority legend, and status legend. Successful creates, updates, and deletes SHALL be observable through subsequent `GET /api/tasks/backlog` responses rather than a separate task read contract.

#### Scenario: OpenAPI includes the task write surface
- **WHEN** the backend OpenAPI document is fetched
- **THEN** it includes `/api/tasks` and `/api/tasks/{id}` write operations alongside `GET /api/tasks/backlog`

#### Scenario: Subsequent backlog reads reflect successful writes
- **WHEN** a task is created, updated, or deleted successfully
- **THEN** a subsequent `GET /api/tasks/backlog` response reflects the resulting task rows and derived summary counts

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