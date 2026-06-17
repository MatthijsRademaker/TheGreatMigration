## 1. Composable — Extend useDailySchedule with pagination

- [ ] 1.1 Add `UseDailyScheduleOptions` interface with optional `page` (default 1), `daysPerPage` (default 4), and `start` fields
- [ ] 1.2 Import and call `usePlanningWindow()` to resolve the canonical planning window
- [ ] 1.3 Compute `totalDays` from `planningWindow.planWindowDays.value.length`, falling back to `daysPerPage` when unavailable
- [ ] 1.4 Compute `totalPages` as `Math.max(1, Math.ceil(totalDays / daysPerPage))`
- [ ] 1.5 Compute `startParam` from planning window start date + `(page-1) * daysPerPage` when `page > 1`, exactly mirroring `usePeopleAvailability` lines 85-100
- [ ] 1.6 Pass computed `startParam` and `daysPerPage` as `query.start` and `query.days` to `getDashboardDailyScheduleQuery()`
- [ ] 1.7 Defer query until planning window resolves (when no explicit `start` is provided), mirroring `queryEnabled` pattern
- [ ] 1.8 Add planning-window-change watcher with `oldDates.length > 0` guard that resets `page` to 1 when the window changes and `page > 1`
- [ ] 1.9 Add `goToPrevPage()` and `goToNextPage()` helpers with boundary guards (`page > 1` and `page < totalPages`)
- [ ] 1.10 Expose `page`, `totalPages`, `daysPerPage`, `totalDays`, `goToPrevPage`, `goToNextPage` from the composable return value

## 2. Route — Wire pagination in HomeView

- [ ] 2.1 Destructure `page`, `totalPages`, `daysPerPage`, `totalDays`, `goToPrevPage`, `goToNextPage` from `useDailySchedule()`
- [ ] 2.2 Render pagination navigation bar between the state-driven cards and the `<DailySchedule>` component, showing:
  - Date range label: first and last day from `scheduleData.days` (e.g., "5 Jul (Sun) – 8 Jul (Wed)")
  - Page indicator: "Page X of Y"
  - Previous button: disabled when `page <= 1`
  - Next button: disabled when `page >= totalPages`
- [ ] 2.3 Ensure pagination bar is only visible in the success state (after loading/error/empty state cards)
- [ ] 2.4 Preserve the existing loading, error, and empty state cards unchanged

## 3. Route — Wire pagination in CalendarView

- [ ] 3.1 Destructure `page`, `totalPages`, `goToPrevPage`, `goToNextPage` from `useDailySchedule()`
- [ ] 3.2 Render pagination navigation bar above the `<DailySchedule>` component in the success template, showing:
  - Date range label from `scheduleData.days`
  - Page indicator: "Page X of Y"
  - Previous button: disabled when `page <= 1`
  - Next button: disabled when `page >= totalPages`
- [ ] 3.3 Ensure pagination bar is only visible in the success state
- [ ] 3.4 Preserve existing loading, error, and empty state cards unchanged

## 4. Tests — Composable page derivation and navigation

- [ ] 4.1 Extend `frontend/tests/calendar/useDailySchedule.test.ts` with client-side mount tests (mirroring `usePeopleAvailability.integration.test.ts`)
  - Page initializes at 1
  - `goToNextPage` advances page by 1
  - `goToPrevPage` retreats page by 1
  - `goToPrevPage` does not go below 1
  - `goToNextPage` does not go beyond `totalPages`
  - `totalPages` reflects planning window days divided by `daysPerPage`
  - `totalPages` is at least 1
  - Page resets to 1 when planning window changes and `page > 1`
  - Page does NOT reset on initial planning window load
- [ ] 4.2 Verify composable returns all expected pagination properties (`page`, `totalPages`, `daysPerPage`, `totalDays`, `goToPrevPage`, `goToNextPage`)
- [ ] 4.3 Verify existing composable tests still pass (adaptation, empty state, loading state)

## 5. Tests — Route-level pagination controls

- [ ] 5.1 Extend `frontend/tests/app-routes-render.test.ts` for `/` route:
  - Mock planning-window response returning 40 days
  - Assert pagination date range label, page indicator, Previous/Next buttons render
  - Assert Previous button has `disabled` attribute on page 1
  - Assert DailySchedule remains read-only (no Edit/Delete/Add task controls)
- [ ] 5.2 Extend `frontend/tests/app-routes-render.test.ts` for `/calendar` route:
  - Mock planning-window response returning 40 days
  - Assert pagination date range label, page indicator, Previous/Next buttons render
  - Assert calendar-specific controls (Add task, Edit, Delete) still render
- [ ] 5.3 Add client-side mount tests for HomeView and CalendarView pagination (mirroring `PeopleView.test.ts` pattern):
  - Previous button disabled on page 1
  - Next button disabled on last page
  - Clicking Next advances page and re-fetches data

## 6. Verification

- [ ] 6.1 Run `scripts/precommit-run` and verify all lint, type-check, and test checks pass without errors
- [ ] 6.2 Verify `DailySchedule.test.ts` passes unchanged (component contract preserved)
- [ ] 6.3 Manually verify: navigate to `/` and `/calendar`, confirm schedule shows 4 days, click Next to see next 4-day slice, confirm schedule CRUD works for dates on any page
- [ ] 6.4 Manually verify: change planning window in Settings → schedule page resets to page 1