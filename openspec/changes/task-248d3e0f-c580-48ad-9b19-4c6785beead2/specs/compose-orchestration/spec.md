# compose-orchestration Specification (Delta)

## Purpose

Update the canonical `openspec/specs/compose-orchestration/spec.md` to add a Postgres database service and wire it to the backend.

## ADDED Requirements

### Requirement: Compose SHALL include a Postgres database service

`compose.yml` SHALL define a `db` service using the `postgres:16-alpine` image. The service SHALL expose port 5432 on the `app-network`. The service SHALL set `POSTGRES_DB=the_great_migration`, `POSTGRES_USER=app`, and `POSTGRES_PASSWORD=app`. The service SHALL include a healthcheck using `pg_isready -U app -d the_great_migration` with interval 5s, timeout 3s, retries 5, and start_period 5s. The service SHALL use a named volume `pgdata` mounted at `/var/lib/postgresql/data` for data persistence across container restarts.

#### Scenario: Postgres service starts and becomes healthy
- **WHEN** `docker compose up` is executed from the repo root
- **THEN** the `db` service reaches a healthy state as reported by `docker compose ps`

#### Scenario: Postgres data persists across restarts
- **WHEN** the compose stack is stopped and restarted
- **THEN** the seed data and any later modifications are preserved

### Requirement: Backend SHALL depend on the database service

The backend service in `compose.yml` SHALL receive a `DATABASE_URL` environment variable set to `postgres://app:app@db:5432/the_great_migration?sslmode=disable`. The backend service SHALL declare `depends_on` with `db: service_healthy` so that the backend waits for Postgres to be ready before starting. The backend service SHALL also receive `DB_AUTO_MIGRATE=true`.

#### Scenario: Backend waits for database before starting
- **WHEN** `docker compose up` is executed
- **THEN** the backend service does not start until the `db` service passes its healthcheck

#### Scenario: Backend can connect to database over compose network
- **WHEN** the compose stack is running
- **THEN** the backend service can resolve and connect to `db:5432`

### Requirement: Existing compose behavior SHALL remain intact

The frontend service, its `depends_on` the backend, `VITE_API_BASE_URL` wiring, volume mounts, and port mappings SHALL remain unchanged. The backend's existing healthcheck (`wget /api/hello`) and port mapping SHALL remain unchanged.

#### Scenario: Frontend still waits for backend
- **WHEN** `docker compose up` is executed
- **THEN** the frontend service waits for the backend healthcheck to pass before starting

#### Scenario: All three services become healthy
- **WHEN** `docker compose up` is executed
- **THEN** the `db`, `backend`, and `frontend` services all reach a healthy state