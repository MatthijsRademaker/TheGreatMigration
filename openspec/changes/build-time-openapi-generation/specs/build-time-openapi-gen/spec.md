## ADDED Requirements

### Requirement: The openapi-gen tool SHALL emit the complete OpenAPI 3.1 specification to stdout
The backend SHALL provide a Go command at `cmd/openapi-gen/main.go` that registers all API endpoints, calls `api.OpenAPI()`, and writes the resulting JSON to stdout. The tool SHALL NOT start an HTTP server, connect to a database, or require any environment variables. It SHALL exit zero on success and non-zero on failure.

#### Scenario: Tool outputs valid OpenAPI JSON
- **WHEN** `go run ./cmd/openapi-gen` is executed from the `backend/` directory
- **THEN** stdout contains valid JSON matching the OpenAPI 3.1 schema with all registered endpoint paths, and exit code is 0

#### Scenario: Tool includes all registered endpoints
- **WHEN** `go run ./cmd/openapi-gen` is executed
- **THEN** the output includes paths for `/api/hello`, `/api/dashboard/people-availability`, `/api/dashboard/daily-schedule`, `/api/planning-window`, `/api/tasks/backlog`, `/api/people`, `/api/people/{id}`, `/api/people/{id}/availability/{date}`, `/api/rooms`, and `/api/rooms/{id}`

#### Scenario: Tool requires no database connection
- **WHEN** `go run ./cmd/openapi-gen` is executed without a running Postgres instance and without `DATABASE_URL` set
- **THEN** the tool still produces valid OpenAPI JSON and exits zero

#### Scenario: Tool requires no environment variables
- **WHEN** `go run ./cmd/openapi-gen` is executed in a clean environment (no `DATABASE_URL`, `DB_AUTO_MIGRATE`, or other project env vars)
- **THEN** the tool exits zero and produces valid output

### Requirement: Endpoint registration SHALL be shared between server and openapi-gen
The backend SHALL provide a `backend/api/` package containing a `RegisterAll(api huma.API, store Store)` function. The server binary (`backend/main.go`) and the openapi-gen tool (`backend/cmd/openapi-gen/main.go`) SHALL both call this function. The store parameter SHALL accept `nil` for spec-only generation.

#### Scenario: Server uses shared registration
- **WHEN** `go build ./...` is executed from the `backend/` directory
- **THEN** the server binary imports `backend/api` and calls `api.RegisterAll()` with a real store

#### Scenario: openapi-gen uses shared registration
- **WHEN** `go run ./cmd/openapi-gen` is executed
- **THEN** the tool calls `api.RegisterAll()` with `nil` store and produces output identical in structure to the running server's `/openapi.json`

### Requirement: The refresh-openapi-snapshot script SHALL use the Go tool
The `frontend/scripts/refresh-openapi-snapshot.mjs` script SHALL be replaced. The `refresh:openapi-snapshot` npm script SHALL invoke `go run ./cmd/openapi-gen` from the backend directory and pipe stdout to `frontend/openapi-snapshot.json`.

#### Scenario: npm script regenerates snapshot
- **WHEN** `npm run refresh:openapi-snapshot` is executed from the `frontend/` directory with Go on PATH and the `backend/` module present
- **THEN** `frontend/openapi-snapshot.json` is overwritten with output matching what the backend would serve at `/openapi.json`

#### Scenario: npm script fails cleanly when Go is missing
- **WHEN** `npm run refresh:openapi-snapshot` is executed without `go` on PATH
- **THEN** the script exits non-zero with a clear error message
