# task-backlog-api Specification

## Purpose
Defines the first read-only backend API contract for the task backlog, providing seeded task rows and derived summary counts for the homepage dashboard and `/tasks` route.

## ADDED Requirements

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

The endpoint handler SHALL use hardcoded in-memory seed data with at least 10 tasks. The seed SHALL exercise all three priority values, all three status values, at least 4 distinct `room` values, and a mix of `assignedTo` states (empty, partially filled, fully filled relative to `peopleNeeded`). The seed SHALL include at least one task with an empty `assignedTo` array and at least one task where `len(assignedTo) > 0` and `len(assignedTo) < peopleNeeded`.

#### Scenario: Seed data includes at least 10 tasks with stable IDs
- **WHEN** the endpoint is called
- **THEN** `tasks` contains at least 10 entries, each with a non-empty `id` starting with `"task-"`

#### Scenario: Seed data exercises all priority and status values
- **WHEN** the endpoint is called
- **THEN** across all tasks in the response, the priorities `high`, `medium`, and `low` each appear at least once
- **AND** the statuses `backlog`, `ready`, and `assigned` each appear at least once

#### Scenario: Seed data exercises assignment variety
- **WHEN** the endpoint is called
- **THEN** at least one task has an empty `assignedTo` array
- **AND** at least one task has a non-empty `assignedTo` with `len(assignedTo) < peopleNeeded`

#### Scenario: Seed data includes at least 4 distinct rooms
- **WHEN** the endpoint is called
- **THEN** the set of distinct `room` values across all tasks has at least 4 entries

### Requirement: Existing behavior SHALL remain intact

The new endpoint SHALL be purely additive. `GET /api/hello` SHALL continue to return `{"message": "Hello from the backend!"}` with status 200. `GET /api/dashboard/people-availability` and `GET /api/planning-window` SHALL continue to return their expected responses. CORS SHALL continue to allow origins `http://localhost:5173` and `http://frontend:5173`. `GET /openapi.json` SHALL include all endpoints.

#### Scenario: Hello endpoint is unchanged
- **WHEN** `GET /api/hello` is called after the new endpoint is registered
- **THEN** the response is 200 with body `{"message": "Hello from the backend!"}` and `Content-Type: application/json`

#### Scenario: People-availability endpoint is unchanged
- **WHEN** `GET /api/dashboard/people-availability` is called after the new endpoint is registered
- **THEN** the response is 200 with `range`, `summary`, `people`, and `statuses` top-level fields

#### Scenario: Planning-window endpoint is unchanged
- **WHEN** `GET /api/planning-window` is called after the new endpoint is registered
- **THEN** the response is 200 with `startDate`, `endDate`, and `days` fields

#### Scenario: OpenAPI includes all endpoints
- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document contains paths `/api/hello`, `/api/dashboard/people-availability`, `/api/planning-window`, and `/api/tasks/backlog`

### Requirement: Backend tests SHALL cover the new endpoint

Backend tests in `backend/main_test.go` SHALL include a test function `TestTaskBacklog` that:

- Sends `GET /api/tasks/backlog` and asserts 200 OK with `Content-Type: application/json`.
- Unmarshals the response and verifies top-level fields `summary`, `tasks`, `priorities`, `statuses` are present.
- Asserts `summary.totalTasks` equals `len(tasks)`.
- Asserts `summary.highPriorityTasks` equals the actual count of tasks with `priority == "high"`.
- Asserts `summary.unassignedTasks` equals the actual count of tasks with empty `assignedTo`.
- Asserts `summary.understaffedTasks` equals the actual count where `len(assignedTo) > 0` and `len(assignedTo) < peopleNeeded`.
- Asserts all priority and status values across tasks are canonical.
- Asserts existing `TestHelloEndpoint`, `TestDashboardPeopleAvailability`, and `TestPlanningWindowEndpoint` still pass.

#### Scenario: Task backlog test passes
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** `TestTaskBacklog` passes

#### Scenario: Existing tests still pass
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** `TestHelloEndpoint`, `TestDashboardPeopleAvailability`, and `TestPlanningWindowEndpoint` all pass
