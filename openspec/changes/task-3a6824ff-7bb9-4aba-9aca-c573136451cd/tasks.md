## 1. Backend: Planning-Window Endpoint

- [ ] 1.1 Create `backend/planning_window.go` with request/response types (`PlanningWindowInput`, `PlanningWindowOutput`, `PlanningWindowBody`) and seed data (Go constants for `2026-07-05` through `2026-08-13`).
- [ ] 1.2 Register `GET /api/planning-window` in `backend/main.go` via `huma.Register` with `operationId: "get-planning-window"`, tags `["Planning"]`, path `/api/planning-window`, and description documenting the global move range contract.
- [ ] 1.3 Implement the handler to return `{startDate: string, endDate: string, days: number}` where `startDate` and `endDate` are ISO 8601 strings and `days` is the inclusive day count.
- [ ] 1.4 Verify `GET /openapi.json` includes the new endpoint with correct response schema.

## 2. Backend: Contract Tests

- [ ] 2.1 Extend `newTestAPI()` in `backend/main_test.go` to register the planning-window endpoint alongside hello and dashboard.
- [ ] 2.2 Add `TestPlanningWindowEndpoint` asserting 200 OK, `Content-Type: application/json`, presence of `startDate`/`endDate`/`days` fields, correct default values (`"2026-07-05"`, `"2026-08-13"`, `40`), and that `startDate < endDate`.
- [ ] 2.3 Verify existing `TestHelloEndpoint` and `TestDashboardPeopleAvailability` still pass.

## 3. Frontend: OpenAPI Snapshot and Client Regeneration

- [ ] 3.1 Start the backend and run `refresh-openapi-snapshot.mjs` to capture the updated `/openapi.json` into `frontend/openapi-snapshot.json`.
- [ ] 3.2 Run `generate-api.mjs` to regenerate `frontend/src/client/` (SDK functions, types, Pinia Colada query artifacts) including the new `get-planning-window` operation.
- [ ] 3.3 Verify the generated `@pinia/colada.gen.ts` exports `getPlanningWindowQuery`, `getPlanningWindowQueryKey`, and related types.

## 4. Frontend: Refactor planWindow.ts

- [ ] 4.1 Remove `PLAN_WINDOW_START`, `PLAN_WINDOW_END`, `generatePlanWindowDays()`, `planWindowDays`, and `planWindowDayCount` from `frontend/src/shared/lib/planWindow.ts`.
- [ ] 4.2 Retain `formatPlanDayLabel(date: Date): string` and the `PlanWindowDay` interface as pure utility exports.
- [ ] 4.3 Update the module JSDoc comment to document the new composable-based consumer pattern and remove references to eager evaluation and rebuild constraints.

## 5. Frontend: Create usePlanningWindow Composable

- [ ] 5.1 Create `frontend/src/shared/composables/usePlanningWindow.ts`.
- [ ] 5.2 Import `getPlanningWindowQuery` from `@/client/@pinia/colada.gen` and use `useQuery(getPlanningWindowQuery())` to fetch backend data.
- [ ] 5.3 Derive `planWindowDays: ComputedRef<PlanWindowDay[]>` from the fetched `startDate`/`endDate` using the same inclusive day-generation algorithm (moved from the removed `generatePlanWindowDays`).
- [ ] 5.4 Derive `planWindowDayCount: ComputedRef<number>` as `planWindowDays.value.length`.
- [ ] 5.5 Expose `isLoading: ComputedRef<boolean>` and `error: ComputedRef<boolean>` from the Pinia Colada query state.
- [ ] 5.6 Expose the shared query key for cache coherency across consumers.

## 6. Frontend: Update HomeView

- [ ] 6.1 Replace `import { planWindowDayCount } from '@/shared/lib/planWindow'` with `import { usePlanningWindow } from '@/shared/composables/usePlanningWindow'`.
- [ ] 6.2 Replace the static `String(planWindowDayCount)` in the summaries array with a computed ref that reads from the composable's reactive `planWindowDayCount`.
- [ ] 6.3 Add loading and error states for the Move days card, following the exact pattern of the hello card (`v-if` branches for loading, error, and success).

## 7. Frontend: Update CalendarView

- [ ] 7.1 Replace `import { planWindowDays } from '@/shared/lib/planWindow'` with `import { usePlanningWindow } from '@/shared/composables/usePlanningWindow'`.
- [ ] 7.2 Call `usePlanningWindow()` and bind the reactive `planWindowDays` to the template `v-for`.
- [ ] 7.3 Add a loading skeleton: when `isLoading` is true, render placeholder day-column elements with the same responsive grid classes (`grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7`, `min-h-36`, border, bg-muted/40, rounded-lg, animate-pulse).
- [ ] 7.4 Add an error state: when `error` is true, render an error message (e.g., "Planning window unavailable — check backend") inside the CardContent.

## 8. Frontend: Update Tests

- [ ] 8.1 Refactor `frontend/tests/planWindow.test.ts` to test the composable by mocking the fetch response: test success (correct day count, first/last dates), loading state, and error state.
- [ ] 8.2 Retain existing `formatPlanDayLabel` tests (these are pure utility tests and should still pass).
- [ ] 8.3 Update `frontend/tests/app-routes-render.test.ts` to install Pinia and Pinia Colada before SSR render (they are already installed per the hello-world-integration change) and mock the planning-window endpoint response.
- [ ] 8.4 Update the `/calendar` route assertion: the day column count must match the mocked response's `days` value instead of the removed `planWindowDayCount` constant.

## 9. Verification

- [ ] 9.1 Run `scripts/precommit-run` and confirm all lint, type-check, and test checks pass.
- [ ] 9.2 Manually verify in the browser that HomeView displays the Move days count from the backend and CalendarView renders the correct day columns.
- [ ] 9.3 Manually verify that stopping the backend shows error states in both views.
- [ ] 9.4 Verify `GET /api/hello` and `GET /api/dashboard/people-availability` continue to return correct responses.