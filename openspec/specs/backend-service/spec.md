# backend-service Specification

## Purpose
TBD - created by archiving change task-30e902a9-a751-4be8-b4ea-149a16a8740d. Update Purpose after archive.
## Requirements
### Requirement: A Go backend service SHALL serve a hello-world endpoint with auto-generated OpenAPI
The change SHALL create a `backend/` Go module using Huma v2 with the chi adapter. The service SHALL expose `GET /api/hello` returning `{"message": "Hello from the backend!"}` and SHALL serve auto-generated OpenAPI 3.1 specification at `/openapi.json`. The service SHALL include a Dockerfile for compose orchestration.

#### Scenario: Hello endpoint returns valid JSON
- **WHEN** an HTTP client sends `GET /api/hello` to the running backend
- **THEN** the response status is 200 and the body is `{"message": "Hello from the backend!"}` with `Content-Type: application/json`

#### Scenario: OpenAPI specification is auto-generated
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` returns a valid OpenAPI 3.1 specification document that includes the `/api/hello` endpoint

#### Scenario: CORS allows frontend origins
- **WHEN** the frontend dev server at `http://localhost:5173` or the compose frontend service sends a cross-origin request to the backend
- **THEN** the response includes `Access-Control-Allow-Origin` matching the request origin and CORS preflight (`OPTIONS`) requests succeed

#### Scenario: Dockerfile builds and runs the backend
- **WHEN** `docker build backend/` is executed using `golang:1.26.2-alpine`
- **THEN** the resulting image starts the Go service and responds to `GET /api/hello` on port 8080

#### Scenario: Go module structure is conventional
- **WHEN** `go vet ./...` and `go build ./...` are executed in `backend/`
- **THEN** both commands exit zero without warnings or errors

### Requirement: Backend SHALL initialize a Postgres connection pool at startup

The backend SHALL read `DATABASE_URL` from the environment and create a `pgxpool.Pool` before registering HTTP handlers. If the pool cannot be created after a bounded retry window (default 30 seconds), the backend SHALL log a diagnostic message and exit with a non-zero status code. The pool SHALL be passed to a `PgStore` constructor and injected into handlers.

#### Scenario: Backend starts successfully with valid DATABASE_URL
- **WHEN** `DATABASE_URL` points to a reachable Postgres instance and the backend is started
- **THEN** the backend creates a connection pool and starts serving HTTP on port 8080

#### Scenario: Backend exits cleanly with unreachable database
- **WHEN** `DATABASE_URL` points to an unreachable host
- **THEN** the backend prints a diagnostic message and exits with a non-zero status

### Requirement: Backend SHALL run goose migrations at startup

The backend SHALL embed migration SQL files via `//go:embed` and run `goose.Up()` after establishing the connection pool but before registering HTTP handlers. Migration execution SHALL be gated by `DB_AUTO_MIGRATE` (default `true`). If migrations fail, the backend SHALL log the error and exit.

#### Scenario: Migrations apply successfully on first run
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=true` against an empty database
- **THEN** all migration files are applied and the backend proceeds to serve HTTP

#### Scenario: Migrations are skipped when disabled
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=false`
- **THEN** goose.Up() is not called and the backend proceeds to serve HTTP

#### Scenario: Already-applied migrations are a no-op
- **WHEN** the backend starts against a database where all migrations are already applied
- **THEN** goose.Up() completes without errors and the backend proceeds normally

### Requirement: Backend SHALL inject Store into handlers

The `main()` function SHALL construct a `Store` implementation (either `PgStore` connected to the pool or the mock for testing) and pass it to `registerPlanningWindow` and `registerDashboardPeopleAvailability`. The `registerTasksBacklog` handler registration SHALL NOT require a `Store` parameter as it remains in-memory.

#### Scenario: Handlers receive Store at registration
- **WHEN** the backend starts normally
- **THEN** the planning-window and dashboard handlers are registered with a non-nil Store

#### Scenario: Tasks backlog handler does not require Store
- **WHEN** the backend starts normally
- **THEN** `registerTasksBacklog` is called without a Store parameter

### Requirement: Existing backend behavior SHALL remain intact

`GET /api/hello` SHALL continue to return `{"message": "Hello from the backend!"}`. CORS SHALL continue to allow `http://localhost:5173` and `http://frontend:5173`. The Dockerfile SHALL continue to build from `golang:1.26.2-alpine` and produce a working image. `go vet ./...` and `go build ./...` SHALL continue to pass.

#### Scenario: Hello endpoint returns valid JSON
- **WHEN** an HTTP client sends `GET /api/hello` to the running backend
- **THEN** the response status is 200 and the body is `{"message": "Hello from the backend!"}` with `Content-Type: application/json`

#### Scenario: OpenAPI specification is auto-generated
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` returns a valid OpenAPI 3.1 specification document

#### Scenario: Dockerfile builds and runs the backend
- **WHEN** `docker build backend/` is executed using `golang:1.26.2-alpine`
- **THEN** the resulting image starts the Go service and responds to `GET /api/hello` on port 8080

### Requirement: Backend SHALL run the demo seed dataset only when DB_SEED is enabled

After the schema `goose.Up()` completes, the backend SHALL conditionally apply the demo seed
dataset embedded from `backend/seed/`. Seed application SHALL be gated by the `DB_SEED`
environment variable, defaulting to `false` (opt-in). When `DB_SEED` is truthy and auto-migrate
is enabled, the backend SHALL run a second `goose.Up()` over the seed dataset using a distinct
version table (`goose_seed_version`) and SHALL restore the schema version table name afterward.
If the seed pass fails, the backend SHALL log the error and exit. The seed pass SHALL NOT run when
`DB_AUTO_MIGRATE` is disabled.

#### Scenario: Seed runs when enabled
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=true` and `DB_SEED=true` against an empty database
- **THEN** the schema migrations apply, the seed dataset applies under `goose_seed_version`, and the backend serves HTTP with seeded demo data

#### Scenario: Seed is skipped by default
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=true` and `DB_SEED` unset against an empty database
- **THEN** the schema migrations apply, no seed dataset is applied, and the backend serves HTTP with empty domain tables

#### Scenario: Seed does not run when migrations are disabled
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=false` and `DB_SEED=true`
- **THEN** neither the schema migrations nor the seed dataset are applied

#### Scenario: Seed pass is idempotent across restarts
- **WHEN** the backend restarts with `DB_SEED=true` against a database where the seed dataset is already applied
- **THEN** the seed `goose.Up()` completes without re-inserting rows and the backend proceeds normally
