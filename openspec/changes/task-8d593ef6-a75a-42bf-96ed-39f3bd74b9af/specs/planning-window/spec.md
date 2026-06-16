## ADDED Requirements

### Requirement: The backend SHALL expose a PUT endpoint to update the planning window

The backend SHALL register `PUT /api/planning-window` via Huma with `operationId: "put-planning-window"`. The request body SHALL contain `startDate` and `endDate` as ISO 8601 date strings (`YYYY-MM-DD`). The endpoint SHALL validate that both dates are present and valid, and that `endDate >= startDate`. On validation failure, the endpoint SHALL return `422 Unprocessable Entity` with a descriptive error body. On success, the endpoint SHALL persist the singleton planning-window row using an UPSERT and return the updated `PlanningWindowBody` (`startDate`, `endDate`, `days` with `days` recalculated server-side). The endpoint SHALL appear in the auto-generated `/openapi.json`.

#### Scenario: PUT endpoint persists a valid date range and returns updated body

- **WHEN** `PUT /api/planning-window` is called with `{"startDate": "2026-07-10", "endDate": "2026-07-20"}`
- **THEN** the response is `200 OK` with `Content-Type: application/json`
- **AND** the body contains `startDate: "2026-07-10"`, `endDate: "2026-07-20"`, and `days: 11`
- **AND** a subsequent `GET /api/planning-window` returns the same values

#### Scenario: PUT endpoint rejects endDate before startDate

- **WHEN** `PUT /api/planning-window` is called with `{"startDate": "2026-08-01", "endDate": "2026-07-01"}`
- **THEN** the response is `422 Unprocessable Entity` with a descriptive error

#### Scenario: PUT endpoint rejects malformed dates

- **WHEN** `PUT /api/planning-window` is called with `{"startDate": "not-a-date", "endDate": "2026-07-01"}`
- **THEN** the response is `422 Unprocessable Entity` with a descriptive error

#### Scenario: PUT endpoint returns 500 on store failure

- **WHEN** `PUT /api/planning-window` is called and the Store returns an error
- **THEN** the response is `500 Internal Server Error`

#### Scenario: PUT endpoint appears in OpenAPI specification

- **WHEN** the backend is running with the PUT endpoint registered
- **THEN** `GET /openapi.json` includes the `/api/planning-window` path with a `put` operation
- **AND** the operation includes the request body schema with `startDate` and `endDate` properties
- **AND** the operation includes the success response schema with `startDate`, `endDate`, and `days`

### Requirement: The Store interface SHALL include an update method for the planning window

The `Store` interface in `backend/store.go` SHALL include `UpdatePlanningWindow(ctx context.Context, startDate, endDate time.Time) (*PlanningWindowBody, error)`. All four implementations (`PgStore`, `mockStore`, `failingStore`, `partialFailingStore`) SHALL implement this method.

#### Scenario: PgStore UpdatePlanningWindow persists via UPSERT

- **WHEN** `PgStore.UpdatePlanningWindow(ctx, startDate, endDate)` is called with valid dates
- **THEN** the `planning_windows` singleton row is updated or inserted
- **AND** the returned `PlanningWindowBody` contains the new `startDate`, `endDate`, and recalculated `days`

#### Scenario: mockStore UpdatePlanningWindow mutates its planning window

- **WHEN** `mockStore.UpdatePlanningWindow(ctx, startDate, endDate)` is called
- **THEN** subsequent `mockStore.GetPlanningWindow(ctx)` returns the updated values with recalculated `days`

#### Scenario: failingStore UpdatePlanningWindow returns an error

- **WHEN** `failingStore.UpdatePlanningWindow(ctx, startDate, endDate)` is called
- **THEN** the call returns `nil, errTestFailure`

### Requirement: Backend tests SHALL cover the PUT planning-window endpoint

Backend tests in `backend/main_test.go` SHALL include test functions for the PUT endpoint: success (200 with updated body), validation failure (422), and store failure (500). Backend integration tests in `backend/main_integration_test.go` SHALL include a persistence round-trip test.

#### Scenario: PUT success test passes

- **WHEN** `go test ./...` runs in `backend/`
- **THEN** a test function exists that sends `PUT /api/planning-window` with valid data, asserts 200, and verifies the response body matches the submitted range with correct `days`

#### Scenario: PUT validation failure test passes

- **WHEN** `go test ./...` runs in `backend/`
- **THEN** a test function exists that sends `PUT /api/planning-window` with `endDate < startDate`, asserts 422

#### Scenario: PUT store failure test passes

- **WHEN** `go test ./...` runs in `backend/`
- **THEN** a test function exists that sends `PUT /api/planning-window` against a failing store, asserts 500

#### Scenario: Integration test verifies persistence round-trip

- **WHEN** integration tests run against a real database
- **THEN** a test sends PUT with new dates, then GET, and verifies the GET response reflects the updated range

### Requirement: The frontend SHALL provide a mutation composable for updating the planning window

The frontend SHALL define `useUpdatePlanningWindow()` in `frontend/src/shared/composables/useUpdatePlanningWindow.ts`. The composable SHALL use the generated Pinia Colada mutation for `putPlanningWindow` (from `@/client/@pinia/colada.gen`). On successful mutation, the composable SHALL invalidate the `getPlanningWindow` query key via `useQueryCache().invalidate` so that all `usePlanningWindow()` consumers automatically refetch. The composable SHALL return `{ mutate, isPending, error }`.

#### Scenario: Mutation composable invalidates query on success

- **WHEN** the mutation is called with valid `{startDate, endDate}` and the backend responds 200
- **THEN** the `getPlanningWindow` query key is invalidated
- **AND** subsequent reads from `usePlanningWindow()` return the updated planning window data

#### Scenario: Mutation composable exposes error state

- **WHEN** the mutation is called and the backend responds with an error
- **THEN** `error.value` is non-null
- **AND** `isPending.value` transitions from `true` to `false`

### Requirement: Frontend tests SHALL cover the mutation composable and SettingsView form

Frontend tests SHALL cover the `useUpdatePlanningWindow()` mutation composable (success and error states) and the SettingsView form behavior (prefill from composable, validation, save, error display, reset). The existing route render test for `/settings` SHALL be updated to assert the presence of the planning-window form elements.

#### Scenario: SettingsView form prefill test verifies date pickers show composable values

- **WHEN** the SettingsView renders with mocked planning-window data (`startDate: "2026-07-05", endDate: "2026-08-13"`)
- **THEN** the start date picker trigger displays a date matching `"2026-07-05"`
- **AND** the end date picker trigger displays a date matching `"2026-08-13"`

#### Scenario: SettingsView blocks save when endDate < startDate

- **WHEN** the user sets end date before start date
- **THEN** the save button is disabled

#### Scenario: SettingsView blocks save when dates are unchanged

- **WHEN** the user has not modified either date from the loaded values
- **THEN** the save button is disabled

#### Scenario: SettingsView shows error on mutation failure

- **WHEN** the save mutation fails
- **THEN** an error message is displayed in the SettingsView card

#### Scenario: Route render test asserts planning-window form

- **WHEN** the route render test runs for `/settings`
- **THEN** the rendered HTML no longer contains `"Feature coming soon"`
- **AND** the rendered HTML contains form elements for the planning window (date picker triggers and save button)

#### Scenario: Precommit checks pass

- **WHEN** `scripts/precommit-run` is executed
- **THEN** all lint, type-check, and test checks pass without errors