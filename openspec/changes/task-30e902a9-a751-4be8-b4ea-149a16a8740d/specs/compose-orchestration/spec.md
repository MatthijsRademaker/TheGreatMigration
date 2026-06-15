## ADDED Requirements

### Requirement: A root Docker Compose file SHALL orchestrate frontend and backend together
The change SHALL add a `compose.yml` at the repo root that starts both the frontend and backend services on a shared network. The frontend service SHALL be configured to reach the backend via a `VITE_API_BASE_URL` environment variable. The compose setup SHALL target dev-mode with volume mounts for frontend hot-reload.

#### Scenario: Both services start successfully
- **WHEN** `docker compose up` is executed from the repo root
- **THEN** both the frontend and backend services reach a healthy state without errors

#### Scenario: Frontend can reach backend over compose network
- **WHEN** the compose stack is running
- **THEN** the frontend service can resolve and connect to `http://backend:8080/api/hello` over the shared compose network

#### Scenario: Frontend hot-reload works via volume mounts
- **WHEN** a source file in `frontend/src/` is modified while the compose stack is running
- **THEN** the Vite dev server detects the change and triggers HMR without requiring a container rebuild

#### Scenario: Backend healthcheck prevents premature frontend startup
- **WHEN** `docker compose up` is executed
- **THEN** the frontend service waits for the backend healthcheck to pass before starting, preventing failed fetch attempts from a not-yet-ready backend