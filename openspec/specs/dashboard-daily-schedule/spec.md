# dashboard-daily-schedule Specification

## Purpose
TBD - created by archiving change task-6d5fb009-6081-4806-98e4-2fee95f68972. Update Purpose after archive.
## Requirements
### Requirement: Backend SHALL expose a Huma-registered endpoint for homepage daily schedule

The backend SHALL register `GET /api/dashboard/daily-schedule` via Huma v2 so the endpoint appears in the auto-generated `/openapi.json`. The endpoint SHALL be tagged consistently with existing dashboard endpoints. It SHALL accept optional query parameters `start` (ISO 8601 date string, `YYYY-MM-DD`) and `days` (positive integer). If omitted, `start` SHALL default to the canonical planning-window start date `"2026-07-05"` and `days` SHALL default to `4`. The `days` parameter SHALL include the start date.

#### Scenario: Endpoint appears in OpenAPI specification
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` includes the `/api/dashboard/daily-schedule` endpoint with its query parameters, response schema, and description

#### Scenario: Default query parameters produce the default four-day planning window
- **WHEN** `GET /api/dashboard/daily-schedule` is called with no query parameters
- **THEN** the response `range.startDate` is `"2026-07-05"`
- **AND** the response `range.days` is `4`
- **AND** the response `range.endDate` is `"2026-07-08"`

#### Scenario: Explicit query parameters override defaults
- **WHEN** `GET /api/dashboard/daily-schedule?start=2026-07-10&days=3` is called
- **THEN** `range.startDate` is `"2026-07-10"`
- **AND** `range.endDate` is `"2026-07-12"`
- **AND** `range.days` is `3`

#### Scenario: Malformed start dates are rejected
- **WHEN** `GET /api/dashboard/daily-schedule?start=2026-13-99` is called
- **THEN** the response status is `400`

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

`availablePeopleCount` for each day SHALL be derived by counting entries in the existing `seedPeople` data whose status for that date equals `available`, using the same offset-based status logic as the people-availability endpoint. Task assignees in `assignedPeople` SHALL reuse the same stable person IDs, names, and initials already used by `GET /api/dashboard/people-availability`.

#### Scenario: Daily helper counts use the same availability source
- **WHEN** the schedule response is returned for a given date window
- **THEN** each day's `availablePeopleCount` equals the number of seeded people whose availability status for that same date is `available`

#### Scenario: Assignee identities match seeded people identities
- **WHEN** a task card includes assigned people
- **THEN** each assigned person `id` is one of the existing seeded person IDs (`p1` through `p8`)
- **AND** the returned `name` and `initials` match that seeded person identity

### Requirement: First slice SHALL use deterministic in-memory seeded schedule data

The endpoint handler SHALL use deterministic in-memory seed data for schedule tasks. The seeded task set SHALL be generated by day offset from the requested start date rather than being limited to one hard-coded historical calendar week. For the default four-day request, the seed SHALL cover multiple priorities and staffing states, including at least one fully staffed `2/2` task card, at least one fully staffed `1/1` task card, and at least one understaffed task card.

#### Scenario: Default seeded window covers the design-backed variety
- **WHEN** `GET /api/dashboard/daily-schedule` is called with the default four-day window
- **THEN** the returned tasks include `high`, `medium`, and `low` priorities across the window
- **AND** at least one task card is fully staffed at `2/2`
- **AND** at least one task card is fully staffed at `1/1`
- **AND** at least one task card is understaffed

#### Scenario: Identical requests return the same seeded task order
- **WHEN** two identical requests are made for the same `start` and `days`
- **THEN** the returned day ordering and task ordering are the same in both responses

### Requirement: Existing backend behavior SHALL remain intact

The daily-schedule endpoint SHALL be purely additive. `GET /api/hello` SHALL continue to return `{"message": "Hello from the backend!"}` with status `200`. `GET /api/dashboard/people-availability` SHALL continue to return its existing dashboard payload. `GET /api/planning-window` SHALL continue to return the canonical planning window. CORS behavior and `/openapi.json` serving SHALL continue to work unchanged.

#### Scenario: Existing endpoints remain available
- **WHEN** the backend is running after the daily-schedule endpoint is added
- **THEN** `GET /api/hello`, `GET /api/dashboard/people-availability`, and `GET /api/planning-window` each still return successful responses matching their existing contracts

#### Scenario: OpenAPI includes all backend endpoints
- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document contains paths for `/api/hello`, `/api/dashboard/people-availability`, `/api/planning-window`, and `/api/dashboard/daily-schedule`

### Requirement: Backend tests SHALL cover the daily-schedule contract

Backend tests SHALL include coverage that:

- sends `GET /api/dashboard/daily-schedule` and asserts `200 OK` with `Content-Type: application/json`
- verifies default-window range values and explicit `start` / `days` behavior
- verifies malformed `start` returns `400`
- verifies the response shape and one-day-per-requested-date behavior
- verifies `priority` and `staffingStatus` only use canonical values
- verifies `assignedCount`, `assignedPeople`, and `peopleNeeded` invariants
- verifies `/openapi.json` includes `/api/dashboard/daily-schedule`
- preserves existing endpoint regression tests

#### Scenario: Daily-schedule contract test passes
- **WHEN** backend tests are run
- **THEN** the daily-schedule endpoint test passes with all contract assertions satisfied

#### Scenario: Existing backend tests still pass
- **WHEN** backend tests are run after the new endpoint is added
- **THEN** the existing hello, people-availability, and planning-window tests still pass

#### Scenario: Repository verification passes
- **WHEN** `scripts/precommit-run` is executed
- **THEN** the repository verification checks complete without failures
