## MODIFIED Requirements

### Requirement: `/tasks` SHALL own the minimal CRUD interactions and refresh from the canonical backlog query

The `/tasks` route SHALL provide add, edit, delete, and assignment update flows using the typed generated frontend client for the task write endpoints. Create and edit interactions SHALL use the shared add-operation modal, `Select`, `Input`, `Checkbox`, and `Button` primitives for a focused form instead of the existing shared `Sheet`. The Room / Area field SHALL use a `Select` populated from the `listRoomsQuery` (`GET /api/rooms`) instead of a free-form text `Input`. While the rooms query is pending, the room Select SHALL show a loading placeholder. If the rooms query fails, the room Select SHALL display an error message with a retry control. `TasksView.vue` SHALL continue to own task form state, mutation handlers, query invalidation, assignment toggling, and modal open-state control. The add and edit flows SHALL both use explicit `v-model:open` state management, with the Add Task handler opening the modal directly instead of relying on trigger timing. After each successful create or update mutation, the task backlog query SHALL be invalidated or refetched so the rendered rows, summary-derived information, and legends reflect backend-confirmed state, and the modal SHALL close only after route-owned reset logic runs. Failed create or update mutations SHALL keep the modal open and preserve user-entered values while showing the error inside the modal body. Delete behavior SHALL remain outside the add-operation modal and SHALL continue to refresh from `GET /api/tasks/backlog` after success.

#### Scenario: Room select shows managed rooms on success

- **WHEN** a user opens the task create or edit modal and the rooms query resolves
- **THEN** the Room / Area field renders as a Select populated with room names from the backend response
- **AND** the user can select a room from the dropdown instead of typing free-form text

#### Scenario: Room select shows loading state while rooms fetch

- **WHEN** a user opens the task modal and the rooms query is pending
- **THEN** the Room / Area field shows a loading placeholder
- **AND** the rest of the form remains interactive

#### Scenario: Room select shows error state on failure

- **WHEN** the rooms query fails
- **THEN** the Room / Area field displays an error message indicating rooms could not be loaded
- **AND** the field provides a retry control for the rooms query

#### Scenario: Selected room name is submitted as room value

- **WHEN** a user selects a room from the dropdown and submits the form
- **THEN** the submitted `room` value is the selected room's name string

#### Scenario: Add Task opens the reusable modal and creates a task

- **WHEN** a user activates `Add Task` on `/tasks`
- **THEN** the route opens the shared add-operation modal with empty task form values and add-oriented title/description/submit labeling
- **AND** submitting valid data calls the create mutation using the same trimmed task payload fields as before

#### Scenario: Edit opens the reusable modal with prefilled values

- **WHEN** a user activates `Edit` for an existing task on `/tasks`
- **THEN** the route opens the shared add-operation modal with the current task values prefilled
- **AND** submitting the form calls the update mutation for that task ID

#### Scenario: Successful create or update refreshes and closes the modal

- **WHEN** a task create or update mutation succeeds
- **THEN** the route invalidates or refetches the canonical backlog query
- **AND** the route resets task form state
- **AND** the modal closes after the success path completes

#### Scenario: Failed create or update keeps user state visible

- **WHEN** a task create or update mutation fails
- **THEN** the modal remains open
- **AND** the user-entered task values remain populated
- **AND** the route displays the mutation error inside the modal body

#### Scenario: Cancel or close resets route-owned task form state

- **WHEN** a user cancels or closes the modal from either add or edit mode
- **THEN** the route resets task form state consistently
- **AND** the next add flow starts from a clean form

#### Scenario: Delete remains outside the add-operation modal

- **WHEN** a user deletes a task from `/tasks`
- **THEN** the route performs the existing delete mutation path without routing that action through the add-operation modal
- **AND** successful delete continues to refresh from `GET /api/tasks/backlog`
