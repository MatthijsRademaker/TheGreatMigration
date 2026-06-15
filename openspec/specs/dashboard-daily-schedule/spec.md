# dashboard-daily-schedule Specification

## Purpose
TBD - created by archiving change task-6d5fb009-6081-4806-98e4-2fee95f68972. Update Purpose after archive.
## Requirements
### Requirement: Backend SHALL expose a Huma-registered endpoint for homepage daily schedule

The backend SHALL continue to register `GET /api/dashboard/daily-schedule` via Huma v2 with optional `start` and `days` query parameters. If `start` is omitted, the endpoint SHALL default to the current planning-window `startDate` returned by the Store-backed planning-window data source. With the seeded demo data, that default remains `2026-07-05`. `days` SHALL continue to default to `4` and remain inclusive of the start date.

#### Scenario: Default query parameters use the planning-window default start
- **WHEN** `GET /api/dashboard/daily-schedule` is called with no query parameters against the seeded database
- **THEN** the response `range.startDate` is `2026-07-05`
- **AND** the response `range.days` is `4`
- **AND** the response `range.endDate` is `2026-07-08`

#### Scenario: Explicit query parameters override defaults
- **WHEN** `GET /api/dashboard/daily-schedule?start=2026-07-10&days=3` is called
- **THEN** `range.startDate` is `2026-07-10`
- **AND** `range.endDate` is `2026-07-12`
- **AND** `range.days` is `3`

### Requirement: Response SHALL contain range metadata and per-day schedule columns

The JSON response body SHALL include a top-level `range` object and a top-level `days` array.

- `range` SHALL contain `startDate` (string, ISO date), `endDate` (string, ISO date), and `days` (integer).
- `days` SHALL contain exactly one object per date in the requested inclusive window, ordered from `startDate` through `endDate`.
- Each day object SHALL contain `date` (ISO date string), `label` (human-readable day header), `availablePeopleCount` (integer), and `tasks` (ordered array).

#### Scenario: Response shape matches the contract
- **WHEN** `GET /api/dashboard/daily-schedule` returns `200`
- **THEN** the JSON body contains `range` and `days` top-level fields with the specified sub-fields

#### Scenario: Every requested date produces one day object
- **WHEN** the response range spans `D` days
- **THEN** `days` contains exactly `D` day objects
- **AND** each day object's `date` value is unique within the response

### Requirement: Each schedule task card SHALL expose canonical dashboard fields

Every task object in `days[].tasks[]` SHALL contain:

- `id` (string, stable for identical requests)
- `title` (string)
- `priority` (one of `high`, `medium`, `low`)
- `roomArea` (string)
- `assignedPeople` (array of `{id, name, initials}` objects)
- `peopleNeeded` (integer, greater than or equal to `1`)
- `assignedCount` (integer)
- `staffingStatus` (one of `fullyStaffed`, `underStaffed`)

Derived invariants SHALL hold for every task card:

- `assignedCount` equals `len(assignedPeople)`
- `assignedCount` never exceeds `peopleNeeded`
- `staffingStatus` is `fullyStaffed` when `assignedCount == peopleNeeded`
- `staffingStatus` is `underStaffed` when `assignedCount < peopleNeeded`

#### Scenario: Task-card fields are complete and canonical
- **WHEN** a task card is returned in the response
- **THEN** it includes all required fields
- **AND** `priority` is one of `high`, `medium`, or `low`
- **AND** `staffingStatus` is one of `fullyStaffed` or `underStaffed`

#### Scenario: Staffing invariants are internally consistent
- **WHEN** the response is returned
- **THEN** every task card's `assignedCount` equals the number of entries in `assignedPeople`
- **AND** no task card has `assignedCount` greater than `peopleNeeded`

### Requirement: Daily helper counts and assignee identities SHALL stay consistent with the people-availability endpoint

`availablePeopleCount` for each returned day SHALL be derived from the persisted availability data by counting rows whose status is `available` for that date. Task assignees in `assignedPeople` SHALL resolve through persisted person records so the returned `id`, `name`, and `initials` remain consistent with the people identities already exposed by `GET /api/dashboard/people-availability`.

#### Scenario: Daily helper counts use the availability table
- **WHEN** the schedule response is returned for a given date window
- **THEN** each day's `availablePeopleCount` equals the count of persisted availability rows for that date whose status is `available`

#### Scenario: Assignee identities match persisted people identities
- **WHEN** a task card includes assigned people
- **THEN** each assigned person `id` references a persisted person row
- **AND** the returned `name` and `initials` match that person identity

### Requirement: First slice SHALL use deterministic in-memory seeded schedule data

The endpoint handler SHALL read schedule task-card data from a `Store` interface backed by Postgres tables and sqlc-generated queries. Seed data inserted through migrations SHALL reproduce the current deterministic daily-schedule output, including:

- the current four seeded day groups used by the schedule board;
- exact current task-card IDs for identical requests;
- multiple priorities and staffing states across the default four-day response; and
- the current assignee variety and ordering needed for task-card invariants.

The persisted schedule read model SHALL remain separate from the backlog task read model.

#### Scenario: Database-backed schedule data preserves the default design-backed variety
- **WHEN** `GET /api/dashboard/daily-schedule` is called with the default four-day window against seeded Postgres data
- **THEN** the returned tasks include `high`, `medium`, and `low` priorities across the window
- **AND** at least one task card is fully staffed at `2/2`
- **AND** at least one task card is fully staffed at `1/1`
- **AND** at least one task card is understaffed

#### Scenario: Identical requests preserve deterministic task-card IDs and ordering
- **WHEN** two identical requests are made for the same `start` and `days`
- **THEN** the returned day ordering, task ordering, and task-card IDs are the same in both responses

### Requirement: Existing backend behavior SHALL remain intact

The daily-schedule endpoint SHALL be purely additive. `GET /api/hello` SHALL continue to return `{"message": "Hello from the backend!"}` with status `200`. `GET /api/dashboard/people-availability` SHALL continue to return its existing dashboard payload. `GET /api/planning-window` SHALL continue to return the canonical planning window. CORS behavior and `/openapi.json` serving SHALL continue to work unchanged.

#### Scenario: Existing endpoints remain available
- **WHEN** the backend is running after the daily-schedule endpoint is added
- **THEN** `GET /api/hello`, `GET /api/dashboard/people-availability`, and `GET /api/planning-window` each still return successful responses matching their existing contracts

#### Scenario: OpenAPI includes all backend endpoints
- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document contains paths for `/api/hello`, `/api/dashboard/people-availability`, `/api/planning-window`, and `/api/dashboard/daily-schedule`

### Requirement: Backend tests SHALL cover the daily-schedule contract

Backend tests SHALL continue to verify the existing daily-schedule contract and SHALL add Store-backed failure-path coverage for the handler. Integration tests SHALL validate the daily-schedule contract against seeded Postgres data and assert `/api/dashboard/daily-schedule` appears in OpenAPI path coverage.

#### Scenario: MockStore-backed daily-schedule tests pass
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** the daily-schedule tests verify the existing contract invariants using a Store-backed handler registration

#### Scenario: Real-Postgres daily-schedule integration test passes
- **WHEN** `go test -tags=integration ./...` runs in `backend/`
- **THEN** the integration suite verifies the seeded daily-schedule payload and OpenAPI path registration through the Postgres-backed endpoint

### Requirement: The Store interface SHALL be injected into the daily-schedule handler

The `registerDailySchedule` function SHALL accept a `Store` parameter. The handler SHALL:
- parse and validate explicit `start` values as it does today;
- call `store.GetPlanningWindow(ctx)` when `start` is omitted;
- call `store.GetDailySchedule(ctx, startDate, days)` to build the response; and
- return a `huma.Error500InternalServerError` if either Store call fails.

#### Scenario: Daily-schedule handler delegates to Store on success
- **WHEN** the handler is called and the Store returns valid planning-window and schedule data
- **THEN** the response body matches the Store-backed schedule payload

#### Scenario: Daily-schedule handler returns 500 on Store failure
- **WHEN** the handler is called and either Store call returns an error
- **THEN** the response is a 500 Internal Server Error
