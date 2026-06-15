# Tasks: Date range planning configurable

## 1. Create the planning window module

- [x] 1.1 Create `frontend/src/shared/lib/planWindow.ts` exporting `PLAN_WINDOW_START` and `PLAN_WINDOW_END` as ISO date string constants (`"2026-07-05"` and `"2026-08-13"`).
- [x] 1.2 Implement `formatPlanDayLabel(date: Date): string` using `Intl.DateTimeFormat` with `en-US` locale, producing labels like `"Sun 5 Jul"`.
- [x] 1.3 Implement `planWindowDays` as a derived getter that generates the inclusive day array `{date: Date, label: string, dateString: string}[]` between `PLAN_WINDOW_START` and `PLAN_WINDOW_END`.
- [x] 1.4 Implement `planWindowDayCount` as a derived getter returning `planWindowDays.length`.
- [x] 1.5 Add a module header comment documenting the inclusive contract, the default date range, and the intended consumer pattern.

## 2. Refactor CalendarView

- [x] 2.1 Import `planWindowDays` from `@/shared/lib/planWindow` in `CalendarView.vue`.
- [x] 2.2 Remove the hard-coded `moveDays` array (the Mon–Fri entries with `day`, `focus`, `helpers`).
- [x] 2.3 Replace the template `v-for="day in moveDays"` with `v-for="planDay in planWindowDays"` using `planDay.dateString` as the key.
- [x] 2.4 Update the column header from `{{ day.day }}` to `{{ planDay.label }}`.
- [x] 2.5 Replace the per-column `Badge` (showing `day.helpers`) and focus text (`day.focus`) with a generic placeholder description like "Planned move day".
- [x] 2.6 Update the grid class from `md:grid-cols-5` to `grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7` to accommodate windows of varying sizes.

## 3. Refactor HomeView

- [x] 3.1 Import `planWindowDayCount` from `@/shared/lib/planWindow` in `HomeView.vue`.
- [x] 3.2 Change the "Move days" entry in the `summaries` array to use `String(planWindowDayCount)` as its `value` instead of the hard-coded `"5"`.
- [x] 3.3 Verify the card label ("Move days"), description ("Scheduled working days in the plan"), and icon (`CalendarDaysIcon`) remain unchanged.

## 4. Document PeopleView deferral

- [x] 4.1 Add a comment in `PeopleView.vue` above the `helpers` array noting that the weekday availability strings are placeholder data and that future integration with `planWindowDays` is expected when real availability CRUD is built.
- [x] 4.2 Do not modify the PeopleView template or data structure in this change.

## 5. Amend hello-world-integration spec

- [x] 5.1 Create `openspec/changes/task-b9e23168-60c8-4b6a-9a30-6b2dabbda869/specs/hello-world-integration/spec.md` as a spec delta.
- [x] 5.2 Add a `## MODIFIED` requirement that narrows the "Other HomeView sections are preserved" scenario: the "Move days" card value is now derived from the planning window, while "Available today" and "Under-staffed" remain static.
- [x] 5.3 Verify the delta does not alter any other scenario in the canonical spec.

## 6. Create planning-window capability spec

- [x] 6.1 Create `openspec/changes/task-b9e23168-60c8-4b6a-9a30-6b2dabbda869/specs/planning-window/spec.md` with requirements for the shared module, CalendarView consumption, HomeView consumption, PeopleView deferral, and test coverage.

## 7. Add unit tests for planWindow

- [x] 7.1 Create `frontend/tests/planWindow.test.ts` importing the module under test.
- [x] 7.2 Test: `PLAN_WINDOW_START` equals `"2026-07-05"` and `PLAN_WINDOW_END` equals `"2026-08-13"`.
- [x] 7.3 Test: `planWindowDayCount` equals 40 for the default range (inclusive: 27 July days + 13 August days).
- [x] 7.4 Test: `planWindowDays.length` equals `planWindowDayCount`.
- [x] 7.5 Test: `planWindowDays[0].dateString` equals `PLAN_WINDOW_START` and `planWindowDays[last].dateString` equals `PLAN_WINDOW_END`.
- [x] 7.6 Test: `formatPlanDayLabel(new Date("2026-07-05"))` produces a string containing "5" and "Jul" (format example: "Sun 5 Jul").
- [x] 7.7 Test: `formatPlanDayLabel(new Date("2026-08-13"))` produces a string containing "13" and "Aug".

## 8. Extend route render test

- [x] 8.1 In `frontend/tests/app-routes-render.test.ts`, import `planWindowDayCount` from the planWindow module.
- [x] 8.2 Extend the `/calendar` route test case to assert the rendered HTML contains the expected number of day-column elements (e.g., by checking for a specific CSS class repeated `planWindowDayCount` times, or by asserting the grid contains the expected child count).

## 9. Verify

- [x] 9.1 Run `scripts/precommit-run` from the repository root and confirm all checks pass.
- [x] 9.2 Confirm `frontend/` unit and render tests pass (`cd frontend && npx vitest run`).
- [x] 9.3 Manually verify that `HomeView` shows `40` for Move days (or the count matching the configured window) and that `CalendarView` renders 40 day columns with date labels.
