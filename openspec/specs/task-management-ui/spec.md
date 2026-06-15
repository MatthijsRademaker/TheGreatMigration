# task-management-ui Specification

## Purpose
TBD - created by archiving change task-f6503ce2-ff41-47e9-84fc-639a69c67a8a. Update Purpose after archive.
## Requirements
### Requirement: `/tasks` SHALL render a Task Management panel instead of the placeholder

The `/tasks` route SHALL render a feature-local Task Management panel composed through `frontend/src/tasks/TasksView.vue` in place of the existing `Task foundation` placeholder. The panel SHALL include the section title `Task Management`, a Filter control, a `+ Add Task` control, a five-column task list/table area, and a separate priority legend.

#### Scenario: Route renders the task management shell

- **WHEN** the `/tasks` route is rendered
- **THEN** the page contains `Task Management`
- **AND** the toolbar contains `Filter` and `Add Task` controls
- **AND** the task list/table shows the columns `Task`, `Priority`, `People Needed`, `Room / Area`, and `Status`
- **AND** a separate priority legend with `High`, `Medium`, and `Low` pills is visible

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
