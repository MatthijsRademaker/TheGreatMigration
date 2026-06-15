## Context

The repository already has Postgres wiring for `GET /api/planning-window` and `GET /api/dashboard/people-availability` through `PgStore`, but `GET /api/tasks/backlog` still derives its response from `seedTasks` in `backend/tasks.go` and `GET /api/dashboard/daily-schedule` still derives its response from `seedTasksForDay`, `seedPeople`, and the hardcoded `planWindowStart` constant in `backend/daily_schedule.go`. Existing canonical specs for backend persistence, task backlog, and daily schedule still encode the earlier in-memory design and now need to be superseded by a persistence-backed version of the same contracts.

## Goals / Non-Goals

### Goals
- Move `GET /api/tasks/backlog` onto Postgres-backed storage without changing its response shape, summary semantics, legends, or canonical vocabularies.
- Move `GET /api/dashboard/daily-schedule` onto Postgres-backed storage without changing its response shape, range semantics, task-card invariants, or assignee identity shape.
- Preserve the existing Postgres-backed implementation and response contracts for `GET /api/planning-window` and `GET /api/dashboard/people-availability`.
- Extend migrations, sqlc queries, generated code, `Store`, `PgStore`, and `MockStore` to support all four read endpoints.
- Seed deterministic demo data that reproduces the current backlog and daily-schedule outputs, including exact current seed IDs, assignment variety, and default four-day schedule variety.
- Expand unit and integration coverage so all four endpoints and the generated OpenAPI paths are validated against the Postgres-backed implementation.

### Non-Goals
- No create, edit, delete, auth, tenancy, or collaboration workflows.
- No frontend redesign, route rewiring, or broader domain remodel beyond the read models needed for the existing contracts.
- No replacement of the existing pgx/sqlc/goose persistence stack.
- No change to the public JSON contracts for planning window, people availability, task backlog, or daily schedule.

## Decisions

### D1: Keep separate read-model persistence for backlog tasks and daily-schedule cards

Use separate persistence tables for backlog tasks and daily-schedule task cards rather than trying to normalize both responses into one shared task table. The two endpoints expose different field sets and vocabularies (`room` vs. `roomArea`, `status` vs. `staffingStatus`, `assignedTo` vs. `assignedPeople`), so separate read-model families preserve the current contracts without nullable cross-wiring.

### D2: Preserve the current deterministic demo output and exact seed IDs

Seed the database from the current in-memory task backlog and daily-schedule templates so first-run responses remain deterministic. Backlog tasks keep the exact IDs `task-1` through `task-11`. Daily-schedule persistence must preserve the current deterministic task-card IDs for identical requests and reproduce the existing four day-group schedule variety used by the default four-day window.

### D3: Default daily-schedule start from the Store-backed planning window

When `GET /api/dashboard/daily-schedule` is called without `start`, the handler must load the planning window from `Store.GetPlanningWindow()` and use that `startDate` instead of the hardcoded `planWindowStart` constant. This keeps the schedule default aligned with the canonical planning-window data already stored in Postgres.

### D4: Extend the existing Store/PgStore/sqlc pattern instead of adding a new data path

Add `GetTaskBacklog(ctx)` and `GetDailySchedule(ctx, startDate, days)` to `Store`, implement them in `PgStore`, and update `MockStore` accordingly. Refactor `registerTasksBacklog` and `registerDailySchedule` to accept `Store` and delegate all data loading and response shaping through Store methods, while retaining legend constants in Go code.

### D5: Derive `availablePeopleCount` from availability rows, not duplicated schedule data

Daily-schedule `availablePeopleCount` must be derived from the existing availability data for each returned date. The Postgres-backed implementation should use a query that counts availability rows with status `available` for each day rather than relying on the old in-memory helper or adding a duplicated count column.

### D6: Validate through both mock-based unit tests and real-Postgres integration tests

Unit tests should cover success and Store-failure paths for backlog and daily-schedule handlers using `MockStore`, while integration tests should run all four Postgres-backed endpoint contracts against a migrated database and assert `/openapi.json` contains all five backend paths.

## Conflict Resolution

The accepted refinement evidence resolves two conflicts:

1. **In-memory spec requirements vs. the new task scope**: existing canonical specs for `backend-persistence`, `task-backlog-api`, and `dashboard-daily-schedule` still require in-memory data. This change explicitly supersedes those requirements so task backlog and daily schedule become Postgres-backed while preserving their public contracts.
2. **Shared-table vs. separate-table modeling**: accepted architectural guidance chose separate read-model tables for backlog tasks and schedule task cards. The specification follows that decision and does not normalize the two response models into one table family.

## Risks

| Risk | Mitigation |
| --- | --- |
| Seed SQL drifts from the current deterministic Go seed behavior. | Require migrations to reproduce the current backlog rows, schedule day-group variety, assignments, summary invariants, and exact seed IDs. |
| Daily-schedule defaults drift from planning-window data. | Require omitted `start` to resolve through `Store.GetPlanningWindow()` instead of a duplicated constant. |
| Refactor leaves hidden in-memory reads in handlers or tests. | Require `registerTasksBacklog` and `registerDailySchedule` to accept `Store` and add unit tests for Store success/failure plus real-Postgres integration coverage. |
| OpenSpec traceability breaks because canonical specs still mention in-memory storage. | Include spec deltas that supersede those in-memory requirements as part of this change. |

## Traceability

- Task: `331e0016-59fa-4e1b-a550-ff7205113ec1`
- Dossier: `2026-06-15T20:47:25.531Z`
- Accepted decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Validated round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `task-331e0016-59fa-4e1b-a550-ff7205113ec1 / initial`
