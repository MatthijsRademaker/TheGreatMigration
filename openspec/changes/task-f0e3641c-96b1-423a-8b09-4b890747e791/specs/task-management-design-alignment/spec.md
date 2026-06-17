# task-management-design-alignment Specification

## Purpose

Defines the frontend design-alignment requirements that unify the task backlog under a single authoritative shell (`TaskManagementPanel.vue`), standardize the heading to "Tasks Backlog", surface the total-task count chip, preserve the priority legend footer, and ensure the home dashboard remains read-only while `/tasks` reuses the unified shell with full CRUD capabilities. All changes are frontend-presentational only; no backend, data-contract, or vocabulary changes are introduced.

## ADDED Requirements

### Requirement: TaskManagementPanel heading SHALL be "Tasks Backlog"

The `CardTitle` in `TaskManagementPanel.vue` SHALL render the text "Tasks Backlog" instead of "Task Management". This heading SHALL render identically in both read-only (home dashboard) and editable (`/tasks`) variants.

#### Scenario: Heading renders on home dashboard

- **WHEN** `HomeView.vue` renders `<TaskManagementPanel read-only />`
- **THEN** the rendered output contains "Tasks Backlog"
- **AND** the rendered output does NOT contain "Task Management"

#### Scenario: Heading renders on /tasks route

- **WHEN** `TasksView.vue` renders `<TaskManagementPanel>` (non-readOnly)
- **THEN** the rendered output contains "Tasks Backlog"
- **AND** the rendered output does NOT contain "Task backlog"

### Requirement: TaskManagementPanel SHALL surface a total-task count chip

The panel header SHALL render a compact count indicator adjacent to the heading, displaying `{summary.totalTasks} tasks` (e.g., "12 tasks"). The count chip SHALL use an existing badge variant (e.g., `secondary`) and SHALL render in both read-only and editable variants. The count SHALL derive from `data.summary.totalTasks` as provided by `useTaskBacklog()`.

#### Scenario: Count chip renders on home dashboard

- **WHEN** `GET /api/tasks/backlog` returns `summary.totalTasks: 3` and the home route renders
- **THEN** the rendered output contains "3 tasks" in the panel header area

#### Scenario: Count chip renders on /tasks route

- **WHEN** `GET /api/tasks/backlog` returns `summary.totalTasks: 3` and the `/tasks` route renders
- **THEN** the rendered output contains "3 tasks" in the panel header area

#### Scenario: Count chip updates with data

- **WHEN** the backlog query returns a different `totalTasks` value
- **THEN** the count chip displays the updated count

### Requirement: TaskManagementPanel SHALL preserve the priority legend footer

The existing footer containing the priority label and three priority badges (High, Medium, Low) SHALL remain in `TaskManagementPanel.vue`. The footer SHALL render in both read-only and editable variants.

#### Scenario: Priority legend renders on home dashboard

- **WHEN** the home route renders with task data
- **THEN** the panel footer contains the text "Priority:"
- **AND** the footer contains badge elements for "High", "Medium", and "Low"

#### Scenario: Priority legend renders on /tasks route

- **WHEN** the `/tasks` route renders with task data
- **THEN** the panel footer contains the text "Priority:"
- **AND** the footer contains badge elements for "High", "Medium", and "Low"

### Requirement: TaskManagementPanel SHALL support per-row edit and delete actions in editable mode

`TaskManagementPanel.vue` SHALL emit `edit-task` and `delete-task` events. `TaskRow.vue` SHALL accept an optional named slot (e.g., `actions`) for per-row action buttons. When the panel's `readOnly` prop is false, the consumer's Edit/Delete buttons SHALL render via the slot; when `readOnly` is true, no action buttons SHALL render. The Edit button SHALL emit `edit-task` with the task row object. The Delete button SHALL emit `delete-task` with the task ID.

#### Scenario: Edit/Delete buttons render in editable mode

- **WHEN** `TaskManagementPanel` renders with `readOnly=false` and the consumer provides action buttons via the actions slot
- **THEN** each task row displays Edit and Delete action buttons
- **AND** activating Edit emits `edit-task` with the corresponding task row
- **AND** activating Delete emits `delete-task` with the corresponding task ID

#### Scenario: Edit/Delete buttons absent in read-only mode

- **WHEN** `TaskManagementPanel` renders with `readOnly=true`
- **THEN** no task row displays Edit or Delete buttons
- **AND** no `edit-task` or `delete-task` events are emitted from row interactions

### Requirement: TasksView SHALL reuse TaskManagementPanel as its primary shell

`TasksView.vue` SHALL remove its standalone `<Card>` wrapper and inline `<ul>` task list from the template. It SHALL render `<TaskManagementPanel>` (non-readOnly) inside its existing `<section>` wrapper. The panel's `add-task` emit SHALL be wired to `startNewTask()`, `edit-task` to `startEdit(task)`, and `delete-task` to `handleDelete(id)`. The `AddOperationModal` and all CRUD form state, mutation handlers, room-select loading/error/retry states, and query invalidation SHALL remain owned by `TasksView.vue`.

#### Scenario: TasksView renders the unified panel

- **WHEN** the `/tasks` route renders
- **THEN** the rendered output contains the `TaskManagementPanel` heading "Tasks Backlog"
- **AND** the rendered output contains the total-task count chip
- **AND** the rendered output contains the 5-column table with TaskRow sub-components
- **AND** the rendered output does NOT contain a separate inline `<ul>` task list
- **AND** the rendered output does NOT contain two nested Card shells

#### Scenario: Add Task on /tasks opens the CRUD modal

- **WHEN** a user activates the "Add Task" button on the `/tasks` route
- **THEN** the `AddOperationModal` opens with empty form values
- **AND** submitting the form calls the create mutation and invalidates the backlog query

#### Scenario: Edit on /tasks opens the CRUD modal with prefilled values

- **WHEN** a user activates the "Edit" button on a task row on `/tasks`
- **THEN** the `AddOperationModal` opens with that task's current values prefilled

#### Scenario: Delete on /tasks removes the task

- **WHEN** a user activates the "Delete" button on a task row on `/tasks`
- **THEN** the delete mutation is called and the backlog query is invalidated on success

#### Scenario: TasksView CRUD state survives the refactor

- **WHEN** the `/tasks` route renders after the panel unification
- **THEN** room-select loading state still shows "Loading rooms…" while rooms query is pending
- **AND** room-select error state still shows "Could not load rooms." with a Retry button on failure
- **AND** mutation errors still render inside the modal body without closing the modal
- **AND** cancel or close still resets the form state

### Requirement: Home dashboard SHALL remain strictly read-only

`HomeView.vue` SHALL continue to render `<TaskManagementPanel read-only />`. The home route SHALL NOT render Filter, Add Task, Edit, or Delete controls. The heading and count chip changes SHALL propagate to the home dashboard automatically since they live inside the shared panel component.

#### Scenario: Home dashboard shows backlog data without management controls

- **WHEN** the home route `/` renders
- **THEN** the rendered output contains "Tasks Backlog" heading
- **AND** the rendered output contains the total-task count chip
- **AND** the rendered output contains task row content (e.g., "Pack kitchen boxes")
- **AND** the rendered output does NOT contain "Filter"
- **AND** the rendered output does NOT contain "Add Task"
- **AND** the rendered output does NOT contain per-row Edit or Delete buttons

### Requirement: TaskRow SHALL NOT add a grip/drag handle

The `TaskRow.vue` template SHALL continue to start directly with the task title cell. No leading grip, drag, or handle affordance SHALL be added. The existing five-column anatomy (title, priority badge, people-needed icon+count, room, status badge) SHALL remain unchanged except for the optional actions slot added for editable mode.

#### Scenario: Row anatomy unchanged aside from actions slot

- **WHEN** `TaskRow` renders without the actions slot provided
- **THEN** the row contains exactly five cells: title, priority badge, people-needed icon+count, room text, and status badge
- **AND** no grip/drag handle is rendered before the title cell

### Requirement: No canonical task vocabulary or data contract SHALL change

All task types in `frontend/src/tasks/types.ts`, display helpers in `frontend/src/tasks/helpers.ts`, and the `useTaskBacklog` composable SHALL remain unchanged. The `priority`, `status`, and `assignedTo` vocabulary SHALL stay exactly as currently defined. The `GET /api/tasks/backlog` query and its response shape SHALL not be modified.

#### Scenario: Task types and helpers are untouched

- **WHEN** the change is implemented
- **THEN** `frontend/src/tasks/types.ts` has no diff from the baseline
- **AND** `frontend/src/tasks/helpers.ts` has no diff from the baseline
- **AND** `frontend/src/tasks/composables/useTaskBacklog.ts` has no diff from the baseline

### Requirement: Loading, error, and empty states SHALL be preserved

The `TaskManagementPanel` SHALL continue to render explicit loading ("Loading tasks…"), error ("Could not load tasks. Please try again."), and empty ("No tasks yet.") states. These states SHALL render identically in both read-only and editable variants.

#### Scenario: Loading state renders

- **WHEN** the backlog query is pending
- **THEN** the panel renders "Loading tasks…"
- **AND** the panel does not render the task table or footer

#### Scenario: Error state renders

- **WHEN** the backlog query fails
- **THEN** the panel renders "Could not load tasks. Please try again."
- **AND** the panel does not render the task table or footer

#### Scenario: Empty state renders

- **WHEN** the backlog query succeeds but returns zero tasks
- **THEN** the panel renders "No tasks yet."
- **AND** the panel does not render the task table or footer

### Requirement: Frontend route-render tests SHALL assert the refined design structure

`frontend/tests/app-routes-render.test.ts` SHALL update assertions for both `/` and `/tasks` routes to cover: the "Tasks Backlog" heading, the total-task count chip text, the 5-column table headers, and route-specific control visibility.

#### Scenario: Home route tests assert read-only design structure

- **WHEN** `app-routes-render.test.ts` renders the home route `/`
- **THEN** the test asserts `expect(html).toContain('Tasks Backlog')`
- **AND** the test asserts `expect(html).toContain('3 tasks')` (count chip)
- **AND** the test asserts column headers: "Task", "Priority", "People Needed", "Room / Area", "Status"
- **AND** the test asserts `expect(html).not.toContain('Filter')`
- **AND** the test asserts `expect(html).not.toContain('Add Task')`
- **AND** the test does NOT assert the old heading "Task Management"

#### Scenario: /tasks route tests assert editable design structure

- **WHEN** `app-routes-render.test.ts` renders the `/tasks` route
- **THEN** the test asserts `expect(html).toContain('Tasks Backlog')`
- **AND** the test asserts `expect(html).toContain('3 tasks')` (count chip)
- **AND** the test asserts column headers: "Task", "Priority", "People Needed", "Room / Area", "Status"
- **AND** the test asserts `expect(html).toContain('Add Task')`
- **AND** the test asserts `expect(html).toContain('Edit')`
- **AND** the test asserts `expect(html).toContain('Delete')`
- **AND** the test does NOT assert the old heading "Task backlog"

### Requirement: A TaskManagementPanel component test suite SHALL be added

A new test file at `frontend/src/tasks/components/__tests__/TaskManagementPanel.spec.ts` SHALL provide component-level coverage for the panel's design-aligned anatomy. Tests SHALL cover: heading text, count chip rendering, readOnly control visibility, non-readOnly control visibility, per-row action button visibility, loading/error/empty state text, column header labels, and priority legend footer presence.

#### Scenario: Component tests cover heading and count chip

- **WHEN** `TaskManagementPanel.spec.ts` mounts the panel with mock backlog data
- **THEN** a test asserts the heading element contains "Tasks Backlog"
- **AND** a test asserts the count chip displays the correct total from mock `summary.totalTasks`

#### Scenario: Component tests cover readOnly control visibility

- **WHEN** `TaskManagementPanel.spec.ts` mounts the panel with `readOnly=true`
- **THEN** a test asserts Filter and Add Task buttons are absent
- **AND** a test asserts per-row Edit and Delete buttons are absent
- **WHEN** mounted with `readOnly=false`
- **THEN** a test asserts Filter and Add Task buttons are present
- **AND** a test asserts per-row Edit and Delete buttons are present

#### Scenario: Component tests cover state rendering

- **WHEN** `TaskManagementPanel.spec.ts` mounts the panel in loading/error/empty states
- **THEN** tests assert the correct state text renders for each: "Loading tasks…", "Could not load tasks. Please try again.", "No tasks yet."

#### Scenario: Component tests cover column headers and footer

- **WHEN** `TaskManagementPanel.spec.ts` mounts the panel with task data
- **THEN** tests assert the five column headers render: "Task", "Priority", "People Needed", "Room / Area", "Status"
- **AND** tests assert the priority legend footer renders with "Priority:" text and High/Medium/Low badges
