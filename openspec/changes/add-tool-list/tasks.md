## 1. Database & query layer

- [x] 1.1 Add migration `backend/migrations/012_create_tools_table.sql` creating `tools` (`id` text PK, `name` text not null, `brought_by` text null references `person(id)` ON DELETE SET NULL, `sort_order` int not null) plus an ID sequence/default consistent with existing `tool-`-style prefixing
- [x] 1.2 (Optional) Seed a few demo tools in the migration or a seed migration, matching the demo-data pattern in `002_seed_demo_data.sql`
- [x] 1.3 Add `backend/queries/tools.sql` with sqlc queries: list tools ordered by sort_order, get by id, create (append sort_order), update name+sort_order, delete, set bringer, clear bringer, and clear-bringer-by-person (or rely on FK ON DELETE SET NULL)
- [x] 1.4 Regenerate the db layer via the Docker sqlc path so `backend/db/tools.sql.go` and `models.go` are produced
- [x] 1.5 Verify generation with `scripts/build-go`

## 2. Backend API

- [x] 2.1 Create `backend/api/tools.go` with request/response types: `ToolsBody` (`summary{total,claimed,open}`, `tools[]{id,name,broughtBy}`), create/update bodies, and bringer body `{personId}`
- [x] 2.2 Implement `GET /api/tools` computing the derived coverage summary server-side (mirror `TaskBacklogBody`/`TaskSummary`)
- [x] 2.3 Implement `POST /api/tools` (validate non-empty name, assign `tool-` ID + append sort order, null bringer), `PUT /api/tools/{id}` (update name/sort, 404 unknown), `DELETE /api/tools/{id}` (404 unknown)
- [x] 2.4 Implement `PUT /api/tools/{id}/bringer` (validate tool exists → 404, person exists → 400, replace existing bringer, no availability check) and idempotent `DELETE /api/tools/{id}/bringer`
- [x] 2.5 Add the corresponding `Store` interface methods and their implementations, reusing `store.PersonExists`
- [x] 2.6 Register the tools endpoints alongside the existing registrations (mirror `registerTasksBacklog`/`registerTasksEndpoints`)
- [x] 2.7 Add Go tests covering: summary counts, create validation, 404s, claim/replace/unclaim idempotency, unknown-person rejection, and person-delete reverting claimed tools to open
- [x] 2.8 Verify with `scripts/check-go` and `scripts/test-go`

## 3. OpenAPI & generated client

- [x] 3.1 Regenerate the OpenAPI spec from the updated Huma registrations (build-time generation path)
- [x] 3.2 Regenerate the frontend client (`frontend/src/client/`) — types and Pinia Colada query/mutation helpers for the new tools operations
- [x] 3.3 Verify with `scripts/check-dashboard`

## 4. Frontend — tools feature

- [x] 4.1 Create `frontend/src/tools/types.ts` (component-local `Tool`, summary, props) mirroring `frontend/src/tasks/types.ts`
- [x] 4.2 Create a `useTools` composable wrapping the generated `GET /api/tools` query plus claim/unclaim/create/delete mutations, exposing `isLoading`/`isError`/`isEmpty` and a refresh, modeled on `useTaskBacklog`
- [x] 4.3 Create `ToolsView.vue` listing tools, visually distinguishing crossed-off (showing bringer) from open, with loading/error/empty states
- [x] 4.4 Add a bringer picker fed by the existing people query (the same source as `usePeopleAvailability`), with no availability filtering; wire select → claim, and an unclaim control on crossed-off tools
- [x] 4.5 Add organizer add-tool and remove-tool controls wired to create/delete
- [x] 4.6 Add component/composable tests under `frontend/tests/tools/` or `frontend/src/tools/__tests__/` following the people/tasks test patterns

## 5. Routing, navigation & KPI

- [x] 5.1 Register the `/tools` route in `frontend/src/app/routes.ts`
- [x] 5.2 Add the **Tools** nav item (`WrenchIcon`, route `/tools`) to the Plan group in `AppSidebar`, keeping the seven-item / icon ordering from the sidebar spec
- [x] 5.3 Add the **Tools covered** KPI card (fifth, `border-info` / `bg-info-soft text-info` / `WrenchIcon`, `claimed / total` fraction) to `frontend/src/home/components/KpiCards.vue`, wired to the tools read
- [x] 5.4 Update `KpiCards.spec.ts` for five cards and the new card's assertions; update any sidebar nav test for the new item

## 6. Verification

- [x] 6.1 Run `scripts/precommit-run` and resolve any failures
- [x] 6.2 Confirm all spec scenarios in `tool-list-api`, `tool-list-ui`, `kpi-summary-cards`, and `sidebar-navigation` are satisfied
