## Why

The OpenAPI snapshot (`frontend/openapi-snapshot.json`) is currently fetched from a running backend via `refresh-openapi-snapshot.mjs`. This requires a live backend + database, making it impossible to validate in CI and a constant merge-conflict magnet. A recent merge stripped people and rooms endpoint definitions from the snapshot without detection, causing runtime crashes in Vue components that import now-nonexistent generated exports. Build-time generation from Go types eliminates the external dependency and enables CI validation that the committed snapshot matches the backend code.

## What Changes

- **New** `backend/cmd/openapi-gen/main.go`: Build-time tool that registers all Huma endpoints and writes OpenAPI 3.1 JSON to stdout — no server, no database required
- **Refactor** `backend/main.go`: Extract a shared `registerAllEndpoints(api, store)` function so `main.go` and `openapi-gen` use identical endpoint registration
- **Replace** `frontend/scripts/refresh-openapi-snapshot.mjs`: Instead of `fetch("http://localhost:8080/openapi.json")`, pipe `go run ./cmd/openapi-gen` to `openapi-snapshot.json`
- **Add** precommit validation to `scripts/check`: Verify committed `openapi-snapshot.json` matches output of `go run ./cmd/openapi-gen` — fails CI if they diverge
- **Update** `frontend/package.json` scripts: `refresh:openapi-snapshot` and `regen:api` use the Go tool instead of HTTP fetch

## Capabilities

### New Capabilities
- `build-time-openapi-gen`: A Go command in `backend/cmd/openapi-gen/` that registers all API endpoints and emits the full OpenAPI 3.1 specification to stdout without a running server or database connection

### Modified Capabilities
- `backend-service`: Refactor endpoint registration in `main.go` into a shared `registerAllEndpoints(api, store)` function callable by both the server binary and the openapi-gen tool
- `verification-refit`: Add snapshot-consistency check to `scripts/check` that fails precommit when committed `openapi-snapshot.json` diverges from what `go run ./cmd/openapi-gen` produces

## Impact

- **Backend**: `main.go` minor refactor (extract `registerAllEndpoints`), new `cmd/openapi-gen/main.go` with zero new dependencies
- **Frontend**: `scripts/refresh-openapi-snapshot.mjs` replaced, `package.json` scripts updated, no dependency changes
- **CI/Precommit**: `scripts/check` gains a Go build + diff step; Go image must be available (already is)
- **DX**: Developers no longer need a running backend to regenerate the API client — just `go run`
- **No breaking changes**: Generated `sdk.gen.ts`, `colada.gen.ts`, `types.gen.ts` are identical when backend code is unchanged
