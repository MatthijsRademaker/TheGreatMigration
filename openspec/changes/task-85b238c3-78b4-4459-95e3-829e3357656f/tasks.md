## 1. Backend — Dashboard endpoint default start date

- [ ] 1.1 Modify `backend/api/dashboard.go`: In the `registerDashboardPeopleAvailability` handler, when `input.Start` is empty, call `store.GetPlanningWindow(ctx)` to get the planning window start date and use it as the default. If no planning window exists (store returns nil or zero-value), return `400 Bad Request` with a descriptive error. If the store call fails, return `500 Internal Server Error`. Return errors via `huma.Error400BadRequest` / `huma.Error500InternalServerError`.
- [ ] 1.2 Update `backend/api/dashboard.go`: Add `isZeroDate` helper or use a zero `time.Time` check to detect missing planning window (comparing against `time.Time{}`).
- [ ] 1.3 Verify no changes needed in `backend/api/register.go` — `registerDashboardPeopleAvailability` already receives `store` as a parameter.
- [ ] 1.4 Verify no changes needed in `backend/store.go` — `GetPlanningWindow` already exists on the `Store` interface.
- [ ] 1.5 Regenerate OpenAPI snapshot and frontend client with `go run ./cmd/openapi-gen > ../frontend/openapi-snapshot.json && cd ../frontend && npm run regen:api` to capture any schema changes (there should be none since only handler logic changed, not the API contract).

## 2. Frontend — Composable `usePeopleAvailability` changes

- [ ] 2.1 Update `frontend/src/shared/composables/usePeopleAvailability.ts`: Accept an optional `start` parameter in the composable signature. When `start` is not provided, call `usePlanningWindow()` to get the planning window's `startDate` and pass it as the `start` query parameter to the dashboard query. Use an `enabled` guard (or `dependsOn`) to defer the dashboard query until the planning window query has resolved.
- [ ] 2.2 Update `adaptToComponentProps`: Add a `daysISO` output — build a parallel `string[]` of ISO date strings (`YYYY-MM-DD`) in the same UTC-based iteration loop that builds the day labels. Return `daysISO` alongside the existing adapted props.
- [ ] 2.3 Update the composable's return value: Expose `daysISO` as a `ComputedRef<string[]>` alongside `data`, `rawData`, `isLoading`, `isError`, `isEmpty`.
- [ ] 2.4 Verify the composable still works when called without arguments in `frontend/src/home/HomeView.vue` and `frontend/src/people/PeopleView.vue` — both should continue to work as the new `start` parameter defaults to auto-derivation from the planning window.

## 3. Frontend — PeopleView date derivation and error handling

- [ ] 3.1 Update `frontend/src/people/PeopleView.vue` replacement: Destructure `daysISO` from `usePeopleAvailability()` return value.
- [ ] 3.2 Remove the existing `getISODate` function from `PeopleView.vue`.
- [ ] 3.3 Update `handleCellUpdate`: Replace `getISODate(dayIndex)` with `daysISO.value[dayIndex]`. Add a defensive guard: if `dayIndex >= daysISO.value.length`, set `updateError` to a clear message indicating the cell cannot be mapped to a date, and return without making an API call.
- [ ] 3.4 Update error handling in the catch block of `handleCellUpdate`: Inspect the error body's `detail` field via `(err as any)?.cause?.body?.detail` and produce distinct messages for each 400 subtype (date outside planning window, invalid date format, invalid status value). For unrecognized 400s, display the detail text as a fallback. Keep 404 handling unchanged.
- [ ] 3.5 Clean up any unused imports or variables after removing `getISODate` and related code (e.g., `rawData` if it's no longer referenced elsewhere).

## 4. Backend tests

- [ ] 4.1 Open `backend/main_test.go` and update the existing dashboard tests: the "Default query parameters produce a 4-day window" test should now expect `range.startDate` to equal the planning window's `startDate` (`"2026-07-05"`) instead of `time.Now()`. Ensure the mock store returns a properly seeded planning window.
- [ ] 4.2 Add a test scenario in `backend/main_test.go`: dashboard endpoint returns 400 when no planning window exists and no explicit `start` is provided. Use the `failingStore` or a store variant that returns nil for `GetPlanningWindow`.
- [ ] 4.3 Open `backend/main_integration_test.go` and update any integration tests that depend on dashboard default behavior against the real database. The seeded planning window (`2026-07-05`) should now be the default, so assertions on `range.startDate` need updating.
- [ ] 4.4 Run `go test ./...` in `backend/` and verify all tests pass.

## 5. Frontend tests

- [ ] 5.1 Open `frontend/src/shared/composables/__tests__/usePeopleAvailability.test.ts` (or create if it doesn't exist): Add tests for the new planning-window integration:
  - Composable passes planning window `startDate` as query parameter when no explicit `start` given
  - Dashboard query is deferred until planning window resolves
  - Explicit `start` parameter bypasses planning window
  - `daysISO` array matches `days` label array length and order
  - `daysISO` values are valid ISO 8601 date strings
- [ ] 5.2 Open `frontend/src/people/__tests__/PeopleView.spec.ts`: Add tests for:
  - `handleCellUpdate` derives correct date from `daysISO` by index
  - Out-of-range `dayIndex` produces clear error without API call
  - Error handling parses `detail` field for distinct 400 subtypes
  - Error handling shows fallback for unrecognized 400 with detail text
- [ ] 5.3 Update route render test (`frontend/tests/app-routes-render.test.ts`) if needed to account for the planning-window dependency in the people-availability query.
- [ ] 5.4 Run frontend tests and verify all pass.

## 6. Verification

- [ ] 6.1 Run `scripts/precommit-run` and verify all lint, type-check, and test checks pass without errors.
- [ ] 6.2 Verify the fix manually by starting the backend and frontend, navigating to `/people`, clicking a status pill, and confirming the upsert succeeds (no "Invalid status or date" error).
