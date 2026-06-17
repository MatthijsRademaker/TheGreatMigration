# daily-schedule-design-alignment Specification

## Purpose
TBD - created by archiving change task-29534381-9c4d-4aa3-9fb0-4ecff9210a15. Update Purpose after archive.
## Requirements
### Requirement: The Daily Schedule panel SHALL render a single compact header row

The Daily Schedule component SHALL render a single compact header row replacing the current three-layer stacked hierarchy (pagination bar above CardHeader, CardHeader with CardTitle, CardContent controls row). The compact header SHALL contain the panel title ("Daily Schedule") on the left; on the right it SHALL contain the date-range label, page indicator ("Page X of Y"), Previous and Next buttons, the decorative "View by: Day" label, and a conditional "Add task" button when `readOnly` is false.

#### Scenario: Compact header renders with all elements in one row

- **WHEN** DailySchedule.vue renders with `days` prop, `page=2`, `totalPages=5`, `dateRangeLabel="5 Jul – 8 Jul"`, and `readOnly=false`
- **THEN** the rendered output contains a single header row
- **AND** the header row contains "Daily Schedule"
- **AND** the header row contains the date-range label "5 Jul – 8 Jul"
- **AND** the header row contains "Page 2 of 5"
- **AND** the header row contains Previous and Next buttons
- **AND** the header row contains "View by: Day"
- **AND** the header row contains "Add task"
- **AND** the rendered output does NOT contain separate "data-slot=\"card-title\"" or "data-slot=\"card-header\"" wrapper elements

#### Scenario: Compact header hides Add task when readOnly is true

- **WHEN** DailySchedule.vue renders with `readOnly=true`
- **THEN** the header row contains "Daily Schedule" and "View by: Day"
- **AND** the header row does NOT contain "Add task"

#### Scenario: Compact header omits pagination when page is 0

- **WHEN** DailySchedule.vue renders with `page=0` and `totalPages=0`
- **THEN** the header row does NOT contain "Page", "Previous", or "Next"
- **AND** the header row still contains "Daily Schedule" and "View by: Day"

### Requirement: Task cards SHALL use simplified anatomy matching design density

Each task card SHALL render: a priority accent left border (`border-l-4`), the task title, a right-aligned priority badge, compact assignee text (comma-separated initials or names), staffing ratio text (`X / Y`), and a conditional under-staffed indicator ("— needs help"). On the editable variant (`readOnly=false`), each card SHALL also render compact Edit and Delete ghost buttons. Task cards SHALL NOT render a "From backlog" badge or individual Avatar components per assignee.

#### Scenario: Task card renders simplified anatomy

- **WHEN** DailySchedule.vue renders a task card with `priority="high"`, `title="Painting hall"`, `assignedPeople=[{name:"Taylor",initials:"T"}]`, `assignedCount=1`, `peopleNeeded=2`, `staffingStatus="underStaffed"`, and `readOnly=false`
- **THEN** the card contains the title "Painting hall"
- **AND** the card contains a priority badge with `data-variant="priorityHigh"`
- **AND** the card contains staffing ratio "1 / 2"
- **AND** the card contains "— needs help"
- **AND** the card contains Edit and Delete buttons
- **AND** the card does NOT contain a "From backlog" badge
- **AND** the card does NOT render an Avatar component with initials "T"
- **AND** the card displays assignee information in compact text format

#### Scenario: Task card hides Edit/Delete when readOnly is true

- **WHEN** DailySchedule.vue renders a task card with `readOnly=true`
- **THEN** the card contains the title and priority badge
- **AND** the card does NOT contain Edit or Delete buttons
- **AND** the card does NOT contain "From backlog"

#### Scenario: Task card does not show From backlog badge regardless of taskId

- **WHEN** DailySchedule.vue renders a task card with `taskId="task-1"` (non-null)
- **THEN** the card renders normally with title, priority, and staffing
- **AND** the card does NOT contain a "From backlog" badge or label

### Requirement: Day headers SHALL use compact horizontal layout

Each day column header SHALL render the day label (e.g., "1 Aug (Sat)") and the available-people count (e.g., "3 available") in a compact horizontal layout rather than a vertical stack of separate heading and paragraph elements.

#### Scenario: Day header renders compact horizontal format

- **WHEN** DailySchedule.vue renders a day column for `label="1 Aug (Sat)"` with `availablePeopleCount=3`
- **THEN** the rendered output contains "1 Aug (Sat)"
- **AND** the rendered output contains "3 available"
- **AND** the label and count appear in a horizontal layout within the day header

### Requirement: The four-column board SHALL use tightened grid spacing

The schedule board SHALL render day columns in a horizontal flex layout with tightened spacing so that the default 4-day slice fits without horizontal scroll on standard desktop viewports. The layout SHALL preserve `overflow-x-auto` and `shrink-0` behavior for non-default day counts.

#### Scenario: Four columns fit without scroll on standard desktop

- **WHEN** DailySchedule.vue renders 4 day columns on a viewport >= 1280px wide
- **THEN** all 4 columns are visible without horizontal scroll
- **AND** each column uses a minimum width appropriate for readable card content

#### Scenario: More than four columns activate overflow scroll

- **WHEN** DailySchedule.vue renders more than 4 day columns
- **THEN** the board container provides horizontal scroll
- **AND** no columns collapse below their readable minimum width

### Requirement: Add task placeholder SHALL be visible only in editable mode

The per-column dashed-border Add task placeholder ("+ Add task") SHALL render only when `readOnly` is false. When `readOnly` is true, all per-column placeholders SHALL be omitted.

#### Scenario: Placeholder visible in editable mode

- **WHEN** DailySchedule.vue renders with `readOnly=false`
- **THEN** each day column contains a dashed-border placeholder card
- **AND** each placeholder contains "+ Add task"

#### Scenario: Placeholder hidden in read-only mode

- **WHEN** DailySchedule.vue renders with `readOnly=true`
- **THEN** no day column contains a "+ Add task" placeholder
- **AND** no day column contains a dashed-border placeholder card

### Requirement: Component SHALL preserve existing prop and emit contracts

The DailySchedule component SHALL continue to accept all existing props (`days`, `readOnly`, `page`, `totalPages`, `dateRangeLabel`) and emit all existing events (`add-task`, `edit-task`, `delete-task`, `prev-page`, `next-page`). No props or emits SHALL be added or removed.

#### Scenario: All existing props are still accepted

- **WHEN** DailySchedule.vue is inspected
- **THEN** its `defineProps` declares `days`, `readOnly`, `page`, `totalPages`, and `dateRangeLabel`

#### Scenario: All existing emits are still declared

- **WHEN** DailySchedule.vue is inspected
- **THEN** its `defineEmits` declares `add-task`, `edit-task`, `delete-task`, `prev-page`, and `next-page`

### Requirement: "View by: Day" SHALL remain a static decorative label

The "View by: Day" control SHALL render as a static decorative label (current behavior). It SHALL NOT activate a dropdown, modal, or any interactive behavior. Its presence matches the design artifacts.

#### Scenario: View by Day is present but non-interactive

- **WHEN** DailySchedule.vue renders with any props
- **THEN** the output contains "View by: Day"
- **AND** "View by: Day" does not have click handlers or dropdown behavior

### Requirement: Pagination controls SHALL remain inside DailySchedule.vue

Pagination controls (date-range label, page indicator, Previous/Next buttons) SHALL remain inside DailySchedule.vue, continuing the current component-layer pagination architecture. The pagination bar SHALL be visually integrated into the compact header row rather than rendered above the Card shell. Route-level consumers (HomeView, CalendarView) SHALL continue to pass pagination props and bind pagination emits into DailySchedule.vue.

#### Scenario: Pagination renders inside the component header

- **WHEN** DailySchedule.vue renders with `page=2`, `totalPages=5`, `dateRangeLabel="5 Jul – 8 Jul"`
- **THEN** the compact header row contains "Page 2 of 5"
- **AND** the compact header row contains Previous and Next buttons
- **AND** the compact header row contains the date-range label "5 Jul – 8 Jul"

#### Scenario: Route consumers still pass pagination props to the component

- **WHEN** HomeView.vue and CalendarView.vue render DailySchedule
- **THEN** both routes continue to pass `page`, `totalPages`, and `dateRangeLabel` props
- **AND** both routes continue to bind `@prev-page` and `@next-page` emits

### Requirement: Frontend tests SHALL assert the refined design structure

Component tests for DailySchedule SHALL assert the compact header row structure, simplified card anatomy (no "From backlog" badge, compact assignee text), compact day header format, and read-only vs. editable control visibility. Route-render tests SHALL continue to assert home read-only behavior and calendar editable behavior with updated structure expectations.

#### Scenario: Component tests assert compact header

- **WHEN** `DailySchedule.test.ts` renders the component with `days`, `page`, `totalPages`, and `readOnly=false`
- **THEN** tests assert the compact header contains "Daily Schedule", "View by: Day", "Add task", and pagination controls in one row
- **AND** tests do NOT assert on `data-slot="card-title"` or separate CardHeader/CardTitle wrappers

#### Scenario: Component tests assert simplified card anatomy

- **WHEN** `DailySchedule.test.ts` renders the component with task cards
- **THEN** tests assert cards contain titles, priority badges, staffing ratios, and compact assignee text
- **AND** tests assert cards do NOT contain "From backlog"
- **AND** tests assert cards do NOT contain Avatar components for each assignee

#### Scenario: Route tests assert home read-only and calendar editable

- **WHEN** `app-routes-render.test.ts` renders the `/` home route
- **THEN** the rendered output does NOT contain "Add task"
- **AND** the rendered output does NOT contain "+ Add task"
- **AND** the rendered output does NOT contain ">Edit<"
- **AND** the rendered output does NOT contain ">Delete<"
- **WHEN** `app-routes-render.test.ts` renders the `/calendar` route
- **THEN** the rendered output contains "Add task"
- **AND** the rendered output contains "+ Add task"
- **AND** the rendered output contains Edit and Delete controls
