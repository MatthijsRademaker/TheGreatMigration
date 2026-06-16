## Context

The frontend already has a generated query for `GET /api/dashboard/daily-schedule`, but `DailySchedule.vue`, `HomeView.vue`, and `CalendarView.vue` still render demo data. The backend already serves the read model, yet schedule persistence is incomplete: `schedule_task_cards` use demo-oriented `day_group` mapping, the read response emits synthetic IDs, and no schedule-card create/update/delete API exists.

## Goals

- Wire the Daily Schedule board on `/` and `/calendar` to the existing BFF read model.
- Preserve `DailySchedule.vue` as a typed, prop-driven presentation component.
- Add schedule-specific CRUD backed by `schedule_task_cards` and `schedule_task_assignments`.
- Make successful writes visible through subsequent `GET /api/dashboard/daily-schedule` reads and frontend query refresh.
- Keep people-availability-derived helper counts and assignee identity resolution intact.

## Non-Goals

- Authentication, authorization, invitations, or collaboration features.
- Replacing backlog-task CRUD or people-availability management beyond required invalidation and compatibility work.
- Drag-and-drop scheduling, automatic staffing optimization, notifications, or calendar export.
- Reusing backlog-task CRUD for the separate schedule-card model.

## Decisions

1. Deliver the change in two layers within one proposal: first wire the existing read path, then add the schedule write surface after its schema and spec prerequisites are in place.
2. Introduce a shared `useDailySchedule` composable that follows the established `usePeopleAvailability` and `useTaskBacklog` adapter pattern, including null-safe array handling and canonical enum narrowing.
3. Keep `DailySchedule.vue` presentational. The home dashboard remains read-only, while `/calendar` owns create, update, and delete interactions.
4. Implement schedule CRUD on `schedule_task_cards` and `schedule_task_assignments`, not `backlog_tasks`.
5. Replace demo-only `day_group` addressing with persisted `scheduled_date` storage and expose stable `sched-<id>` card identifiers in the read model.
6. Regenerate and commit the OpenAPI snapshot and generated frontend client artifacts when the write contract is added.

## Conflict Resolution

Refinement disagreed on whether CRUD should be deferred. This proposal keeps CRUD in scope because the task explicitly asks to fix forward missing CRUD endpoints and the validated room output converged on the required preconditions: a schedule-specific write model, a `scheduled_date` migration, stable persisted IDs, and an explicit spec delta. To preserve the existing component boundary, persistence is owned by the `/calendar` route rather than by `DailySchedule.vue` itself.

## Risks

- The `scheduled_date` migration must preserve the seeded four-day schedule variety and ordering expected by current daily-schedule tests.
- Frontend tests that assert demo schedule fixtures must be rewritten to use mocked BFF-shaped data and explicit state coverage.
- Schedule mutations must invalidate/refetch daily-schedule queries reliably so route state does not go stale across affected date windows.

## Traceability

- Task source: `task:48c52b0a-4b52-4413-83fa-0f7fba30036f`
- Dossier: `dossier:2026-06-16T13:53:59.281Z`
- Accepted decisions: `decision:1-swarm-architect-recommendation`, `decision:1-swarm-lead-dev-recommendation`, `decision:1-swarm-reviewer-recommendation`
- Validated room evidence: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
