## 1. Backend Write Path

- [ ] 1.1 Add `UpsertPlanningWindow` sqlc query in `backend/queries/planning_window.sql` using `INSERT INTO planning_windows (id, start_date, end_date) VALUES (1, $1, $2) ON CONFLICT (id) DO UPDATE SET start_date = $1, end_date = $2, updated_at = NOW() RETURNING id, start_date, end_date, created_at, updated_at`
- [ ] 1.2 Regenerate `backend/db/planning_window.sql.go` by running sqlc
- [ ] 1.3 Add `UpdatePlanningWindow(ctx context.Context, startDate, endDate time.Time) (*PlanningWindowBody, error)` to the `Store` interface in `backend/store.go`
- [ ] 1.4 Implement `UpdatePlanningWindow` on `PgStore` in `backend/store.go`: parse dates, call `UpsertPlanningWindow`, return `PlanningWindowBody` with recalculated `days`
- [ ] 1.5 Add request/response types (`UpdatePlanningWindowInput`, `UpdatePlanningWindowOutput`) and register `PUT /api/planning-window` handler in `backend/planning_window.go` with `operationId: "put-planning-window"`, validation (both dates required, `endDate >= startDate`, return 422 on failure), and `huma.Error500InternalServerError` on store failure
- [ ] 1.6 Implement `UpdatePlanningWindow` on `mockStore` in `backend/store_mock_test.go`: accept params, update `m.planningWindow` fields, recalculate `days`
- [ ] 1.7 Implement `UpdatePlanningWindow` on `failingStore` in `backend/main_test.go`: return `errTestFailure`
- [ ] 1.8 Implement `UpdatePlanningWindow` on `partialFailingStore` in `backend/main_test.go`: succeed to support daily-schedule tests that need a valid planning window
- [ ] 1.9 Add backend unit tests in `backend/main_test.go`: `TestUpdatePlanningWindowSuccess` (200, verify response body matches updated range + recalculated days), `TestUpdatePlanningWindowValidationFailure` (422, endDate < startDate, malformed dates), `TestUpdatePlanningWindowStoreFailure` (500 via failingStore)
- [ ] 1.10 Add backend integration test in `backend/main_integration_test.go`: PUT new range, then GET to verify persistence round-trip

## 2. OpenAPI and Frontend Client Refresh

- [ ] 2.1 Start backend with new PUT endpoint exposed
- [ ] 2.2 Regenerate OpenAPI snapshot: `npm run refresh:openapi-snapshot` (from `frontend/`)
- [ ] 2.3 Regenerate frontend client: `npm run generate:api` (from `frontend/`)
- [ ] 2.4 Verify `frontend/src/client/sdk.gen.ts` contains an `updatePlanningWindow` function and `frontend/src/client/@pinia/colada.gen.ts` contains mutation exports

## 3. Frontend Settings Form

- [ ] 3.1 Create `frontend/src/shared/composables/useUpdatePlanningWindow.ts`: use the generated Pinia Colada mutation for `putPlanningWindow`; on success, invalidate `getPlanningWindowQueryKey()` via `useQueryCache`; return `{ mutate, isPending, error }`
- [ ] 3.2 Replace `frontend/src/settings/SettingsView.vue` placeholder with interactive "Planning window" card:
  - Import `usePlanningWindow`, `useUpdatePlanningWindow`, `DatePicker`, `Button`, `Card`
  - Display current range summary from `usePlanningWindow()`
  - Two `DatePicker` controls for start/end dates with `v-model` bound to local refs
  - Convert between ISO strings (from composable) and `CalendarDate` (for DatePicker) using `@internationalized/date`
  - Form validation: block save when `endDate < startDate` or dates are unchanged from loaded values
  - Save button with loading/disabled states from mutation composable
  - Reset button to revert to loaded values
  - Error display area for mutation errors
  - Card description referencing planning-window management
- [ ] 3.3 Update `frontend/tests/app-routes-render.test.ts`: replace `/settings` content assertion from `"Feature coming soon"` to assertions for planning-window form elements (date picker triggers, save button, card title matching "Planning window")
- [ ] 3.4 Create `frontend/tests/settings-view.test.ts` covering: form prefill from mocked composable data, validation (disallowed endDate < startDate, detection of unchanged values), successful save triggers mutation and query invalidation, error display on mutation failure, reset restores original values

## 4. Spec Amendments

- [ ] 4.1 Amend `openspec/specs/planning-window/spec.md`: add ADDED requirements for PUT /api/planning-window (operationID, request body schema, validation rules, response contract, error responses, Store method)
- [ ] 4.2 Amend `openspec/specs/sidebar-navigation/spec.md`: MODIFIED — remove placeholder-only restriction on SettingsView; reference planning-window spec for content contract; update scenario to allow interactive planning-window card

## 5. Verification

- [ ] 5.1 Run `cd src && go test ./...` — all backend tests pass
- [ ] 5.2 Run `npm run test` from `frontend/` — all frontend tests pass
- [ ] 5.3 Run `scripts/precommit-run` — lint, type-check, and test checks pass
- [ ] 5.4 Manually verify: save a new date range in settings, refresh page, confirm dates persist; confirm GET /api/planning-window returns updated values