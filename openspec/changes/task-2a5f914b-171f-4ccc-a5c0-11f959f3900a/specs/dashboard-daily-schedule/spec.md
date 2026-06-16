## MODIFIED Requirements

### Requirement: Each schedule task card SHALL expose canonical dashboard fields

Every task object in `days[].tasks[]` SHALL continue to include `title`, `priority`, `roomArea`, `assignedPeople`, `peopleNeeded`, `assignedCount`, and `staffingStatus`, and `assignedCount` SHALL continue to equal `len(assignedPeople)` without exceeding `peopleNeeded`. The `id` field SHALL continue to expose a stable persisted schedule-card identifier derived from `schedule_task_cards.id` and formatted as a string prefixed `sched-`; synthetic per-request IDs SHALL NOT be used. When a schedule card references a backlog task, the response SHALL additionally include a `taskId` field containing the referenced backlog task's ID (e.g., `"task-3"`). If the card has no backlog reference, `taskId` SHALL be `null`.

#### Scenario: Response task cards keep canonical fields and invariants

- **WHEN** a task card is returned in the daily-schedule response
- **THEN** it includes all canonical dashboard fields
- **AND** its assignment and staffing invariants remain internally consistent

#### Scenario: Task card identifiers are stable across identical reads

- **WHEN** two identical daily-schedule requests are made for the same persisted data
- **THEN** the returned schedule card identifiers are identical in both responses
- **AND** each returned `id` is unique within the response

#### Scenario: Created cards appear in later reads with their persisted identifier

- **WHEN** a schedule card is created successfully
- **THEN** a subsequent `GET /api/dashboard/daily-schedule` for the affected window returns that card
- **AND** the returned task card `id` matches the persisted `sched-<id>` value created by the backend

#### Scenario: Referenced cards include the linked taskId

- **WHEN** a schedule card was created with a `taskId` reference
- **THEN** the daily-schedule response includes `"taskId": "task-<id>"` in that card's object
- **AND** the `taskId` matches the backlog task referenced at creation time

#### Scenario: Unreferenced cards return null taskId

- **WHEN** a schedule card has no backlog task reference
- **THEN** the daily-schedule response includes `"taskId": null` in that card's object

### Requirement: Backend SHALL expose Store-backed schedule-card CRUD endpoints

The backend SHALL register schedule-card write operations alongside `GET /api/dashboard/daily-schedule`:

- `POST /api/schedule/cards`
- `PUT /api/schedule/cards/{id}`
- `DELETE /api/schedule/cards/{id}`

Create and update requests SHALL accept the full persisted schedule-card shape: scheduled date, title, priority, room or area, people needed, and the complete assigned-person ID set. Requests SHALL additionally accept an optional `taskId` string field. When `taskId` is provided and no explicit value is supplied for `title`, `priority`, `roomArea`, or `peopleNeeded`, those fields SHALL be inherited from the referenced backlog task at creation time. If explicit values are supplied alongside `taskId`, the explicit values take precedence. Successful writes SHALL be observable through subsequent `GET /api/dashboard/daily-schedule` responses for the affected date window.

#### Scenario: OpenAPI includes the schedule write surface

- **WHEN** the backend OpenAPI document is fetched after this change
- **THEN** it includes `POST /api/schedule/cards`, `PUT /api/schedule/cards/{id}`, and `DELETE /api/schedule/cards/{id}` alongside the existing daily-schedule read path

#### Scenario: Subsequent reads reflect successful schedule writes

- **WHEN** a schedule card is created, updated, or deleted successfully
- **THEN** a subsequent `GET /api/dashboard/daily-schedule` response for the affected window reflects the resulting persisted schedule state

#### Scenario: Unknown schedule-card identifiers return not found

- **WHEN** `PUT /api/schedule/cards/{id}` or `DELETE /api/schedule/cards/{id}` addresses a schedule card that does not exist
- **THEN** the endpoint returns `404`

#### Scenario: Schedule card created with taskId inherits task fields

- **WHEN** `POST /api/schedule/cards` includes `"taskId": "task-3"` and no explicit title, priority, roomArea, or peopleNeeded
- **THEN** the created card's title, priority, roomArea, and peopleNeeded are copied from the referenced backlog task at creation time
- **AND** subsequent reads return the inherited values

#### Scenario: Explicit fields override taskId inheritance

- **WHEN** `POST /api/schedule/cards` includes both `"taskId": "task-3"` and an explicit `"peopleNeeded": 5`
- **THEN** the explicit `peopleNeeded` value is used instead of the referenced task's value
- **AND** inherited fields without explicit overrides still come from the referenced task

### Requirement: Schedule writes SHALL validate canonical fields and planning-window constraints

Create and update requests SHALL reject empty titles, non-canonical priorities, missing room or area values, `peopleNeeded < 1`, unknown assigned person IDs, assignment counts greater than `peopleNeeded`, and scheduled dates outside the planning window. When `taskId` is provided, the backend SHALL additionally validate that the referenced backlog task exists. Validation failures SHALL not partially persist schedule cards or assignment rows.

#### Scenario: Invalid schedule input is rejected

- **WHEN** a create or update request includes an empty title, a missing room or area value, a non-canonical priority, or `peopleNeeded < 1`
- **THEN** the request is rejected
- **AND** no partial card or assignment write is committed

#### Scenario: Unknown assignee identifiers are rejected

- **WHEN** a create or update request includes assigned-person IDs that do not reference persisted people
- **THEN** the request is rejected
- **AND** the persisted assignment set is unchanged

#### Scenario: Out-of-window scheduled dates are rejected

- **WHEN** a create or update request targets a date outside the current planning window
- **THEN** the request is rejected
- **AND** the schedule read model is unchanged

#### Scenario: Unknown taskId is rejected

- **WHEN** a create or update request includes a `taskId` that does not reference an existing backlog task
- **THEN** the request is rejected with a 400 error
- **AND** no partial card or assignment write is committed

### Requirement: Store-backed persistence SHALL create, update, and delete schedule cards transactionally

The Store interface, sqlc queries, and Postgres-backed implementation SHALL support schedule-card create, update, and delete over `schedule_task_cards` and `schedule_task_assignments`. Creates SHALL return the stable persisted card identifier used by subsequent reads. Updates SHALL replace the full assignment set transactionally and MAY change the scheduled date as part of the same persisted card. Deletes SHALL remove the card and its assignment rows without touching backlog tasks. When a `taskId` is provided at creation or update, the Store SHALL persist the `task_id` foreign key reference in `schedule_task_cards`. Successful writes SHALL invalidate the frontend daily-schedule query so refreshed reads show the persisted result.

#### Scenario: Create returns a stable persisted schedule-card identifier

- **WHEN** `POST /api/schedule/cards` succeeds
- **THEN** the created schedule card is addressable by a stable `sched-<id>` identifier
- **AND** later daily-schedule reads return that same identifier

#### Scenario: Update replaces assignments and may move the card to a new date

- **WHEN** `PUT /api/schedule/cards/{id}` succeeds with a new assignee set or scheduled date
- **THEN** the persisted assignment rows are replaced transactionally
- **AND** a later daily-schedule read returns the updated assignments and date placement

#### Scenario: Delete removes the card and its assignment rows

- **WHEN** `DELETE /api/schedule/cards/{id}` succeeds
- **THEN** subsequent daily-schedule reads no longer include the card
- **AND** no orphaned `schedule_task_assignments` rows remain for that card

#### Scenario: Backlog task CRUD remains separate

- **WHEN** schedule-card writes are implemented
- **THEN** backlog task create, update, and delete behavior remains a separate concern
- **AND** schedule writes do not call backlog-task persistence paths

#### Scenario: Created schedule card persists the taskId reference

- **WHEN** a schedule card is created with a `taskId` value
- **THEN** the `task_id` column in `schedule_task_cards` stores the referenced backlog task's ID
- **AND** subsequent daily-schedule reads include that `taskId` in the card's response

#### Scenario: Deleting a referenced task is rejected if schedule cards reference it

- **WHEN** a `DELETE /api/tasks/{id}` targets a backlog task that has referencing schedule cards
- **THEN** the delete is rejected with an error message indicating the task has scheduled cards
- **AND** the backlog task row and its referencing schedule cards remain intact

### Requirement: Generated API artifacts and verification SHALL cover the schedule write surface

When the schedule write contract is added, the committed OpenAPI snapshot and generated frontend client artifacts SHALL be refreshed so the frontend consumes typed schedule queries and mutations without ad-hoc fetches. Backend unit and integration tests SHALL cover read-contract preservation plus schedule CRUD success, validation, not-found, transactional assignment behavior, `taskId` reference creation, and `taskId` validation. Frontend tests SHALL cover mutation-driven refresh behavior on `/calendar` and the task-selection modal workflow.

#### Scenario: Committed API artifacts include schedule mutations

- **WHEN** the OpenAPI snapshot and generated frontend client are refreshed for this change
- **THEN** the committed artifacts include typed schedule-card mutation operations for the new endpoints
- **AND** the `CreateScheduleCardRequestBody` includes an optional `taskId` string field

#### Scenario: Backend tests cover read and CRUD flows together

- **WHEN** backend tests run after this change
- **THEN** they verify the existing daily-schedule read contract and the new schedule-card CRUD flows
- **AND** they cover `taskId` creation, explicit-field override, and unknown-taskId validation

#### Scenario: Calendar mutation flows refresh the rendered schedule

- **WHEN** a `/calendar` schedule mutation succeeds in the frontend
- **THEN** the daily-schedule query is invalidated or refreshed
- **AND** the rendered board updates without a manual page reload

## ADDED Requirements

### Requirement: Backend SHALL include taskId in schedule-card responses

The `GET /api/dashboard/daily-schedule` response's task card objects SHALL include a `taskId` field. When the card's `schedule_task_cards.task_id` column is non-null, the field SHALL contain the referenced backlog task's ID string. When the column is null, the field SHALL be `null` or omitted.

#### Scenario: Task card with taskId returns the reference

- **WHEN** a schedule card has a non-null `task_id` in the database
- **THEN** the daily-schedule response includes `"taskId": "task-<id>"` in that card's object

#### Scenario: Task card without taskId returns null

- **WHEN** a schedule card has a null `task_id` in the database
- **THEN** the daily-schedule response includes `"taskId": null` in that card's object
