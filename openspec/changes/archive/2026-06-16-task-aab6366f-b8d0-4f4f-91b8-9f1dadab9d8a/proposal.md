## Why

The home dashboard (`/` route) currently renders a hello-world card, two hardcoded summary cards, and a planning-window card in the top summary row instead of the four KPI summary cards shown in section 3 of `designs/components.png`. The backend already exposes contract-backed endpoints for people availability (`GET /api/dashboard/people-availability`) and task backlog summaries (`GET /api/tasks/backlog`), but the frontend's committed generated client and OpenAPI snapshot only expose hello, planning-window, and people-availability — the task backlog endpoint is missing from the codegen artifacts. This change wires those existing backend contracts into live KPI cards on the home page, replacing the static placeholder row with data-driven summary cards.

## What Changes

- Regenerate `frontend/openapi-snapshot.json` and all `frontend/src/client/` artifacts to include the `/api/tasks/backlog` endpoint so the frontend can consume `highPriorityTasks` and `unassignedTasks` through committed generated Pinia Colada queries.
- Create a feature-local `frontend/src/home/components/KpiCards.vue` component that renders four KPI cards using existing Card, CardHeader, CardContent, CardTitle, CardDescription, and Badge primitives.
- Wire the **People available today** card to the existing `getDashboardPeopleAvailabilityQuery`, displaying `availableToday` of `totalPeople`.
- Wire the **High priority tasks** card to `highPriorityTasks` from the newly generated `getTasksBacklogQuery`.
- Wire the **Unassigned jobs** card to `unassignedTasks` from the newly generated `getTasksBacklogQuery`.
- Render the **Rooms completed** card as an explicitly isolated placeholder with a documented data-testid, displaying a placeholder value (`—`) and a code comment linking to a future room-progress contract. It must not derive counts from any unrelated backend data.
- Replace the current top-row summary grid in `HomeView.vue` with the new `KpiCards` component, removing the hello-world card from the top row.
- Preserve the existing lower `Today's plan` and `Move notes` sections unchanged.
- Update `frontend/tests/app-routes-render.test.ts` to mock `GET /api/dashboard/people-availability` and `GET /api/tasks/backlog` responses, and assert the four KPI card labels appear in the rendered home-route output.
- Style all cards using existing Card primitives, semantic status tokens, and lucide-vue icons from the Design System v2 foundation. No new global theme tokens or shared primitive directories are introduced.

## Impact

- **Home route rendering**: The top summary row changes from hello-world + 2 static + planning-window cards to four KPI cards driven by backend contracts. The hello-world developer visibility aid moves off the top grid.
- **Generated client artifacts**: `frontend/openapi-snapshot.json` gains the `/api/tasks/backlog` path and schema. `frontend/src/client/sdk.gen.ts`, `types.gen.ts`, and `@pinia/colada.gen.ts` gain task backlog types, SDK function, and query exports.
- **SSR verification**: `frontend/tests/app-routes-render.test.ts` must mock two additional API endpoints for the home route to render successfully.
- **Non-impacted**: Backend services, persistence, sidebar, top bar, people matrix, schedule board, rooms/areas page, settings page, and the lower HomeView sections (Today's plan, Move notes) remain unchanged.
