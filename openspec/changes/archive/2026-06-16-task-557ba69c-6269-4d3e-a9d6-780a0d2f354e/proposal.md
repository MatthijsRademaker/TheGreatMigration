## Why

The SettingsView date-picker Save button permanently shows "Saving..." text and never transitions back to "Save". The button is stuck in a perpetual saving-indicator state because the `useUpdatePlanningWindow` composable derives its pending signal from the mutation's `status` ref, which starts as `"pending"` before any mutation is ever triggered.

## What Changes

- **Fix `useUpdatePlanningWindow` composable**: Replace `isPending` derivation from `status` (data state) with `isLoading` (async state) so that the pending indicator is only `true` while an HTTP request is actually in flight.
- **Update `SettingsView.vue`**: (Minor) Confirm template uses the corrected `isPending` signal correctly.
- **Update `planning-window` spec**: Clarify the contract for `isPending` semantics in the mutation composable requirement.
- **Update `settings-view.test.ts`**: Add a test that verifies `isPending` transitions correctly during a mutation lifecycle (not just the mock being set to a static value).

## Capabilities

### New Capabilities

- *(none — this is a pure bug fix within existing capabilities)*

### Modified Capabilities

- `planning-window`: The `useUpdatePlanningWindow` composable requirement's `isPending` semantics SHALL be clarified — `isPending` SHALL be `false` when no mutation is in flight and SHALL only become `true` during the HTTP request. This corrects the current behavior where `isPending` is `true` from composable mount time.

## Impact

- **Frontend**: `frontend/src/shared/composables/useUpdatePlanningWindow.ts` — one-line change (use `isLoading` instead of `status`-derived `isPending`)
- **Tests**: `frontend/tests/settings-view.test.ts` — add mutation lifecycle test; mock `isLoading` instead of `isPending`
- **Specs**: `openspec/specs/planning-window/spec.md` — clarify `isPending` semantics in the mutation composable requirement
- No backend changes, no API changes, no dependency changes
