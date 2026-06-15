## Why

The repository is currently a single root-level Vue/Vite application with no backend, no Docker Compose orchestration, and verification scripts that assume a single Node project. The task asks for a structural split: move the frontend into its own `frontend/` directory, introduce a new Go backend using Huma v2 for automatic OpenAPI generation, wire both together with Docker Compose, and connect them through a simple hello-world integration.

## What Changes

- **Frontend relocation**: Move all Vue/Vite source (`src/`), tests (`tests/`), and config files (`package.json`, `vite.config.ts`, `tsconfig.json`, `components.json`, `index.html`) into a new `frontend/` service directory. All path aliases, theme references, test imports, and build scripts are updated so the existing dashboard routes render identically after the move.
- **Backend service**: Create a new `backend/` Go module with Huma v2 powering a `GET /api/hello` endpoint that returns a JSON hello-world response with auto-generated OpenAPI output at `/openapi.json`. The service includes a Dockerfile for compose orchestration.
- **Docker Compose orchestration**: Add a root `compose.yml` that starts both services on a shared network with service names, port mappings, and the frontend configured to reach the backend via a `VITE_API_BASE_URL` environment variable. The compose setup targets dev-mode with volume mounts for frontend hot-reload.
- **Hello-world integration**: Replace the static summary-card data in `HomeView.vue` with a live fetch from the backend's `GET /api/hello` endpoint, visible at runtime.
- **Verification refit**: Rewrite `scripts/check` and `scripts/test` to delegate to both `frontend/` (Node via NODE24_IMAGE) and `backend/` (Go via GO_IMAGE) using the existing Docker-backed verification harness. The frontend checks run first, backend checks second, with failures reported clearly.

## Impact

- The repo root becomes an orchestration layer with `compose.yml`, updated `scripts/`, and thin verification wrappers. No root `package.json` remains after relocation.
- All existing Vue tooling, shadcn-vue theming, and design-system contracts continue to work from within `frontend/` with updated relative paths.
- The pre-commit pipeline (`.pre-commit-config.yaml`) continues to run `scripts/check` and `scripts/test`, now covering both services.
- A new runtime dependency is introduced: the HomeView dashboard card depends on the backend being reachable. If the backend is down, the card shows a loading or error state.
- No authentication, persistence, user management, or production deployment hardening is introduced.