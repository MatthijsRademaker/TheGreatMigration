## Context

The household organizer needs to control the canonical move timeline from /settings. The repository already has:

- A `planning_windows` database table with a singleton row (migration `001_create_tables.sql`)
- A `GET /api/planning-window` endpoint backed by a `Store` interface and sqlc query
- A `usePlanningWindow()` Pinia Colada composable with cached reads
- A shared `DatePicker` component (single-date, Reka UI `CalendarRoot`)
- A `PlanningWindowBodyWritable` type in the generated frontend client

The gap: no write path on the backend, no mutation composable, and a SettingsView that is spec-constrained to placeholder-only content.

## Goals / Non-Goals

### Goals

1. Let the organizer view and edit the move start and end dates from /settings.
2. Persist the chosen range as the canonical planning window via `PUT /api/planning-window`.
3. Refresh the settings UI after save without a full page reload (query invalidation).
4. Reuse existing infrastructure: `Store` interface, sqlc queries, `usePlanningWindow()`, `DatePicker`, `Card`.
5. Update canonical OpenSpec specs (planning-window and sidebar-navigation) to reflect the new capabilities.

### Non-Goals

- Do not build the broader future settings areas (notifications, account details, team visibility).
- Do not rework calendar/people/task views to become fully date-driven.
- Do not introduce general availability CRUD or schedule editing.
- Do not extend seed data or fix downstream availability for arbitrary date ranges.
- Do not build a date range picker component; the two single-date `DatePicker`s are sufficient.

## Decisions

### D1: Use PUT (not PATCH) for the singleton planning-window resource

**Decision**: The write endpoint is `PUT /api/planning-window` with full resource replacement semantics. Both `startDate` and `endDate` are required in the request body.

**Rationale**: The planning window is a singleton with exactly two fields. PUT semantics (idempotent full replacement) are simpler than PATCH (partial update). The existing `PlanningWindowBodyWritable` type already models a PUT-style body.

**Sources**: swarm-architect (round 1), swarm-lead-dev (round 1), accepted decisions `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`.

### D2: Accept any valid range (endDate >= startDate) without seed-data constraints

**Decision**: The backend validates only that `endDate >= startDate` (and both are valid ISO 8601 dates). It does not reject dates outside the seeded 2026-07-05..2026-08-13 availability range.

**Rationale**: Downstream endpoints (`people-availability`, `daily-schedule`) already default missing availability rows to `"off"`. No endpoint crashes; the demo experience degrades gracefully. Blocking writes would couple this task to full availability CRUD, which is out of scope.

**Sources**: swarm-architect (round 1), swarm-lead-dev (round 1), accepted decisions `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`.

### D3: Use an UPSERT SQL pattern for robustness

**Decision**: The sqlc query uses `INSERT ... ON CONFLICT ... DO UPDATE` (UPSERT) rather than a plain `UPDATE`.

**Rationale**: The `planning_windows` table has no unique constraint on a singleton marker column; the schema uses `id SERIAL PRIMARY KEY`. An UPSERT on a known sentinel `id = 1` is robust even if the singleton row is missing, which could happen in a fresh database without seed data.

**Sources**: swarm-architect (round 1), accepted decision `1-swarm-architect-recommendation`.

### D4: Invalidate the planning-window query after successful save

**Decision**: The frontend mutation composable invalidates the `getPlanningWindow` query key after a successful PUT, causing `usePlanningWindow()` consumers to refetch.

**Rationale**: Pinia Colada's `useQueryCache().invalidate` is the standard pattern for write-then-refresh. The existing composable already exposes `queryKey`.

**Sources**: swarm-architect (round 1), accepted decision `1-swarm-architect-recommendation`.

### D5: Regenerate OpenAPI snapshot and client as part of backend work

**Decision**: After adding the PUT endpoint to the backend, regenerate `openapi-snapshot.json` and run `npm run generate:api` before writing frontend mutation code.

**Rationale**: The frontend generated client must include the new mutation function before the SettingsView can call it. This is a mechanical sequencing constraint, not an architectural choice.

**Sources**: swarm-lead-dev (round 1), accepted decision `1-swarm-lead-dev-recommendation`.

### D6: Extend all four Store implementations in lockstep

**Decision**: Adding `UpdatePlanningWindow` to the `Store` interface requires implementing it on `PgStore`, `mockStore`, `failingStore`, and `partialFailingStore`.

**Rationale**: All four types implement the full `Store` interface. The compiler enforces this; missing implementations cause build failures.

**Sources**: swarm-lead-dev (round 1), accepted decision `1-swarm-lead-dev-recommendation`.

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| OpenAPI/frontend-client generation pipeline: backend must be running during snapshot refresh | Medium | Documented step; `npm run refresh:openapi-snapshot` requires a running backend |
| Seed data mismatch: dates outside 2026-07-05..2026-08-13 cause sparse availability data | Medium | Accepted; downstream defaults to `"off"`; SettingsView card description can note expected range |
| DatePicker type conversion: ISO strings ↔ `DateValue` (@internationalized/date) conversion required in SettingsView | Medium | Straightforward; use `CalendarDate` from `@internationalized/date`; tested |
| Store interface expansion requires four implementations | Low | Mechanical; compiler-enforced |

## Traceability

- **Task**: `8d593ef6-a75a-42bf-96ed-39f3bd74b9af` — "Add date settings"
- **Dossier**: `2026-06-16T07:32:40.863Z` — exploration dossier
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Rounds**: Round 1 agents: swarm-architect (`019ecf58-d59a-74eb-bff9-b0b1af61c8bc`), swarm-lead-dev (`019ecf5a-60c2-798e-8983-efdb235c3e5b`), swarm-reviewer (`019ecf5b-d9f7-7a13-bb50-6c756a5e2970`)
- **Evidence**: backend/planning_window.go, backend/store.go, backend/queries/planning_window.sql, backend/store_mock_test.go, backend/main_test.go, frontend/src/settings/SettingsView.vue, frontend/src/shared/composables/usePlanningWindow.ts, frontend/tests/app-routes-render.test.ts, openspec/specs/planning-window/spec.md, openspec/specs/sidebar-navigation/spec.md
