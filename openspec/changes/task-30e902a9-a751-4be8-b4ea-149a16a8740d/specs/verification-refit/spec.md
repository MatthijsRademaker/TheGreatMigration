## ADDED Requirements

### Requirement: Root verification scripts SHALL cover both the frontend and backend services
The `scripts/check` and `scripts/test` verification scripts SHALL be rewritten to delegate to both the `frontend/` (Node) and `backend/` (Go) service directories using the existing Docker-backed verification harness. The scripts SHALL run frontend checks first, backend checks second, and report which service failed if either exits non-zero.

#### Scenario: scripts/check runs both frontend and backend verification
- **WHEN** `scripts/check` is executed from the repo root
- **THEN** the script runs `cd frontend && npm install && npm run check` inside NODE24_IMAGE, then runs `cd backend && go vet ./... && go build ./...` inside GO_IMAGE, and both phases pass

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