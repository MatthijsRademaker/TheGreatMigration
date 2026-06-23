## MODIFIED Requirements

### Requirement: Each task row SHALL contain canonical fields

Each object in the `tasks` array SHALL include:

- `id` (string): stable task identifier, prefixed `"task-"`.
- `title` (string): human-readable task description.
- `priority` (string): one of `"high"`, `"medium"`, or `"low"`.
- `peopleNeeded` (integer): number of people required for the task, minimum 1.
- `area` (object): `{ id, name }` resolved from the task's `rooms_areas` foreign key. The prior scalar `room` string SHALL NOT be returned.
- `status` (string): one of `"backlog"`, `"ready"`, or `"assigned"`.
- `assignedTo` (array of strings): person-ID strings for assigned helpers, may be empty.

#### Scenario: Task rows have all required fields with correct types

- **WHEN** the response is returned
- **THEN** every object in `tasks` has non-empty `id`, `title`, `priority`, and `status` string fields
- **AND** every task has `area` as an object with non-empty `id` and `name`
- **AND** no task has a `room` string field
- **AND** every task has `peopleNeeded` as a positive integer
- **AND** every task has `assignedTo` as an array (possibly empty)

#### Scenario: Task IDs follow the stable prefix pattern

- **WHEN** the response is returned
- **THEN** every task `id` starts with `"task-"` and is unique within the response

### Requirement: Task writes SHALL use canonical task fields and validate assignment references

The task create and update endpoints SHALL accept `areaId` (replacing the prior `room` string) and SHALL reject a request whose `areaId` does not reference an existing `rooms_areas` row. `areaId` SHALL remain required alongside the existing canonical-field validations (priority, status, title, peopleNeeded) and the assignment-reference validation.

#### Scenario: Task write with a valid area id succeeds

- **WHEN** a task is created or updated with an `areaId` referencing an existing `rooms_areas` row and otherwise valid fields
- **THEN** the write succeeds and the persisted task references that area

#### Scenario: Task write with an unknown area id is rejected

- **WHEN** a task create or update is submitted with an `areaId` not present in `rooms_areas`
- **THEN** the request is rejected with a client error
- **AND** no task is created or modified

#### Scenario: Task write with a missing area id is rejected

- **WHEN** a task create or update omits `areaId`
- **THEN** the request is rejected with a 400 error
