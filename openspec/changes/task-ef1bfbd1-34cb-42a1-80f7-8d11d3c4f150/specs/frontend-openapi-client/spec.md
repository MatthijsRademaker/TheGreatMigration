## ADDED Requirements

### Requirement: The frontend SHALL generate and commit a typed API client from the backend OpenAPI contract

The frontend SHALL use Hey API configuration rooted in `frontend/openapi-ts.config.ts` with a committed OpenAPI snapshot at `frontend/openapi-snapshot.json` as the default input. The generated output SHALL be committed under `frontend/src/client/` and include the TypeScript types, fetch client, SDK functions, and Pinia Colada query artifacts required for the current backend contract.

#### Scenario: Snapshot-backed code generation is deterministic

- **WHEN** the frontend code generation configuration is inspected
- **THEN** the default OpenAPI input is the committed `frontend/openapi-snapshot.json` file
- **AND** the generated output path is `frontend/src/client/`

#### Scenario: Normal verification does not require a live backend

- **WHEN** `scripts/check`, `scripts/test`, or `scripts/precommit-run` are executed without a running backend
- **THEN** frontend install, type-check, and test workflows complete using committed snapshot and generated client artifacts
- **AND** no implicit code generation step requires contacting `/openapi.json`

#### Scenario: Regeneration is an explicit developer action

- **WHEN** a developer needs to refresh the frontend API client from the backend contract
- **THEN** the frontend provides an explicit generation workflow that documents how to refresh the snapshot from a running backend and regenerate the committed client artifacts

### Requirement: The generated frontend client SHALL preserve the existing browser-to-backend routing contract

The generated fetch client SHALL default to an empty-string base URL so browser requests continue to use the existing same-origin `/api` path through the Vite proxy. The client SHALL allow `VITE_API_BASE_URL` to override the base URL when explicitly provided.

#### Scenario: Same-origin proxy is the default browser behavior

- **WHEN** the generated client runs in the browser without `VITE_API_BASE_URL` set
- **THEN** requests to backend endpoints use same-origin `/api/...` URLs
- **AND** the existing Vite proxy remains the default path to the backend

#### Scenario: Explicit base URL override is supported

- **WHEN** `VITE_API_BASE_URL` is set for the frontend runtime
- **THEN** the generated client uses that value as its base URL for backend requests

### Requirement: Query-backed frontend views SHALL receive Pinia and Pinia Colada context before rendering

The application SHALL install Pinia before Pinia Colada in the frontend bootstrap so generated query artifacts can be consumed by Vue views. Any SSR route render setup used by the frontend verification suite SHALL install the same plugins before rendering views that consume generated queries.

#### Scenario: Runtime bootstrap provides query context

- **WHEN** `frontend/src/main.ts` initializes the Vue application
- **THEN** it installs `createPinia()` before Pinia Colada
- **AND** installs the existing router after the query stack is available

#### Scenario: SSR route rendering provides the same query context

- **WHEN** the SSR route render smoke test creates its app instance
- **THEN** it installs Pinia and Pinia Colada before rendering routes that use generated queries
- **AND** rendering does not fail because query context is missing