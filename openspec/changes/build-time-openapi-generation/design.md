## Context

The backend uses Huma v2 with the chi adapter to auto-generate OpenAPI 3.1 from Go structs and operation metadata. Currently, the frontend obtains this spec by fetching `http://localhost:8080/openapi.json` from a running backend and saving it as `frontend/openapi-snapshot.json`. The `@hey-api/openapi-ts` code generator then consumes this snapshot to produce `sdk.gen.ts`, `colada.gen.ts`, and `types.gen.ts`.

This pipeline has two failure modes:

1. **Requires live infrastructure**: `refresh-openapi-snapshot.mjs` needs a running backend + Postgres. Not feasible in CI or on a dev machine without starting compose services.

2. **Silent drift**: If someone modifies backend endpoints without regenerating the snapshot, or resolves a merge conflict incorrectly (as happened with the rooms/people merge), the committed snapshot diverges from backend code. Generated frontend code then references exports that don't exist, causing runtime crashes. TypeScript catches this (`vue-tsc --noEmit` reports missing exports) but the precommit check may pass if dependencies aren't installed.

Huma v2's `api.OpenAPI()` method builds the spec from registered operation metadata — it never invokes handlers or accesses external resources. This means the full OpenAPI spec can be produced at build time from Go types alone.

## Goals / Non-Goals

**Goals:**
- Generate `openapi-snapshot.json` from Go types without a running server or database
- CI validation that committed snapshot matches backend code
- Eliminate the HTTP fetch step from the regeneration workflow
- Single source of truth for endpoint registration (no duplication between server binary and spec generator)

**Non-Goals:**
- Changing the `@hey-api/openapi-ts` code generation step (it still consumes the snapshot)
- Removing the committed snapshot (stays committed for now; could be `.gitignore`'d later)
- Adding OpenAPI spec validation beyond structural correctness (already handled by Huma)
- Changing the backend framework (Huma v2 stays)

## Decisions

### Decision 1: Extract endpoint registration into `backend/api/` package

**Chosen**: Move all `register*` functions and their input/output types from `package main` into a new `backend/api/` library package with a `RegisterAll(api huma.API, store Store)` entry point.

**Rationale**: Go does not allow importing `main` packages. The `cmd/openapi-gen` binary needs to call the same registration functions as `main.go`. A shared library package is the Go-idiomatic solution. It also improves separation of concerns — request/response types and handler registration live in a dedicated API package rather than the server entrypoint.

**Alternatives considered**:
- Build tags with two `main` functions: messy, fragile, duplicates registration calls
- Keep everything in `main` and shell out to `go run . --openapi-gen`: requires runtime flags and still boots the server

### Decision 2: Pass `nil` store to openapi-gen

**Chosen**: `cmd/openapi-gen/main.go` calls `api.RegisterAll(api, nil)`. Huma's `OpenAPI()` introspects registered Go types — it never invokes handler closures where `store` is actually used.

**Rationale**: All six `register*` functions only reference `store` inside handler closures (the third argument to `huma.Register`). Registration itself (`huma.Register(api, op, handler)`) captures the handler but doesn't execute it. Verified by code audit of `dashboard.go`, `daily_schedule.go`, `people.go`, `planning_window.go`, `rooms_areas.go`, `tasks.go`.

**Risk**: If a future registration function uses `store` during registration (not inside a closure), `nil` would panic. Mitigation: add a comment at `RegisterAll` documenting this constraint.

### Decision 3: CI validation in `scripts/check`

**Chosen**: Add a step after `go build ./...` that runs `go run ./cmd/openapi-gen` and diffs against the committed snapshot. Fail if they diverge.

**Rationale**: Catches drift before merge. The Go image is already available in the check container (used for `go vet` and `go build`). The diff is fast (sub-second for a ~8KB JSON file).

**Format handling**: Huma's `json.Marshal` produces deterministic output (sorted keys in the standard library as of Go 1.21+), but whitespace may differ. Use `diff <(go run ./cmd/openapi-gen | jq .) <(jq . frontend/openapi-snapshot.json)` if `jq` is available, or `diff <(go run ./cmd/openapi-gen) frontend/openapi-snapshot.json` relying on Huma's stable output. The `refresh:openapi-snapshot` script will pipe `go run ./cmd/openapi-gen` directly to the snapshot file, ensuring format consistency.

### Decision 4: Frontend scripts call Go tool from frontend directory

**Chosen**: `refresh:openapi-snapshot` becomes `cd ../backend && go run ./cmd/openapi-gen > ../frontend/openapi-snapshot.json`. The `regen:api` script (which chains refresh + generate) stays the same.

**Rationale**: Scripts run from `frontend/`. Relative path to backend is `../backend`. Using `go run ./cmd/openapi-gen` means developers don't need to pre-compile anything — Go handles it. The `go` binary must be on PATH (already required for backend development).

**Alternative**: Compile a binary during `go build` and invoke it. More complex, no benefit for a small codegen tool.

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| `nil` store panics if future registration code accesses store outside closures | Document constraint in `RegisterAll` godoc; CI catches immediately (diff fails if openapi-gen panics) |
| Huma `json.Marshal` output changes between versions, causing CI false positives | Pin Huma version in `go.mod`; snapshot regeneration script produces matching format |
| `go run` adds latency to `regen:api` (~1-2s for compilation) | Acceptable — previous flow had ~500ms HTTP fetch + backend must be running anyway |
| Moving types to `backend/api/` breaks imports in test files | Test files already import from `package main` (same package); they'll need `api.` prefix updates — one-time mechanical change |
