## MODIFIED Requirements

### Requirement: The frontend SHALL provide a presentational People Availability component

`frontend/src/people/PeopleAvailability.vue` SHALL render a titled People Availability panel as a frontend-only component. The panel SHALL include person rows, day columns, availability status cells covering `available`, `busy`, `partial`, and `off`, and a visible legend. When rendered in editable mode, the component SHALL accept an `editable` boolean prop, make status cells clickable, and emit an `update-cell` event when a cell's status is changed. The component SHALL remain frontend-only, never importing generated query artifacts or client types.

#### Scenario: Standalone component renders the designed matrix shape

- **WHEN** `PeopleAvailability.vue` renders with its default data
- **THEN** it displays a titled card or panel
- **AND** it renders day labels across the matrix
- **AND** it renders one row per person with a status indicator for each day
- **AND** the statuses `Available`, `Busy`, `Partial`, and `Off` are visibly represented in the component

#### Scenario: The component remains frontend-only

- **WHEN** the component source is inspected
- **THEN** it does not fetch `GET /api/dashboard/people-availability`
- **AND** it does not import generated query artifacts or generated client types
- **AND** it can render deterministically without backend access

#### Scenario: Editable mode makes cells clickable

- **WHEN** `editable` prop is `true` and a user clicks a status cell
- **THEN** an inline Popover opens showing the four canonical status options and a Clear action
- **AND** selecting a status emits an `update-cell` event with `personId`, `dayIndex`, and the selected `status`

#### Scenario: Read-only mode renders identically without editable prop

- **WHEN** `editable` prop is `false` or not provided
- **THEN** the component renders status cells as plain badges without click handlers or Popovers
- **AND** no delete buttons appear in person rows
- **AND** the rendered HTML matches the existing read-only output

### Requirement: The People Availability component SHALL expose a data contract supporting both read-only and editable modes

`PeopleAvailability.vue` SHALL accept typed props for its displayed content, including the matrix days and people rows, while providing sensible demo defaults. The contract SHALL extend to support an optional `editable` boolean prop. When `editable` is `true`, status cells SHALL render as interactive triggers for a Popover-based status picker, and each person row SHALL include a delete trigger that emits a `delete-person` event. The component SHALL emit typed events rather than managing mutations itself.

#### Scenario: Default data supports immediate visual verification

- **WHEN** the component is rendered without props
- **THEN** it shows deterministic demo content for the title, day labels, people rows, and legend
- **AND** the output is stable enough for SSR test assertions

#### Scenario: Props can replace the demo content

- **WHEN** a caller passes explicit day and people data through the component props
- **THEN** the component renders the provided values using the same visual structure
- **AND** it does not require template changes to switch from demo data to real data later

#### Scenario: Editable mode emits cell change events

- **WHEN** a user selects a new status in the inline picker
- **THEN** the component emits `update-cell` with `{ personId: string, dayIndex: number, status: AvailabilityStatus }`
- **WHEN** a user selects "Clear" in the inline picker
- **THEN** the component emits `update-cell` with `{ personId: string, dayIndex: number, status: null }`
- **WHEN** a user clicks a delete trigger on a person row
- **THEN** the component emits `delete-person` with the person's `id` string

#### Scenario: The component does not import generated client types

- **WHEN** the component source is inspected
- **THEN** it only imports from `@/shared/ui/*` and `./types`
- **AND** it does not reference `@/client/*` or any generated API artifacts

### Requirement: The People Availability component SHALL reuse existing design-system primitives and dense row treatment

The component SHALL be composed from existing frontend primitives and semantic tokens. In editable mode, the inline status picker SHALL reuse the existing `Popover` component from `@/shared/ui/popover`. Status options in the picker SHALL render as clickable `Badge` instances using the canonical variants. Small screens SHALL preserve readability through horizontal overflow of the matrix rather than collapsing the design into stacked cards.

#### Scenario: Status cells use canonical badge semantics

- **WHEN** availability states are rendered in the matrix, legend, or the inline picker
- **THEN** they use the existing `Badge` semantic variants `available`, `busy`, `partial`, and `off`
- **AND** no raw hex color classes or duplicate status-chip primitives are introduced

#### Scenario: Inline picker reuses Popover primitive

- **WHEN** editable mode is active and a cell is clicked
- **THEN** the picker renders inside a `PopoverContent` from `@/shared/ui/popover`
- **AND** the Popover is anchored to the clicked badge element via `PopoverTrigger`
- **AND** the picker dismisses when the user selects an option or clicks outside

#### Scenario: Small screens keep the matrix readable

- **WHEN** the component renders on a narrow viewport
- **THEN** the day matrix can scroll horizontally
- **AND** labels remain readable instead of being compressed or restacked per person

### Requirement: The people route and frontend tests SHALL adopt the editable component

`frontend/src/people/PeopleView.vue` SHALL pass `editable={true}` to `PeopleAvailability.vue` and wire the emitted `update-cell` and `delete-person` events to the existing mutations. Frontend verification SHALL add tests for the editable interaction, including Popover rendering and cell event emission. Read-only tests for the component SHALL remain unchanged.

#### Scenario: PeopleView passes editable prop and handles events

- **WHEN** the `/people` route renders and data is available
- **THEN** `PeopleAvailability` renders with `editable={true}`
- **AND** the `update-cell` event is wired to the upsert and delete-availability mutations
- **AND** the `delete-person` event is wired to the delete-person mutation

#### Scenario: Tests cover editable mode interactions

- **WHEN** the frontend tests run
- **THEN** an editable-mode test renders the component with `editable={true}` and asserts the presence of Popover triggers
- **AND** a cell interaction test simulates clicking a cell and asserts the emitted event payload matches the expected `{ personId, dayIndex, status }` shape
- **AND** the existing read-only SSR tests continue to pass unchanged
