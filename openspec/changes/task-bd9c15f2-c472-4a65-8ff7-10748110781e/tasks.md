## 1. Compose the home dashboard in `HomeView.vue`

- [x] Import `TaskManagementPanel`, `PeopleAvailability`, and `DailySchedule` from their existing feature locations.
- [x] Remove the `upcomingWork` array, the `Today’s plan` card, and the unused `Badge` import from `frontend/src/home/HomeView.vue`.
- [x] Keep `<KpiCards />` as the first row.
- [x] Add one responsive grid row with `TaskManagementPanel` in the wider column and `PeopleAvailability` in the narrower column.
- [x] Add a second responsive grid row with `DailySchedule` in the wider column and the existing static `Move notes` card in the narrower column.
- [x] Use two independent grid containers and reuse the existing `xl:grid-cols-[1.4fr_0.9fr]` pattern.

## 2. Preserve composition-only scope

- [x] Keep the change frontend-only: no backend wiring, no generated-client regeneration, no route changes, and no `AppShell.vue` or `AppSidebar.vue` edits.
- [x] Reuse existing component internals as-is unless a minimal composition-level layout adjustment is required.
- [x] Keep the `Move notes` content inline in `HomeView.vue`; do not add a reusable notes component or notes CRUD.

## 3. Update route-render verification

- [x] Remove `Today’s plan` assertions from the `/` route checks in `frontend/tests/app-routes-render.test.ts`.
- [x] Update the `/` `routeCases` entry to use `Move notes` as the home-route content anchor.
- [x] Add `/` route assertions for `Task Management`, `People availability`, `Daily Schedule`, and `Move notes`.
- [x] Preserve the existing KPI assertions for `People available today`, `High priority tasks`, `Unassigned jobs`, `Rooms completed`, and the mocked values `6`, `of 8`, `4`, and `3`.
- [x] Preserve the existing route-render assertions for shell chrome, sidebar navigation, and the other routes.

## 4. Run repository verification

- [x] Run `scripts/precommit-run`.
