## Implementation Tasks

- [x] Publish spec deltas for `daily-schedule-component` and `dashboard-daily-schedule`, including the explicit move from demo-only schedule addressing to persisted, date-based schedule CRUD.
- [x] Add a shared `useDailySchedule` composable that wraps the generated daily-schedule query, adapts nullable arrays and canonical enum values, and exposes loading, error, empty, and success states.
- [x] Wire `HomeView.vue` and `CalendarView.vue` to the shared daily-schedule read path, remove demo-data fallback on successful reads, keep the home dashboard read-only, and preserve typed prop support for `DailySchedule.vue` tests.
- [x] Add a backend migration and seed/backfill changes so schedule cards are addressed by persisted `scheduled_date` values and returned with stable `sched-<id>` identifiers.
- [x] Extend the Store interface, sqlc queries, and Huma handlers with schedule-card `POST`, `PUT`, and `DELETE` operations backed by `schedule_task_cards` and `schedule_task_assignments`.
- [x] Implement `/calendar` schedule create/update/delete flows with generated client mutations and daily-schedule query invalidation while leaving `/` as an overview.
- [x] Regenerate and commit `frontend/openapi-snapshot.json` plus generated frontend client artifacts for the schedule write surface.
- [x] Update frontend component/SSR tests and backend unit/integration tests to cover read-state handling, schedule CRUD, validation, not-found paths, and contract preservation.
- [x] Run the project verification flow with `scripts/precommit-run`.
