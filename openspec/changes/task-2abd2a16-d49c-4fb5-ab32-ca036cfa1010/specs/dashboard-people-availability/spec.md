# dashboard-people-availability Specification

## MODIFIED Requirements

### Requirement: Backend SHALL expose a Huma-registered endpoint for homepage people availability

The backend SHALL register `GET /api/dashboard/people-availability` via Huma v2, making it appear in the auto-generated `/openapi.json`. The endpoint SHALL accept optional query parameters `start` (ISO 8601 date string, YYYY-MM-DD), `days` (positive integer), `offset` (non-negative integer), and `limit` (non-negative integer). If omitted, `start` SHALL default to the planning window's start date, `days` SHALL default to 4, `offset` SHALL default to 0, and `limit` SHALL default to 0 (no limit — return all people). When `limit` is 0, all people matching the date range SHALL be returned. When `limit` > 0, at most `limit` people SHALL be returned, starting from `offset`. The `days` parameter SHALL include the start date (i.e., `days=4` returns the start date plus the next 3 dates).

#### Scenario: Endpoint appears in OpenAPI specification

- **WHEN** the backend is running
- **THEN** `GET /openapi.json` includes the `/api/dashboard/people-availability` endpoint with its query parameters (`start`, `days`, `offset`, `limit`), response schema, and description

#### Scenario: Default query parameters use the planning window start date with all people

- **WHEN** `GET /api/dashboard/people-availability` is called with no query parameters
- **THEN** the response `range.startDate` equals the planning window's start date, `range.days` is 4, `summary.totalPeople` equals the total count of people in the database, and `people` contains all people

#### Scenario: Explicit query parameters override defaults

- **WHEN** `GET /api/dashboard/people-availability?start=2024-07-02&days=3` is called
- **THEN** `range.startDate` is `"2024-07-02"`, `range.endDate` is `"2024-07-04"`, and `range.days` is 3

#### Scenario: Paginated query returns a subset of people with metadata

- **WHEN** `GET /api/dashboard/people-availability?offset=0&limit=3` is called against seed data with 8 people
- **THEN** the response `people` contains exactly 3 entries, `summary.totalPeople` is 8, and pagination metadata is present

#### Scenario: Offset skips the specified number of people

- **WHEN** `GET /api/dashboard/people-availability?offset=3&limit=3` is called
- **THEN** the response `people` contains 3 entries and the first entry's `id` differs from the first entry when `offset=0`

#### Scenario: limit=0 returns all people

- **WHEN** `GET /api/dashboard/people-availability` is called with no `limit` parameter (defaults to 0) against seed data with 8 people
- **THEN** `people` contains at least 8 entries and `summary.totalPeople` is at least 8

### Requirement: Response SHALL contain range metadata, summary counts, per-person availability, pagination metadata, and a status legend

The JSON response body SHALL include the following top-level fields:

- `range`: an object with `startDate` (string, ISO date), `endDate` (string, ISO date), `days` (integer), and `selectedDate` (string, ISO date, defaults to `startDate`).
- `summary`: an object with `availableToday` (integer, count of people with status `"available"` on `selectedDate`) and `totalPeople` (integer, total number of people matching the query regardless of pagination).
- `people`: an array of objects, each with `id` (string, stable key), `name` (string), `initials` (string), and `availability` (array of `{date: string, status: string}` objects, one per date in the range).
- `pagination`: an object with `totalPeople` (integer, total count of people regardless of pagination), `page` (integer, 1-indexed current page number, computed from offset/limit), and `perPage` (integer, the limit value, or `totalPeople` when no limit was applied).
- `statuses`: an array of canonical status legend objects, each with `id` (string), `label` (string), and `colorIntent` (string).

#### Scenario: Response shape matches the contract

- **WHEN** `GET /api/dashboard/people-availability` returns 200
- **THEN** the JSON body contains `range`, `summary`, `people`, `pagination`, and `statuses` top-level fields with the specified sub-fields

#### Scenario: pagination.totalPeople equals summary.totalPeople

- **WHEN** the response is returned with pagination
- **THEN** `pagination.totalPeople` equals `summary.totalPeople`

#### Scenario: pagination.page is computed from offset and limit

- **WHEN** `offset=0&limit=10` is requested
- **THEN** `pagination.page` is 1
- **WHEN** `offset=10&limit=10` is requested
- **THEN** `pagination.page` is 2

#### Scenario: pagination.perPage equals the requested limit when paginating

- **WHEN** `limit=5` is requested
- **THEN** `pagination.perPage` is 5

#### Scenario: pagination.perPage equals totalPeople when no limit is applied

- **WHEN** no `limit` parameter is provided
- **THEN** `pagination.perPage` equals `summary.totalPeople`

### Requirement: The Store interface SHALL be injected into the handler

The `registerDashboardPeopleAvailability` function SHALL accept a `Store` parameter. The handler closure SHALL call `store.GetPeopleAvailability(ctx, startDate, days, offset, limit)` to retrieve data instead of constructing it from in-memory slices and closure-based status functions. The handler SHALL return a `huma.Error500InternalServerError` if the Store call fails.

#### Scenario: Handler delegates to Store on success

- **WHEN** the handler is called and the Store returns valid dashboard data
- **THEN** the response body matches the Store's returned data including pagination metadata

#### Scenario: Handler returns error on Store failure

- **WHEN** the handler is called and the Store returns an error
- **THEN** the response is a 500 Internal Server Error

### Requirement: Backend tests SHALL cover pagination scenarios

Backend tests SHALL include test functions that:

- Send `GET /api/dashboard/people-availability?offset=0&limit=3` and assert exactly 3 people are returned with correct pagination metadata.
- Send `GET /api/dashboard/people-availability?offset=10&limit=3` and assert empty people array when offset exceeds total count.
- Send `GET /api/dashboard/people-availability` (no pagination params) and assert all people are returned with `pagination.page` of 1 and `pagination.perPage` equal to total count.

#### Scenario: Pagination limit test passes

- **WHEN** `go test ./...` runs in `backend/`
- **THEN** the pagination limit test passes

#### Scenario: Pagination offset-beyond-count test passes

- **WHEN** `go test ./...` runs in `backend/`
- **THEN** the offset-beyond-count test returns a 200 with an empty people array and correct `summary.totalPeople`

#### Scenario: Existing tests still pass

- **WHEN** `go test ./...` runs in `backend/`
- **THEN** all existing test functions for the dashboard endpoint still pass
