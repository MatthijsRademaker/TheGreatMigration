## 1. Dependency Setup

- [x] 1.1 Add `pgx/v5`, `sqlc`, and `pressly/goose/v3` dependencies to `backend/go.mod` and run `go mod tidy`.

## 2. Database Schema and Migrations

- [x] 2.1 Create `backend/migrations/001_create_tables.sql` with DDL for `planning_windows`, `people`, and `availability` tables including CHECK constraints on status values.
- [x] 2.2 Create `backend/migrations/002_seed_demo_data.sql` with INSERT statements that reproduce the exact in-memory seed data (1 planning window row, 8 people, 320 availability rows for 40 dates).

## 3. sqlc Query Layer

- [x] 3.1 Create `backend/sqlc.yaml` configuring the Postgres engine, schema path, queries path, and Go output directory (`backend/db/`).
- [x] 3.2 Create `backend/queries/planning_window.sql` with a SELECT query for the planning window row.
- [x] 3.3 Create `backend/queries/people_availability.sql` with SELECT queries for people and their availability rows within a date range.
- [x] 3.4 Run `sqlc generate` and commit the generated Go code under `backend/db/`.

## 4. Store Layer

- [x] 4.1 Create `backend/store.go` with the `Store` interface (`GetPlanningWindow`, `GetPeopleAvailability`).
- [x] 4.2 Implement `PgStore` struct with `NewPgStore(ctx, pool)` constructor that wraps sqlc queries.
- [x] 4.3 Implement `GetPlanningWindow` method: query `planning_windows` table, return `PlanningWindowBody` data.
- [x] 4.4 Implement `GetPeopleAvailability` method: query `people` + `availability` tables, build `DashboardBody` with correct range metadata, summary counts, and status legend.

## 5. Backend Startup Wiring

- [x] 5.1 Add `//go:embed migrations/*.sql` directive and goose migration runner to `backend/main.go`.
- [x] 5.2 Add pgxpool initialization with retry logic reading `DATABASE_URL` from env.
- [x] 5.3 Run goose migrations at startup (gated by `DB_AUTO_MIGRATE`, default `true`) before starting the HTTP server.
- [x] 5.4 Refactor handler registration: pass `Store` to `registerPlanningWindow` and `registerDashboardPeopleAvailability`.
- [x] 5.5 Update `backend/dashboard.go` and `backend/planning_window.go` handlers to read from `Store` instead of in-memory constants/slices.
- [x] 5.6 Preserve `/api/hello`, `/api/tasks/backlog`, CORS, and `/openapi.json` behavior unchanged.

## 6. Compose Orchestration

- [x] 6.1 Add `db` service to `compose.yml` using `postgres:16-alpine` with healthcheck, env vars, and volume for data persistence.
- [x] 6.2 Add `DATABASE_URL` env var and `depends_on: db: service_healthy` to the backend service.
- [x] 6.3 Add `DB_AUTO_MIGRATE=true` env var to the backend service.

## 7. Tests

- [x] 7.1 Create `backend/store_mock_test.go` with a mock `Store` implementation for fast unit tests.
- [x] 7.2 Refactor `newTestAPI()` to accept a `Store` parameter (defaults to mock store).
- [x] 7.3 Update existing tests (`TestHelloEndpoint`, `TestDashboardPeopleAvailability`, `TestPlanningWindowEndpoint`) to use the mock store.
- [x] 7.4 Create `backend/main_integration_test.go` (build tag `integration`) with a `TestDBBackedEndpoints` function that starts a Postgres sidecar via `verification_start_postgres_sidecar`, applies migrations, seeds data, and asserts the endpoint contracts.
- [x] 7.5 Ensure `verification_cleanup_postgres_sidecar` is called in `t.Cleanup`.

## 8. OpenSpec Spec Updates

- [x] 8.1 Create `openspec/changes/task-248d3e0f-c580-48ad-9b19-4c6785beead2/specs/backend-persistence/spec.md` defining the persistence architecture.
- [x] 8.2 Create MODIFIED spec delta for `dashboard-people-availability` replacing in-memory requirement with DB-backed query requirement.
- [x] 8.3 Create MODIFIED spec delta for `planning-window` replacing compile-time constant requirement with DB-backed query requirement.
- [x] 8.4 Create MODIFIED spec delta for `compose-orchestration` adding Postgres service and health dependency.
- [x] 8.5 Create MODIFIED spec delta for `backend-service` adding database initialization at startup.

## 9. Verification

- [x] 9.1 Run `scripts/precommit-run` and confirm all lint, vet, build, and test checks pass.
- [x] 9.2 Run `docker compose up` and verify the full stack starts (db → backend → frontend) and all endpoints return correct data.
