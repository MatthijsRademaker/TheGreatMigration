## ADDED Requirements

### Requirement: The frontend SHALL provide a presentational Daily Schedule board component

`frontend/src/calendar/DailySchedule.vue` SHALL render a frontend-only Daily Schedule board. The component SHALL include a titled panel, a compact header row with static `View by: Day` and `Add task` controls, four day columns, task cards, and a dashed `Add task` placeholder at the bottom of each day column.

#### Scenario: Standalone component renders the designed board anatomy
- **WHEN** `DailySchedule.vue` renders with its default data
- **THEN** it displays the title `Daily Schedule`
- **AND** it displays the header controls `View by: Day` and `Add task`
- **AND** it renders the day labels `2 Jul (Tue)`, `3 Jul (Wed)`, `4 Jul (Thu)`, and `5 Jul (Fri)`
- **AND** it renders the availability labels `6 available`, `7 available`, `7 available`, and `5 available`
- **AND** it renders representative task titles including `Painting hall`, `Steam walls`, `Clean up`, `Sanding`, `Bedroom painting`, `Touch up woodwork`, `Living room finishing`, `2nd floor walls`, `Kitchen painting`, and `Final clean`
- **AND** it renders one `Add task` placeholder per day column

#### Scenario: The component remains frontend-only
- **WHEN** the component source is inspected
- **THEN** it does not fetch `GET /api/dashboard/daily-schedule`
- **AND** it does not import generated query artifacts or generated client types
- **AND** it renders deterministically without backend access

### Requirement: The Daily Schedule component SHALL expose a future-ready local data contract

`DailySchedule.vue` SHALL accept typed props backed by feature-local interfaces that mirror the backend daily-schedule field names: day objects with `date`, `label`, `availablePeopleCount`, and `tasks`; task objects with `id`, `title`, `priority`, `roomArea`, `assignedPeople`, `peopleNeeded`, `assignedCount`, and `staffingStatus`; and assigned-person objects with `id`, `name`, and `initials`. The component SHALL provide deterministic demo defaults using the design-image dates and task titles, while treating those titles as demo fixtures only.

#### Scenario: Default props provide stable design-backed fixtures
- **WHEN** the component renders without props
- **THEN** its default data covers `high`, `medium`, and `low` priorities
- **AND** its rendered staffing metadata includes representative `2 / 2`, `1 / 1`, and under-staffed coverage through the same contract
- **AND** the rendered output is stable enough for SSR assertions

#### Scenario: Props can replace the demo fixtures without template changes
- **WHEN** a caller passes explicit day, task, and assigned-person data through component props
- **THEN** the component renders the provided values using the same day-column and task-card structure
- **AND** it preserves the backend-compatible field names inside the local contract
- **AND** no generated client types or API adapters are required inside the component

### Requirement: The Daily Schedule component SHALL reuse existing design-system primitives and responsive board behavior

The component SHALL compose existing frontend primitives and semantic tokens only. It SHALL use `Card` for the outer panel, `Button` for the header controls and dashed placeholder affordances, `Badge` variants `priorityHigh`, `priorityMedium`, and `priorityLow` for priority labels, and `Avatar` for compact assignee metadata. Each task card SHALL be a plain `div` with card-like styling and a `border-l-4` priority accent using semantic destructive, warning, or success tokens. Small screens SHALL preserve readability through horizontal overflow with a flex row and day columns that keep a readable minimum width rather than collapsing into a compressed grid.

#### Scenario: Priority treatments use canonical badge and border semantics
- **WHEN** a task card renders a `high`, `medium`, or `low` priority
- **THEN** the badge uses `priorityHigh`, `priorityMedium`, or `priorityLow`
- **AND** the card shows a matching semantic left accent for that priority
- **AND** no raw hex color classes or duplicate priority primitives are introduced

#### Scenario: Small screens keep the board readable
- **WHEN** the component renders on a narrow viewport
- **THEN** the day columns can scroll horizontally within the panel
- **AND** each day column keeps a minimum readable width
- **AND** the layout does not rely on compressing task cards into a dense multi-column grid

### Requirement: The calendar route and frontend tests SHALL adopt the new Daily Schedule board

`frontend/src/calendar/CalendarView.vue` SHALL replace its current placeholder content with `DailySchedule.vue` while preserving the route metadata and app-shell behavior. Frontend verification SHALL add a dedicated SSR component test and SHALL update the `/calendar` route render test so it asserts Daily Schedule content instead of the removed placeholder and planning-window column count.

#### Scenario: The calendar route renders the board instead of the placeholder
- **WHEN** the `/calendar` route renders after this change
- **THEN** it shows the Daily Schedule component
- **AND** it no longer shows the text `Schedule board foundation`

#### Scenario: SSR tests cover the board contract
- **WHEN** the frontend SSR tests run
- **THEN** a dedicated Daily Schedule component test asserts the panel title, header controls, day labels, availability counts, representative task titles, priority labels, assignee metadata, staffing counts, and `Add task` placeholders
- **AND** the `/calendar` route render test asserts new board content while preserving the existing app-shell and sidebar checks
