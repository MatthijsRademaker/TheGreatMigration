## ADDED Requirements

### Requirement: A standalone production compose file SHALL orchestrate the stack in production shape

The repository SHALL provide a `compose.prod.yml` at the repo root that orchestrates the `db`,
`backend`, and `frontend` services on a shared network using the `pgdata` named volume. It SHALL be
usable as a standalone entrypoint via `docker compose -f compose.prod.yml up` without layering the
dev `compose.yml`. The dev `compose.yml` and the dev Dockerfiles SHALL remain unchanged.

#### Scenario: Production stack starts from the standalone file
- **WHEN** `docker compose -f compose.prod.yml up` is executed from the repo root
- **THEN** the `db`, `backend`, and `frontend` services all reach a healthy state without referencing `compose.yml`

#### Scenario: Dev compose is unaffected
- **WHEN** `docker compose up` is executed against the unchanged `compose.yml`
- **THEN** the dev stack still runs the Vite dev server with HMR and seeded demo data

### Requirement: The production frontend SHALL be an optimized static build served by nginx

The production frontend image SHALL be built from `frontend/Dockerfile.prod` as a multi-stage
build: a Node stage producing an optimized `vite build` (`dist/`), and an `nginx` runtime serving
those static assets. The dev-server, source bind mount, `node_modules` volume, and HMR SHALL NOT be
present in the production frontend service.

#### Scenario: Frontend serves built assets, not a dev server
- **WHEN** the production frontend container is running
- **THEN** it serves the prebuilt static `dist/` assets via nginx
- **AND** no Vite dev server, source bind mount, or HMR is active

### Requirement: The production frontend SHALL proxy /api to the backend same-origin

nginx in the production frontend image SHALL reverse-proxy requests under `/api` to
`http://backend:8080` over the compose network, so the relative-URL API client works same-origin
with no build-time API base URL.

#### Scenario: API requests reach the backend through the proxy
- **WHEN** the running app issues a request to a relative `/api/...` URL
- **THEN** nginx forwards it to `backend:8080` and returns the backend response

### Requirement: The production frontend SHALL serve the SPA history fallback

Because the router uses history mode, nginx SHALL fall back to `index.html` for non-asset,
non-`/api` routes so that deep links resolve to the SPA.

#### Scenario: Deep link resolves to the SPA
- **WHEN** a client requests a client-side route (e.g. `/tools`) directly from the server
- **THEN** nginx returns `index.html` and the SPA renders the route (not a 404)

### Requirement: The production backend SHALL NOT seed demo data

The production `backend` service SHALL run schema migrations (`DB_AUTO_MIGRATE=true`) but SHALL NOT
enable seeding (`DB_SEED` unset/false), so a fresh production database contains the schema with no
demo rows.

#### Scenario: Fresh production database is schema-only
- **WHEN** the production stack starts against a fresh `pgdata` volume
- **THEN** the schema migrations are applied
- **AND** the `people`, `backlog_tasks`, `rooms_areas`, and `schedule_task_cards` tables contain no demo rows

### Requirement: The production backend SHALL be a slim runtime image

The production backend image SHALL be built from `backend/Dockerfile.prod` as a multi-stage build
that compiles the `CGO_ENABLED=0` binary and copies it into a minimal runtime image, rather than
shipping the Go toolchain and source. The runtime image SHALL retain the tooling needed for the
existing `wget`-based healthcheck.

#### Scenario: Backend healthcheck passes on the slim image
- **WHEN** the production backend container is running
- **THEN** its `wget`-based healthcheck against `/api/hello` succeeds and the service is reported healthy

### Requirement: Only the frontend SHALL be exposed to the host in production

In `compose.prod.yml`, only the `frontend` service SHALL publish a host port (`80`). The `backend`
and `db` services SHALL NOT publish host ports; they SHALL be reachable only over the compose
network. Service startup ordering SHALL be preserved via `depends_on` healthchecks (db → backend →
frontend).

#### Scenario: Internal services are not host-exposed
- **WHEN** the production stack is running
- **THEN** `http://localhost:80` serves the app
- **AND** `backend:8080` and `db:5432` are not reachable from the host

#### Scenario: Startup ordering is preserved
- **WHEN** `docker compose -f compose.prod.yml up` is executed
- **THEN** the backend starts only after `db` is healthy, and the frontend starts only after the backend is healthy
