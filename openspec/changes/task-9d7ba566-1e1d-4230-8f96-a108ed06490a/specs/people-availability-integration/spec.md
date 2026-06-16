## ADDED Requirements

### Requirement: The frontend SHALL adapt dashboard availability data into the People Availability component contract

The frontend SHALL introduce a shared read path for People Availability that calls the generated `GET /api/dashboard/people-availability` query and adapts the generated `DashboardBody` payload into the local props expected by `PeopleAvailability.vue`. The adapter SHALL handle nullable generated arrays safely, convert ISO date strings into readable day labels for the matrix, and validate availability statuses against the canonical set `available`, `busy`, `partial`, and `off` before passing them to the component.

#### Scenario: Adapter converts the dashboard response into component props
- **WHEN** the dashboard availability query succeeds
- **THEN** the adapted output provides day labels, person rows, availability cells, `availableToday`, and `totalPeople` values that `PeopleAvailability.vue` can render without template changes
- **AND** null or missing generated arrays are treated as empty collections
- **AND** only canonical status values are passed through to the component

#### Scenario: Summary text stays aligned with backend data
- **WHEN** the adapted People Availability props are rendered on a route
- **THEN** the visible `availableToday of totalPeople available today` summary matches the backend response `summary.availableToday` and `summary.totalPeople`

### Requirement: The homepage and people route SHALL render the same BFF-backed availability surface

`frontend/src/home/HomeView.vue` and `frontend/src/people/PeopleView.vue` SHALL use the same shared People Availability read path so both routes render the backend-derived availability window. The homepage SHALL keep the existing dashboard layout and remain read-only. The `/people` route SHALL show explicit loading, backend-unavailable, and empty states instead of silently falling back to demo defaults.

#### Scenario: HomeView shows backend-derived availability in the existing dashboard slot
- **WHEN** the homepage renders and the availability query succeeds
- **THEN** the People Availability panel shows backend-derived people rows and day columns in the existing home dashboard layout
- **AND** the homepage does not render create, update, or delete controls for availability management

#### Scenario: PeopleView handles loading, error, and empty states explicitly
- **WHEN** the `/people` route renders around the shared availability query
- **THEN** it shows a loading state while the query is pending
- **AND** it shows a backend-unavailable state if the query fails
- **AND** it shows an empty-state prompt when no people are returned
- **AND** on success it renders the backend-derived People Availability matrix

### Requirement: The `/people` route SHALL own the minimal availability management interactions

The `/people` route SHALL be the surface that exposes the first management controls for this feature: adding a person, updating a single person’s status for a single date, and removing a person. Those interactions SHALL use the typed frontend client generated from the backend OpenAPI contract, while the homepage remains a read-only overview.

#### Scenario: Creating or deleting a person updates the `/people` matrix
- **WHEN** a person is created or an unreferenced person is deleted from `/people`
- **THEN** the route refreshes to show the resulting people list from the backend-backed matrix

#### Scenario: Updating one status cell refreshes the rendered matrix
- **WHEN** a single-date availability update succeeds from `/people`
- **THEN** the route refreshes to show the backend-confirmed status in the matching matrix cell

### Requirement: Frontend verification SHALL cover BFF-backed People Availability rendering

Frontend tests SHALL cover the shared adapter/composable, component rendering with backend-shaped props, `/people` route state handling, and SSR route rendering for both `/` and `/people` with mocked BFF data.

#### Scenario: Frontend tests assert backend-shaped People Availability behavior
- **WHEN** the frontend component and route tests run
- **THEN** they assert adapted backend-shaped props, canonical status rendering, and the explicit loading, error, and empty states
- **AND** the shared-shell route render suite still verifies the app shell and sidebar behavior for `/` and `/people`
