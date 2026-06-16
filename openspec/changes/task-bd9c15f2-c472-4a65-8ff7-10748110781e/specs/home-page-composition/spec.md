## ADDED Requirements

### Requirement: The home route SHALL compose the dashboard from existing frontend panels

The `/` route SHALL render `frontend/src/home/HomeView.vue` as a composition shell with `KpiCards` as the first row, a second row containing `TaskManagementPanel` in the wider column beside `PeopleAvailability` in the narrower column, and a third row containing `DailySchedule` in the wider column beside the static `Move notes` card in the narrower column. The route SHALL keep the existing app shell, sidebar, and route metadata unchanged.

#### Scenario: Home route renders the composed dashboard hierarchy

- **WHEN** the `/` route renders
- **THEN** the page contains the KPI labels `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`
- **AND** it contains `Task Management`
- **AND** it contains `People availability`
- **AND** it contains `Daily Schedule`
- **AND** it contains `Move notes`

#### Scenario: Each composed dashboard row keeps the intended wide-left / narrow-right layout

- **WHEN** `HomeView.vue` is inspected
- **THEN** the Tasks/People row and the Schedule/Move Notes row are rendered as two independent grid containers
- **AND** each grid uses the existing `xl:grid-cols-[1.4fr_0.9fr]` layout pattern
- **AND** the wider left column hosts `TaskManagementPanel` or `DailySchedule`
- **AND** the narrower right column hosts `PeopleAvailability` or the static `Move notes` card

### Requirement: HomeView SHALL remain a frontend-only composition surface

`frontend/src/home/HomeView.vue` SHALL reuse existing feature components from their current locations and SHALL not introduce new backend wiring, generated-client work, route changes, shell changes, or new domain components. The removed `Today’s plan` placeholder SHALL not be replaced with new custom data structures in this slice.

#### Scenario: Placeholder-only content is removed from HomeView

- **WHEN** `HomeView.vue` is updated for this change
- **THEN** it no longer defines the `upcomingWork` array
- **AND** it no longer renders the `Today’s plan` card
- **AND** it no longer imports `Badge` for the removed placeholder content

#### Scenario: Composition reuses existing components and retains static Move Notes copy

- **WHEN** `HomeView.vue` is inspected
- **THEN** it imports `TaskManagementPanel` from `@/tasks/components/TaskManagementPanel.vue`
- **AND** it imports `PeopleAvailability` from `@/people/PeopleAvailability.vue`
- **AND** it imports `DailySchedule` from `@/calendar/DailySchedule.vue`
- **AND** the `Move notes` panel remains a static `Card` with the existing reminder copy
- **AND** no notes API, reusable notes component, backend handler, or generated query is added in this change

### Requirement: Frontend route-render tests SHALL assert the composed home dashboard

`frontend/tests/app-routes-render.test.ts` SHALL update the `/` route assertions to match the composed dashboard while preserving the existing KPI mock coverage and the existing assertions for other routes.

#### Scenario: Home-route smoke assertions stop depending on the removed placeholder

- **WHEN** the shared-shell home-route test runs
- **THEN** it does not assert `Today’s plan`
- **AND** the `routeCases` entry for `/` uses a home-specific visible string from the composed dashboard such as `Move notes`

#### Scenario: Home-route assertions cover the new dashboard sections and existing KPI output

- **WHEN** `frontend/tests/app-routes-render.test.ts` renders `/`
- **THEN** it asserts `Task Management`, `People availability`, `Daily Schedule`, and `Move notes`
- **AND** it still asserts `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`
- **AND** it still verifies the mocked KPI output values `6`, `of 8`, `4`, and `3`

#### Scenario: Other route assertions remain valid

- **WHEN** the route-render test suite runs for `/tasks`, `/calendar`, `/people`, `/rooms`, and `/settings`
- **THEN** their existing route metadata and content assertions continue to pass
