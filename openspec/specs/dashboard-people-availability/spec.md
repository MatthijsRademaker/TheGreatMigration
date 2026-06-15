# dashboard-people-availability Specification

## Purpose
TBD - created by archiving change task-7ea46653-dc77-4474-b5e1-3579c35ebc0b. Update Purpose after archive.
## Requirements
### Requirement: Backend SHALL expose a Huma-registered endpoint for homepage people availability

The backend SHALL register `GET /api/dashboard/people-availability` via Huma v2, making it appear in the auto-generated `/openapi.json`. The endpoint SHALL accept optional query parameters `start` (ISO 8601 date string, YYYY-MM-DD) and `days` (positive integer). If omitted, `start` SHALL default to the server-local current date and `days` SHALL default to 4. The `days` parameter SHALL include the start date (i.e., `days=4` returns the start date plus the next 3 dates).

#### Scenario: Endpoint appears in OpenAPI specification
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` includes the `/api/dashboard/people-availability` endpoint with its query parameters, response schema, and description

#### Scenario: Default query parameters produce a 4-day window
- **WHEN** `GET /api/dashboard/people-availability` is called with no query parameters
- **THEN** the response `range.days` is 4 and `range.endDate` is `range.startDate` plus 3 days

#### Scenario: Explicit query parameters override defaults
- **WHEN** `GET /api/dashboard/people-availability?start=2024-07-02&days=3` is called
- **THEN** `range.startDate` is `"2024-07-02"`, `range.endDate` is `"2024-07-04"`, and `range.days` is 3

### Requirement: Response SHALL contain range metadata, summary counts, per-person availability, and a status legend

The JSON response body SHALL include the following top-level fields:

- `range`: an object with `startDate` (string, ISO date), `endDate` (string, ISO date), `days` (integer), and `selectedDate` (string, ISO date, defaults to `startDate`).
- `summary`: an object with `availableToday` (integer, count of people with status `"available"` on `selectedDate`) and `totalPeople` (integer, total number of people in the response).
- `people`: an array of objects, each with `id` (string, stable key), `name` (string), `initials` (string), and `availability` (array of `{date: string, status: string}` objects, one per date in the range).
- `statuses`: an array of canonical status legend objects, each with `id` (string), `label` (string), and `colorIntent` (string).

#### Scenario: Response shape matches the contract
- **WHEN** `GET /api/dashboard/people-availability` returns 200
- **THEN** the JSON body contains `range`, `summary`, `people`, and `statuses` top-level fields with the specified sub-fields

#### Scenario: Summary card values are consistent with per-person data
- **WHEN** the response is returned
- **THEN** `summary.availableToday` equals the count of entries in `people[].availability` where `date` equals `range.selectedDate` and `status` equals `"available"`

#### Scenario: Every person has an availability entry for each date in the range
- **WHEN** the response contains N people and the range spans D days
- **THEN** each person's `availability` array has exactly D entries, one per date in `[startDate, endDate]`

### Requirement: Availability status values SHALL be constrained to the design-backed vocabulary

Every `status` value in `people[].availability` entries and in `statuses[]` SHALL be one of exactly four values: `"available"`, `"busy"`, `"partial"`, `"off"`. The `statuses` legend SHALL include all four statuses with their design-system metadata:

| id | label | colorIntent |
| --- | --- | --- |
| `available` | `Available` | `success` |
| `busy` | `Busy` | `destructive` |
| `partial` | `Partial` | `warning` |
| `off` | `Off` | `muted` |

#### Scenario: All status values are canonical
- **WHEN** the response is returned
- **THEN** every `status` value in every person's `availability` array is one of `available`, `busy`, `partial`, or `off`

#### Scenario: Status legend includes all four statuses with metadata
- **WHEN** the response is returned
- **THEN** `statuses` is an array of exactly four objects with `id`, `label`, and `colorIntent` fields matching the canonical design-system vocabulary

### Requirement: First slice SHALL use in-memory seeded data

The endpoint handler SHALL use hardcoded in-memory seed data. The seed SHALL include at least 8 people with stable, unique `id` values. The seed SHALL exercise all four canonical statuses (`available`, `busy`, `partial`, `off`) across the returned date range so that backend tests can verify all status states appear.

#### Scenario: Seed data includes at least 8 people with stable IDs
- **WHEN** the endpoint is called with the default 4-day window
- **THEN** `people` contains at least 8 entries, each with a non-empty `id`, `name`, and `initials`

#### Scenario: Seed data exercises all four statuses
- **WHEN** the endpoint is called
- **THEN** across all people and dates in the response, the statuses `available`, `busy`, `partial`, and `off` each appear at least once

### Requirement: Existing behavior SHALL remain intact

The new endpoint SHALL be purely additive. `GET /api/hello` SHALL continue to return `{"message": "Hello from the backend!"}` with status 200. CORS SHALL continue to allow origins `http://localhost:5173` and `http://frontend:5173`. `GET /openapi.json` SHALL include both `/api/hello` and `/api/dashboard/people-availability`.

#### Scenario: Hello endpoint is unchanged
- **WHEN** `GET /api/hello` is called after the new endpoint is registered
- **THEN** the response is 200 with body `{"message": "Hello from the backend!"}` and `Content-Type: application/json`

#### Scenario: CORS still allows frontend origins
- **WHEN** a cross-origin `OPTIONS` preflight request is sent with `Origin: http://localhost:5173`
- **THEN** the response includes `Access-Control-Allow-Origin: http://localhost:5173`

#### Scenario: OpenAPI includes both endpoints
- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document contains paths `/api/hello` and `/api/dashboard/people-availability`

### Requirement: Backend tests SHALL cover the new endpoint

Backend tests SHALL include a test function that:

- Sends `GET /api/dashboard/people-availability` and asserts 200 OK with `Content-Type: application/json`.
- Unmarshals the response as JSON and verifies the top-level fields (`range`, `summary`, `people`, `statuses`) are present.
- Asserts that `summary.availableToday` matches the actual count of people with status `"available"` on `range.selectedDate`.
- Asserts that no status values outside the canonical set appear.

#### Scenario: Happy-path test passes
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** the new test function for the dashboard endpoint passes

#### Scenario: Hello test still passes
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** the existing `TestHelloEndpoint` function still passes
