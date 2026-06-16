## Context

The SettingsView at `frontend/src/settings/SettingsView.vue` provides a date-picker form for updating the planning window (start date, end date). Users report that the Save button permanently shows "Saving..." text and never transitions back to "Save", making the form appear stuck.

Investigation traced the issue to the `useUpdatePlanningWindow` composable at `frontend/src/shared/composables/useUpdatePlanningWindow.ts`. The composable derives `isPending` from the Pinia Colada `useMutation()` return value `status` ref, which initializes as `"pending"` — the data state "the mutation has not completed" — rather than the async state "a mutation is currently in flight."

The frontend uses Pinia Colada v1.3.1 for data fetching and mutations. The `useMutation()` composable provides two distinct signals:

| Signal | Initial value | Purpose |
|---|---|---|
| `status` | `"pending"` | Data state — whether the mutation has ever completed |
| `asyncStatus` | `"idle"` | Async state — whether a mutation request is in flight |
| `isLoading` | `false` | Convenience computed for `asyncStatus === "loading"` |

The existing code uses `status`, which is the wrong signal for "is saving in progress."

## Goals / Non-Goals

**Goals:**

- Fix the Save button indicator so it shows "Save" when idle and "Saving..." only during an actual HTTP request
- Add a test that validates the pending signal lifecycle

**Non-Goals:**

- No backend changes — the PUT endpoint, store, and validation are all correct
- No visual or layout changes to the SettingsView beyond the button text fix
- No changes to the backend OpenAPI contract or type generation

## Decisions

### Decision: Use `isLoading` instead of `status`-derived `isPending`

The `useMutation` composable from Pinia Colada exposes an `isLoading` computed ref that is `true` only while the mutation function is executing (between `asyncStatus` transitioning to `"loading"` and back to `"idle"`). This maps exactly to the UX requirement.

**Alternatives considered:**

- **Use `asyncStatus` directly**: Equivalent but requires an additional `computed(() => asyncStatus.value === "loading")`. `isLoading` is the idiomatic API provided by the library.
- **Use a local `saving` ref**: Would work but duplicates state that the mutation library already manages correctly. Adds unnecessary complexity.
- **Use `status` with initial-value check**: Possible but fragile — would need to distinguish "never called" from "in flight". The library already provides the right abstraction.

### Decision: Pinia Colada's `isLoading` semantics are correct for this use case

The `mutate()` function in Pinia Colada sets `entry.asyncStatus.value = "loading"` before calling the mutation function and resets it to `"idle"` in the `finally` block. This means `isLoading` (and `asyncStatus === "loading"`) is `true` while waiting for the HTTP response, which is exactly the UX we want.

### Decision: Update the mock in `settings-view.test.ts` to expose `isLoading`

The existing mock for `useUpdatePlanningWindow` returns a `controlledIsPending` ref. Since the fix uses `isLoading` instead of the derived `isPending`, the mock should provide `isLoading` directly. The composable return type changes from `{ mutate, isPending, error }` to using `isLoading` internally, but the public API surface (return values) stays the same since we rename `isLoading` to `isPending` in the destructure:

```ts
const { mutate, isLoading: isPending, error } = useMutation({ ... });
```

## Risks / Trade-offs

- **[Low] The `isLoading`/`asyncStatus` behavior could differ across Pinia Colada versions**: Pinia Colada 1.3.1 is the current version and the behavior is well-defined in the source (see `mutate()` function in the mutation cache). If an upgrade changes `asyncStatus` semantics, a test will catch the regression.
- **[None] No breaking changes**: The public API of the composable (`{ mutate, isPending, error }`) is unchanged — only the internal derivation changes.
