## 1. Add frontend API generation dependencies and configuration

- [ ] Add pinned frontend dependencies for `pinia`, `@pinia/colada`, and the Hey API packages required for code generation and the generated fetch/SDK client.
- [ ] Create `frontend/openapi-ts.config.ts` with the committed snapshot as input and `frontend/src/client/` as the generated output.
- [ ] Add and document an explicit `generate:api` script that regenerates the snapshot-backed client artifacts.
- [ ] Commit `frontend/openapi-snapshot.json` as the deterministic default codegen input.
- [ ] Generate and commit the client artifacts required for `GET /api/hello`, including types, fetch client, SDK functions, and Pinia Colada query artifacts.

## 2. Bootstrap Pinia and Pinia Colada for generated queries

- [ ] Update `frontend/src/main.ts` to install `createPinia()`, then Pinia Colada, then the existing router before mounting the app.
- [ ] Ensure generated client runtime configuration preserves the current same-origin `/api` behavior by default and supports `VITE_API_BASE_URL` when explicitly set.

## 3. Replace the HomeView manual fetch with the generated query

- [ ] Refactor `frontend/src/home/HomeView.vue` to remove the manual `onMounted` fetch and use the generated Pinia Colada query for `GET /api/hello`.
- [ ] Preserve the existing loading, error, and success card states in the first summary card.
- [ ] Keep the remaining summary cards and the "Today's plan" / "Move notes" sections unchanged.

## 4. Keep SSR render verification working

- [ ] Update `frontend/tests/app-routes-render.test.ts` so the SSR app installs Pinia and Pinia Colada before rendering routes that consume generated queries.
- [ ] Verify the route render smoke test continues to pass without requiring a live backend.

## 5. Verify the change through the existing Docker-backed workflow

- [ ] Run `scripts/check` from the repo root.
- [ ] Run `scripts/test` from the repo root.
- [ ] Run `scripts/precommit-run` from the repo root.