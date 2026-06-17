## Why

The People Availability panel displays an editable matrix of people × days for assigning availability statuses. The Settings page lets users configure a planning window of any length (e.g., 40 days), but the availability matrix is hardcoded to a fixed 4-day view. Users can only see and assign availability for the first 4 days of the planning window; the remaining days are unreachable through the UI. This makes the planning window configuration functionally incomplete — users set a date range but cannot act on it in the only panel designed for day-by-day availability assignment.

## What Changes

- **People Availability panel gains day pagination**: The matrix now shows a configurable number of days per page (default 7). Navigation controls (prev/next, page indicator) let users move through the full planning window.
- **The `usePeopleAvailability` composable exposes a `days` parameter**: Callers can specify the number of days to display, replacing the hardcoded default of 4. The composable derives the current page's date range from the planning window and the page state.
- **Backend `GET /api/dashboard/people-availability` gains person-level pagination**: New optional query parameters `offset` and `limit` let the frontend request a subset of people. The response gains pagination metadata (`totalPeople`, `page`, `perPage`). This prevents unbounded payload growth when the planning window spans many days and the organization has many people — a 40-day × 100-person matrix would otherwise produce 4000 cells per request.
- **PeopleView passes the visible date range as `start`/`days`**: Instead of relying on the default 4-day window, PeopleView computes the current page's start date and page size and passes them explicitly to the API. This makes the pagination explicit and testable.
- **The PeopleAvailability component accepts pagination metadata**: New props for `startDate`, `endDate`, `currentPage`, `totalPages`, `perPage`, and `totalDays` let the component render navigation controls with awareness of where the user is in the planning window.

## Capabilities

### New Capabilities

*(none — this change extends existing capabilities, it doesn't introduce new ones)*

### Modified Capabilities

- `dashboard-people-availability`: The endpoint adds optional `offset` and `limit` query parameters for person-level pagination, and returns `totalPeople`, `page`, and `perPage` fields in the response body. Existing behavior (no offset/limit returns all people as before) is preserved for backward compatibility.
- `people-availability-integration`: The frontend composable (`usePeopleAvailability`) exposes `days` and `page` parameters, computes the effective `start` date from the planning window and page state, and passes `start`, `days`, `offset`, and `limit` to the dashboard query. The PeopleView component adds day-navigation controls and renders the matrix for the current page. The PeopleAvailability component accepts new pagination props and renders prev/next controls.

## Impact

| Area | Impact |
|------|--------|
| `backend/api/dashboard.go` | `DashboardInput` gains `offset` and `limit` query params. Handler passes them to store. |
| `backend/api/store.go` | `GetPeopleAvailability` signature gains `offset` and `limit` parameters. |
| `backend/store.go` | `GetPeopleAvailability` implementation adds SQL `LIMIT/OFFSET` to the people query and returns pagination metadata. |
| `backend/db/people_availability.sql` / sqlc | New query variant for paginated people fetch; generated Go code updates. |
| `backend/main_test.go` | Tests for paginated response scenarios. |
| Frontend `usePeopleAvailability.ts` | Accepts `days` and `page` options; computes `start` from planning window + page; passes `offset`/`limit` to API. |
| Frontend `PeopleView.vue` | Adds day-navigation controls; passes page state to composable; delegates visible range to API params. |
| Frontend `PeopleAvailability.vue` | Accepts and renders pagination props (prev/next, page indicator, date range label). |
| Frontend `types.ts` | New types for pagination metadata. |
| OpenAPI snapshot | Regenerated to capture `offset`/`limit` params and `totalPeople`/`page`/`perPage` response fields. |
