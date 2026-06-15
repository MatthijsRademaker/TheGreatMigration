# dashboard-people-availability Specification (Delta)

## Purpose

Update the canonical `openspec/specs/dashboard-people-availability/spec.md` to reflect the transition from in-memory seeded data to Postgres-backed persistence while preserving the exact API contract.

## MODIFIED Requirements

### Requirement: The endpoint SHALL read from Postgres-backed queries

The endpoint handler SHALL read people and availability data from a `Store` interface backed by a Postgres database via sqlc-generated queries. The requirement "First slice SHALL use in-memory seeded data" is replaced by this requirement. The endpoint SHALL accept the same query parameters (`start`, `days`), return the same response shape (`DashboardBody` with `range`, `summary`, `people`, `statuses`), and enforce the same validation semantics. The `availableToday` count, status legend, and four-status vocabulary SHALL remain unchanged.

#### Scenario: Endpoint returns DB-backed data matching the contract
- **WHEN** `GET /api/dashboard/people-availability` is called with the default 4-day window and seed data is present in the database
- **THEN** the response contains `range`, `summary`, `people`, and `statuses` top-level fields with the specified sub-fields
- **AND** `summary.availableToday` equals the count of people with status `available` on `range.selectedDate`
- **AND** every person's `availability` array has exactly D entries
- **AND** all status values are canonical (`available`, `busy`, `partial`, `off`)

#### Scenario: Explicit query parameters work with DB data
- **WHEN** `GET /api/dashboard/people-availability?start=2026-07-05&days=3` is called against seeded data
- **THEN** `range.startDate` is `"2026-07-05"`, `range.endDate` is `"2026-07-07"`, and `range.days` is 3

#### Scenario: Seed data includes at least 8 people
- **WHEN** the endpoint is called with any valid date range
- **THEN** `people` contains at least 8 entries with non-empty `id`, `name`, and `initials`

#### Scenario: Seed data exercises all four statuses
- **WHEN** the endpoint is called
- **THEN** across all people and dates in the response, the statuses `available`, `busy`, `partial`, and `off` each appear at least once

### Requirement: The Store interface SHALL be injected into the handler

The `registerDashboardPeopleAvailability` function SHALL accept a `Store` parameter. The handler closure SHALL call `store.GetPeopleAvailability(ctx, startDate, days)` to retrieve data instead of constructing it from in-memory slices and closure-based status functions. The handler SHALL return a `huma.Error500InternalServerError` if the Store call fails.

#### Scenario: Handler delegates to Store on success
- **WHEN** the handler is called and the Store returns valid dashboard data
- **THEN** the response body matches the Store's returned data

#### Scenario: Handler returns error on Store failure
- **WHEN** the handler is called and the Store returns an error
- **THEN** the response is a 500 Internal Server Error