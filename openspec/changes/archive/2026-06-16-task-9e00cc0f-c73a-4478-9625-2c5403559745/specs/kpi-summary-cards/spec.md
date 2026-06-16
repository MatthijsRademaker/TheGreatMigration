## ADDED Requirements

### Requirement: KPI card elements exist as Dashboard sub-elements

The LikeC4 model SHALL define four `component` sub-elements inside the `Dashboard` feature area, one per KPI summary card. Each SHALL have a unique identifier and a description that specifies the card's purpose and data source.

- `peopleKpi` — "People available today. Displays the number of helpers with confirmed availability. Backed by `peopleData`."
- `highPriorityKpi` — "High priority tasks. Displays the count of tasks flagged as high priority. Backed by `taskData`."
- `unassignedKpi` — "Unassigned jobs. Displays the count of tasks not yet assigned to any helper. Backed by `taskData`."
- `roomsKpi` — "Rooms completed. Displays the number of fully completed rooms. Currently a static placeholder — no backend data source."

#### Scenario: Four KPI card sub-elements defined

- **WHEN** the LikeC4 model is parsed
- **THEN** the `Dashboard` feature area SHALL have exactly four child elements: `peopleKpi`, `highPriorityKpi`, `unassignedKpi`, and `roomsKpi`

### Requirement: KPI cards have summarizes relationships to their data sources

Each KPI card that consumes data from a planning data store SHALL declare a `summarizes` relationship from the card element to the corresponding data store element. The relationship title SHALL match the card's purpose.

- `dashboard.peopleKpi -[summarizes]-> peopleData` with title "Confirmed helpers count"
- `dashboard.highPriorityKpi -[summarizes]-> taskData` with title "High priority task count"
- `dashboard.unassignedKpi -[summarizes]-> taskData` with title "Unassigned job count"
- `dashboard.roomsKpi` SHALL have NO `summarizes` relationship (placeholder)

#### Scenario: Data-backed cards have summarizes relationships

- **WHEN** inspecting the model relationships for `dashboard.peopleKpi`
- **THEN** SHALL find a `summarizes` relationship to `peopleData`
- **WHEN** inspecting the model relationships for `dashboard.highPriorityKpi`
- **THEN** SHALL find a `summarizes` relationship to `taskData`
- **WHEN** inspecting the model relationships for `dashboard.unassignedKpi`
- **THEN** SHALL find a `summarizes` relationship to `taskData`
- **WHEN** inspecting the model relationships for `dashboard.roomsKpi`
- **THEN** SHALL find NO `summarizes` relationships

### Requirement: Scoped dashboard KPI view

The LikeC4 model SHALL include a scoped component view `dashboard-kpis of migration.spa.dashboard` that renders the dashboard, its four KPI card children, their data store dependencies, and the `summarizes` relationships between them.

The view SHALL:
- Have title "Dashboard KPI Cards"
- Have description "How the four KPI summary cards derive their data from planning stores."
- Include `*` (dashboard and its direct children)
- Include the three planning data stores (`taskData`, `peopleData`, `scheduleData` — the latter via its existing relationship)
- Exclude external actors, the app shell, routing, and unrelated feature areas

#### Scenario: Dashboard KPI view renders correctly

- **WHEN** the LikeC4 model is validated with `likec4 validate --json --no-layout`
- **THEN** the validation SHALL pass with no errors in the model
- **WHEN** the `dashboard-kpis` view is processed
- **THEN** SHALL include the dashboard element, four KPI card elements, and the three planning data elements
- **AND** SHALL show the `summarizes` relationship from each data-consuming card to its data store

### Requirement: Existing views continue to render without errors

All existing views (`system-context`, `container-overview`, `feature-map`, `planning-flow`) SHALL remain valid after the model changes. The new sub-elements and relationships SHALL NOT break auto-layout or view-level style rules.

#### Scenario: Existing views remain valid

- **WHEN** the LikeC4 model is validated
- **THEN** all four existing views SHALL produce no validation errors
- **AND** the `feature-map` view SHALL still apply its green `style` rule to the `dashboard` element
