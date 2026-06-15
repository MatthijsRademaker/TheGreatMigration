## Why

The backend currently serves three API endpoints — `/api/hello`, `/api/planning-window`, and `/api/dashboard/people-availability` — but all data is hardcoded as in-memory constants and seed slices. The archive design docs explicitly acknowledged this as a first-slice decision to defer persistence. With the endpoint contracts now stable and frontend consumers depending on them, the next step is to introduce a real persistence layer using the recommended `pgx` + `sqlc` + `goose` Postgres stack. This closes the persistence gap and gives the application durable data that can survive process restarts and evolve toward the product vision.

## What Changes

- **Go dependencies**: Add `pgx/v5`, `sqlc`, and `pressly/goose/v3` to `backend/go.mod`.
- **Database schema and migrations**: Create goose migrations under `backend/migrations/` that define `planning_windows`, `people`, and `availability` tables, plus a seed migration that inserts the canonical demo data (planning window `2026-07-05`..`2026-08-13`, 8 people with all four availability statuses).
- **sqlc query layer**: Add `backend/sqlc.yaml` and query SQL files under `backend/queries/`. Generate Go query code and commit it so that `sqlc` CLI is not required at Docker build time.
- **Store abstraction**: Introduce a `Store` interface in `backend/store.go` with methods `GetPlanningWindow()` and `GetPeopleAvailability(startDate, days)`. Implement a `PgStore` backed by a `pgxpool.Pool` and sqlc-generated queries.
- **Handler wiring**: Refactor `backend/main.go` to initialize a database pool, run embedded goose migrations at startup, create the `PgStore`, and inject it into the planning-window and dashboard handlers. The `/api/tasks/backlog` endpoint stays in-memory.
- **Compose orchestration**: Add a `db` service (`postgres:16-alpine`) to `compose.yml` with a healthcheck (`pg_isready`). Add `DATABASE_URL` environment variable to the backend service with `depends_on` the `db` service being healthy.
- **Tests**: Add a mock `Store` for fast handler unit tests. Add integration tests (gated by Go build tag `integration`) that use the existing `verification_start_postgres_sidecar` / `verification_cleanup_postgres_sidecar` helpers from `scripts/lib/docker-verification.sh` to test DB-backed handlers against a disposable Postgres instance.
- **OpenSpec specs**: Update the canonical specs for `dashboard-people-availability`, `planning-window`, `compose-orchestration`, and `backend-service` to reflect the persistence refactor. Add a new `backend-persistence` spec defining the persistence architecture.

## Impact

- **Affected specs**: `backend-service`, `dashboard-people-availability`, `planning-window`, `compose-orchestration`, plus new `backend-persistence`.
- **Affected code**: `backend/main.go`, `backend/dashboard.go`, `backend/planning_window.go`, `backend/main_test.go`, `backend/go.mod`, `backend/Dockerfile`, `compose.yml`, plus new files (`backend/store.go`, `backend/queries/`, `backend/migrations/`, `backend/sqlc.yaml`, `backend/db/` for generated code).
- **Frontend**: No changes required — the existing API contracts are preserved exactly.
- **Dev workflow**: `docker compose up` now starts Postgres alongside the backend and frontend. Goose migrations run automatically on backend startup. Developers can run `go test -short` for fast unit tests without Docker, or `scripts/test` / `go test -tags=integration` for full DB-backed integration coverage.
