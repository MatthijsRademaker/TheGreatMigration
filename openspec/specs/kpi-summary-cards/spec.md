# kpi-summary-cards Specification

## Purpose
Define the KPI summary cards feature across both the LikeC4 architecture model (element decomposition, data-flow relationships, scoped views) and the frontend application (component implementation, data fetching, tests, and design-system compliance).

## Requirements

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

### Requirement: The home page SHALL render a KPI summary card grid driven by contract-backed backend data

The `/` route SHALL render a feature-local `KpiCards` component in place of the current static summary row. The component SHALL display four KPI cards structurally aligned to section 3 of `designs/components.png`, using existing Card, CardHeader, CardContent, CardTitle, and CardDescription primitives from the shared UI library. The component SHALL live at `frontend/src/home/components/KpiCards.vue`.

#### Scenario: Home route renders four KPI cards

- **WHEN** the `/` route renders
- **THEN** the top summary row contains exactly four cards with labels `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`

#### Scenario: KPI cards use Card primitives

- **WHEN** the KpiCards component renders
- **THEN** each card uses Card, CardHeader, CardContent, CardTitle, and CardDescription from `@/shared/ui/card`

#### Scenario: Lower HomeView sections are preserved

- **WHEN** the `/` route renders after the KPI card change
- **THEN** the `Today's plan` and `Move notes` sections remain present below the KPI row

### Requirement: The People available today card SHALL display live availability data from the backend contract

The card SHALL consume `getDashboardPeopleAvailabilityQuery` (already generated in `frontend/src/client/@pinia/colada.gen.ts`). The card value SHALL display `availableToday` of `totalPeople` (e.g., "6 of 8"). The card SHALL handle loading and error states.

#### Scenario: Card displays availability count from API

- **WHEN** `getDashboardPeopleAvailabilityQuery` resolves with `summary.availableToday = 6` and `summary.totalPeople = 8`
- **THEN** the `People available today` card displays the value `6` and the context `8` in a format communicating "6 of 8 available"

#### Scenario: Card shows loading state while query is pending

- **WHEN** `getDashboardPeopleAvailabilityQuery` is in-flight
- **THEN** the `People available today` card displays a loading indicator rather than a stale or zero value

#### Scenario: Card shows error state on query failure

- **WHEN** `getDashboardPeopleAvailabilityQuery` fails
- **THEN** the `People available today` card displays a graceful error state (e.g., "Backend unavailable")

### Requirement: The High priority tasks and Unassigned jobs cards SHALL display live task-backlog data through the generated client

The cards SHALL consume `getTasksBacklogQuery` (to be generated by refreshing `frontend/openapi-snapshot.json` and regenerating `frontend/src/client/`). The High priority tasks card SHALL display `highPriorityTasks`. The Unassigned jobs card SHALL display `unassignedTasks`. Both cards SHALL handle loading and error states.

#### Scenario: High priority tasks card displays count from API

- **WHEN** `getTasksBacklogQuery` resolves with `summary.highPriorityTasks = 4`
- **THEN** the `High priority tasks` card displays the value `4`

#### Scenario: Unassigned jobs card displays count from API

- **WHEN** `getTasksBacklogQuery` resolves with `summary.unassignedTasks = 3`
- **THEN** the `Unassigned jobs` card displays the value `3`

#### Scenario: Task-backlog cards show loading state while query is pending

- **WHEN** `getTasksBacklogQuery` is in-flight
- **THEN** both the `High priority tasks` and `Unassigned jobs` cards display a loading indicator

#### Scenario: Task-backlog cards show error state on query failure

- **WHEN** `getTasksBacklogQuery` fails
- **THEN** both task-backlog KPI cards display a graceful error state

### Requirement: The Rooms completed card SHALL be an isolated placeholder with no derived business logic

The fourth card SHALL use the same Card primitive layout as the other three cards. It SHALL display a non-numeric placeholder value (`—`) labeled `Rooms completed`. It SHALL include `data-testid="kpi-placeholder-rooms-completed"` and a code comment documenting it as a placeholder for a future room-progress contract. It SHALL NOT derive values from task backlog, schedule, people availability, or any other backend data.

#### Scenario: Rooms completed card renders as placeholder

- **WHEN** the KpiCards component renders
- **THEN** the fourth card displays the label `Rooms completed` and the value `—`
- **AND** the card renders `data-testid="kpi-placeholder-rooms-completed"`

#### Scenario: Rooms completed card does not query any backend endpoint

- **WHEN** the KpiCards component renders
- **THEN** no network request is made for room-progress data
- **AND** the card content is purely static

### Requirement: The KPI card styling SHALL use existing design system tokens and primitives

The cards SHALL use semantic icon/accent classes from the Design System v2 token surface exposed in `frontend/src/app/styles.css`. Card icon circles SHALL use semantic background/foreground classes (not raw hex values). No new global theme tokens, shared primitive directories, or page-specific raw color hacks SHALL be introduced.

#### Scenario: Card icons use semantic token classes

- **WHEN** the KpiCards component renders
- **THEN** card icon backgrounds use semantic classes (e.g., `bg-secondary text-secondary-foreground`) rather than raw color values

#### Scenario: No new shared component directories are created

- **WHEN** the KPI card implementation is complete
- **THEN** no new directories are added under `frontend/src/shared/ui/` for this change

### Requirement: The SSR route-render test SHALL mock the new KPI card data sources

The `renderRoute` helper in `frontend/tests/app-routes-render.test.ts` SHALL provide conditional mock responses for `GET /api/dashboard/people-availability` and `GET /api/tasks/backlog`. The home-route assertion SHALL verify that the four KPI card labels appear in the rendered output. Existing route-render assertions for shell chrome, sidebar, and lower HomeView sections SHALL continue to pass.

#### Scenario: Home route SSR test asserts KPI card labels

- **WHEN** the SSR route-render test runs for `/`
- **THEN** the rendered HTML contains `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`

#### Scenario: People-availability mock returns valid contract shape

- **WHEN** the SSR mocks `/api/dashboard/people-availability`
- **THEN** the mock returns a 200 response with `summary.availableToday`, `summary.totalPeople`, `range`, `people`, and `statuses` fields

#### Scenario: Task-backlog mock returns valid contract shape

- **WHEN** the SSR mocks `/api/tasks/backlog`
- **THEN** the mock returns a 200 response with `summary.highPriorityTasks`, `summary.unassignedTasks`, `tasks`, `priorities`, and `statuses` fields

#### Scenario: Existing route-render assertions still pass

- **WHEN** the SSR route-render test runs for all routes
- **THEN** all existing assertions for shell chrome, sidebar navigation, page titles, and descriptions continue to pass
