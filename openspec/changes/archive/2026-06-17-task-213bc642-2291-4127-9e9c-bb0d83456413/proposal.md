## Why

The Daily Schedule panel is locked to a fixed 4-day view inherited from the backend's default `days=4` parameter. Users configure a planning window that can span many weeks (e.g., 40 days), but can only see the first 4 days of schedule data in the dashboard and calendar views. People Availability already solved this by deriving page-based `start` and `days` values from the configured planning window and adding prev/next navigation. The Daily Schedule now lags behind that product behavior, making the configured date range feel only partially usable.

## What Changes

- **`useDailySchedule` composable gains pagination**: Accepts optional `page` (1-indexed, default 1) and `daysPerPage` (default 4) options. Derives explicit `start` and `days` query parameters from the planning window and current page state, replacing the zero-option call that relied on backend defaults. Exposes `page`, `totalPages`, `daysPerPage`, `totalDays`, `goToPrevPage`, and `goToNextPage`.
- **Route-level pagination controls**: Both `HomeView` (read-only dashboard) and `CalendarView` (editable schedule board) render prev/next navigation controls outside the shared `DailySchedule.vue` component, exactly mirroring the `PeopleView → PeopleAvailability` pattern where the route owns pagination state and the presentational component stays prop-driven.
- **`DailySchedule.vue` remains presentational**: The shared component does NOT acquire pagination props, emits, data-fetching, or navigation rendering. Its `days` and `readOnly` props continue unchanged.
- **In-memory page state only**: Page position is composable-local; no URL sync, no localStorage persistence.
- **No backend changes**: The existing `GET /api/dashboard/daily-schedule` endpoint already supports optional `start` and `days` query parameters with the correct contract. The frontend now derives and passes these explicitly instead of relying on backend defaults.

## Impact

| Area | Impact |
|------|--------|
| `frontend/src/calendar/composables/useDailySchedule.ts` | Extended with page/daysPerPage options, planning-window-derived start/days, totalPages/goToPrevPage/goToNextPage, and planning-window-change page-reset watcher |
| `frontend/src/home/HomeView.vue` | Wires paginated composable; renders prev/next navigation bar alongside read-only DailySchedule |
| `frontend/src/calendar/CalendarView.vue` | Wires paginated composable; renders prev/next navigation bar alongside editable DailySchedule |
| `frontend/src/calendar/DailySchedule.vue` | Unchanged — stays presentational with `days`/`readOnly` props only |
| `frontend/tests/calendar/useDailySchedule.test.ts` | Extended with page derivation, boundary navigation, and planning-window-change reset assertions |
| `frontend/tests/calendar/DailySchedule.test.ts` | Unchanged — component contract is preserved |
| `frontend/tests/app-routes-render.test.ts` | Extended with pagination control assertions for `/` and `/calendar` routes |
| Backend | No changes required |
