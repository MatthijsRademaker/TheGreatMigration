## MODIFIED Requirements

### Requirement: Backend SHALL inject Store into handlers

The `main()` function SHALL construct a `Store` implementation and pass it to `api.RegisterAll(api, store)`. The `RegisterAll` function in the `backend/api/` package SHALL call every individual `register*` function with the provided store. The openapi-gen tool SHALL call `api.RegisterAll(api, nil)` for spec-only generation without a database connection.

#### Scenario: Server handler registration uses shared RegisterAll
- **WHEN** the backend starts normally
- **THEN** `main()` calls `api.RegisterAll(api, store)` with a non-nil `Store` and all handlers are registered

#### Scenario: openapi-gen uses nil store
- **WHEN** `go run ./cmd/openapi-gen` is executed
- **THEN** `api.RegisterAll(api, nil)` is called and produces a valid OpenAPI spec without panicking

### Requirement: A Go backend service SHALL serve a hello-world endpoint with auto-generated OpenAPI

The service SHALL expose `GET /api/hello` returning `{"message": "Hello from the backend!"}` and SHALL serve auto-generated OpenAPI 3.1 specification at `/openapi.json`. The service SHALL include a Dockerfile for compose orchestration. Endpoint registration types and functions SHALL live in the `backend/api/` package, shared between the server binary and the openapi-gen tool.

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

#### Scenario: Go module structure builds including openapi-gen
- **WHEN** `go build ./...` is executed in `backend/`
- **THEN** the command compiles both the server binary and `cmd/openapi-gen` without errors

### Requirement: Existing backend behavior SHALL remain intact

`GET /api/hello` SHALL continue to return `{"message": "Hello from the backend!"}`. CORS SHALL continue to allow `http://localhost:5173` and `http://frontend:5173`. The Dockerfile SHALL continue to build from `golang:1.26.2-alpine` and produce a working image. `go vet ./...` and `go build ./...` SHALL continue to pass. All existing API endpoints SHALL respond identically to pre-refactor behavior.

#### Scenario: Hello endpoint returns valid JSON
- **WHEN** an HTTP client sends `GET /api/hello` to the running backend
- **THEN** the response status is 200 and the body is `{"message": "Hello from the backend!"}` with `Content-Type: application/json`

#### Scenario: OpenAPI specification is auto-generated
- **WHEN** the backend is running
- **THEN** `GET /openapi.json` returns a valid OpenAPI 3.1 specification document

#### Scenario: Dockerfile builds and runs the backend
- **WHEN** `docker build backend/` is executed using `golang:1.26.2-alpine`
- **THEN** the resulting image starts the Go service and responds to `GET /api/hello` on port 8080
