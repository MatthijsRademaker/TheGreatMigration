# task-management-ui Specification

## Purpose
TBD - created by archiving change task-f6503ce2-ff41-47e9-84fc-639a69c67a8a. Update Purpose after archive.
## Requirements
### Requirement: `/tasks` SHALL render a Task Management panel instead of the placeholder

The frontend SHALL replace static task fixtures with a shared `useTaskBacklog` read path that calls the generated `GET /api/tasks/backlog` query and adapts its payload into the Task Management panel contract. The rendered surface SHALL use backend-derived task rows, summary-derived information, priority legend, and status legend. The read path SHALL expose explicit loading, backend-unavailable, and empty states instead of silently falling back to fixture data.

#### Scenario: `/tasks` renders backlog data from the BFF

- **WHEN** the backlog query succeeds on `/tasks`
- **THEN** the panel renders backend-derived task rows and legends rather than local fixture rows

#### Scenario: Task routes show explicit read states

- **WHEN** the shared backlog query is pending, fails, or returns no tasks
- **THEN** the route renders explicit loading, backend-unavailable, or empty states respectively

### Requirement: Task rows SHALL use feature-local typed fixture data aligned to the future backend contract

The Task Management implementation SHALL define feature-local task data matching the future backend row shape: `id`, `title`, `priority`, `peopleNeeded`, `room`, `status`, and `assignedTo`. The first slice SHALL use dedicated static fixture data under `frontend/src/tasks/` and SHALL NOT make API calls. The fixtures SHALL include the design example row `Painting hall` and enough rows to exercise `high`, `medium`, and `low` priorities, multiple rooms/areas, and empty `assignedTo` values.

#### Scenario: Static fixtures drive the initial route content

- **WHEN** the `/tasks` route renders in this slice
- **THEN** the task list is populated from local fixture data
- **AND** no backend fetch or mutation is required to display the rows
- **AND** one rendered row contains `Painting hall`, `High`, `2`, and `Hall`

### Requirement: Status and staffing display SHALL remain forward-compatible with later API wiring

Task rows SHALL be rendered through a dedicated `TaskRow` sub-component. The displayed `Unassigned` status pill SHALL be derived from empty `assignedTo` values through a centralized mapping helper instead of becoming a new canonical status value. The People Needed column SHALL render an inline icon-plus-count presentation using a verified `@lucide/vue` people icon and the numeric `peopleNeeded` value.

#### Scenario: Empty assignments render the design status label

- **WHEN** a fixture row has an empty `assignedTo` array
- **THEN** the row displays an `Unassigned` status pill
- **AND** the underlying typed row data still preserves the canonical task fields for later backend mapping

#### Scenario: Staffing count is rendered as icon plus number

- **WHEN** a task row is rendered
- **THEN** the `People Needed` cell shows a people/user icon next to the numeric staffing count

### Requirement: The panel SHALL reuse existing shared primitives and semantic badge variants only

The Task Management panel SHALL compose existing shared UI primitives and semantic styles only. Priority pills in both rows and the legend SHALL use the existing `priorityHigh`, `priorityMedium`, and `priorityLow` badge variants. The `Unassigned` status pill SHALL use the existing `secondary` badge variant. The implementation SHALL NOT introduce raw hex colors, inline color styles, component-scoped color hacks, or a new shared table primitive for this slice.

#### Scenario: Priority and status pills use existing semantic variants

- **WHEN** the panel renders priority and status badges
- **THEN** `High`, `Medium`, and `Low` use the existing priority badge variants
- **AND** `Unassigned` uses the existing `secondary` badge variant
- **AND** the rendered component source does not require raw color values or inline color styling

### Requirement: Toolbar controls SHALL expose placeholder events without data behavior

The Filter and `+ Add Task` controls SHALL remain visibly interactive in the UI but SHALL NOT perform filtering, task creation, or backend mutation work in this slice. The panel SHALL emit placeholder `filter` and `add-task` events so later wiring can attach behavior without changing the panel contract.

#### Scenario: Toolbar clicks do not require backend integration

- **WHEN** a user activates the Filter or `+ Add Task` control
- **THEN** the panel emits the corresponding placeholder event
- **AND** no API call, data mutation, or disabled-state workaround is required

### Requirement: Route-render tests SHALL verify the new static task-management content

The frontend test suite SHALL update `/tasks` route assertions to target the new Task Management content instead of `Task foundation`. The route-render test SHALL cover the panel title, representative row content, column labels, toolbar controls, and the displayed `Unassigned` state.

#### Scenario: Route test matches the new panel

- **WHEN** `frontend/tests/app-routes-render.test.ts` renders `/tasks`
- **THEN** the assertions include `Task Management`
- **AND** include representative task content such as `Painting hall`, `People Needed`, `Room / Area`, `Unassigned`, `Filter`, and `Add Task`

### Requirement: The task panel SHALL remain presentational and support a read-only home variant

`frontend/src/tasks/components/TaskManagementPanel.vue` SHALL accept a documented `readOnly` prop or equivalent variant. `frontend/src/home/HomeView.vue` SHALL render the panel in read-only mode with no filter, add, edit, or delete controls. `frontend/src/tasks/TasksView.vue` SHALL render the management variant and own the create, edit, delete, and assignment interactions.

#### Scenario: HomeView stays read-only

- **WHEN** the home dashboard renders the task panel
- **THEN** the panel shows backend-derived task data without filter, add, edit, or delete controls

#### Scenario: `/tasks` exposes the management variant

- **WHEN** the `/tasks` route renders successfully
- **THEN** the task panel shows the management controls needed for task creation and maintenance

### Requirement: `/tasks` SHALL own the minimal CRUD interactions and refresh from the canonical backlog query

The `/tasks` route SHALL provide add, edit, delete, and assignment update flows using the typed generated frontend client for the task write endpoints. Create and edit interactions SHALL use the existing shared `Sheet`, `Input`, `Select`, and `Button` primitives for a focused form. After each successful mutation, the task backlog query SHALL be invalidated or refetched so the rendered rows, summary-derived information, and legends reflect backend-confirmed state.

#### Scenario: Successful mutations refresh the management view

- **WHEN** a task create, update, assignment change, or delete succeeds on `/tasks`
- **THEN** the route refreshes from `GET /api/tasks/backlog`
- **AND** the rendered rows and derived summary information match the backend-confirmed response

#### Scenario: The home dashboard remains read-only after shared data refresh

- **WHEN** backlog data changes and the home dashboard re-renders
- **THEN** the home panel reflects the updated backend data
- **AND** it still does not expose destructive controls

### Requirement: Task row typing and display SHALL stay aligned with the backend contract

Frontend task rows SHALL use the canonical priority and status vocabulary from the generated client or matching local union types rather than loose `string` status fields. The displayed `Unassigned` pill SHALL continue to be derived from an empty `assignedTo` array and SHALL NOT become a new canonical backend status value.

#### Scenario: Canonical task vocabulary is preserved in frontend types

- **WHEN** task rows are adapted for display
- **THEN** task priority and status values are constrained to the canonical backend vocabulary

#### Scenario: Empty assignments still render the design label

- **WHEN** a task row has an empty `assignedTo` array
- **THEN** the UI displays `Unassigned`
- **AND** the underlying task `status` remains one of `backlog`, `ready`, or `assigned`

### Requirement: Frontend verification SHALL cover BFF-backed route rendering and mutation behavior

Frontend tests SHALL cover the shared task backlog composable/adapter, explicit loading/error/empty states, create/edit/delete mutation flows with query invalidation, and SSR route rendering for `/tasks` and `/` using realistic BFF-shaped task backlog mocks. Route tests SHALL stop relying on fixture-derived content leaking through when `/api/tasks/backlog` is mocked empty.

#### Scenario: Route render tests use backlog-shaped mocks

- **WHEN** the SSR route-render suite mocks `/api/tasks/backlog`
- **THEN** the mock payload matches the backlog contract shape used by the BFF-backed panel
- **AND** `/tasks` assertions validate backend-driven rendering rather than fixture leakage

#### Scenario: Frontend tests cover explicit route states and management refresh

- **WHEN** the frontend task test suite runs
- **THEN** it verifies loading, error, and empty states plus successful management mutations that invalidate or refetch the backlog query
