# backend-persistence Specification

## Purpose
TBD - created by archiving change task-248d3e0f-c580-48ad-9b19-4c6785beead2. Update Purpose after archive.
## Requirements
### Requirement: Backend SHALL depend on pgx, sqlc, and goose for database access

The backend Go module SHALL add dependencies on `github.com/jackc/pgx/v5` (for connection pooling), `github.com/sqlc-dev/sqlc` (for type-safe query generation tooling), and `github.com/pressly/goose/v3` (for schema migrations). The `sqlc` CLI SHALL NOT be required at Docker build time; generated Go code SHALL be committed to the repository.

#### Scenario: Go module includes persistence dependencies
- **WHEN** `go mod tidy` runs in `backend/`
- **THEN** the module graph includes `pgx/v5`, `sqlc`, and `goose/v3` as direct or indirect dependencies

#### Scenario: Generated sqlc code is committed
- **WHEN** `go build ./...` runs in `backend/` without `sqlc` installed
- **THEN** the build succeeds because all sqlc-generated `.go` files are present in the repository

### Requirement: Goose migrations SHALL create the persistence schema

The existing migrations for `planning_windows`, `people`, and `availability` SHALL remain intact. Additional goose migrations SHALL extend the schema with separate read-model tables for:

- backlog tasks, including stable task identifiers, titles, canonical priority values, required people counts, room names, canonical backlog statuses, and deterministic ordering;
- backlog task assignments, including foreign keys to both backlog tasks and `people(id)`;
- daily-schedule task cards, including stable seeded identifiers, titles, canonical priority values, room-area labels, required people counts, a scheduled date, and deterministic ordering; and
- daily-schedule task assignments, including foreign keys to both schedule task cards and `people(id)`.

The schema migrations SHALL enforce the existing canonical vocabularies at the database level. The schema migration set SHALL contain **structure only** — it SHALL NOT insert demo rows. Demo seed data is provided by a separate, flag-gated dataset (see "Demo seed data SHALL be a separate, flag-gated goose dataset"). Schema migration version numbers SHALL NOT be renumbered when the seed migrations are removed; goose tolerates the resulting version gaps on a fresh database.

#### Scenario: Existing persistence tables are preserved while backlog and schedule tables are added
- **WHEN** the schema goose migrations are applied to an empty Postgres database
- **THEN** the `planning_windows`, `people`, and `availability` tables still exist with their current columns
- **AND** separate persisted read models exist for backlog tasks and daily-schedule task cards with assignment tables referencing `people(id)`

#### Scenario: Schema-only migration leaves domain tables empty
- **WHEN** the schema goose migrations are applied to an empty Postgres database and the seed dataset is NOT applied
- **THEN** the `people`, `backlog_tasks`, `rooms_areas`, and `schedule_task_cards` tables contain zero rows
- **AND** the backend starts and serves HTTP

### Requirement: Demo seed data SHALL be a separate, flag-gated goose dataset

Demo seed data SHALL live in `backend/seed/` as its own goose dataset, embedded via a dedicated `//go:embed` directive and tracked in a goose version table distinct from the schema table (`goose_seed_version`). The seed dataset SHALL seed the exact backlog IDs `task-1` through `task-11`, the eight demo people (`p1`–`p8`) and their availability, the eight rooms/areas (`room-1`–`room-8`), and the deterministic daily-schedule cards and assignments needed to reproduce the current dashboard contract. Because the seed dataset is applied after the full schema, schedule-card seed rows SHALL set the `scheduled_date` column directly (the `day_group` column no longer exists). After inserting, the seed dataset SHALL advance the `people_id_seq`, `rooms_areas_id_seq`, and backlog-task ID sequences past the seeded maximum so subsequently created entities receive non-colliding sequential IDs.

#### Scenario: Seed dataset reproduces the deterministic demo responses
- **WHEN** the schema migrations are applied and then the seed dataset is applied
- **THEN** the backlog task read model contains the same 11 seeded tasks and assignment variety
- **AND** the schedule read model contains the deterministic seeded task-card data and assignments needed to reproduce the current default schedule behavior
- **AND** eight demo people with their availability and eight rooms/areas exist

#### Scenario: Seed dataset advances ID sequences past seeded rows
- **WHEN** the seed dataset has been applied and a new person, room, and backlog task are created
- **THEN** they receive the next sequential identifiers (`p9`, `room-9`, `task-12`) with no collision against seeded IDs

#### Scenario: Seed dataset is tracked independently of schema
- **WHEN** the backend starts twice against a database where schema and seed have both been applied
- **THEN** neither the schema pass nor the seed pass re-applies any migration, and startup succeeds both times

### Requirement: Backend SHALL embed and run migrations at startup

The backend SHALL use Go's `//go:embed` directive to bundle migration `.sql` files into the binary. The backend SHALL call `goose.Up()` with the embedded filesystem before starting the HTTP server. Migration execution SHALL be gated by the `DB_AUTO_MIGRATE` environment variable, defaulting to `"true"`. If `DB_AUTO_MIGRATE` is `"false"`, the backend SHALL skip migrations and proceed directly to serving.

#### Scenario: Migrations run at startup by default
- **WHEN** the backend starts with `DB_AUTO_MIGRATE` unset and `DATABASE_URL` pointing to an empty Postgres database
- **THEN** the backend applies all pending goose migrations before the HTTP server begins accepting connections

#### Scenario: Migrations can be skipped via env var
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=false` and `DATABASE_URL` pointing to an already-migrated Postgres database
- **THEN** the backend skips goose migration execution and starts serving HTTP immediately

### Requirement: Backend SHALL connect to Postgres via a single DATABASE_URL env var

The backend SHALL read the `DATABASE_URL` environment variable as the sole database connectivity configuration. The value SHALL be a standard Postgres connection URI (`postgres://user:pass@host:port/dbname?...`). The backend SHALL use `pgxpool.ParseConfig` to create a connection pool. If `DATABASE_URL` is empty or the pool cannot be created after a configurable retry window, the backend SHALL exit with a clear error message.

#### Scenario: Backend connects with a valid DATABASE_URL
- **WHEN** `DATABASE_URL` is set to a reachable Postgres instance
- **THEN** the backend creates a `pgxpool.Pool` and proceeds through startup

#### Scenario: Backend exits on unreachable database
- **WHEN** `DATABASE_URL` points to a non-existent or unreachable host
- **THEN** the backend retries for a bounded period and exits with a diagnostic message if the pool cannot be established

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

### Requirement: sqlc SHALL generate type-safe query code from SQL files

A `backend/sqlc.yaml` configuration SHALL define the Postgres engine, schema file paths, query file paths, and Go output directory (`backend/db/`). Query `.sql` files under `backend/queries/` SHALL contain named SQL queries (e.g., `-- name: GetPlanningWindow :one`) that sqlc uses to generate Go functions. The generated Go code SHALL be committed to the repository.

#### Scenario: sqlc configuration is valid
- **WHEN** `sqlc generate` runs with the committed `sqlc.yaml`
- **THEN** it produces Go files under `backend/db/` without errors

#### Scenario: Generated code compiles
- **WHEN** `go build ./...` runs in `backend/`
- **THEN** the sqlc-generated code compiles without errors

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
