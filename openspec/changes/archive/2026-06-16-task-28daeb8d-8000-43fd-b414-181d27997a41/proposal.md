## Why

The daily planner and task backlog forms both use free-form text Inputs for the Room / Area field, even though the backend manages rooms as a first-class persisted entity with its own CRUD API (`GET /api/rooms` returns seeded Kitchen, Living Room, Bedroom 1, etc.). This creates a data-integrity gap: users can type typos, invent room names that don't exist in the managed list, and drift between what the backend knows and what the forms accept. Replacing the free-form Input with a Select populated from `GET /api/rooms` keeps room names consistent across the product and leverages the existing backend resource.

## What Changes

- Replace the free-form text `Input` for "Room / Area" in `CalendarView.vue` (daily planner schedule-card create/edit modal) with a `Select` populated from `listRoomsQuery`.
- Replace the free-form text `Input` for "Room / Area" in `TasksView.vue` (task backlog create/edit modal) with a `Select` populated from `listRoomsQuery`.
- Replace the free-form text `Input` for "Scheduled date" in the CalendarView modal with the existing shared `DatePicker` component.
- Add a local `useRooms()` composable or inline query usage that fetches `listRoomsQuery` and provides the room list with loading/error states to both views.
- Update frontend tests to assert room selection from backend-shaped room data instead of free-form room entry.
- No backend changes — `GET /api/rooms` already exists and returns the required data.

## Capabilities

### New Capabilities

*(none — this change modifies existing capabilities, it doesn't introduce new ones)*

### Modified Capabilities

- `daily-schedule-component`: The schedule-card create/edit form in the `/calendar` route changes the Room / Area field from a free-form text `Input` to a `Select` populated by `GET /api/rooms`. The Scheduled Date field changes from a free-form text `Input` to the existing shared `DatePicker` component.
- `task-management-ui`: The task create/edit form in the `/tasks` route changes the Room / Area field from a free-form text `Input` to a `Select` populated by `GET /api/rooms`.

## Impact

| Area | Impact |
|------|--------|
| `frontend/src/calendar/CalendarView.vue` | Replace `Input#form-room` with `Select` fetching from `listRoomsQuery`; replace `Input#form-date` with `DatePicker`; add room query loading/error states to the modal |
| `frontend/src/tasks/TasksView.vue` | Replace `Input#form-room` with `Select` fetching from `listRoomsQuery`; add room query loading/error states to the modal |
| Backend | No change — `GET /api/rooms` already exists and returns all 8 seeded rooms |
| Frontend client | No regeneration needed — `listRoomsQuery` already exists in generated client |
| Tests | Update CalendarView and TasksView modal tests to assert room Select behavior and DatePicker usage instead of free-form Input assertions |
