## 1. Relocate the frontend into `frontend/`

- [ ] Move `src/`, `tests/`, `package.json`, `vite.config.ts`, `tsconfig.json`, `components.json`, `index.html`, and `env.d.ts` from repo root into `frontend/` using `git mv`.
- [ ] Verify `vite.config.ts` resolves `@` to `./src` relative to `frontend/vite.config.ts` (file-relative, should work as-is).
- [ ] Verify `tsconfig.json` resolves `@/*` to `./src/*` and `include` globs cover `src/` and `tests/` relative to `frontend/tsconfig.json`.
- [ ] Verify `components.json` points `css` to `src/app/styles.css` relative to `frontend/components.json`.
- [ ] Verify `index.html` script src `/src/main.ts` resolves correctly within the Vite project root at `frontend/index.html`.
- [ ] Update `docs/design-system-v2.md` to reference `frontend/src/app/styles.css` as the global theme surface.
- [ ] Verify test imports (`tests/*.test.ts` → `../src/...`) continue to resolve correctly after the move.

## 2. Create the Go backend with Huma v2

- [ ] Initialize `backend/go.mod` with module path `github.com/user/the-great-migration/backend` and Go 1.24.
- [ ] Add Huma v2 dependency: `github.com/danielgtaylor/huma/v2`.
- [ ] Add Huma chi adapter: `github.com/danielgtaylor/huma/v2/adapters/humago` and `github.com/go-chi/chi/v5`.
- [ ] Create `backend/main.go` with a Huma v2 server exposing `GET /api/hello` returning `{"message": "Hello from the backend!"}`.
- [ ] Configure CORS middleware to allow `http://localhost:5173` and the compose frontend origin.
- [ ] Configure Huma to serve auto-generated OpenAPI at `/openapi.json`.
- [ ] Create `backend/Dockerfile` using `golang:1.26.2-alpine` (aligned with `scripts/.versions` GO_IMAGE) for building and running the Go service.
- [ ] Add `go vet ./...` and `go build ./...` as the backend verification commands.

## 3. Add Docker Compose orchestration

- [ ] Create `compose.yml` at repo root with `frontend` and `backend` services on a shared network.
- [ ] Frontend service: build from `frontend/Dockerfile` with `VITE_API_BASE_URL=http://backend:8080` environment variable, expose port 5173, mount `frontend/` as a volume for HMR.
- [ ] Backend service: build from `backend/Dockerfile`, expose port 8080, configure healthcheck on `GET /api/hello`.
- [ ] Add `depends_on` with `condition: service_healthy` on the frontend service to prevent premature fetch attempts.

## 4. Wire hello-world into the frontend

- [ ] Replace the first summary card in `frontend/src/home/HomeView.vue` (currently "High priority / 4") with a live fetch to `GET /api/hello`.
- [ ] Display the returned `message` field in the card's value position with a loading state while the fetch is in-flight and an error state if the backend is unreachable.
- [ ] Preserve the remaining three summary cards and the "Today's plan" / "Move notes" sections unchanged.

## 5. Refit root verification scripts

- [ ] Rewrite `scripts/check` to:
  - Run frontend checks: `cd frontend && npm install && npm run check` inside NODE24_IMAGE container.
  - Run backend checks: `cd backend && go vet ./... && go build ./...` inside GO_IMAGE container.
  - Report which service failed if either phase exits non-zero.
- [ ] Rewrite `scripts/test` to:
  - Run frontend tests: `cd frontend && npm install && npm test` inside NODE24_IMAGE container.
  - Run backend tests: `cd backend && go test ./...` inside GO_IMAGE container.
  - Report which service failed if either phase exits non-zero.
- [ ] Verify `.pre-commit-config.yaml` hooks continue to invoke `scripts/check` and `scripts/test` without changes.

## 6. Verify the change

- [ ] Run `scripts/check` from repo root — both frontend and backend checks pass.
- [ ] Run `scripts/test` from repo root — both frontend and backend tests pass.
- [ ] Run `docker compose up` from repo root — both services start and the frontend's HomeView renders the hello-world message from the backend.
- [ ] Verify the four remaining summary cards and the "Today's plan" / "Move notes" sections still render correctly.