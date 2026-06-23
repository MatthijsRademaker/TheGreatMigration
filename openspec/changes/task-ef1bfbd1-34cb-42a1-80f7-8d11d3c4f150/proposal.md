## Why

The frontend currently proves backend connectivity with a hand-written fetch in `frontend/src/home/HomeView.vue`, while the backend already exposes a Huma-generated OpenAPI document at `/openapi.json`. This change turns that ad hoc integration into a typed generated API layer and Pinia Colada query integration without making normal frontend builds, checks, or tests depend on a live backend.

## What Changes

- Add Hey API code generation to the frontend using a committed OpenAPI snapshot sourced from the backend's `/openapi.json` document.
- Add `pinia` and `@pinia/colada` to the frontend bootstrap so generated Pinia Colada query artifacts can be consumed by Vue views.
- Create `frontend/openapi-ts.config.ts`, commit `frontend/openapi-snapshot.json`, and generate committed client artifacts under `frontend/src/client/`.
- Configure the generated fetch client to default to an empty-string base URL so browser requests continue to flow through the existing Vite `/api` proxy, while allowing `VITE_API_BASE_URL` to override the base URL when explicitly provided.
- Refactor `frontend/src/home/HomeView.vue` to replace the manual hello fetch with the generated `GET /api/hello` Pinia Colada query while preserving the current loading, error, and success states in the first summary card.
- Update the SSR route render test to install the same Pinia and Pinia Colada plugins required by query-backed views.
- Provide an explicit `generate:api` workflow so developers can refresh the snapshot and generated client from a running backend when needed.

## Impact

- Frontend dependencies and bootstrap change to include the query stack required for generated API usage.
- The repository gains committed OpenAPI/codegen artifacts so `scripts/check`, `scripts/test`, and `scripts/precommit-run` continue to work without starting Compose or a live backend.
- The hello-world dashboard integration becomes type-safe and generated from the backend contract instead of manually coded fetch logic.
- Existing non-hello dashboard cards, "Today's plan", and "Move notes" remain unchanged.
- No new backend endpoints, persistence, authentication, roles, or unrelated UI redesign work are introduced.
