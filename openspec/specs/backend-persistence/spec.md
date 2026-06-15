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

A set of goose migration files under `backend/migrations/` SHALL define the database schema. Migration `001_create_tables.sql` SHALL create three tables:
- `planning_windows`: `id` (SERIAL PRIMARY KEY), `start_date` (DATE NOT NULL), `end_date` (DATE NOT NULL), `created_at` (TIMESTAMPTZ DEFAULT NOW()), `updated_at` (TIMESTAMPTZ DEFAULT NOW()).
- `people`: `id` (TEXT PRIMARY KEY), `name` (TEXT NOT NULL), `initials` (TEXT NOT NULL), `created_at` (TIMESTAMPTZ DEFAULT NOW()).
- `availability`: `id` (SERIAL PRIMARY KEY), `person_id` (TEXT NOT NULL REFERENCES people(id)), `date` (DATE NOT NULL), `status` (TEXT NOT NULL CHECK (status IN ('available','busy','partial','off'))), `created_at` (TIMESTAMPTZ DEFAULT NOW()), UNIQUE(person_id, date).

Migration `002_seed_demo_data.sql` SHALL insert:
- One planning window row with `start_date='2026-07-05'` and `end_date='2026-08-13'`.
- Eight people with IDs `p1` through `p8`, matching the names and initials from the current in-memory `seedPeople`.
- Availability rows for every person on every date from `2026-07-05` through `2026-08-13` (40 dates × 8 people) with status values that exactly reproduce the current closure-based logic (p1–p6 always `available`, p7 always `busy`, p8 cycling through `off`/`partial`/`busy`/`available`).

#### Scenario: Migration files are present and valid
- **WHEN** goose migrations are applied to an empty Postgres database
- **THEN** the `planning_windows`, `people`, and `availability` tables exist with the specified columns

#### Scenario: Seed data matches the in-memory defaults
- **WHEN** the seed migration is applied
- **THEN** `SELECT start_date, end_date FROM planning_windows` returns exactly one row with `'2026-07-05'` and `'2026-08-13'`
- **AND** `SELECT COUNT(*) FROM people` returns 8
- **AND** `SELECT COUNT(*) FROM availability` returns 320

#### Scenario: Availability statuses match the closure logic
- **WHEN** the seed data is queried for person `p7`
- **THEN** every availability row for `p7` has status `'busy'`
- **WHEN** the seed data is queried for person `p8`
- **THEN** the status cycles through `'off'`, `'partial'`, `'busy'`, `'available'` in date order

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

A `Store` interface in `backend/store.go` SHALL define the data access contract:
- `GetPlanningWindow(ctx context.Context) (PlanningWindowData, error)` — returns the singleton planning window row (`startDate`, `endDate`, `days`).
- `GetPeopleAvailability(ctx context.Context, startDate time.Time, days int) (DashboardData, error)` — returns people, their availability within the date window, and the status legend.

An implementation (`PgStore`) SHALL satisfy this interface using a `*pgxpool.Pool` and sqlc-generated query functions. A mock implementation (`MockStore`) SHALL also satisfy this interface for use in handler unit tests.

#### Scenario: PgStore satisfies the Store interface
- **WHEN** `PgStore` is compiled
- **THEN** it implements all methods of the `Store` interface without compilation errors

#### Scenario: MockStore satisfies the Store interface
- **WHEN** `MockStore` is compiled
- **THEN** it implements all methods of the `Store` interface without compilation errors

### Requirement: sqlc SHALL generate type-safe query code from SQL files

A `backend/sqlc.yaml` configuration SHALL define the Postgres engine, schema file paths, query file paths, and Go output directory (`backend/db/`). Query `.sql` files under `backend/queries/` SHALL contain named SQL queries (e.g., `-- name: GetPlanningWindow :one`) that sqlc uses to generate Go functions. The generated Go code SHALL be committed to the repository.

#### Scenario: sqlc configuration is valid
- **WHEN** `sqlc generate` runs with the committed `sqlc.yaml`
- **THEN** it produces Go files under `backend/db/` without errors

#### Scenario: Generated code compiles
- **WHEN** `go build ./...` runs in `backend/`
- **THEN** the sqlc-generated code compiles without errors

### Requirement: Backend handlers SHALL read from the Store interface

The `registerDashboardPeopleAvailability` and `registerPlanningWindow` handler registration functions SHALL accept a `Store` parameter. The handler functions SHALL call `store.GetPeopleAvailability()` and `store.GetPlanningWindow()` respectively instead of reading from in-memory constants or seed slices. The response types (`DashboardBody`, `PlanningWindowBody`) and their JSON serialization SHALL remain unchanged.

#### Scenario: Planning window handler reads from Store
- **WHEN** `GET /api/planning-window` is called and the Store returns planning window data
- **THEN** the response contains `startDate`, `endDate`, and `days` matching the Store's data

#### Scenario: Dashboard handler reads from Store
- **WHEN** `GET /api/dashboard/people-availability` is called and the Store returns people and availability data
- **THEN** the response contains `range`, `summary`, `people`, and `statuses` matching the Store's data

### Requirement: Existing endpoints SHALL remain intact

The persistence refactor SHALL NOT change the behavior of `GET /api/hello`, `GET /api/tasks/backlog`, CORS handling, or `/openapi.json` generation. The tasks backlog endpoint SHALL continue to use in-memory seed data as it is explicitly out of scope for this change.

#### Scenario: Hello endpoint is unchanged after persistence refactor
- **WHEN** `GET /api/hello` is called after the refactor
- **THEN** the response is 200 with body `{"message": "Hello from the backend!"}`

#### Scenario: Tasks backlog endpoint remains in-memory
- **WHEN** `GET /api/tasks/backlog` is called after the refactor
- **THEN** the response contains at least 10 tasks, all three priority values, all three task statuses, and at least 4 distinct rooms

#### Scenario: CORS still allows frontend origins
- **WHEN** a cross-origin `OPTIONS` preflight request is sent with `Origin: http://localhost:5173`
- **THEN** the response includes `Access-Control-Allow-Origin: http://localhost:5173`

#### Scenario: OpenAPI includes all endpoints
- **WHEN** `GET /openapi.json` is called
- **THEN** the OpenAPI document contains paths `/api/hello`, `/api/dashboard/people-availability`, `/api/planning-window`, and `/api/tasks/backlog`

### Requirement: Backend tests SHALL cover DB-backed handlers

Backend unit tests SHALL use the `MockStore` to test handler logic without requiring a Postgres instance. Backend integration tests (build tag `integration`) SHALL start a disposable Postgres sidecar using `verification_start_postgres_sidecar()` from `scripts/lib/docker-verification.sh`, apply goose migrations, seed data, and run the same endpoint contract tests against a real database. Integration tests SHALL call `verification_cleanup_postgres_sidecar()` in `t.Cleanup` to ensure resources are released.

#### Scenario: Unit tests pass with mock store
- **WHEN** `go test -short ./...` runs in `backend/`
- **THEN** all existing endpoint tests pass using the mock Store

#### Scenario: Integration tests pass with real Postgres
- **WHEN** `go test -tags=integration ./...` runs in `backend/` with Docker available
- **THEN** the DB-backed endpoint tests start a Postgres sidecar, apply migrations, and assert the planning-window and dashboard contracts hold

#### Scenario: Precommit checks pass
- **WHEN** `scripts/precommit-run` is executed
- **THEN** all lint, vet, build, and test checks pass without errors
