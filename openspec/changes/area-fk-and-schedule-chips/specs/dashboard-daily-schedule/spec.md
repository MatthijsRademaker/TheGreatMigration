## MODIFIED Requirements

### Requirement: Each schedule task card SHALL expose canonical dashboard fields

Every task object in `days[].tasks[]` SHALL include `title`, `priority`, `area`, `assignedPeople`, `peopleNeeded`, `assignedCount`, and `staffingStatus`, and `assignedCount` SHALL continue to equal `len(assignedPeople)` without exceeding `peopleNeeded`. The `area` field SHALL be an object `{ id, name }` resolved from the card's `rooms_areas` foreign key; the prior scalar `roomArea` string SHALL NOT be returned. The `id` field SHALL expose a stable persisted schedule-card identifier derived from `schedule_task_cards.id` and formatted as a string prefixed `sched-`; synthetic per-request IDs SHALL NOT be used.

#### Scenario: Response task cards keep canonical fields and invariants

- **WHEN** a task card is returned in the daily-schedule response
- **THEN** it includes all canonical dashboard fields
- **AND** it includes `area` as an object with non-empty `id` and `name`
- **AND** it does not include a `roomArea` string field
- **AND** its assignment and staffing invariants remain internally consistent

#### Scenario: Task card identifiers are stable across identical reads

- **WHEN** two identical daily-schedule requests are made for the same persisted data
- **THEN** the returned schedule card identifiers are identical in both responses
- **AND** each returned `id` is unique within the response

#### Scenario: Created cards appear in later reads with their persisted identifier

- **WHEN** a schedule card is created successfully
- **THEN** a subsequent `GET /api/dashboard/daily-schedule` for the affected window returns that card
- **AND** the returned task card `id` matches the persisted `sched-<id>` value created by the backend

### Requirement: Schedule writes SHALL validate canonical fields and planning-window constraints

The schedule-card create and update endpoints SHALL accept `areaId` (replacing the prior `roomArea` string) and SHALL reject a request whose `areaId` does not reference an existing `rooms_areas` row. `areaId` SHALL remain required alongside the existing canonical-field and planning-window validations.

#### Scenario: Write with a valid area id succeeds

- **WHEN** a schedule card is created or updated with an `areaId` that references an existing `rooms_areas` row and otherwise valid fields
- **THEN** the write succeeds and the persisted card references that area

#### Scenario: Write with an unknown area id is rejected

- **WHEN** a schedule-card create or update is submitted with an `areaId` not present in `rooms_areas`
- **THEN** the request is rejected with a client error
- **AND** no card is created or modified

#### Scenario: Write with a missing area id is rejected

- **WHEN** a schedule-card create or update omits `areaId`
- **THEN** the request is rejected with a 400 error
