# dashboard-daily-schedule Specification (Delta)

## Purpose

Extend the daily-schedule backend capability from a read-only demo-aligned surface into a persisted, date-based schedule contract that still powers the existing dashboard board.

## MODIFIED Requirements

### Requirement: Each schedule task card SHALL expose canonical dashboard fields

Every task object in `days[].tasks[]` SHALL continue to include `title`, `priority`, `roomArea`, `assignedPeople`, `peopleNeeded`, `assignedCount`, and `staffingStatus`, and `assignedCount` SHALL continue to equal `len(assignedPeople)` without exceeding `peopleNeeded`. The `id` field SHALL now expose a stable persisted schedule-card identifier derived from `schedule_task_cards.id` and formatted as a string prefixed `sched-`; synthetic per-request IDs SHALL NOT be used.

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

### Requirement: First slice SHALL use deterministic in-memory seeded schedule data

The schedule read model SHALL remain separate from `backlog_tasks`, but `schedule_task_cards` SHALL be anchored by a persisted `scheduled_date` rather than demo-only modulo day grouping. Backend migration and seeded data SHALL preserve the current default four-day variety, ordering expectations, and people-availability or assignee invariants while making cards addressable on real dates.

#### Scenario: Seeded default window still returns date-specific schedule coverage

- **WHEN** `GET /api/dashboard/daily-schedule` is called for the default seeded planning window
- **THEN** the response still contains the seeded multi-day mix of priorities, staffing states, and assignee variety
- **AND** each task card is anchored to a persisted schedule date

#### Scenario: Explicit date windows filter by persisted scheduled dates

- **WHEN** `GET /api/dashboard/daily-schedule` is called with explicit `start` and `days` values
- **THEN** the returned cards are selected from persisted `scheduled_date` values in that inclusive window
- **AND** modulo grouping is not required to place cards on dates

#### Scenario: Backlog reads remain separate from schedule reads and writes

- **WHEN** schedule cards are created, updated, or deleted
- **THEN** those changes are reflected through `GET /api/dashboard/daily-schedule`
- **AND** they do not reuse or overwrite the `backlog_tasks` read or write model

## ADDED Requirements

### Requirement: Backend SHALL expose Store-backed schedule-card CRUD endpoints

The backend SHALL register schedule-card write operations alongside `GET /api/dashboard/daily-schedule`:

- `POST /api/schedule/cards`
- `PUT /api/schedule/cards/{id}`
- `DELETE /api/schedule/cards/{id}`

Create and update requests SHALL accept the full persisted schedule-card shape: scheduled date, title, priority, room or area, people needed, and the complete assigned-person ID set. Successful writes SHALL be observable through subsequent `GET /api/dashboard/daily-schedule` responses for the affected date window.

#### Scenario: OpenAPI includes the schedule write surface

- **WHEN** the backend OpenAPI document is fetched after this change
- **THEN** it includes `POST /api/schedule/cards`, `PUT /api/schedule/cards/{id}`, and `DELETE /api/schedule/cards/{id}` alongside the existing daily-schedule read path

#### Scenario: Subsequent reads reflect successful schedule writes

- **WHEN** a schedule card is created, updated, or deleted successfully
- **THEN** a subsequent `GET /api/dashboard/daily-schedule` response for the affected window reflects the resulting persisted schedule state

#### Scenario: Unknown schedule-card identifiers return not found

- **WHEN** `PUT /api/schedule/cards/{id}` or `DELETE /api/schedule/cards/{id}` addresses a schedule card that does not exist
- **THEN** the endpoint returns `404`

### Requirement: Schedule writes SHALL validate canonical fields and planning-window constraints

Create and update requests SHALL reject empty titles, non-canonical priorities, missing room or area values, `peopleNeeded < 1`, unknown assigned person IDs, assignment counts greater than `peopleNeeded`, and scheduled dates outside the planning window. Validation failures SHALL not partially persist schedule cards or assignment rows.

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

### Requirement: Store-backed persistence SHALL create, update, and delete schedule cards transactionally

The Store interface, sqlc queries, and Postgres-backed implementation SHALL support schedule-card create, update, and delete over `schedule_task_cards` and `schedule_task_assignments`. Creates SHALL return the stable persisted card identifier used by subsequent reads. Updates SHALL replace the full assignment set transactionally and MAY change the scheduled date as part of the same persisted card. Deletes SHALL remove the card and its assignment rows without touching backlog tasks. Successful writes SHALL invalidate the frontend daily-schedule query so refreshed reads show the persisted result.

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

### Requirement: Generated API artifacts and verification SHALL cover the schedule write surface

When the schedule write contract is added, the committed OpenAPI snapshot and generated frontend client artifacts SHALL be refreshed so the frontend consumes typed schedule queries and mutations without ad-hoc fetches. Backend unit and integration tests SHALL cover read-contract preservation plus schedule CRUD success, validation, not-found, and transactional assignment behavior. Frontend tests SHALL cover mutation-driven refresh behavior on `/calendar`.

#### Scenario: Committed API artifacts include schedule mutations

- **WHEN** the OpenAPI snapshot and generated frontend client are refreshed for this change
- **THEN** the committed artifacts include typed schedule-card mutation operations for the new endpoints

#### Scenario: Backend tests cover read and CRUD flows together

- **WHEN** backend tests run after this change
- **THEN** they verify the existing daily-schedule read contract and the new schedule-card CRUD flows

#### Scenario: Calendar mutation flows refresh the rendered schedule

- **WHEN** a `/calendar` schedule mutation succeeds in the frontend
- **THEN** the daily-schedule query is invalidated or refreshed
- **AND** the rendered board updates without a manual page reload
