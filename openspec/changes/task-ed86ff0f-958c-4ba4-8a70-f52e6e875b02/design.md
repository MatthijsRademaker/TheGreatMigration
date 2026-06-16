## Context

The codebase has strong, well-documented patterns for backend persistence (Store interface, Huma handler registration, sqlc-generated queries, goose migrations) and frontend API consumption (committed OpenAPI snapshot, Hey API codegen, Pinia Colada queries). The room/area CRUD feature fits cleanly into these patterns. The main architectural delta is the introduction of mutation endpoints — POST, PUT, and DELETE — which the backend already supports via its Chi CORS configuration but the frontend has never consumed. The frontend currently has only GET-only generated SDK functions and query-only Pinia Colada hooks, so room CRUD will establish the mutation and cache invalidation pattern for the entire project.

Existing room references in `backlog_tasks.room` and `schedule_task_cards.room_area` are free-form TEXT strings with no foreign key constraints. This means the new rooms catalog can remain independent of historical task and schedule data without schema changes — a deliberate simplification that avoids the complexity of foreign key migrations and cascading updates.

## Goals

- Add a persisted `rooms_areas` table with columns `id` (TEXT PK), `name` (TEXT NOT NULL), `type` (TEXT NOT NULL with CHECK constraint for 'room'/'area'), `created_at`, and `updated_at`.
- Expose full CRUD endpoints through Huma v2 handlers backed by Store methods and sqlc-generated queries.
- Seed the table with 8 demo records derived from existing room labels in backlog and schedule seed data.
- Replace the `/rooms` placeholder view with a single-page management screen using existing shared primitives (Card, Input, Button, Select).
- Refresh the committed OpenAPI snapshot and regenerate frontend client artifacts to include room CRUD SDK functions and query/mutation hooks.
- Establish the first Pinia Colada `useMutation` + `invalidateCache` pattern in the frontend.
- Cover the new backend handlers with unit tests (mockStore/failingStore) and integration tests (Postgres-backed).
- Update the SSR route render test to mock the rooms endpoint and stop asserting "Feature coming soon".
- Keep the dashboard "Rooms completed" KPI placeholder unchanged.

## Non-Goals

- Building a floor-plan designer, checklist system, or other advanced room-planning tooling.
- Changing the home-page "Rooms completed" KPI from its current static placeholder behavior.
- Refactoring tasks or daily-schedule storage into foreign-key relationships to the rooms table.
- Adding authentication, permissions, or collaboration workflows.
- Deriving room progress or completion metrics from the CRUD model.
- Creating separate views, modals, wizards, or multi-step forms for room management.

## Decisions

### 1. Entity model: single rooms_areas table with type discriminator

Use a single `rooms_areas` table with a `type` column constrained to 'room' or 'area' rather than separate tables. This matches the dossier assumption that a single entity represents both concepts, keeps the schema and CRUD surface minimal, and follows the product's simplicity-first philosophy.

### 2. ID scheme: TEXT primary keys with deterministic seed IDs

Use TEXT primary keys (`room-1`, `room-2`, …) matching the existing pattern for people (`p1`-`p8`) and tasks (`task-1`-`task-11`). This ensures deterministic seed data, makes seed IDs predictable in tests, and avoids introducing a new ID convention (UUID or SERIAL) that would be inconsistent with the existing codebase.

### 3. Delete semantics: allow deletion without cascade

Deleting a room from the catalog does not modify or cascade to existing `backlog_tasks.room` or `schedule_task_cards.room_area` TEXT values. Those columns are free-form strings with no foreign key constraints, and the dossier non-goals explicitly exclude refactoring task storage into FK relationships. Orphaned room-name strings in task and schedule records are acceptable and represent a deliberate simplification.

### 4. Seed data strategy: seed 8 demo rooms from existing labels

Seed 6-8 demo records (Kitchen, Living Room, Bedroom 1, Bedroom 2, Garage, Bedroom, Storage, Office) derived from the room labels already used in `004_seed_backlog_schedule_data.sql`. This maintains the product's deterministic demo pattern, ensures the feature is immediately usable on first load, and avoids the need for a novel empty-state design.

### 5. Frontend mutation pattern: configure Hey API for full SDK generation

The existing `frontend/openapi-ts.config.ts` only generates GET SDK functions and query helpers. To support room create/update/delete, the configuration must be updated to produce POST, PUT, and DELETE SDK functions. The generated Pinia Colada plugin will produce `defineMutationOptions` for these operations. The frontend will use `useMutation` with `invalidateCache` to refresh the room list after mutations — establishing the pattern for all future mutations in the project.

### 6. Frontend architecture: single management screen

All room CRUD stays on a single `RoomsView.vue` screen. A compact form at the top handles create and edit (reusing the same form fields). A list/table of existing records below shows each room with Edit and Delete action buttons. Clicking Edit populates the form with the selected room's data. Delete removes the record and refreshes the list via cache invalidation. All UI uses existing shared primitives (Card, Input, Button, Select) — no new shared-ui components are created.

### 7. SSR rendering: mock rooms endpoint in renderRoute helper

The SSR `renderRoute` helper must mock `GET /api/rooms` (returning an empty or seeded list) so that `RoomsView.vue` renders its management screen rather than a loading state during SSR. The test assertion changes from "Feature coming soon" to asserting room-management affordances (form fields, action buttons).

### 8. Dashboard KPI: explicitly frozen

The "Rooms completed" KPI card remains an isolated static placeholder displaying "—" with `data-testid="kpi-placeholder-rooms-completed"`. It must not query the rooms endpoint or derive any value from the CRUD model. Both the dossier non-goals and the canonical `kpi-summary-cards` spec freeze this behavior.

## Risks

- **HIGH — Novel frontend mutation pattern**: The codebase has zero mutation examples. `useMutation` + `invalidateCache` from Pinia Colada v1.3.1 is untested in this repo. Mitigation: generate the client first to inspect mutation option shapes, then wire the UI. Consider a light integration spike before full frontend work.
- **MEDIUM — Store interface expansion**: Adding 4 methods to the `Store` interface requires updating `mockStore`, `failingStore`, and `partialFailingStore`. The `partialFailingStore` pattern currently only returns success for `GetPlanningWindow` — room methods must return `errTestFailure` unless custom partial behavior is needed. Mitigation: follow the existing failingStore pattern exactly; add room method stubs that return `errTestFailure`.
- **MEDIUM — OpenAPI snapshot drift**: If the backend `/openapi.json` is regenerated but the committed `frontend/openapi-snapshot.json` is not refreshed, the generated client will be stale and the frontend will fail to compile or fetch from wrong URLs. Mitigation: run `npm run regen:api` after backend handler registration is complete.
- **LOW — SSR test mock gap**: The existing `renderRoute` helper does not mock a rooms endpoint. After the change, `RoomsView.vue` will use query hooks that fetch on mount. SSR rendering will show a loading state unless a mock is added. Mitigation: add a rooms API mock to the `renderRoute` helper returning a valid room list response.
- **LOW — Pre-existing room string coupling in tests**: `store_mock_test.go` hardcodes room strings in `seedTasks` and `seedTasksForDay`. New room-specific tests must not accidentally couple to these backlog seed rooms. Mitigation: keep room CRUD test fixtures independent of backlog seed data.

## Conflict Resolution

- **Delete semantics (reviewer blocker)**: Resolved to allow deletion without cascade. The reviewer identified this as a blocker requiring explicit commitment. Both architect and lead-dev confirmed the approach, and the dossier non-goals explicitly exclude FK refactoring. The spec now explicitly states: "Deleting a room from the catalog SHALL NOT cascade to or modify existing backlog_tasks.room or schedule_task_cards.room_area TEXT values."
- **Seed data strategy (reviewer blocker)**: Resolved to seed 8 demo rooms from existing labels. All three agents agreed on seeding from demo labels to maintain the deterministic demo pattern. The spec now explicitly requires the seed migration.
- **Mutation pattern gap (reviewer risk)**: Resolved to configure `@hey-api/openapi-ts` to generate POST/PUT/DELETE SDK functions, establishing the mutation pattern for the project. Manual fetch calls are rejected in favor of generated SDK consistency.

## Traceability

- Task: `ed86ff0f-958c-4ba4-8a70-f52e6e875b02`
- Dossier: `2026-06-16T07:15:43.803Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial` snapshot for `task-ed86ff0f-958c-4ba4-8a70-f52e6e875b02`