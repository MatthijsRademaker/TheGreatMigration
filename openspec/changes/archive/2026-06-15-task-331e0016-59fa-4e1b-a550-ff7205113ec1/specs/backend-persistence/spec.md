# backend-persistence Specification (Delta)

## Purpose

Extend the existing Postgres persistence layer so the remaining read endpoints (`GET /api/tasks/backlog` and `GET /api/dashboard/daily-schedule`) use the same pgx/sqlc/goose/Store foundation already used by planning window and people availability.

## MODIFIED Requirements

### Requirement: Goose migrations SHALL create the persistence schema

The existing migrations for `planning_windows`, `people`, and `availability` SHALL remain intact. Additional goose migrations SHALL extend the schema with separate read-model tables for:

- backlog tasks, including stable task identifiers, titles, canonical priority values, required people counts, room names, canonical backlog statuses, and deterministic ordering;
- backlog task assignments, including foreign keys to both backlog tasks and `people(id)`;
- daily-schedule task cards, including stable seeded identifiers, titles, canonical priority values, room-area labels, required people counts, deterministic ordering, and seeded grouping metadata sufficient to reproduce the current four day-group schedule behavior; and
- daily-schedule task assignments, including foreign keys to both schedule task cards and `people(id)`.

The new migrations SHALL enforce the existing canonical vocabularies at the database level, seed the exact backlog IDs `task-1` through `task-11`, and seed the deterministic daily-schedule card data and assignee relationships needed to reproduce the current dashboard contract.

#### Scenario: Existing persistence tables are preserved while backlog and schedule tables are added
- **WHEN** goose migrations are applied to an empty Postgres database
- **THEN** the `planning_windows`, `people`, and `availability` tables still exist with their current columns
- **AND** separate persisted read models exist for backlog tasks and daily-schedule task cards with assignment tables referencing `people(id)`

#### Scenario: Seed data reproduces the current deterministic demo responses
- **WHEN** the seed migrations are applied
- **THEN** the backlog task read model contains the same 11 seeded tasks and assignment variety currently defined in `backend/tasks.go`
- **AND** the schedule read model contains the deterministic seeded task-card groups and assignments needed to reproduce the current default four-day schedule behavior

### Requirement: Backend SHALL define a Store interface for data access

The `Store` interface in `backend/store.go` SHALL include:
- `GetPlanningWindow(ctx context.Context) (*PlanningWindowBody, error)`
- `GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (*DashboardBody, error)`
- `GetTaskBacklog(ctx context.Context) (*TaskBacklogBody, error)`
- `GetDailySchedule(ctx context.Context, startDate time.Time, days int) (*DailyScheduleBody, error)`

`PgStore` SHALL satisfy all four methods using sqlc-generated queries, and `MockStore` SHALL satisfy all four methods for unit-test coverage.

#### Scenario: PgStore satisfies the expanded Store interface
- **WHEN** `PgStore` is compiled
- **THEN** it implements all methods of the expanded `Store` interface without compilation errors

#### Scenario: MockStore satisfies the expanded Store interface
- **WHEN** `MockStore` is compiled
- **THEN** it implements all methods of the expanded `Store` interface without compilation errors

### Requirement: Backend handlers SHALL read from the Store interface

The `registerDashboardPeopleAvailability` and `registerPlanningWindow` handlers SHALL remain Store-backed. `registerTasksBacklog` and `registerDailySchedule` SHALL also accept a `Store` parameter and SHALL stop reading directly from `seedTasks`, `seedTasksForDay`, `seedPeople`, `findPersonByID`, `countAvailableForDay`, or `planWindowStart`.

For daily schedule, when the `start` query parameter is omitted, the handler SHALL call `store.GetPlanningWindow()` and use the returned `startDate` as the default before calling `store.GetDailySchedule()`.

#### Scenario: Task backlog handler reads from Store
- **WHEN** `GET /api/tasks/backlog` is called and the Store returns backlog data
- **THEN** the response body matches the Store-backed backlog payload

#### Scenario: Daily schedule handler resolves its default start through Store
- **WHEN** `GET /api/dashboard/daily-schedule` is called without a `start` query parameter
- **THEN** the handler obtains the default `startDate` from `store.GetPlanningWindow()` before loading schedule data from `store.GetDailySchedule()`

### Requirement: Existing endpoints SHALL remain intact

The persistence refactor SHALL NOT change the behavior of `GET /api/hello`, CORS handling, or `/openapi.json` generation. `GET /api/planning-window` and `GET /api/dashboard/people-availability` SHALL remain Postgres-backed. `GET /api/tasks/backlog` and `GET /api/dashboard/daily-schedule` SHALL move from in-memory seed data to Postgres-backed storage while preserving their current response contracts, legends, and deterministic seeded output.

#### Scenario: Public endpoint contracts are preserved across the persistence refactor
- **WHEN** the backend is running after the refactor
- **THEN** `GET /api/planning-window`, `GET /api/dashboard/people-availability`, `GET /api/tasks/backlog`, and `GET /api/dashboard/daily-schedule` each return their existing response shapes and canonical vocabularies

#### Scenario: OpenAPI includes all backend paths
- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document contains `/api/hello`, `/api/planning-window`, `/api/dashboard/people-availability`, `/api/tasks/backlog`, and `/api/dashboard/daily-schedule`

### Requirement: Backend tests SHALL cover DB-backed handlers

Backend unit tests SHALL use `MockStore` to test Store-backed handler logic without requiring Postgres. Backend integration tests (build tag `integration`) SHALL run against a migrated Postgres database and cover planning window, people availability, task backlog, daily schedule, and OpenAPI path registration.

#### Scenario: Unit tests cover Store-backed handler success and failure paths
- **WHEN** `go test -short ./...` runs in `backend/`
- **THEN** the handler tests cover Store-backed success and error paths for task backlog and daily schedule without requiring Postgres

#### Scenario: Integration tests cover all Postgres-backed read endpoints
- **WHEN** `go test -tags=integration ./...` runs in `backend/` with Docker available
- **THEN** the tests validate planning-window, people-availability, task-backlog, and daily-schedule contracts against a migrated database
- **AND** they assert `/openapi.json` includes all five backend paths

#### Scenario: Precommit checks pass
- **WHEN** `scripts/precommit-run` is executed
- **THEN** all lint, build, and test checks pass without errors
