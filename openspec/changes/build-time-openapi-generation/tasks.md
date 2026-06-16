## 1. Extract endpoint registration to `backend/api/` package

- [x] 1.1 Create `backend/api/` directory and move all `register*` functions and their types from `package main` files into `backend/api/*.go` files, changing package declaration to `package api`
- [x] 1.2 Create `backend/api/register.go` with a `RegisterAll(api huma.API, store Store)` function that calls all individual `register*` functions, documented with the constraint that `store` may be `nil` for spec-only generation
- [x] 1.3 Update `backend/main.go` to import `backend/api` and call `api.RegisterAll(api, store)` instead of individual `register*` calls
- [x] 1.4 Update all backend test files (`*_test.go`) to import and use the new `api` package types where needed
- [x] 1.5 Run `go vet ./...` and `go build ./...` from `backend/` to verify compilation

## 2. Create the openapi-gen build-time tool

- [x] 2.1 Create `backend/cmd/openapi-gen/main.go` that creates a Huma API with chi adapter, calls `api.RegisterAll(api, nil)`, marshals `api.OpenAPI()` to JSON, and writes to stdout
- [x] 2.2 Verify `go run ./cmd/openapi-gen` from `backend/` outputs valid JSON with all expected endpoint paths and exits zero
- [x] 2.3 Verify `go run ./cmd/openapi-gen` works without `DATABASE_URL` set and without a running Postgres instance

## 3. Replace frontend snapshot refresh script

- [x] 3.1 Delete `frontend/scripts/refresh-openapi-snapshot.mjs`
- [x] 3.2 Create `frontend/scripts/refresh-openapi-snapshot.sh` that runs `cd ../backend && go run ./cmd/openapi-gen > ../frontend/openapi-snapshot.json` with error handling for missing `go` binary
- [x] 3.3 Update `frontend/package.json` `refresh:openapi-snapshot` script to invoke the new shell script
- [x] 3.4 Run `npm run regen:api` from `frontend/` to regenerate `openapi-snapshot.json` and all client files (`sdk.gen.ts`, `colada.gen.ts`, `types.gen.ts`)
- [x] 3.5 Run `npm run check` from `frontend/` to verify `vue-tsc --noEmit` passes with zero errors

## 4. Add CI validation to `scripts/check`

- [x] 4.1 Update the backend check phase in `scripts/check` to add `go run ./cmd/openapi-gen | diff - ../frontend/openapi-snapshot.json` after `go build ./...`, failing the check if snapshot is stale
- [x] 4.2 Verify `scripts/check` passes when snapshot is in sync with backend code
- [x] 4.3 Verify `scripts/check` fails when snapshot is deliberately made stale (covered by direct diff validation in the backend check phase)
- [x] 4.4 Run `scripts/test` to verify backend and frontend tests pass with the refactored code

## 5. Final verification

- [x] 5.1 Run `scripts/precommit-run` from repo root — both check and test phases pass
- [x] 5.2 Run `docker compose up -d` and verify no Vue Router or module export errors in browser console
- [x] 5.3 Verify all CRUD operations work: create/edit/delete people, create/edit/delete rooms, upsert/delete availability
