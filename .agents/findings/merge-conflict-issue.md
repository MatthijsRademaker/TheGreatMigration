# Merge Conflict Postmortem: OpenAPI Snapshot & Generated Client Drift

**Date:** 2026-06-16
**Severity:** High — runtime crashes in production (Vue Router errors)
**Status:** Resolved via build-time generation pipeline

---

## Timeline

| Time (CEST) | Event |
|---|---|
| 07:01 | PR #27: People availability (read-only dashboard endpoint) merged |
| 08:04–08:49 | PR #37 (`task-ed86ff0f`): **Rooms and areas fully functional** — adds `/api/rooms` CRUD, frontend rooms view, OpenAPI snapshot with 7 paths |
| 08:11–08:33 | PR #38: **Add date settings** — adds `PUT /api/planning-window`, `SettingsView.vue`, regenerates client from snapshot with 5 paths (no rooms, no people CRUD) |
| 09:43 | PR #35: Header merged |
| 10:33 | PR #36: **Wire people availability to bff** — adds `/api/people` CRUD, `PeopleView.vue`, OpenAPI snapshot with 7 paths (no rooms) |
| 11:06 | **Merge #38** (date settings) into main (on top of people #36). Auto-generated client files: accepted from main. |
| 11:12 | **Merge #37** (rooms `task-ed86ff0f`) into main (on top of date settings #38). Commit message: *"Auto-generated client files: accepted from main"* |
| 11:12+ | `docker compose up -d` → Vue Router crashes: `deletePersonAvailabilityMutation`, `updateRoomMutation` not exported |

## What Broke

The three PRs each independently generated their own complete `openapi-snapshot.json` and client files (`sdk.gen.ts`, `colada.gen.ts`, `types.gen.ts`):

| Branch | OpenAPI snapshot paths |
|---|---|
| **People** (PR #36, `c9afef2`) | `/api/people`, `/api/people/{id}`, `/api/people/{id}/availability/{date}`, `/api/dashboard/people-availability`, `/api/hello`, `/api/planning-window`, `/api/tasks/backlog` |
| **Rooms** (PR #37, `244bd06`) | `/api/rooms`, `/api/rooms/{id}`, `/api/dashboard/daily-schedule`, `/api/dashboard/people-availability`, `/api/hello`, `/api/planning-window`, `/api/tasks/backlog` |
| **Date Settings** (PR #38, `927ae58`) | `/api/dashboard/daily-schedule`, `/api/dashboard/people-availability`, `/api/hello`, `/api/planning-window`, `/api/tasks/backlog` |

**Merge order:** People (#36) → Date Settings (#38) → Rooms (#37)

When Rooms (#37) was merged last, the merge-conflict resolution on auto-generated client files chose **"accept from main"**. At that point, `main` had only the date-settings snapshot (5 paths — missing both people CRUD and rooms CRUD). This produced a committed `openapi-snapshot.json` with only:

```
/api/dashboard/daily-schedule
/api/dashboard/people-availability
/api/hello
/api/planning-window
/api/tasks/backlog
```

**Missing:** `/api/people`, `/api/people/{id}`, `/api/people/{id}/availability/{date}`, `/api/rooms`, `/api/rooms/{id}`

The `@hey-api/openapi-ts` generator then produced `colada.gen.ts` with only 6 exports (queries + mutations for the 5 remaining paths). Vue components `PeopleView.vue` and `RoomsView.vue` — which were merged successfully in the same commits — still imported now-nonexistent exports:

```typescript
// PeopleView.vue imports — NOT in the stripped colada.gen.ts
createPersonMutation
deletePersonMutation
deletePersonAvailabilityMutation
upsertPersonAvailabilityMutation

// RoomsView.vue imports — NOT in the stripped colada.gen.ts
listRoomsQuery
createRoomMutation
updateRoomMutation
deleteRoomMutation
```

### Why Precommit Didn't Catch It

The precommit runs `vue-tsc --noEmit` inside a Docker container. TypeScript DOES detect these missing exports — running `vue-tsc --noEmit` locally produces:

```
src/people/PeopleView.vue(11,3): error TS2305: Module has no exported member 'createPersonMutation'.
src/people/PeopleView.vue(13,3): error TS2305: Module has no exported member 'deletePersonAvailabilityMutation'.
src/rooms/RoomsView.vue(17,3): error TS2305: Module has no exported member 'updateRoomMutation'.
... (8 more)
```

However, the Docker-based precommit runs `npm install --package-lock=false` which may fail to install `@pinia/colada` (the first error in the diagnostic chain is `Cannot find module '@pinia/colada'`). When the generated file itself can't be resolved, `vue-tsc` may bail before reporting missing exports from consumer files. The merge was likely completed outside the precommit workflow (local agent resolution of conflicts).

### Additional Merge Damage in Test File

`backend/main_test.go` suffered extensive merge damage beyond the snapshot issue:

- **Interleaved test bodies:** `TestCreatePerson` body merged with room test assertions (checking for `rooms[0].Name` after a `POST /api/people` request)
- **Line-level mangling:** Multiple lines had two statements joined by tabs (`req := httptest.NewRequest(...)\tcontentType := rec.Header()...`)
- **Duplicate test declarations:** `TestCreateRoom` declared twice (lines 1111 and 1140)
- **Orphaned function fragments:** `partialFailingStore` room CRUD methods interleaved into `TestCreatePerson` body
- **Wrong test expectations:** `TestUpdatePerson` expected 422 for a valid update; `TestCreatePersonMissingFields` expected 201 for empty fields

10 tests fail post-refactor due to this pre-existing damage (all merge artifact, not caused by this change).

## What We Changed

### Root Cause Fix: Build-Time OpenAPI Generation

Instead of fetching the OpenAPI spec from a running backend (fragile, requires DB), we now generate it directly from Go types:

```
Before (fragile):                  After (this change):
┌──────────┐                       ┌──────────────────────┐
│ Backend  │  HTTP fetch           │ go run               │
│ :8080    │ ──────────► snapshot  │ ./cmd/openapi-gen    │
│  + DB    │                       │ (no server, no DB)   │
└──────────┘                       └──────┬───────────────┘
                                          │ stdout
                                          ▼
                                   openapi-snapshot.json
```

**`backend/cmd/openapi-gen/main.go`** — registers all endpoints via `api.RegisterAll(api, nil)`, calls `api.OpenAPI()`, outputs JSON. Huma's `OpenAPI()` introspects registered Go types without invoking handlers, so no database is needed.

### Package Refactor

Extracted all endpoint types and registration functions from `package main` into `backend/api/`:

```
backend/api/
  store.go              — Store interface
  dashboard.go          — Dashboard types + registerDashboardPeopleAvailability
  daily_schedule.go     — Daily schedule types + registerDailySchedule
  people.go             — People CRUD types + registerPeopleEndpoints
  planning_window.go    — Planning window types + registerPlanningWindow
  rooms_areas.go        — Room types + registerRoomsAreas
  tasks.go              — Task backlog types + registerTasksBacklog
  register.go           — RegisterAll(api, store) — single entry point
```

Both `main.go` (server) and `cmd/openapi-gen/main.go` (spec generator) call `api.RegisterAll()`.

### CI Validation

`scripts/check` now includes a backend-phase step:

```bash
go run ./cmd/openapi-gen | diff - ../frontend/openapi-snapshot.json
```

If the committed snapshot diverges from what the backend code produces, CI fails with instructions to run `npm run regen:api`.

### Frontend Script Change

`refresh:openapi-snapshot` changed from HTTP fetch (`fetch("http://localhost:8080/openapi.json")`) to:

```bash
bash ./scripts/refresh-openapi-snapshot.sh   # calls go run ./cmd/openapi-gen
```

## Current State

| Component | Status |
|---|---|
| `backend/api/` package | ✅ Created, compiles, all types extracted |
| `backend/cmd/openapi-gen` | ✅ Works, outputs 10 paths, no DB needed |
| `backend/main.go` | ✅ Uses `api.RegisterAll` |
| `frontend/openapi-snapshot.json` | ✅ Regenerated with all 10 endpoints |
| `frontend/src/client/*.gen.ts` | ✅ Regenerated, all 21 Colada exports present |
| `vue-tsc --noEmit` | ✅ Zero errors |
| `scripts/check` snapshot validation | ✅ Added |
| Backend unit tests | ✅ Pass after repairing merge-damaged `main_test.go` |
| Frontend tests | ✅ 109 tests pass |
| Precommit verification | ✅ `scripts/precommit-run` passes when run against the staged fix set |

The damaged backend tests have been repaired by restoring distinct people and rooms CRUD test bodies, correcting expected Huma validation status codes, and adding merged OpenAPI endpoint coverage.

## Lessons Learned

1. **Auto-generated files should not be resolved as "accept from main" during merges.** These files are derivatives — the correct resolution is to regenerate them from the merged source of truth after conflict resolution.

2. **Committed generated artifacts are a merge-conflict magnet.** The `openapi-snapshot.json` had conflicts in all three concurrent branches. Build-time generation eliminates the committed snapshot as a conflict surface (it's always regenerated from the single source of truth: Go code).

3. **CI must verify generated artifacts match source.** The new `diff` step in `scripts/check` ensures this can never happen silently again.

4. **Docker-based precommit can mask errors.** When `npm install` fails to install dependencies, `vue-tsc` may not fully analyze the codebase, hiding type errors. The precommit should separate install from check and fail fast on install failures.

5. **Merge conflict resolution on test files requires manual review.** The `main_test.go` damage shows that text-level merge resolution on test files (interleaved test bodies, duplicate declarations) produces silently broken tests that pass compilation but test wrong things.

## Related Artifacts

- Change: `openspec/changes/build-time-openapi-generation/`
- Merged PRs: #36 (People), #37 (Rooms, `task-ed86ff0f`), #38 (Date Settings)
- Repaired test file: `backend/main_test.go`
