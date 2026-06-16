## MODIFIED Requirements

### Requirement: `/tasks` and the home dashboard SHALL render a shared BFF-backed Task Management surface

The frontend SHALL replace static task fixtures with a shared `useTaskBacklog` read path that calls the generated `GET /api/tasks/backlog` query and adapts its payload into the Task Management panel contract. The rendered surface SHALL use backend-derived task rows, summary-derived information, priority legend, and status legend. The read path SHALL expose explicit loading, backend-unavailable, and empty states instead of silently falling back to fixture data.

#### Scenario: `/tasks` renders backlog data from the BFF
- **WHEN** the backlog query succeeds on `/tasks`
- **THEN** the panel renders backend-derived task rows and legends rather than local fixture rows

#### Scenario: Task routes show explicit read states
- **WHEN** the shared backlog query is pending, fails, or returns no tasks
- **THEN** the route renders explicit loading, backend-unavailable, or empty states respectively

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