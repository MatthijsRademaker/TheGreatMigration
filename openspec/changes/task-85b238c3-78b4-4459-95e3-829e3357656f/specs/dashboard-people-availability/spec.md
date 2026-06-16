## MODIFIED Requirements

### Requirement: Backend SHALL expose a Huma-registered endpoint for homepage people availability

The backend SHALL register `GET /api/dashboard/people-availability` via Huma v2, making it appear in the auto-generated `/openapi.json`. The endpoint SHALL accept optional query parameters `start` (ISO 8601 date string, YYYY-MM-DD) and `days` (positive integer). If omitted, `start` SHALL default to the **planning window's start date** (fetched from `store.GetPlanningWindow(ctx)`) and `days` SHALL default to 4. If `start` is omitted and no planning window exists in the database, the endpoint SHALL return `400 Bad Request` with a descriptive message. The `days` parameter SHALL include the start date (i.e., `days=4` returns the start date plus the next 3 dates). When an explicit `start` value is provided, it SHALL be used as-is without consulting the planning window.

#### Scenario: Endpoint appears in OpenAPI specification

- **WHEN** the backend is running
- **THEN** `GET /openapi.json` includes the `/api/dashboard/people-availability` endpoint with its query parameters, response schema, and description

#### Scenario: Default query parameters use the planning window start date

- **WHEN** `GET /api/dashboard/people-availability` is called with no query parameters and a planning window exists with `startDate=2026-07-05`
- **THEN** the response `range.startDate` is `"2026-07-05"`, `range.endDate` is `"2026-07-08"`, and `range.days` is 4

#### Scenario: Explicit `start` query parameter overrides the planning window

- **WHEN** `GET /api/dashboard/people-availability?start=2024-07-02&days=3` is called
- **THEN** `range.startDate` is `"2024-07-02"`, `range.endDate` is `"2024-07-04"`, and `range.days` is 3
- **AND** the endpoint does not consult the planning window

#### Scenario: Missing planning window with no explicit `start` returns 400

- **WHEN** `GET /api/dashboard/people-availability` is called with no query parameters and no planning window row exists in the database
- **THEN** the endpoint returns `400 Bad Request` with an error message indicating that no planning window is configured

### Requirement: Backend handler status legend, summary invariants, and availability padding SHALL remain unchanged

All response invariants from the existing specification — four canonical status values, availability padding for missing DB rows (`"off"` default), summary consistency, and the status legend — SHALL remain identical to the pre-change contract. Only the default-start resolution logic changes.

#### Scenario: Summary card values are consistent with per-person data

- **WHEN** the response is returned
- **THEN** `summary.availableToday` equals the count of entries in `people[].availability` where `date` equals `range.selectedDate` and `status` equals `"available"`

#### Scenario: Every person has an availability entry for each date in the range

- **WHEN** the response contains N people and the range spans D days
- **THEN** each person's `availability` array has exactly D entries, one per date in `[startDate, endDate]`

#### Scenario: All status values are canonical

- **WHEN** the response is returned
- **THEN** every `status` value in every person's `availability` array is one of `available`, `busy`, `partial`, or `off`

#### Scenario: Status legend includes all four statuses with metadata

- **WHEN** the response is returned
- **THEN** `statuses` is an array of exactly four objects with `id`, `label`, and `colorIntent` fields matching the canonical design-system vocabulary

## ADDED Requirements

### Requirement: Dashboard endpoint handler SHALL propagate store errors gracefully

When `store.GetPlanningWindow(ctx)` returns an error and no explicit `start` was provided, the handler SHALL return `500 Internal Server Error`. When `store.GetPlanningWindow(ctx)` succeeds but returns `nil` (no row found) and no explicit `start` was provided, the handler SHALL return `400 Bad Request` with a descriptive error message. Existing store errors from `store.GetPeopleAvailability(ctx, ...)` continue to produce `500 Internal Server Error`.

#### Scenario: Planning-window store error returns 500

- **WHEN** the store's `GetPlanningWindow` returns an error (e.g., database connection failure)
- **THEN** the endpoint returns `500 Internal Server Error`

#### Scenario: Missing planning window returns 400 with descriptive message

- **WHEN** no planning-window row exists and no explicit `start` is provided
- **THEN** the endpoint returns `400 Bad Request` with `"detail"` containing a message explaining that a planning window must be configured
