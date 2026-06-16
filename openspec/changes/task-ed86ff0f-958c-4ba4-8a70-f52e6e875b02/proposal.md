## Why

The app exposes a `/rooms` route in navigation, but it is still a placeholder and the backend has no room/area model or CRUD endpoints. Room labels already appear across seeded backlog and schedule data, and the dashboard carries a separate static "Rooms completed" placeholder — users can see the concept of spaces without being able to create, edit, list, or delete them anywhere in the product. This change makes rooms and areas a first-class persisted feature with full CRUD on both the backend and frontend.

## What Changes

- **Backend**: Add a persisted `rooms_areas` table via a new goose migration, sqlc-generated queries, Store interface methods, and Huma v2 CRUD handlers (`GET /api/rooms`, `POST /api/rooms`, `PUT /api/rooms/{id}`, `DELETE /api/rooms/{id}`) following the existing handler registration pattern.
- **Seed data**: Add a seed migration populating the rooms_areas table with 8 demo records (Kitchen, Living Room, Bedroom 1, Bedroom 2, Garage, Bedroom, Storage, Office) derived from existing room labels in backlog and schedule seed data.
- **Frontend**: Replace the placeholder `RoomsView.vue` with a single management screen using existing shared Card, Input, Button, and Select primitives. The screen provides list display, a compact create/edit form, inline edit affordances, and delete actions — all on one page.
- **Client generation**: Refresh the committed `frontend/openapi-snapshot.json` from the backend `/openapi.json` and regenerate client artifacts under `frontend/src/client/`. Configure Hey API to produce POST/PUT/DELETE SDK functions alongside the query helpers.
- **Mutation pattern**: Establish the first `useMutation` + cache invalidation pattern in the frontend using Pinia Colada, so the room list refreshes after create, update, and delete without a manual browser reload.
- **Tests**: Add backend unit tests for room handler success/failure paths (mockStore/failingStore pattern), backend integration tests for Postgres-backed CRUD, and update the SSR route render test to mock the rooms endpoint and stop asserting "Feature coming soon".
- **Spec updates**: Amend the sidebar-navigation spec to require a functional room management view instead of placeholder content, and extend the backend-persistence spec to cover room CRUD Store methods and the rooms_areas table.

## Impact

- The backend gains 4 new CRUD handlers registered in `main.go` and 4 new Store interface methods; `PgStore`, `mockStore`, `failingStore`, and `partialFailingStore` all require updates.
- The Store interface grows from 4 to 8 methods; every implementing struct must be updated.
- A new goose migration (005) creates the `rooms_areas` table; a seed migration (006) populates demo data.
- New sqlc query file (`queries/rooms_areas.sql`) and regenerated code under `backend/db/`.
- The committed `frontend/openapi-snapshot.json` gains `/api/rooms` paths; frontend client artifacts are regenerated.
- The `/rooms` route test assertion changes from placeholder text to room-management UI assertions.
- Existing non-room endpoints, the dashboard "Rooms completed" KPI placeholder, and free-form room strings in backlog/schedule data remain unchanged.
- Delete operations on rooms do not cascade to or modify existing `backlog_tasks.room` or `schedule_task_cards.room_area` TEXT values.