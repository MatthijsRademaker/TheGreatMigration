## ADDED Requirements

### Requirement: Backend SHALL expose persisted room/area CRUD via Store-backed Huma v2 handlers

The backend SHALL register four Huma v2 handlers in `backend/main.go` for room/area management:
- `GET /api/rooms` — list all room/area records
- `POST /api/rooms` — create a new room/area record
- `PUT /api/rooms/{id}` — update an existing room/area record by ID
- `DELETE /api/rooms/{id}` — delete a room/area record by ID

Each handler SHALL accept a `Store` parameter and delegate to the corresponding Store method. Create and update handlers SHALL validate that `type` is one of `room` or `area`. Update and delete handlers SHALL return a 404 response when the specified ID does not exist.

#### Scenario: List rooms returns all records

- **WHEN** `GET /api/rooms` is called and rooms exist in the database
- **THEN** the response is 200 OK with a JSON array of room objects, each containing `id`, `name`, `type`, `createdAt`, and `updatedAt`

#### Scenario: Create room persists and returns the new record

- **WHEN** `POST /api/rooms` is called with valid `name` and `type` ("room" or "area")
- **THEN** the response is 201 with a JSON object containing the new room's `id`, `name`, `type`, `createdAt`, and `updatedAt`
- **AND** the record is persisted and visible in subsequent `GET /api/rooms` calls

#### Scenario: Create room rejects invalid type

- **WHEN** `POST /api/rooms` is called with a `type` value other than "room" or "area"
- **THEN** the response is a 422 validation error

#### Scenario: Update room modifies an existing record

- **WHEN** `PUT /api/rooms/{id}` is called with a valid existing `id` and updated `name` and `type` fields
- **THEN** the response is 200 with the updated room object
- **AND** subsequent `GET /api/rooms` reflects the changes

#### Scenario: Update room returns 404 for unknown ID

- **WHEN** `PUT /api/rooms/{id}` is called with an `id` that does not exist
- **THEN** the response is 404

#### Scenario: Delete room removes the record

- **WHEN** `DELETE /api/rooms/{id}` is called with a valid existing `id`
- **THEN** the response is 204 No Content
- **AND** subsequent `GET /api/rooms` no longer includes the deleted room

#### Scenario: Delete room returns 404 for unknown ID

- **WHEN** `DELETE /api/rooms/{id}` is called with an `id` that does not exist
- **THEN** the response is 404

#### Scenario: Room endpoints appear in OpenAPI

- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document includes paths for `/api/rooms` (GET, POST), `/api/rooms/{id}` (PUT, DELETE) with complete request and response schemas

### Requirement: Room/area records SHALL persist across application restarts via a goose migration

A new goose migration SHALL create a `rooms_areas` table with columns:
- `id` (TEXT PRIMARY KEY)
- `name` (TEXT NOT NULL)
- `type` (TEXT NOT NULL with a CHECK constraint limiting values to `room` or `area`)
- `created_at` (TIMESTAMPTZ NOT NULL DEFAULT NOW())
- `updated_at` (TIMESTAMPTZ NOT NULL DEFAULT NOW())

A seed migration SHALL populate the table with 8 demo records derived from existing room labels in backlog and schedule seed data: Kitchen, Living Room, Bedroom 1, Bedroom 2, Garage, Bedroom, Storage, and Office. Each record SHALL be assigned an appropriate `type` (room or area).

#### Scenario: Migration creates the rooms_areas table

- **WHEN** goose migrations are applied to an empty Postgres database
- **THEN** the `rooms_areas` table exists with the expected columns and constraints
- **AND** inserting an invalid `type` value is rejected by the database

#### Scenario: Seed migration populates demo rooms

- **WHEN** the seed migration is applied
- **THEN** the `rooms_areas` table contains 8 records
- **AND** each record has a non-empty `name` and a `type` of either "room" or "area"

#### Scenario: Seeded data persists across restarts

- **WHEN** the backend is restarted after seed data is applied
- **THEN** `GET /api/rooms` returns the 8 seeded records

### Requirement: Deleting a room SHALL NOT cascade to existing task or schedule data

Room/area delete operations SHALL NOT modify or cascade to existing `backlog_tasks.room` or `schedule_task_cards.room_area` TEXT values. Deleted room names may continue appearing as historical text in task and schedule records. The `rooms_areas` table SHALL NOT define foreign key relationships to the backlog or schedule tables.

#### Scenario: Deleted room name persists in task data

- **WHEN** a room (e.g., "Kitchen") is deleted from the rooms catalog
- **THEN** `GET /api/tasks/backlog` continues to return tasks with `room` equal to "Kitchen"
- **AND** `GET /api/dashboard/daily-schedule` continues to return task cards with `roomArea` equal to "Kitchen"

#### Scenario: No foreign key constraint blocks room deletion

- **WHEN** a room is deleted that matches a name used in backlog or schedule records
- **THEN** the delete succeeds without database errors

### Requirement: The Store interface SHALL include room CRUD methods

The `Store` interface in `backend/store.go` SHALL be extended with:
- `ListRooms(ctx context.Context) ([]Room, error)`
- `CreateRoom(ctx context.Context, input CreateRoomInput) (*Room, error)`
- `UpdateRoom(ctx context.Context, id string, input UpdateRoomInput) (*Room, error)`
- `DeleteRoom(ctx context.Context, id string) error`

`PgStore` SHALL implement all four methods using sqlc-generated queries. `mockStore` SHALL implement them with in-memory state backed by a map. `failingStore` and `partialFailingStore` SHALL return `errTestFailure` for all room methods.

#### Scenario: PgStore satisfies the expanded Store interface

- **WHEN** `PgStore` is compiled
- **THEN** it implements all room methods of the expanded `Store` interface without compilation errors

#### Scenario: mockStore satisfies the expanded Store interface

- **WHEN** `mockStore` is compiled
- **THEN** it implements all room methods of the expanded `Store` interface without compilation errors

#### Scenario: Room unit tests can use mockStore for success paths

- **WHEN** a room handler unit test uses `mockStore`
- **THEN** list, create, update, and delete operations behave deterministically based on the mock's in-memory state

#### Scenario: Room unit tests can use failingStore for error paths

- **WHEN** a room handler unit test uses `failingStore`
- **THEN** all room operations return errors, enabling failure-path coverage

### Requirement: Backend tests SHALL cover room handler success and failure paths

Backend unit tests SHALL cover the room CRUD handlers using `mockStore` (success) and `failingStore` (failure). Integration tests (build tag `integration`) SHALL cover the persisted CRUD flow against a migrated Postgres database.

#### Scenario: Unit tests cover room list handler

- **WHEN** `go test -short ./...` runs in `backend/`
- **THEN** a test exercises `GET /api/rooms` with mockStore returning rooms and asserts the response body matches

#### Scenario: Unit tests cover room create handler success and validation failure

- **WHEN** `go test -short ./...` runs in `backend/`
- **THEN** a test exercises `POST /api/rooms` with valid input and asserts 201
- **AND** a test exercises `POST /api/rooms` with invalid type and asserts 422

#### Scenario: Unit tests cover room update handler success and not-found

- **WHEN** `go test -short ./...` runs in `backend/`
- **THEN** a test exercises `PUT /api/rooms/{id}` with valid input and asserts 200
- **AND** a test exercises `PUT /api/rooms/{id}` with unknown ID and asserts 404

#### Scenario: Unit tests cover room delete handler success and not-found

- **WHEN** `go test -short ./...` runs in `backend/`
- **THEN** a test exercises `DELETE /api/rooms/{id}` with valid ID and asserts 204
- **AND** a test exercises `DELETE /api/rooms/{id}` with unknown ID and asserts 404

#### Scenario: Unit tests cover handler failure paths

- **WHEN** `go test -short ./...` runs in `backend/`
- **THEN** tests using `failingStore` exercise error paths for list, create, update, and delete handlers

#### Scenario: Integration tests cover Postgres-backed CRUD

- **WHEN** `go test -tags=integration ./...` runs in `backend/` with Docker available
- **THEN** tests validate create, list, update, and delete operations against a migrated Postgres database

### Requirement: Frontend SHALL replace the /rooms placeholder with a functional management screen

The `RoomsView.vue` component SHALL render a single management screen using existing shared primitives (Card, Input, Button, Select). The screen SHALL include:
- A list of all room/area records loaded via the generated `listRooms` Pinia Colada query
- A compact create/edit form with `name` (Input) and `type` (Select) fields
- Create, Edit, and Delete action buttons per record
- Loading, empty, and error states for the room list

#### Scenario: Room list loads and displays records

- **WHEN** the `/rooms` route is rendered and the backend returns room records
- **THEN** the view displays each room's `name` and `type` in a list or card layout
- **AND** each record shows Edit and Delete action buttons

#### Scenario: Create form adds a new room

- **WHEN** a user fills in the create form with name and type and submits
- **THEN** the new room appears in the list without a manual page reload
- **AND** the create form fields are cleared

#### Scenario: Edit action populates the form

- **WHEN** a user clicks Edit on an existing room
- **THEN** the form is populated with that room's current `name` and `type` values
- **AND** submitting the form updates the room in the list without a manual page reload

#### Scenario: Delete action removes a room

- **WHEN** a user clicks Delete on an existing room
- **THEN** the room is removed from the list without a manual page reload

#### Scenario: Empty state displays when no rooms exist

- **WHEN** the backend returns an empty room list
- **THEN** the view displays an empty-state message indicating no rooms have been created yet

#### Scenario: Loading state displays while data is fetching

- **WHEN** the room list query is in-flight
- **THEN** the view displays a loading indicator rather than stale content

#### Scenario: Error state displays on query failure

- **WHEN** the room list query fails
- **THEN** the view displays a graceful error message

### Requirement: Frontend SHALL use Pinia Colada mutations with cache invalidation for room CRUD

Room create, update, and delete actions SHALL use `useMutation` from Pinia Colada with `invalidateCache` to refresh the room list after each successful mutation. No manual fetch calls, manual state management, or ad-hoc list manipulation SHALL be used.

#### Scenario: Create mutation invalidates room list cache

- **WHEN** a room is successfully created via mutation
- **THEN** the room list query is automatically re-fetched
- **AND** the updated list includes the newly created room

#### Scenario: Update mutation invalidates room list cache

- **WHEN** a room is successfully updated via mutation
- **THEN** the room list query is automatically re-fetched
- **AND** the updated list reflects the changes

#### Scenario: Delete mutation invalidates room list cache

- **WHEN** a room is successfully deleted via mutation
- **THEN** the room list query is automatically re-fetched
- **AND** the updated list no longer includes the deleted room

### Requirement: Frontend SHALL consume room endpoints through the committed OpenAPI snapshot and generated client

Room CRUD operations SHALL use SDK functions generated from the committed `frontend/openapi-snapshot.json`, not ad-hoc fetch calls. The generated SDK SHALL include `listRooms`, `createRoom`, `updateRoom`, and `deleteRoom` functions. The generated Pinia Colada plugin SHALL produce query options for list and mutation options for create, update, and delete.

#### Scenario: Generated SDK includes room functions

- **WHEN** `npm run generate:api` completes after snapshot refresh
- **THEN** `frontend/src/client/sdk.gen.ts` exports `listRooms`, `createRoom`, `updateRoom`, and `deleteRoom` functions

#### Scenario: Generated Pinia Colada plugin includes room hooks

- **WHEN** `npm run generate:api` completes after snapshot refresh
- **THEN** `frontend/src/client/@pinia/colada.gen.ts` exports `listRoomsQuery` and mutation options for create, update, and delete

#### Scenario: Frontend verification does not depend on a live backend

- **WHEN** `scripts/check`, `scripts/test`, or `scripts/precommit-run` are executed without a running backend
- **THEN** frontend install, type-check, and test workflows complete using the committed snapshot and generated client artifacts

### Requirement: SSR route render test SHALL mock the rooms endpoint

The `renderRoute` helper in `frontend/tests/app-routes-render.test.ts` SHALL provide a mock response for `GET /api/rooms` returning a valid room list (at minimum an empty array). The `/rooms` route test assertion SHALL change from "Feature coming soon" to assertions for room-management UI elements (form fields, action buttons, or empty state).

#### Scenario: Rooms endpoint is mocked for SSR rendering

- **WHEN** the SSR route render test runs for `/rooms`
- **THEN** the `renderRoute` helper returns a valid 200 mock response for `GET /api/rooms`
- **AND** the rendered output does not contain "Feature coming soon"

#### Scenario: /rooms route test asserts management UI

- **WHEN** the SSR route render test runs for `/rooms`
- **THEN** the rendered output contains room-management affordances (form input, select, button elements)
- **AND** the route metadata title and description are still present

#### Scenario: Existing route render assertions continue to pass

- **WHEN** the SSR route render test runs for all routes
- **THEN** all existing assertions for shell chrome, sidebar navigation, page titles, descriptions, and non-room route content continue to pass

### Requirement: The dashboard Rooms completed KPI SHALL remain an isolated static placeholder

The "Rooms completed" KPI card on the home page SHALL continue to display a static placeholder value (`—`) with `data-testid="kpi-placeholder-rooms-completed"`. It SHALL NOT query the rooms endpoint or derive any value from the room CRUD model. The existing `kpi-summary-cards` specification for this card SHALL remain unchanged.

#### Scenario: Rooms completed KPI is not data-driven

- **WHEN** the KPI cards component renders after room CRUD is implemented
- **THEN** the "Rooms completed" card displays `—` as its value
- **AND** no network request for room data is made from the KPI component

#### Scenario: Rooms completed KPI render output is unchanged

- **WHEN** the home route SSR test runs
- **THEN** the rendered output still contains `data-testid="kpi-placeholder-rooms-completed"`
- **AND** the card label is still "Rooms completed"

## MODIFIED Requirements

### Requirement: The /rooms route SHALL render a functional room management view instead of placeholder content

The `/rooms` route defined in `frontend/src/app/routes.ts` SHALL continue to resolve to `RoomsView.vue` and render through the shared `AppShell` layout. The view SHALL replace its placeholder content with a functional room management screen that fetches real data, supports CRUD operations, and handles loading/empty/error states. The route's `meta.title` ("Rooms / Areas") and `meta.description` SHALL remain unchanged.

#### Scenario: RoomsView renders through AppShell with management UI

- **WHEN** the RoomsView component is rendered
- **THEN** the content is wrapped by the `AppShell` layout component
- **AND** no additional sidebar, header, or navigation chrome is rendered inside the view
- **AND** the view contains room management affordances (list, form, action buttons) rather than placeholder text

#### Scenario: Route metadata is preserved

- **WHEN** the application routes to `/rooms`
- **THEN** the route resolves to `RoomsView.vue`
- **AND** the rendered HTML contains the route's `meta.title` ("Rooms / Areas") and `meta.description`

### Requirement: Existing non-room endpoints and behaviors SHALL remain unchanged

All existing backend endpoints (`/api/hello`, `/api/planning-window`, `/api/dashboard/people-availability`, `/api/tasks/backlog`, `/api/dashboard/daily-schedule`) SHALL continue to return their current response contracts. The dashboard "Rooms completed" KPI SHALL continue to render as a static placeholder. The `backlog_tasks.room` and `schedule_task_cards.room_area` columns SHALL remain free-form TEXT fields with no foreign key constraints.

#### Scenario: Existing endpoints are unaffected

- **WHEN** the backend is running after room CRUD is registered
- **THEN** `GET /api/hello` returns `{"message": "Hello from the backend!"}`
- **AND** all other existing endpoints return their current response contracts

#### Scenario: Backlog and schedule room fields remain free-form

- **WHEN** the backend schema is inspected after room CRUD migration
- **THEN** `backlog_tasks.room` remains a `TEXT NOT NULL` column with no foreign key to `rooms_areas`
- **AND** `schedule_task_cards.room_area` remains a `TEXT NOT NULL` column with no foreign key to `rooms_areas`

#### Scenario: Dashboard Rooms completed KPI is unchanged

- **WHEN** the home page renders after room CRUD is implemented
- **THEN** the "Rooms completed" card renders with `data-testid="kpi-placeholder-rooms-completed"` and value `—`