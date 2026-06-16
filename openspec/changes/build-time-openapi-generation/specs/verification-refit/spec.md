## MODIFIED Requirements

### Requirement: Root verification scripts SHALL cover both the frontend and backend services
The `scripts/check` and `scripts/test` verification scripts SHALL delegate to both the `frontend/` (Node) and `backend/` (Go) service directories using the existing Docker-backed verification harness. The scripts SHALL run frontend checks first, backend checks second, and report which service failed if either exits non-zero. The backend check phase SHALL additionally validate that the committed OpenAPI snapshot matches the backend code.

#### Scenario: scripts/check runs both frontend and backend verification
- **WHEN** `scripts/check` is executed from the repo root
- **THEN** the script runs `cd frontend && npm install && npm run check` inside NODE24_IMAGE, then runs `cd backend && go mod tidy && go vet ./... && go build ./...` inside GO_IMAGE, and both phases pass

#### Scenario: scripts/test runs both frontend and backend tests
- **WHEN** `scripts/test` is executed from the repo root
- **THEN** the script runs `cd frontend && npm install && npm test` inside NODE24_IMAGE, then runs `cd backend && go test ./...` inside GO_IMAGE, and both phases pass

#### Scenario: Frontend failure is reported clearly
- **WHEN** `scripts/check` is executed and the frontend phase fails (e.g., vue-tsc type error)
- **THEN** the script exits non-zero with output indicating the frontend check failed

#### Scenario: Backend failure is reported clearly
- **WHEN** `scripts/check` is executed and the backend phase fails (e.g., go vet warning)
- **THEN** the script exits non-zero with output indicating the backend check failed

#### Scenario: Pre-commit hooks remain functional
- **WHEN** a commit triggers `.pre-commit-config.yaml` hooks
- **THEN** the `check` and `test` hooks invoke `scripts/check` and `scripts/test` from the repo root, which now cover both services

### Requirement: scripts/check SHALL verify the committed OpenAPI snapshot matches backend code

The backend check phase in `scripts/check` SHALL, after `go build ./...` succeeds, run `go run ./cmd/openapi-gen` and diff its stdout against the committed `frontend/openapi-snapshot.json`. If the files differ, the check SHALL fail with a message instructing the developer to run `npm run regen:api` from the frontend directory.

#### Scenario: Snapshot matches backend code
- **WHEN** `scripts/check` is executed and `frontend/openapi-snapshot.json` is consistent with the backend endpoint definitions
- **THEN** the diff step passes and the overall check succeeds

#### Scenario: Snapshot is stale
- **WHEN** `scripts/check` is executed and backend endpoints have been added, removed, or modified without regenerating the snapshot
- **THEN** the diff step fails, the check exits non-zero, and the error message indicates that `npm run regen:api` should be run

#### Scenario: openapi-gen tool fails to compile
- **WHEN** `scripts/check` is executed and `go run ./cmd/openapi-gen` fails (e.g., compilation error in registration code)
- **THEN** the check exits non-zero with the Go compilation error visible in output
