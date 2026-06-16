## MODIFIED Requirements

### Requirement: The frontend SHALL provide a mutation composable for updating the planning window

The frontend SHALL define `useUpdatePlanningWindow()` in `frontend/src/shared/composables/useUpdatePlanningWindow.ts`. The composable SHALL use the generated Pinia Colada mutation for `putPlanningWindow` (from `@/client/@pinia/colada.gen`). On successful mutation, the composable SHALL invalidate the `getPlanningWindow` query key via `useQueryCache().invalidateQueries` so that all `usePlanningWindow()` consumers automatically refetch. The composable SHALL return `{ mutate, isPending, error }`. The `isPending` return value SHALL be `false` when no mutation is in flight and SHALL only become `true` while the mutation HTTP request is executing.

#### Scenario: Mutation composable invalidates query on success

- **WHEN** the mutation is called with valid `{startDate, endDate}` and the backend responds 200
- **THEN** the `getPlanningWindow` query key is invalidated
- **AND** subsequent reads from `usePlanningWindow()` return the updated planning window data

#### Scenario: Mutation composable exposes error state

- **WHEN** the mutation is called and the backend responds with an error
- **THEN** `error.value` is non-null

#### Scenario: isPending is false when no mutation is in flight

- **WHEN** the composable is first instantiated (before any mutation call)
- **THEN** `isPending.value` is `false`

#### Scenario: isPending is true during mutation execution

- **WHEN** the mutation is called and the HTTP request is in-flight
- **THEN** `isPending.value` is `true`

#### Scenario: isPending returns to false after mutation completes

- **WHEN** the mutation request succeeds
- **THEN** `isPending.value` returns to `false`
