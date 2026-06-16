## 1. Regenerate the frontend OpenAPI snapshot and client to include the task backlog endpoint

- [ ] Refresh `frontend/openapi-snapshot.json` to include the `/api/tasks/backlog` path and its `TaskBacklogBody` response schema by either running the backend and capturing `/openapi.json`, or manually appending the tasks backlog schema fragment.
- [ ] Regenerate `frontend/src/client/` artifacts (`types.gen.ts`, `sdk.gen.ts`, `@pinia/colada.gen.ts`) from the updated snapshot so `getTasksBacklogQuery` is available.
- [ ] Verify the generated artifacts pass `vue-tsc` type-checking (`npm run check` from `frontend/`).

## 2. Create the feature-local KpiCards component

- [ ] Create `frontend/src/home/components/KpiCards.vue` as a Vue 3 `<script setup>` SFC.
- [ ] Render four cards using existing Card, CardHeader, CardContent, CardTitle, and CardDescription primitives.
- [ ] Wire the **People available today** card to `getDashboardPeopleAvailabilityQuery`, displaying `availableToday` of `totalPeople` with a `UsersRoundIcon` icon.
- [ ] Wire the **High priority tasks** card to `highPriorityTasks` from `getTasksBacklogQuery` with a suitable lucide icon.
- [ ] Wire the **Unassigned jobs** card to `unassignedTasks` from `getTasksBacklogQuery` with a suitable lucide icon.
- [ ] Render the **Rooms completed** card as an explicit placeholder: value `—`, label `Rooms completed`, `data-testid="kpi-placeholder-rooms-completed"`, and a code comment marking it for future room-progress contract work.
- [ ] Handle query loading and error states for the availability and backlog cards (loading indicator, error fallback).
- [ ] Use semantic accent classes from the Design System v2 token surface for card icon backgrounds and accent treatments.

## 3. Integrate KpiCards into HomeView

- [ ] Replace the existing top-row `<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">` block in `frontend/src/home/HomeView.vue` with the `KpiCards` component.
- [ ] Remove the hello-world card (helloQuery logic), the hardcoded `summaries` array, and the planning-window card from the top grid.
- [ ] Remove unused imports (`helloQuery`, `getHelloQuery`, `usePlanningWindow`, `planWindowDayCount`, `moveDaysLoading`, `moveDaysError`) if they are no longer used on the page.
- [ ] Preserve the lower `<div class="grid flex-1 gap-4 xl:grid-cols-[1.4fr_0.9fr]">` block with "Today's plan" and "Move notes" cards unchanged.

## 4. Update the SSR route-render test

- [ ] Add mock responses for `GET /api/dashboard/people-availability` and `GET /api/tasks/backlog` in the `renderRoute` helper inside `frontend/tests/app-routes-render.test.ts`.
- [ ] The people-availability mock must return `summary.availableToday` and `summary.totalPeople` matching the contract shape.
- [ ] The tasks backlog mock must return `summary.highPriorityTasks` and `summary.unassignedTasks` matching the contract shape.
- [ ] Add assertions for the four KPI card labels (`People available today`, `High priority tasks`, `Unassigned jobs`, `Rooms completed`) in the home-route test block.
- [ ] Verify that all existing route-render assertions continue to pass.

## 5. Verify the change through the repository verification flow

- [ ] Run `scripts/check` from the repo root.
- [ ] Run `scripts/test` from the repo root.
- [ ] Run `scripts/precommit-run` from the repo root.