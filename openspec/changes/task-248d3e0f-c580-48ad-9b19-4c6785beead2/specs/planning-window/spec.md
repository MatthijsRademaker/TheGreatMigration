# planning-window Specification (Delta)

## Purpose

Update the canonical `openspec/specs/planning-window/spec.md` to reflect the transition from compile-time Go constants to Postgres-backed persistence while preserving the exact API contract.

## MODIFIED Requirements

### Requirement: The planning window SHALL be served from a Postgres-backed query

The requirement "The planning window SHALL be defined by a backend endpoint as the canonical source of truth" is modified: instead of being seeded with compile-time Go constants, the endpoint SHALL read the planning window from a `Store` interface backed by the `planning_windows` database table. The response contract (`startDate`, `endDate`, `days`) and all validation semantics SHALL remain unchanged. The default seed data SHALL preserve the canonical values (`startDate="2026-07-05"`, `endDate="2026-08-13"`, `days=40`).

#### Scenario: Endpoint returns the canonical planning window from DB
- **WHEN** `GET /api/planning-window` is called and seed data is present
- **THEN** the response is 200 OK with `Content-Type: application/json`
- **AND** the body contains `startDate="2026-07-05"`, `endDate="2026-08-13"`, and `days=40`

#### Scenario: Endpoint validates response contract
- **WHEN** the endpoint is called
- **THEN** `startDate` and `endDate` are valid ISO 8601 date strings
- **AND** `startDate` lexicographically precedes `endDate`
- **AND** `days` is a positive integer equal to the inclusive day count

#### Scenario: Endpoint appears in OpenAPI specification
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` includes the `/api/planning-window` endpoint with its response schema

### Requirement: The Store interface SHALL be injected into the handler

The `registerPlanningWindow` function SHALL accept a `Store` parameter. The handler closure SHALL call `store.GetPlanningWindow(ctx)` to retrieve the planning window data instead of reading from compile-time constants `planWindowStart` and `planWindowEnd`. The handler SHALL return a `huma.Error500InternalServerError` if the Store call fails.

#### Scenario: Handler delegates to Store on success
- **WHEN** the handler is called and the Store returns planning window data
- **THEN** the response body matches the Store's returned data

#### Scenario: Handler returns error on Store failure
- **WHEN** the handler is called and the Store returns an error
- **THEN** the response is a 500 Internal Server Error