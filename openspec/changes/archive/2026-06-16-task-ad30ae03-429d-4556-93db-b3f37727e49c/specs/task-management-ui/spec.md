## ADDED Requirements

### Requirement: Shared add-operation modal SHALL provide a reusable centered CRUD shell

The frontend SHALL provide a shared add-operation modal for task create/edit flows and future add-operation callers. The shell SHALL:

- support `v-model:open`
- render an accessible title and description
- accept caller-provided body content through a default slot
- provide default submit and cancel actions with configurable labels, disabled state, and submitting state
- allow callers to override the default footer through a named footer slot
- emit open-state changes, submit actions, and cancel actions back to caller-owned handlers
- size to approximately 60% of viewport width and height on `sm` and above with a stable header/footer and a scrollable body region
- fill the viewport on mobile widths below `sm`

Domain-specific validation and mutation errors SHALL remain caller-rendered content inside the modal body rather than modal-owned behavior.

#### Scenario: Desktop modal sizing keeps long forms usable

- **WHEN** the modal is opened on `sm` and larger viewports
- **THEN** it renders as a centered modal sized to roughly 60% of viewport width and height
- **AND** the body region scrolls independently so the header and footer remain accessible

#### Scenario: Mobile modal fills the viewport

- **WHEN** the modal is opened below the `sm` breakpoint
- **THEN** it fills the viewport rather than rendering as a narrow centered dialog

#### Scenario: Submit, cancel, and open-state changes propagate to callers

- **WHEN** a caller opens the modal and a user activates submit, cancel, or close interactions
- **THEN** the shell emits those interactions back to the caller-owned handlers without branching on task-specific behavior

#### Scenario: Callers render domain-specific errors inside the modal body

- **WHEN** a task create or update mutation fails while the modal is open
- **THEN** the caller can render the error inside the modal body content
- **AND** the shared modal does not need task-specific error props or task-specific error formatting

## MODIFIED Requirements

### Requirement: `/tasks` SHALL own the minimal CRUD interactions and refresh from the canonical backlog query

The `/tasks` route SHALL provide add, edit, delete, and assignment update flows using the typed generated frontend client for the task write endpoints. Create and edit interactions SHALL use the shared add-operation modal, `Input`, `Select`, `Checkbox`, and `Button` primitives for a focused form instead of the existing shared `Sheet`. `TasksView.vue` SHALL continue to own task form state, mutation handlers, query invalidation, assignment toggling, and modal open-state control. The add and edit flows SHALL both use explicit `v-model:open` state management, with the Add Task handler opening the modal directly instead of relying on trigger timing. After each successful create or update mutation, the task backlog query SHALL be invalidated or refetched so the rendered rows, summary-derived information, and legends reflect backend-confirmed state, and the modal SHALL close only after route-owned reset logic runs. Failed create or update mutations SHALL keep the modal open and preserve user-entered values while showing the error inside the modal body. Delete behavior SHALL remain outside the add-operation modal and SHALL continue to refresh from `GET /api/tasks/backlog` after success.

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
