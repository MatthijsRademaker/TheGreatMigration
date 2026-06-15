## Context

The backend exposes `/api/planning-window` and `/api/dashboard/people-availability` — both read endpoints whose data currently lives in Go-level compile-time constants and in-memory slices. The archive design docs for these endpoints explicitly chose in-memory seed data because no persistence model existed. Now that the API contracts are stable (frontend consumers are built, `openspec/specs/` exist, and `frontend/openapi-snapshot.json` is committed), the repository is ready for a real database layer.

The recommended stack is `pgx` (Postgres driver) + `sqlc` (type-safe query generation) + `goose` (schema migrations). The verification harness (`scripts/lib/docker-verification.sh`) already has reusable Postgres sidecar helpers, and `scripts/.versions` already pins `POSTGRES_TEST_IMAGE=postgres:16-alpine`.

## Goals

- Introduce a persistence foundation: pgx connection pooling, sqlc query generation, goose migration framework.
- Replace in-memory data sources for `/api/planning-window` and `/api/dashboard/people-availability` with Postgres-backed queries.
- Add a Postgres service to `compose.yml` with healthcheck and `DATABASE_URL` wiring.
- Run goose migrations at backend startup via the Go library (embed.FS), gated by `DB_AUTO_MIGRATE` env var (default `true`).
- Preserve the exact API contracts for all existing endpoints, including `/api/hello`, `/api/tasks/backlog`, CORS, and `/openapi.json`.
- Add DB-backed integration tests reusing the existing verification Postgres sidecar helpers.

## Non-Goals

- Designing or implementing new product endpoints beyond the existing three.
- Adding authentication, multi-user behavior, or production-grade deployment concerns.
- Building full task CRUD, people CRUD, or write/edit UI flows.
- Wiring the `/api/tasks/backlog` endpoint to Postgres — it stays in-memory.
- Redesigning frontend views or generated client contracts.
- Adding a separate migration container or migration step in CI/compose.

## Conflict Resolution

One refinement participant recommended wiring all three endpoints (including `/api/tasks/backlog`) to Postgres. Two others explicitly recommended scoping to only the planning-window and dashboard-people-availability endpoints, keeping tasks-backlog in-memory. The dossier acceptance criteria name only the two dashboard/planning endpoints, and the non-goals exclude "full task CRUD." Resolution: follow the majority recommendation and the acceptance criteria — only `/api/planning-window` and `/api/dashboard/people-availability` get DB wiring. The tasks backlog stays in-memory as a deliberate scope boundary.

## Decisions

### D1: Store interface abstraction

**Decision:** Introduce a `Store` interface in `backend/store.go` with methods `GetPlanningWindow(ctx) (*PlanningWindowData, error)` and `GetPeopleAvailability(ctx, startDate, days) (*DashboardData, error)`. Implement a `PgStore` struct backed by pgxpool and sqlc queries. Wire handlers to accept the `Store` interface via function parameters (DI at registration time).

**Rationale:** The interface decouples handler logic from the data source, enabling testability via mocks for fast unit tests and real pgx pools for integration tests. The existing handler registration pattern (`registerDashboardPeopleAvailability(api)` → `registerDashboardPeopleAvailability(api, store)`) is a minimal refactor.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-reviewer`

### D2: Schema — three tables, minimal surface

**Decision:** Create `planning_windows` (id, start_date, end_date, created_at, updated_at), `people` (id, name, initials, created_at), and `availability` (id, person_id FK, date, status VARCHAR with CHECK constraint, created_at). The schema covers exactly the data needed by the two endpoints and nothing more.

**Rationale:** The dossier and acceptance criteria scope to planning-window and dashboard-people-availability data only. Adding a `tasks` table would expand scope beyond the acceptance criteria. The `availability.status` uses a CHECK constraint (`IN ('available','busy','partial','off')`) to enforce the design system vocabulary at the database level.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-reviewer`, design-system-v2.md

### D3: Migrations — goose Go library, embed.FS, startup timing

**Decision:** Use `github.com/pressly/goose/v3` as a Go library (not CLI). Embed migration SQL files via `//go:embed migrations/*.sql`. Run migrations in `main()` after establishing a pgx pool but before starting the HTTP server. Gate migrations behind `DB_AUTO_MIGRATE` env var (default `"true"`).

**Rationale:** Running migrations at startup via the Go library avoids adding a goose CLI binary to the Docker image, keeps the Dockerfile simple, and matches the local-dev convenience pattern. The `DB_AUTO_MIGRATE` gate provides an escape hatch for environments where migrations are managed externally.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`

### D4: Seed data — materialized SQL rows

**Decision:** The seed migration (`002_seed_demo_data.sql`) inserts one planning window row (`2026-07-05`..`2026-08-13`), 8 people (p1–p8 with names matching the current `seedPeople`), and availability rows that reproduce the exact closure-based statuses for every person on every date in the planning window range. The seed covers exactly `2026-07-05` through `2026-08-13` (40 days × 8 people = 320 availability rows).

**Rationale:** The current Go code uses `always("available")` for p1–p6, `always("busy")` for p7, and `cycleStatuses` (off/partial/busy/available by dayOffset % 4) for p8. Materializing the full cross-product of 8 people × 40 dates in SQL preserves the exact endpoint behavior. Seeding only the planning window range (not an unbounded buffer) keeps the seed migration deterministic and reviewable.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-reviewer`

### D5: Env var contract — DATABASE_URL

**Decision:** The backend reads a single `DATABASE_URL` environment variable as the Postgres connection string (format: `postgres://user:pass@host:5432/dbname?sslmode=disable`). `pgxpool.ParseConfig` parses it directly.

**Rationale:** Single DSN is simpler to configure, matches `pgx` conventions, and reduces env var surface in `compose.yml`. The existing `verification_postgres_sidecar_url()` helper already returns a DSN in this format.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-reviewer`

### D6: sqlc — committed generated code

**Decision:** Configure `sqlc` via `backend/sqlc.yaml` to read query files from `backend/queries/` and generate Go code into `backend/db/`. Commit the generated Go files. The `sqlc` CLI is NOT required at Docker build time.

**Rationale:** Committing generated code is standard Go monorepo practice and avoids requiring `sqlc` in the Docker build image. The `scripts/check` step (`go vet ./... && go build ./...`) naturally catches drift between queries and generated code.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`

### D7: Test strategy — interface mock + build-tag integration

**Decision:** Handler unit tests use a mock `Store` (in-memory map-based implementation) imported from `backend/store_mock_test.go`. Integration tests use the `//go:build integration` tag and the existing `verification_start_postgres_sidecar` / `verification_cleanup_postgres_sidecar` helpers. The `newTestAPI` function is refactored to accept a `Store` parameter.

**Rationale:** This split preserves fast `go test -short` for handler logic while providing full database integration coverage when Docker is available. The existing `scripts/precommit-run` runs inside Docker with DinD socket access, so integration tests run there naturally.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`

### D8: compose.yml — db service with healthcheck

**Decision:** Add a `db` service to `compose.yml` using `postgres:16-alpine` with `POSTGRES_DB=the_great_migration`, `POSTGRES_USER=app`, `POSTGRES_PASSWORD=app`. Healthcheck via `pg_isready`. Backend gets `DATABASE_URL=postgres://app:app@db:5432/the_great_migration?sslmode=disable` and `depends_on: db: service_healthy`.

**Rationale:** Uses the same image version already pinned in `scripts/.versions`. The backend-on-db health dependency prevents premature connection attempts.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`

## Risks

| Risk | Mitigation |
|------|------------|
| Migration startup delay: goose migrations run before HTTP server starts, adding latency to backend readiness. | Backend retries DB connection with backoff. Compose `depends_on: service_healthy` on the db service ensures Postgres is ready before the backend starts. |
| sqlc generated-code coupling: schema changes cascade into regenerated types. | Keep query files narrowly focused on read models. Committed code surfaces build failures early. |
| Test sidecar lifecycle: orphaned containers if cleanup fails. | `verification_cleanup_postgres_sidecar` called in `t.Cleanup`. Unique run IDs prevent collisions. |
| Seed data range: availability queries outside `2026-07-05`..`2026-08-13` return empty rows. | The default `start` parameter uses server-local current date; queries outside the seed window return 0 availability rows per person. The handler returns empty availability arrays rather than errors, preserving the contract shape. |
| Tasks endpoint asymmetry: `/api/tasks/backlog` stays in-memory while other endpoints become DB-backed. | Documented in code comments and this design record as an intentional scope boundary tied to the acceptance criteria. |
| Backend Dockerfile changes: adding goose/embed.FS requires migration files at build time. | Migration `.sql` files are committed under `backend/migrations/`. The `//go:embed migrations/*.sql` directive picks them up at build time automatically. |

## Traceability

- **Task**: `248d3e0f-c580-48ad-9b19-4c6785beead2`
- **Dossier**: `2026-06-15T19:29:43.926Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Validated round outputs**: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- **Artifact base**: `initial` snapshot for `task-248d3e0f-c580-48ad-9b19-4c6785beead2`