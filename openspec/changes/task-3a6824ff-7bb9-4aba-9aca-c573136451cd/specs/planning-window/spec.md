## ADDED Requirements

### Requirement: The planning window SHALL be defined by a backend endpoint as the canonical source of truth

The planning window SHALL be served by `GET /api/planning-window`, a backend endpoint registered via Huma v2 with `operationId: "get-planning-window"`. The endpoint SHALL accept no query parameters. The response body SHALL contain `startDate` (ISO 8601 date string), `endDate` (ISO 8601 date string), and `days` (inclusive count integer). The endpoint SHALL be seeded with compile-time Go constants defaulting to `startDate="2026-07-05"`, `endDate="2026-08-13"`, and `days=40`. The endpoint SHALL appear in the auto-generated `/openapi.json`.

#### Scenario: Endpoint returns the canonical planning window
- **WHEN** `GET /api/planning-window` is called
- **THEN** the response is 200 OK with `Content-Type: application/json`
- **AND** the body contains `startDate`, `endDate`, and `days` fields
- **AND** `startDate` is `"2026-07-05"` and `endDate` is `"2026-08-13"`
- **AND** `days` equals 40

#### Scenario: Endpoint appears in OpenAPI specification
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` includes the `/api/planning-window` endpoint with its response schema

#### Scenario: Backend validates its own response contract
- **WHEN** the endpoint is called
- **THEN** `startDate` and `endDate` are valid ISO 8601 date strings (format `YYYY-MM-DD`)
- **AND** `startDate` lexicographically precedes `endDate`
- **AND** `days` is a positive integer equal to the inclusive day count between `startDate` and `endDate`

### Requirement: The frontend SHALL consume the planning window through a shared Pinia Colada composable

The frontend SHALL define `usePlanningWindow()` in `frontend/src/shared/composables/usePlanningWindow.ts`. The composable SHALL use the generated Pinia Colada query for `GET /api/planning-window` (from `@/client/@pinia/colada.gen`). The composable SHALL expose reactive computed refs: `planWindowDays` (`PlanWindowDay[]`), `planWindowDayCount` (`number`), `isLoading` (`boolean`), and `isError` (`boolean`). The composable SHALL derive `planWindowDays` from the fetched `startDate`/`endDate` using inclusive day generation, and SHALL derive `planWindowDayCount` as `planWindowDays.length`. The composable SHALL cache the backend result using a shared query key so that multiple concurrent consumers do not trigger duplicate fetches.

#### Scenario: Composable returns derived day list on success
- **WHEN** the backend returns a valid planning-window response
- **THEN** `planWindowDays.value` is an array of `PlanWindowDay` objects
- **AND** `planWindowDays.value.length` equals the backend's `days` value
- **AND** the first entry has `dateString` equal to the backend's `startDate`
- **AND** the last entry has `dateString` equal to the backend's `endDate`

#### Scenario: Composable reflects loading state
- **WHEN** the fetch is in-flight
- **THEN** `isLoading.value` is `true` and `isError.value` is `false`

#### Scenario: Composable reflects error state
- **WHEN** the fetch fails (network error, non-200 response, or malformed response)
- **THEN** `isError.value` is `true` and `isLoading.value` is `false`

### Requirement: HomeView SHALL display the Move days count from the composable with loading and error handling

`frontend/src/home/HomeView.vue` SHALL import and call `usePlanningWindow()`. The "Move days" summary card SHALL display the reactive `planWindowDayCount` value. The card SHALL show a loading indicator while the composable is in the loading state, an error message when the composable is in the error state, and the day count value on success. The card's label ("Move days"), description, and icon SHALL remain unchanged.

#### Scenario: Move days card displays backend value on success
- **WHEN** the composable resolves successfully with the default 40-day planning window
- **THEN** the "Move days" summary card displays `"40"` as its value

#### Scenario: Move days card shows loading state
- **WHEN** the composable is in the loading state
- **THEN** the "Move days" summary card displays a loading indicator

#### Scenario: Move days card shows error state
- **WHEN** the composable is in the error state
- **THEN** the "Move days" summary card displays an error message

### Requirement: CalendarView SHALL render day columns from the composable with loading and error handling

`frontend/src/calendar/CalendarView.vue` SHALL import and call `usePlanningWindow()`. Day columns SHALL be rendered by iterating over the reactive `planWindowDays`. The view SHALL show a loading skeleton while the composable is in the loading state, matching the responsive grid layout. The view SHALL show an error message when the composable is in the error state. On success, the view SHALL render exactly one column per planning-window day with the formatted date label from `planDay.label`.

#### Scenario: CalendarView renders backend-driven day columns on success
- **WHEN** the composable resolves successfully with the default 40-day planning window
- **THEN** the view renders exactly 40 day-column elements, each labeled with a date in the planning window range

#### Scenario: CalendarView grid adapts to arbitrary window sizes
- **WHEN** the backend returns a different number of days
- **THEN** the number of rendered columns changes accordingly without breaking the responsive grid layout

#### Scenario: CalendarView shows loading skeleton
- **WHEN** the composable is in the loading state
- **THEN** the view renders placeholder day-column elements with the same grid classes

#### Scenario: CalendarView shows error state
- **WHEN** the composable is in the error state
- **THEN** the view displays an error message indicating the planning window is unavailable

### Requirement: Existing backend behavior SHALL remain intact

The new endpoint SHALL be purely additive. `GET /api/hello` SHALL continue to return `{"message": "Hello from the backend!"}` with status 200. `GET /api/dashboard/people-availability` SHALL continue to return its combined payload with correct range, summary, people, and statuses. CORS SHALL continue to allow origins `http://localhost:5173` and `http://frontend:5173`. `GET /openapi.json` SHALL include `/api/hello`, `/api/dashboard/people-availability`, and `/api/planning-window`.

#### Scenario: Hello endpoint is unchanged
- **WHEN** `GET /api/hello` is called after the new endpoint is registered
- **THEN** the response is 200 with body `{"message": "Hello from the backend!"}` and `Content-Type: application/json`

#### Scenario: Dashboard endpoint is unchanged
- **WHEN** `GET /api/dashboard/people-availability` is called after the new endpoint is registered
- **THEN** the response is 200 with the expected `range`, `summary`, `people`, and `statuses` top-level fields

#### Scenario: OpenAPI includes all three endpoints
- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document contains paths `/api/hello`, `/api/dashboard/people-availability`, and `/api/planning-window`

### Requirement: Backend tests SHALL cover the planning-window endpoint contract

Backend tests SHALL include a test function that sends `GET /api/planning-window`, asserts 200 OK with `Content-Type: application/json`, unmarshals the JSON body, and verifies `startDate` equals `"2026-07-05"`, `endDate` equals `"2026-08-13"`, and `days` equals 40. Existing hello and dashboard tests SHALL continue to pass.

#### Scenario: Planning-window contract test passes
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** the new planning-window test function passes

#### Scenario: Existing tests still pass
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** `TestHelloEndpoint` and `TestDashboardPeopleAvailability` still pass

### Requirement: Frontend tests SHALL cover composable behavior

Frontend tests SHALL test the composable's success, loading, and error states by mocking the fetch response. The tests SHALL verify that on success, `planWindowDayCount` matches the mocked `days` value and `planWindowDays` boundaries match the mocked `startDate`/`endDate`. The route render test SHALL mock the planning-window endpoint response and assert the correct number of day columns in CalendarView. The `formatPlanDayLabel` utility tests SHALL be preserved.

#### Scenario: Composable test validates success state
- **WHEN** the composable is tested with a mocked success response
- **THEN** `planWindowDayCount` matches the mocked `days` and the day list boundaries match the mocked date range

#### Scenario: Composable test validates loading and error states
- **WHEN** the composable is tested with a pending fetch and with a failed fetch
- **THEN** `isLoading` and `isError` transition correctly

#### Scenario: Route render test reflects backend-driven day count
- **WHEN** the route render test runs for `/calendar` with a mocked planning-window response
- **THEN** the rendered output contains exactly the mocked `days` number of day-column elements

#### Scenario: Precommit checks pass
- **WHEN** `scripts/precommit-run` is executed
- **THEN** all lint, type-check, and test checks pass without errors

## REMOVED Requirements

### Requirement: The planning window SHALL be defined as a single shared TypeScript module

**Reason**: The planning window is now backed by `GET /api/planning-window`. The frontend TypeScript module (`planWindow.ts`) is refactored to a thin derivation consumer of the backend contract, with the composable (`usePlanningWindow`) as the primary consumer API. The hard-coded constants `PLAN_WINDOW_START`, `PLAN_WINDOW_END`, `planWindowDays`, and `planWindowDayCount` are removed.

**Migration**: Consumers import `usePlanningWindow` from `@/shared/composables/usePlanningWindow` instead of importing constants from `@/shared/lib/planWindow`. The `formatPlanDayLabel` utility and `PlanWindowDay` interface remain available from `@/shared/lib/planWindow`.
