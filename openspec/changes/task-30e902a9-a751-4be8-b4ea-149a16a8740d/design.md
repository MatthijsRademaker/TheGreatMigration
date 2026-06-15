## Context

The repository is a single-root Vue 3 + Vite + shadcn-vue + Tailwind v4 application with Docker-backed verification scripts, pre-commit hooks, and a design system rooted at `src/app/styles.css`. There is no existing backend, no Docker Compose orchestration, and no frontend API integration. The task requires splitting this monolith into a multi-service layout with `frontend/` and `backend/` service directories, orchestrated by root-level Compose, with a minimal hello-world integration as the first runtime connection.

## Goals

- Relocate the Vue app into `frontend/` with all path aliases, theme config, shadcn wiring, and tests updated so the existing dashboard routes render identically after the move.
- Create a minimal Go backend in `backend/` using Huma v2 that exposes `GET /api/hello` returning `{"message": "Hello from the backend!"}` and generates OpenAPI at `/openapi.json`.
- Wire both services through a root `compose.yml` on a shared network, targeting dev-mode with Vite HMR and Go service reachable by the frontend.
- Integrate a live hello-world fetch into HomeView, replacing static placeholder data with a runtime backend call.
- Refit root verification scripts (`scripts/check`, `scripts/test`) to verify both frontend Node checks and backend Go checks using the existing Docker-backed harness.

## Non-Goals

- Designing or implementing real move-planning backend domain APIs beyond the hello-world proof of integration.
- Adding authentication, persistence, user management, or production deployment hardening.
- Redesigning the existing frontend information architecture or design system as part of the repo split.
- Supporting production-optimized compose profiles (multi-stage builds, nginx, compiled binaries) — the compose target is dev-mode only.
- Adding hot-reload tooling for the Go backend (e.g., air) in this change.
- Preserving root `package.json` — all Node dependencies and scripts move entirely into `frontend/`.

## Decisions

### 1. Directory layout

The repo root serves as the orchestration layer with `compose.yml` and updated `scripts/`. The Vue app moves entirely into `frontend/` (including `src/`, `tests/`, `package.json`, `vite.config.ts`, `tsconfig.json`, `components.json`, `index.html`, and `env.d.ts`). The Go backend lives in `backend/` as a new service directory.

### 2. Frontend relocation mechanics

All moved files use paths relative to their own location, which means `vite.config.ts` (`@` → `./src`), `tsconfig.json` (`@/*` → `./src/*`), and `components.json` (`css: src/app/styles.css`) work correctly as-is after relocation to `frontend/` because the paths are file-relative. The following specific updates are required:

- `index.html` moves to `frontend/index.html` and its `<script type="module" src="/src/main.ts">` path resolves correctly within the Vite project root at `frontend/`.
- Test imports (`tests/app-routes-render.test.ts`, `tests/design-system-primitives.test.ts`) use `../src/...` relative paths that continue to work after both `src/` and `tests/` move under `frontend/`.
- `docs/design-system-v2.md` reference to `src/app/styles.css` must be updated to `frontend/src/app/styles.css`.
- Files are relocated using `git mv` to preserve commit history.

### 3. Frontend-backend communication strategy

The frontend communicates with the backend using a Vite dev proxy in local dev and a `VITE_API_BASE_URL` environment variable in the compose stack. CORS middleware is configured on the backend as a safety net for direct access. In the compose dev-mode workflow, the frontend fetches `/api/hello` which Vite proxies to the backend service, or uses `VITE_API_BASE_URL=http://backend:8080` for direct cross-service calls.

### 4. Backend technology stack

The backend is a Go module at `backend/go.mod` with module path `github.com/user/the-great-migration/backend`. It uses Huma v2 (`github.com/danielgtaylor/huma/v2`) with the chi adapter (`github.com/danielgtaylor/huma/v2/adapters/humago` and `github.com/go-chi/chi/v5`) for HTTP routing. The server exposes:

- `GET /api/hello` returning `{"message": "Hello from the backend!"}`
- `GET /openapi.json` serving the auto-generated OpenAPI 3.1 specification

CORS is configured to allow the Vite dev server origin (`http://localhost:5173`) and the compose frontend service.

### 5. Docker Compose dev-mode target

The `compose.yml` targets a dev-mode workflow: the frontend service runs the Vite dev server with HMR via volume mounts of the `frontend/` directory, and the backend service runs the Go application from the source directory. Both services connect on a shared compose network. The frontend receives `VITE_API_BASE_URL=http://backend:{BACKEND_PORT}` as a build argument or environment variable.

### 6. Verification refit design

`scripts/check` is rewritten to run two verification phases sequentially:
1. Frontend phase: `cd frontend && npm install && npm run check` inside the NODE24_IMAGE container.
2. Backend phase: `cd backend && go vet ./... && go build ./...` inside the GO_IMAGE container.

`scripts/test` follows the same pattern:
1. Frontend phase: `cd frontend && npm install && npm test` inside the NODE24_IMAGE container.
2. Backend phase: `cd backend && go test ./...` inside the GO_IMAGE container.

Both scripts exit with error if either phase fails, reporting which service failed. The `.pre-commit-config.yaml` hooks remain unchanged — they continue to invoke `scripts/check` and `scripts/test` from the repo root.

### 7. Hello-world integration placement

The hello-world response replaces the first summary card in `HomeView.vue` ("High priority / 4"). The card becomes a live-fetch component that calls `GET /api/hello` on mount and displays the returned message in place of the static value. The remaining three summary cards and the "Today's plan" / "Move notes" sections are preserved unchanged.

## Risks

- **Path alias breakage**: Five-plus files reference `src/` paths. Vite and TypeScript aliases are file-relative and should survive the move, but any missed reference (e.g., CI config, IDE settings, or ad-hoc scripts) will break silently. Mitigation: exhaustive search for `src/` and `./src` references in all non-ignored files after the move.
- **Verification script timing**: The Docker-backed harness runs one container per invocation. Sequential frontend-then-backend checks will roughly double CI time. Running in parallel would require more complex orchestration. Mitigation: accept sequential execution for correctness; optimize later if needed.
- **CORS misconfiguration**: If the backend CORS origin list does not match the frontend origin in each environment (localhost:5173 vs compose service name), the hello-world call fails silently in the browser. Mitigation: configure CORS for both `http://localhost:5173` and the compose frontend origin.
- **Module path choice**: The Go module path `github.com/user/the-great-migration/backend` is a placeholder. If the repo is later published under a different path, the module path must be updated. Mitigation: document this as a deliberate placeholder.
- **Compose startup ordering**: The frontend container may attempt to fetch the backend before it's ready. Mitigation: use depends_on with a healthcheck on the backend service.

## Conflict Resolution

- **Communication strategy**: Resolved to Vite dev proxy for local dev + `VITE_API_BASE_URL` env var for compose + CORS on backend as safety net. This combines the architect's Vite-proxy recommendation with the lead-dev's env-var-plus-CORS approach.
- **Root package.json**: Resolved to moving the entire `package.json` into `frontend/` with no root-level Node metadata remaining. The task is a structural split, not a monorepo workspace setup.
- **Compose profile**: Resolved to dev-mode only (volume mounts, Vite dev server, no production builds). A production profile can be added in a future change.
- **Verification execution order**: Resolved to sequential (frontend, then backend) for simplicity and clarity of failure reporting.
- **OpenAPI exposure**: Resolved to exposing `/openapi.json` from the backend since Huma generates it for free and it satisfies the auto-generation requirement.
- **Huma router adapter**: Resolved to chi (`humago` + `go-chi/chi/v5`) as the standard choice for Huma v2 HTTP services.

## Traceability

- Task: `30e902a9-a751-4be8-b4ea-149a16a8740d`
- Dossier: `2026-06-15T16:28:40.671Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial` snapshot for `task-30e902a9-a751-4be8-b4ea-149a16a8740d`