## Context

The backend manages rooms and areas as a first-class persisted entity through `GET /api/rooms` (and full CRUD), seeded with 8 records: Kitchen, Living Room, Bedroom 1, Bedroom 2, Garage, Bedroom, Storage, Office. The frontend already has a dedicated `/rooms` management page (`RoomsView.vue`) using the generated `listRoomsQuery` from the Pinia Colada client.

However, the two primary creation forms — the schedule-card modal in `/calendar` (CalendarView.vue) and the task modal in `/tasks` (TasksView.vue) — both use free-form text `<Input>` fields for Room / Area. This means:

1. Users can type names that don't match any managed room (typos, misspellings, new names that were never created via `/rooms`)
2. Backlog tasks and schedule cards can silently reference non-existent rooms
3. The Rooms management page exists but its data isn't consumed by the forms that need it most
4. The Scheduled Date field in CalendarView also uses a free-form text Input even though a `DatePicker` shared UI component already exists

Both views already consume backend data for people assignment (checkboxes from `usePeopleAvailability`), so the pattern of querying backend reference data in forms is established.

## Goals / Non-Goals

**Goals:**

- Replace the free-form Room / Area Input with a `Select` populated from `GET /api/rooms` in both CalendarView and TasksView
- Replace the free-form Scheduled Date Input with the shared `DatePicker` component in CalendarView
- Add loading and error states for the rooms query in both modals (room data may lag behind the schedule/task data)
- Update tests to assert room-select behavior instead of free-form room entry

**Non-Goals:**

- No backend changes — `GET /api/rooms` already returns the canonical room list
- No client regeneration — `listRoomsQuery` and its mutation variants already exist in the generated client
- No schema changes — the `room`/`roomArea` field stays a string on the backend; the Select merely constrains the input at the UI layer
- No FK enforcement at the database level — rooms are a soft reference, not a hard constraint
- No changes to HomeView (it's read-only and uses TaskManagementPanel, which doesn't have its own form)
- No changes to DailySchedule.vue itself (it's presentational — the form lives in CalendarView)

## Decisions

| Decision | Choice | Alternative | Rationale |
|----------|--------|-------------|-----------|
| Room query sharing | Inline `useQuery(listRoomsQuery())` in each view | Shared composable `useRooms()` | Both views only need the raw room list for one Select. A shared composable adds indirection for a single query. If a third consumer emerges, extract then. |
| Room selection model | Select by room `name` string (matching the existing `room`/`roomArea` backend field) | Select by room `id` | The backend accepts room as a string, not a FK. Changing the backend contract to accept room IDs would be a schema/API change out of scope. The Select displays room names and sets the value to the name. |
| Date picker | Reuse existing `DatePicker` shared UI component | Keep free-form Input, or build a new date widget | The shared DatePicker already exists at `frontend/src/shared/ui/date-picker/DatePicker.vue` and is registered in the UI barrel. Using it avoids a new component and provides a consistent UX. |
| Loading/error handling | Rooms query runs in parallel with the primary query; modal shows a loading state for the Select while rooms fetch | Require rooms before enabling the form | The schedule/task query may resolve before rooms. Showing a loading spinner inside the room field while rooms resolve is less disruptive than blocking the entire modal. If rooms fail, the room Select shows an error state with a retry affordance. |

## Risks / Trade-offs

| Risk | Severity | Mitigation |
|------|----------|------------|
| Rooms query fails but user still needs to create a task/schedule card | Medium | Show an error message in the room Select area with a retry button. The Select is a single field — users can still fill the rest of the form while rooms load. |
| A room is created in `/rooms` after the form loads (stale room list) | Low | Room management is infrequent. The rooms query runs once per modal open. If a room was just created, the user can close and reopen the modal to refresh. A future enhancement could use query invalidation on room mutations. |
| Room name changes in `/rooms` after a task references it | Low | Room names are soft references. If a room is renamed, existing tasks keep the old name string. This matches the current behavior with free-form inputs and is not regressed. |
| DatePicker component has different UX than a simple text field | Low | The DatePicker provides a calendar popover that is strictly better than typing YYYY-MM-DD. If users preferred typing, the DatePicker's input field still accepts direct text entry. |
