## ADDED Requirements

### Requirement: The frontend SHALL provide a presentational People Availability component

`frontend/src/people/PeopleAvailability.vue` SHALL render a titled People Availability panel as a frontend-only component. The panel SHALL include person rows, four day columns, availability status cells covering `available`, `busy`, `partial`, and `off`, and a visible legend or status row aligned to the design intent from `designs/components.png` and `designs/home-page.png`.

#### Scenario: Standalone component renders the designed matrix shape

- **WHEN** `PeopleAvailability.vue` renders with its default data
- **THEN** it displays a titled card or panel
- **AND** it renders four day labels across the matrix
- **AND** it renders one row per person with a status indicator for each day
- **AND** the statuses `Available`, `Busy`, `Partial`, and `Off` are visibly represented in the component

#### Scenario: The component remains frontend-only

- **WHEN** the component source is inspected
- **THEN** it does not fetch `GET /api/dashboard/people-availability`
- **AND** it does not import generated query artifacts or generated client types
- **AND** it can render deterministically without backend access

### Requirement: The People Availability component SHALL expose a future-ready local data contract

`PeopleAvailability.vue` SHALL accept typed props for its displayed content, including the matrix days and people rows, while providing sensible demo defaults that match the first-slice four-day design. This contract SHALL let a later change pass real backend data into the same structure without rewriting the component template.

#### Scenario: Default data supports immediate visual verification

- **WHEN** the component is rendered without props
- **THEN** it shows deterministic demo content for the title, four day labels, people rows, and legend
- **AND** the output is stable enough for SSR test assertions

#### Scenario: Props can replace the demo content

- **WHEN** a caller passes explicit day and people data through the component props
- **THEN** the component renders the provided values using the same visual structure
- **AND** it does not require template changes to switch from demo data to real data later

### Requirement: The People Availability component SHALL reuse existing design-system primitives and dense row treatment

The component SHALL be composed from existing frontend primitives and semantic tokens. It SHALL reuse the existing `Card` surface, `Badge` availability variants (`available`, `busy`, `partial`, `off`), and `Avatar` initials fallback for row labels. Person rows SHALL use a compact inline Avatar-plus-name treatment rather than the `PersonChip` pill. Small screens SHALL preserve readability through horizontal overflow of the matrix rather than collapsing the design into stacked cards.

#### Scenario: Status cells use canonical badge semantics

- **WHEN** availability states are rendered in the matrix or legend
- **THEN** they use the existing `Badge` semantic variants `available`, `busy`, `partial`, and `off`
- **AND** no raw hex color classes or duplicate status-chip primitives are introduced

#### Scenario: Person rows stay compact enough for matrix density

- **WHEN** a person row is rendered
- **THEN** it shows the person's initials avatar and name inline
- **AND** it does not use the rounded `PersonChip` pill treatment

#### Scenario: Small screens keep the matrix readable

- **WHEN** the component renders on a narrow viewport
- **THEN** the day matrix can scroll horizontally
- **AND** labels remain readable instead of being compressed or restacked per person

### Requirement: The people route and frontend tests SHALL adopt the new component

`frontend/src/people/PeopleView.vue` SHALL replace its current placeholder card with `PeopleAvailability.vue`. Frontend verification SHALL add an SSR component test and update the route-render smoke test so the `/people` route asserts the new component content instead of the old placeholder text.

#### Scenario: PeopleView renders the component instead of the placeholder

- **WHEN** the `/people` route renders
- **THEN** it shows the People Availability component
- **AND** it no longer shows the text `People availability foundation`

#### Scenario: SSR tests cover the component contract

- **WHEN** the frontend SSR tests run
- **THEN** a dedicated People Availability component test asserts the title, sample people, four day labels, all four status labels, and the legend
- **AND** the `/people` route render test asserts new component content while preserving the existing app-shell and sidebar checks
