## Why

The /settings route is still a placeholder, but the household organizer needs to control the canonical move timeline. The repo already has a persisted `planning_windows` table, a `GET /api/planning-window` endpoint, a cached `usePlanningWindow()` composable, and a shared `DatePicker`, but there is no write path, no settings form, and the current specs require SettingsView to remain placeholder-only.

This change adds a date range configurator to /settings, enabling the organizer to view and edit the move start and end dates. The backend write path persists the singleton planning-window row; the frontend form uses the existing composable and DatePicker infrastructure. The existing GET contract remains unchanged, and all seed-data/downstream concerns are addressed with explicit decisions.

## What Changes

### Backend

- **New `PUT /api/planning-window` endpoint** registered via Huma with `operationId: "put-planning-window"`. Accepts `{startDate, endDate}` (ISO 8601 strings), validates `endDate >= startDate`, persists the singleton planning-window row via an UPSERT query, and returns the updated `PlanningWindowBody` with recalculated `days`.
- **New `UpdatePlanningWindow(ctx, startDate, endDate)` method** on the `Store` interface, implemented on `PgStore`, `mockStore`, `failingStore`, and `partialFailingStore`.
- **New `UpsertPlanningWindow` sqlc query** in `backend/queries/planning_window.sql` using `INSERT ... ON CONFLICT ... DO UPDATE`.
- **Backend unit tests** covering PUT success, validation failure (`endDate < startDate`, malformed dates), and store-failure.
- **Backend integration tests** verifying persistence round-trip.

### Frontend

- **New `useUpdatePlanningWindow()` mutation composable** in `frontend/src/shared/composables/` that calls `PUT /api/planning-window` via the generated Pinia Colada SDK and invalidates the `getPlanningWindow` query key on success.
- **SettingsView replaced** — the placeholder "Feature coming soon" card is replaced with an interactive "Planning window" card containing:
  - A summary of the current planning window range loaded from `usePlanningWindow()`.
  - Two `DatePicker` controls for start and end dates, prefilled from the composable.
  - Inline validation blocking save when `endDate < startDate` or dates are unchanged.
  - Save and reset buttons with loading/error feedback via the mutation composable.
- **OpenAPI snapshot regenerated** (`frontend/openapi-snapshot.json`) and **frontend client regenerated** (`npm run generate:api`) after the backend write endpoint is added.

### Specs & Tests

- **`openspec/specs/planning-window/spec.md`**: Extended with ADDED requirements for the PUT write contract (operation ID, request/response schemas, validation rules, error responses).
- **`openspec/specs/sidebar-navigation/spec.md`**: MODIFIED to remove the placeholder-only restriction on SettingsView and reference the planning-window spec for its content contract.
- **Route render tests** updated: `/settings` case replaces "Feature coming soon" assertion with assertions for date pickers, range summary, and save button.
- **New frontend tests** covering SettingsView form behavior: prefill from composable, validation, save, error display, and post-save refresh.

## Impact

### Affected specs

- `openspec/specs/planning-window/spec.md` — ADDED: PUT write-endpoint requirements.
- `openspec/specs/sidebar-navigation/spec.md` — MODIFIED: SettingsView is no longer placeholder-only; references planning-window spec.

### Affected code

- **Backend**: `planning_window.go` (new PUT handler), `store.go` (new Store method + PgStore impl), `store_mock_test.go` (mockStore + failingStore + partialFailingStore updates), `queries/planning_window.sql` (new UPSERT query), `db/planning_window.sql.go` (sqlc regenerated), `main_test.go` (new PUT tests), `main_integration_test.go` (new integration tests).
- **Frontend**: `settings/SettingsView.vue` (replaced placeholder), new `composables/useUpdatePlanningWindow.ts`, `openapi-snapshot.json` (regenerated), `src/client/` (regenerated), `tests/app-routes-render.test.ts` (updated /settings assertions), new `tests/settings-view.test.ts`.

### No breaking changes

- `GET /api/planning-window` contract is unchanged.
- Existing route-render tests are updated, not removed.
- No new database tables or migrations; the existing `planning_windows` table is reused.
- No new shared UI components; `DatePicker` and `Card` are reused.

### Risk: Downstream demo data degradation

Saving dates outside the seeded 2026-07-05..2026-08-13 range causes downstream availability/schedule endpoints to return sparse `"off"`-default data. This is accepted for the demo — the backend already handles missing availability rows gracefully and no endpoint crashes. A note is included in the SettingsView card description.