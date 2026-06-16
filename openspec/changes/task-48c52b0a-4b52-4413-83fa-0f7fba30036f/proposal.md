## Why

`DailySchedule.vue` still renders demo fixtures on `/` and `/calendar` even though the BFF already exposes `GET /api/dashboard/daily-schedule` and the generated frontend client already exposes a typed query for it. The current schedule persistence model also cannot support a truly functional board: `schedule_task_cards` are grouped by demo-only `day_group`, the read endpoint returns synthetic per-request IDs, and there is no schedule-specific CRUD surface.

## What Changes

- Add a shared `useDailySchedule` read path that wraps the generated daily-schedule query, adapts nullable/generated data into the component contract, and exposes loading, error, empty, and success states.
- Wire `HomeView.vue` and `CalendarView.vue` to that shared read path so both routes render backend-derived schedule data while keeping the home dashboard read-only.
- Keep `DailySchedule.vue` presentational and prop-driven: it continues to emit typed add intents, while `/calendar` owns any persisted create/update/delete workflow.
- Evolve the schedule read model to use stable persisted card identifiers and a real `scheduled_date` instead of demo-only modulo day grouping.
- Add schedule-specific CRUD endpoints backed by `schedule_task_cards` and `schedule_task_assignments`, regenerate the OpenAPI snapshot/client artifacts, and invalidate daily schedule queries after successful mutations.
- Update frontend SSR/component tests and backend unit/integration coverage to reflect BFF-backed schedule behavior and the new write surface.

## Impact

- Daily Schedule becomes BFF-backed on both routes with explicit state handling instead of demo fallback.
- Schedule edits persist through the schedule tables rather than reusing backlog task CRUD.
- Existing people availability and backlog contracts remain separate and must stay green.
- Implementation spans Vue route/composable work, backend migration/store/API changes, generated artifacts, and Docker-backed verification.