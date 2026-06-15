## Context

The frontend is a Vue 3 + Vite application that currently demonstrates backend connectivity with a hand-written fetch to `GET /api/hello` in `frontend/src/home/HomeView.vue`. The backend already exposes a Huma-generated OpenAPI 3.1 document at `/openapi.json`, and the existing local/compose runtime uses the Vite `/api` proxy to reach the backend. The repository's Docker-backed verification scripts (`scripts/check`, `scripts/test`, `scripts/precommit-run`) run without starting Compose, so this change must not make normal frontend verification depend on a live backend.

## Goals

- Add Hey API code generation to the frontend using the backend OpenAPI document as the source of truth.
- Add Pinia and Pinia Colada to the app bootstrap so generated Pinia Colada query artifacts can be consumed by Vue views.
- Replace the existing HomeView hello fetch with a generated, type-safe Pinia Colada query while preserving the current loading, error, and success states.
- Provide a repeatable generation command/configuration and commit the generated client artifacts needed for normal check/build/test runs.
- Keep the implementation aligned with the existing Vite proxy, Compose setup, strict TypeScript checks, and Docker-backed verification flow.

## Non-Goals

- Adding new backend endpoints, persistence, authentication, roles, or production collaboration features.
- Redesigning the dashboard or changing unrelated static task, schedule, people, sidebar, or design-system content.
- Requiring the frontend Docker image build, `vue-tsc`, or Vitest to contact a live backend unless an explicit API generation workflow is invoked.
- Introducing optional Pinia Colada devtools as part of this change.

## Decisions

### 1. Code generation input and output

Use a committed static OpenAPI snapshot at `frontend/openapi-snapshot.json` as the default Hey API input so generation and type-checking remain deterministic when the backend is not running. Configure `frontend/openapi-ts.config.ts` to generate committed artifacts under `frontend/src/client/`.

### 2. Generation workflow

Provide an explicit frontend `generate:api` script that runs Hey API against the config file. The workflow must document how developers refresh `frontend/openapi-snapshot.json` from a running backend's `/openapi.json` endpoint and then regenerate the committed client artifacts.

### 3. Generated client composition

Generate and commit the TypeScript types, fetch client, SDK functions, and Pinia Colada query artifacts needed for `GET /api/hello`. This keeps normal `scripts/check`, `scripts/test`, and `scripts/precommit-run` independent of live backend availability.

### 4. Runtime base URL behavior

Configure the generated fetch client to default to an empty-string base URL so browser requests continue to use the existing Vite `/api` proxy. Allow `VITE_API_BASE_URL` to override the base URL when explicitly set.

### 5. Frontend bootstrap

Install `createPinia()` before Pinia Colada in `frontend/src/main.ts`, then install the existing router before mounting the app. Query-backed views depend on this plugin order.

### 6. HomeView integration

Replace the manual `onMounted` fetch in `frontend/src/home/HomeView.vue` with the generated Pinia Colada query for the backend `get-hello` operation. Preserve the current first-card loading, error, and success behaviors and leave all other dashboard content unchanged.

### 7. SSR render compatibility

Update the SSR route render test to install Pinia and Pinia Colada in the SSR app setup before `renderToString`, matching the runtime bootstrap required by query-backed views.

## Risks

- **Generated export shape mismatch**: Exact generated symbol names for the `get-hello` operation must be verified after the first generation run. Mitigation: run generation early and wire `HomeView` to the actual generated exports.
- **Strict TypeScript checks on generated code**: `frontend/tsconfig.json` includes strict `src/**/*.ts` checking with `noUnusedLocals` and `noUnusedParameters`. Mitigation: verify generated output passes existing checks and adjust codegen/check configuration only if required to keep committed artifacts checkable.
- **Snapshot drift**: A committed snapshot can fall behind the running backend contract. Mitigation: document the refresh workflow in scripts or comments so snapshot regeneration is a deliberate developer action.
- **SSR query context failures**: Query-backed views will fail SSR route rendering if Pinia and Pinia Colada are not installed before render. Mitigation: update the SSR render setup to mirror `main.ts` plugin installation order.

## Conflict Resolution

- **Live URL input vs committed snapshot**: Resolved to a committed `frontend/openapi-snapshot.json` as the default codegen input because normal Docker-backed verification does not start a backend service.
- **Generated output location**: Resolved to `frontend/src/client/` per accepted decisions and reviewer guidance.
- **Base URL strategy**: Resolved to an empty-string default with optional `VITE_API_BASE_URL` override so browser requests preserve the existing Vite `/api` proxy contract.
- **Pinia Colada devtools**: Deferred because the accepted decisions and task goals cover only the core Hey API + Pinia Colada integration.

## Traceability

- Task: `ef1bfbd1-34cb-42a1-80f7-8d11d3c4f150`
- Dossier: `2026-06-15T17:09:20.713Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial` snapshot for `task-ef1bfbd1-34cb-42a1-80f7-8d11d3c4f150`