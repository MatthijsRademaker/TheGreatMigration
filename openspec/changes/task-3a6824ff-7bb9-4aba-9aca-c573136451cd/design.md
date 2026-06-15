## Context

The current planning-window source of truth lives entirely in the frontend as hard-coded constants in `frontend/src/shared/lib/planWindow.ts` (`PLAN_WINDOW_START = "2026-07-05"`, `PLAN_WINDOW_END = "2026-08-13"`). The module explicitly documents that changing the window requires a rebuild and foreshadows runtime reconfiguration. The task backlog states this should have been a backend-backed endpoint from the start.

The existing codebase already demonstrates the full pattern for backend-backed features:
- **Backend**: Huma v2 handler registration in `backend/dashboard.go` provides a template for request/response types, OpenAPI metadata, and handler wiring.
- **Frontend client**: `frontend/src/client/@pinia/colada.gen.ts` shows the auto-generated Pinia Colada query pattern for `GET /api/hello`.
- **Frontend consumption**: `HomeView.vue` demonstrates loading, error, and success state handling with Pinia Colada queries.
- **Artifact refresh**: `frontend/scripts/refresh-openapi-snapshot.mjs` and `frontend/scripts/generate-api.mjs` define the established snapshot+codegen workflow.

The canonical spec at `openspec/specs/planning-window/spec.md` currently enshrines the frontend TypeScript module as the source of truth. This change produces a spec delta that replaces that requirement with a backend-backed contract.

The `frontend/src/shared/composables/` directory exists but is empty, ready for the new composable.

## Goals

1. Add a `GET /api/planning-window` backend endpoint following the existing Huma v2 pattern, returning `{startDate, endDate, days}` with validated ISO 8601 dates.
2. Refactor `planWindow.ts` to remove hard-coded constants while preserving `formatPlanDayLabel` and `PlanWindowDay` as pure shared utilities.
3. Create a shared `usePlanningWindow()` composable that wraps the generated Pinia Colada query, derives `planWindowDays` and `planWindowDayCount` as reactive computed refs, and exposes `isLoading`/`error` states.
4. Update `HomeView.vue` and `CalendarView.vue` to consume the composable with proper loading and error handling.
5. Refresh the frontend OpenAPI snapshot and regenerate client artifacts so the new endpoint is available through the generated client.
6. Add backend contract tests for the planning-window endpoint and update frontend tests for the composable-based approach.
7. Produce OpenSpec change artifacts (proposal, design, tasks, spec delta) describing the backend-backed contract.

## Non-Goals

- Building a UI for editing the planning window.
- Adding persistence, authentication, or multi-user configuration management.
- Reworking PeopleView placeholder availability into real per-date CRUD.
- Changing the dashboard people-availability feature beyond what is required to coexist with the new planning-window contract.
- Environment-variable configuration of the planning window in this slice (the backend uses compile-time Go constants matching the current defaults).
- Amending the canonical `hello-world-integration` spec — the existing Move days scenario already references planning-window derivation and does not conflict with the backend-backed approach.

## Decisions

### 1. Dedicated endpoint: `GET /api/planning-window`

Add a standalone endpoint rather than folding the planning range into an existing dashboard contract. The planning window is a global app-level concept (consumed by CalendarView, HomeView, and future views), not a dashboard-specific concern. It takes no query parameters and returns a single global range, keeping the endpoint simple and the contract clear.

**Rationale**: Follows the established Huma pattern, respects separation of concerns, and avoids coupling the planning window to dashboard-specific query parameters (start, days, selectedDate).

### 2. Response shape: `{startDate, endDate, days}`

Return `startDate` and `endDate` as ISO 8601 date strings plus `days` (inclusive count), mirroring the date-range pattern already used by `planWindow.ts`. The backend returns only the range metadata — day-list derivation (`PlanWindowDay[]`) remains in the frontend composable.

**Rationale**: Keeps backend responses small and derivation logic close to the views that consume rendered day objects. The frontend composable can derive the day list once and share it reactively across consumers.

### 3. Shared composable: `usePlanningWindow()`

Create a single composable in `frontend/src/shared/composables/usePlanningWindow.ts` that:
- Calls the generated `getPlanningWindowQuery()` from `@/client/@pinia/colada.gen`.
- Derives `planWindowDays: ComputedRef<PlanWindowDay[]>` and `planWindowDayCount: ComputedRef<number>` from the fetched start/end dates.
- Exposes `isLoading` and `error` reactive states.
- Uses a shared Pinia Colada query key so that HomeView and CalendarView mounting simultaneously do not trigger duplicate fetches.

**Rationale**: Prevents duplicate derivation logic and duplicate network requests. Follows the exact pattern established by the generated `getHelloQuery` in `colada.gen.ts`.

### 4. Error strategy: Fail visible, no silent fallback

When the backend endpoint is unreachable (5xx, network failure) or returns invalid data (start >= end, non-ISO dates, missing fields), both CalendarView and HomeView display a visible error state rather than silently falling back to stale defaults.

**Rationale**: Consistent with the existing HomeView hello-card pattern ("Backend unavailable"). Silent fallback risks masking backend misconfiguration in production. The backend validates its own response (the handler constructs deterministic, valid dates from Go constants), so invalid responses indicate a real deployment problem that should be surfaced.

### 5. `planWindow.ts` refactoring scope

Remove the four hard-coded exports (`PLAN_WINDOW_START`, `PLAN_WINDOW_END`, `planWindowDays`, `planWindowDayCount`) and the internal `generatePlanWindowDays()` function. Retain `formatPlanDayLabel()` and the `PlanWindowDay` interface as pure utility exports. The module-level JSDoc comment is updated to document the new composable-based consumer pattern.

**Rationale**: The acceptance criteria require that consumers stop treating frontend constants as canonical. Retaining the constants alongside a composable would create ambiguity ("import the wrong source"). The pure utilities have no backend dependency and remain valuable for label formatting.

### 6. OpenAPI snapshot and client regeneration

The snapshot refresh (`refresh-openapi-snapshot.mjs`) and code generation (`generate-api.mjs`) are run after the backend endpoint is registered. Both scripts are idempotent — they follow the established workflow and produce committed client artifacts. The generated Pinia Colada query for the new endpoint will follow the same pattern as `getHelloQuery`.

**Rationale**: Mechanical step required by the acceptance criteria. The existing scripts are proven and require no modification.

### 7. CalendarView loading skeleton

CalendarView currently has no loading state — it iterates `planWindowDays` synchronously. With async data, the view must handle the pending state. A minimal loading skeleton (e.g., placeholder day columns matching the responsive grid) prevents layout shift when the backend fetch resolves.

**Rationale**: Identified as a risk by the architect ("CalendarView currently renders day columns without any loading state"). The skeleton uses the same responsive grid classes (`grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7`) for visual stability.

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| **Canonical spec conflict**: `openspec/specs/planning-window/spec.md` Requirement 1 mandates frontend constants as the source of truth, contradicting this change. | High | This change produces a spec delta that replaces Requirement 1 with a backend-backed contract. The canonical spec must be archived and updated in a follow-up archive cycle. |
| **Test breakage**: `planWindow.test.ts` asserts hard-coded constant values and `app-routes-render.test.ts` imports `planWindowDayCount` eagerly. | Medium | Tests are refactored to mock the fetch and test composable success/loading/error states. The route render test is updated to match the new async-aware rendering. |
| **Layout shift**: CalendarView currently renders day columns synchronously at import time. With async data, columns appear after fetch resolves. | Low | A loading skeleton with matching responsive grid classes is added to CalendarView. |
| **Snapshot diff size**: The committed `openapi-snapshot.json` is stale — it lacks the dashboard/people-availability endpoint. After regeneration, the diff includes multiple endpoint additions beyond the planning window. | Low | Expected and acceptable. The snapshot becomes more complete, not less. |
| **Timezone ambiguity**: The dashboard endpoint uses `time.Now()` in server-local timezone. If new planning-window handler accidentally uses server-local time for seed dates, it will differ from the UTC expectations in `formatPlanDayLabel`. | Low | The planning-window handler uses Go constants (`"2026-07-05"`, `"2026-08-13"`) parsed as dates, never `time.Now()`. The backend does not perform timezone-dependent computation for the planning window. |

## Conflict Resolution

No conflicting decisions among refinement participants. All three reviewers (architect, lead-dev, reviewer) independently recommended the same approach: dedicated endpoint, shared composable, preserved pure utilities, and error-on-failure strategy. The reviewer's three blockers (B1: spec conflict, B2: error strategy, B3: async consumer patterns) are addressed as explicit design decisions above rather than left as implicit assumptions.

## Traceability

- Task: `3a6824ff-7bb9-4aba-9aca-c573136451cd`
- Dossier: `2026-06-15T18:19:30.948Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial` snapshot for `task-3a6824ff-7bb9-4aba-9aca-c573136451cd`