## Why

The planning-window feature (task B9E23168) was delivered as frontend-only constants in `frontend/src/shared/lib/planWindow.ts`. Changing the move timeline currently requires a frontend code change and rebuild. The backlog intent is for date-range planning to be backend-configurable, not compile-time-configured. This change moves the canonical planning-window source of truth from hard-coded frontend constants to a backend HTTP endpoint, making the planning window runtime-configurable from the server side without a frontend redeploy.

## What Changes

- **Backend**: Add `GET /api/planning-window` as a new Huma v2 endpoint returning the global move range (`startDate`, `endDate`, `days`). The endpoint is seeded with the current business defaults (2026-07-05 through 2026-08-13) via Go constants, with no persistence, authentication, or multi-user configuration in this slice.
- **Frontend shared utility**: Refactor `frontend/src/shared/lib/planWindow.ts` to remove the hard-coded constants (`PLAN_WINDOW_START`, `PLAN_WINDOW_END`, `planWindowDays`, `planWindowDayCount`) and retain only pure utility exports (`formatPlanDayLabel`, `PlanWindowDay` interface).
- **Frontend composable**: Create `frontend/src/shared/composables/usePlanningWindow.ts` — a shared Pinia Colada composable that wraps the generated `getPlanningWindowQuery`, derives `planWindowDays` and `planWindowDayCount` as reactive computed refs, and exposes `isLoading` and `error` states.
- **Frontend consumers**: Update `HomeView.vue` and `CalendarView.vue` to consume `usePlanningWindow()` instead of directly importing from `planWindow.ts`. Both views gain loading and error handling for the planning-window data, following the existing Pinia Colada pattern demonstrated by the hello card in `HomeView.vue`.
- **OpenAPI artifacts**: Refresh `frontend/openapi-snapshot.json` from the running backend and regenerate `frontend/src/client/` (SDK, types, Pinia Colada query artifacts) via the established `refresh-openapi-snapshot.mjs` and `generate-api.mjs` scripts.
- **Tests**: Add backend contract tests for `GET /api/planning-window` in `backend/main_test.go`. Update frontend tests (`planWindow.test.ts`, `app-routes-render.test.ts`) to test composable behavior instead of asserting hard-coded constants.
- **OpenSpec artifacts**: Publish proposal, design, tasks, and planning-window spec delta describing the backend-backed contract.

## Impact

- **Affected specs**: `planning-window` (spec delta describing backend-backed contract), `hello-world-integration` (the Move days card scenario already references planning-window derivation; no amendment needed in this change).
- **Affected code**: `backend/main.go` (new endpoint registration), `backend/planning_window.go` (new file), `backend/main_test.go` (new test), `frontend/src/shared/lib/planWindow.ts` (refactored), `frontend/src/shared/composables/usePlanningWindow.ts` (new file), `frontend/src/home/HomeView.vue` (consumer update), `frontend/src/calendar/CalendarView.vue` (consumer update), `frontend/openapi-snapshot.json` (regenerated), `frontend/src/client/` (regenerated), `frontend/tests/planWindow.test.ts` (updated), `frontend/tests/app-routes-render.test.ts` (updated).
- **No impact on**: `GET /api/hello`, `GET /api/dashboard/people-availability`, PeopleView, CORS configuration, Vite proxy, Docker Compose setup, existing design-system components.
