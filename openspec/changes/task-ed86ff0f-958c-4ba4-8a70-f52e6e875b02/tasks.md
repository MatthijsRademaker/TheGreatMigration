## 1. Backend persistence layer

- [ ] Create migration `005_create_rooms_areas_table.sql` with columns: `id` (TEXT PRIMARY KEY), `name` (TEXT NOT NULL), `type` (TEXT NOT NULL CHECK IN ('room','area')), `created_at` (TIMESTAMPTZ DEFAULT NOW()), `updated_at` (TIMESTAMPTZ DEFAULT NOW()).
- [ ] Create seed migration `006_seed_rooms_areas.sql` populating 8 demo records: Kitchen, Living Room, Bedroom 1, Bedroom 2, Garage, Bedroom, Storage, Office with appropriate room/area type assignments.
- [ ] Create `backend/queries/rooms_areas.sql` with sqlc query definitions for `ListRooms`, `GetRoomByID`, `CreateRoom`, `UpdateRoom`, `DeleteRoom`.
- [ ] Run `sqlc generate` and verify generated Go code under `backend/db/` compiles.

## 2. Backend Store and handlers

- [ ] Extend the `Store` interface in `backend/store.go` with `ListRooms`, `CreateRoom`, `UpdateRoom`, `DeleteRoom` methods.
- [ ] Implement `PgStore` room CRUD methods using sqlc-generated queries.
- [ ] Implement `mockStore` room CRUD methods using in-memory state backed by a map.
- [ ] Update `failingStore` to return `errTestFailure` for all room methods.
- [ ] Update `partialFailingStore` to return `errTestFailure` for room methods.
- [ ] Create `backend/rooms_areas.go` with Huma v2 handler registration for `GET /api/rooms` (list), `POST /api/rooms` (create), `PUT /api/rooms/{id}` (update), `DELETE /api/rooms/{id}` (delete).
- [ ] Register room handlers in `backend/main.go` alongside existing handler registrations.
- [ ] Define typed input/output structs for each operation with appropriate Huma tags and validation.

## 3. Backend tests

- [ ] Add unit tests in `backend/main_test.go` covering room handler success paths using `mockStore`.
- [ ] Add unit tests covering room handler failure paths using `failingStore`.
- [ ] Add assertion that all four room endpoints appear in `/openapi.json`.
- [ ] Add integration tests in `backend/main_integration_test.go` covering persisted CRUD flow against Postgres.
- [ ] Verify existing handler tests continue to pass without modification.

## 4. Frontend API client refresh

- [ ] Start the backend with room CRUD handlers registered.
- [ ] Run `npm run refresh:openapi-snapshot` to commit the updated `frontend/openapi-snapshot.json` with `/api/rooms` paths.
- [ ] Update `frontend/openapi-ts.config.ts` if needed to produce POST/PUT/DELETE SDK functions.
- [ ] Run `npm run generate:api` to regenerate client artifacts under `frontend/src/client/`.
- [ ] Verify generated SDK includes `createRoom`, `updateRoom`, `deleteRoom` functions and Pinia Colada mutation options.

## 5. Frontend room management view

- [ ] Replace `frontend/src/rooms/RoomsView.vue` placeholder content with a single management screen.
- [ ] Implement room list display using data from the generated `listRooms` query hook.
- [ ] Implement a compact create/edit form with `name` (Input) and `type` (Select with room/area options) fields.
- [ ] Implement create action using `useMutation` with cache invalidation to refresh the list.
- [ ] Implement edit action: populate form with selected room data, save via `useMutation` with cache invalidation.
- [ ] Implement delete action with `useMutation` and cache invalidation.
- [ ] Handle loading, empty, and error states for the room list.
- [ ] Use existing shared primitives only: Card, Input, Button, Select from `@/shared/ui`.

## 6. Test updates

- [ ] Update `frontend/tests/app-routes-render.test.ts` to mock `GET /api/rooms` in the `renderRoute` helper returning a valid room list response.
- [ ] Change the `/rooms` route test assertion from "Feature coming soon" to assertions for room-management UI (form fields, list elements, action buttons).
- [ ] Verify all existing SSR route assertions continue to pass.
- [ ] Verify the dashboard "Rooms completed" KPI placeholder remains unchanged and does not query the rooms endpoint.

## 7. Verification

- [ ] Run `scripts/check` from the repo root.
- [ ] Run `scripts/test` from the repo root.
- [ ] Run `scripts/precommit-run` from the repo root.
- [ ] Verify room CRUD operations persist across backend restarts.
- [ ] Verify room create, update, and delete refresh the displayed list without manual browser reload.