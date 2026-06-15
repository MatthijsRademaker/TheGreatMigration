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
