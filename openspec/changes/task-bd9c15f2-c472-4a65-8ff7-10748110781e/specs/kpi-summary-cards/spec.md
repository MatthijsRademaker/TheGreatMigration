## MODIFIED Requirements

### Requirement: The home page SHALL render a KPI summary card grid driven by contract-backed backend data

The `/` route SHALL continue to render a feature-local `KpiCards` component at the top of the home dashboard. The component SHALL continue to display four KPI cards using existing Card primitives, and the lower portion of `HomeView.vue` SHALL now compose the existing `TaskManagementPanel`, `PeopleAvailability`, `DailySchedule`, and static `Move notes` panel instead of preserving the removed `Today’s plan` placeholder.

#### Scenario: Home route renders four KPI cards

- **WHEN** the `/` route renders
- **THEN** the top summary row contains exactly four cards with labels `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`

#### Scenario: KPI cards use Card primitives

- **WHEN** the KpiCards component renders
- **THEN** each card uses Card, CardHeader, CardContent, CardTitle, and CardDescription from `@/shared/ui/card`

#### Scenario: Lower HomeView sections are composed dashboard panels

- **WHEN** the `/` route renders after the home-page composition change
- **THEN** `Task Management`, `People availability`, `Daily Schedule`, and `Move notes` are present below the KPI row
- **AND** `Today’s plan` is not present

### Requirement: The SSR route-render test SHALL mock the new KPI card data sources

The `renderRoute` helper in `frontend/tests/app-routes-render.test.ts` SHALL continue to provide conditional mock responses for `GET /api/dashboard/people-availability` and `GET /api/tasks/backlog`. The home-route assertion SHALL verify the four KPI card labels and SHALL also verify the composed home-dashboard sections. Existing route-render assertions for shell chrome, sidebar, and the other routes SHALL continue to pass.

#### Scenario: Home route SSR test asserts KPI cards and composed dashboard sections

- **WHEN** the SSR route-render test runs for `/`
- **THEN** the rendered HTML contains `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`
- **AND** it contains `Task Management`, `People availability`, `Daily Schedule`, and `Move notes`

#### Scenario: People-availability mock returns valid contract shape

- **WHEN** the SSR mocks `/api/dashboard/people-availability`
- **THEN** the mock returns a 200 response with `summary.availableToday`, `summary.totalPeople`, `range`, `people`, and `statuses` fields

#### Scenario: Task-backlog mock returns valid contract shape

- **WHEN** the SSR mocks `/api/tasks/backlog`
- **THEN** the mock returns a 200 response with `summary.highPriorityTasks`, `summary.unassignedTasks`, `tasks`, `priorities`, and `statuses` fields

#### Scenario: Existing route-render assertions still pass

- **WHEN** the SSR route-render test runs for all routes
- **THEN** all existing assertions for shell chrome, sidebar navigation, page titles, descriptions, and non-home routes continue to pass
