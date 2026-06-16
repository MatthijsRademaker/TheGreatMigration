## Why

When a user clicks a status pill in the editable People Availability matrix to change someone's availability, the backend rejects the request with "date is outside the planning window" — surfaced to the user as the opaque error "Invalid status or date." The root cause is a date-range mismatch: the dashboard endpoint (`GET /api/dashboard/people-availability`) defaults to `time.Now()` when no `start` query parameter is provided, but the availability upsert endpoint (`PUT /api/people/{id}/availability/{date}`) validates dates against the planning window stored in the database. The dashboard can display dates outside the planning window, making those cells visually editable but functionally broken.

## What Changes

- **Dashboard endpoint** (`GET /api/dashboard/people-availability`): When no explicit `start` parameter is provided, default to the planning window's start date instead of `time.Now()`. This ensures the dashboard only displays dates that are writable.
- **Frontend PeopleView error handling**: Surface the actual backend error message instead of the generic "Invalid status or date" catch-all. Differentiate between invalid status, invalid date format, and date-outside-window errors.
- **Frontend PeopleView date derivation**: Remove the bespoke `getISODate` helper that derives a date from `dayIndex` using `rawData.range.startDate` — instead, derive the date from the `days` array index already available in the component's adapted props, eliminating the dependency on raw API response internals.
- **Preventive guard**: Add a guard in the frontend query to pass the planning window's `startDate` as the `start` query parameter, making the dashboard range explicitly tied to the planning window even when the backend default changes independently.

## Capabilities

### New Capabilities

- *(none — this change modifies existing capabilities, it doesn't introduce new ones)*

### Modified Capabilities

- `dashboard-people-availability`: The default `start` behavior changes — when no `start` query parameter is provided, the endpoint now defaults to the planning window's start date instead of `time.Now()`. The explicit `start` parameter continues to work as before for testing and manual overrides. A `400` error is returned when no planning window exists and no `start` is supplied.
- `people-availability-integration`: The frontend adapter/composable now explicitly passes the planning window `startDate` as the `start` query parameter to the dashboard endpoint. The `handleCellUpdate` function in PeopleView derives dates from the adapted `days` prop index rather than from `rawData.range.startDate`, removing the fragile `getISODate` helper. Error handling differentiates the three 400 subtypes (invalid status, invalid date format, date outside planning window) instead of using a single generic message.

## Impact

| Area | Impact |
|------|--------|
| `backend/api/dashboard.go` | Endpoint handler now reads planning window from store when `start` is empty |
| `backend/api/register.go` | No change — `registerDashboardPeopleAvailability` already receives `store` |
| `backend/store.go` | No change — `GetPlanningWindow` already exists |
| `backend/store_mock_test.go` | Mock store may need a `GetPlanningWindow` returns for dashboard tests |
| `frontend/src/shared/composables/usePeopleAvailability.ts` | Accept optional `start` param, call `usePlanningWindow()` to derive it, pass as query to dashboard endpoint |
| `frontend/src/people/PeopleView.vue` | Replace `getISODate` with date lookup from adapted `days` prop; improve error message differentiation |
| `frontend/src/people/PeopleAvailability.vue` | No change — emits `dayIndex` unchanged |
| `frontend/src/people/types.ts` | No change |
| `frontend/src/home/HomeView.vue` | No change — uses same composable which now passes `start` implicitly |
| Tests | Backend dashboard tests updated for new default behavior; frontend tests for error differentiation and date derivation patterns |
