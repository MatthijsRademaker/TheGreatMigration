## MODIFIED Requirements

### Requirement: The homepage and people route SHALL render the same BFF-backed availability surface

`frontend/src/home/HomeView.vue` and `frontend/src/people/PeopleView.vue` SHALL use the same shared People Availability read path so both routes render the backend-derived availability window. The homepage SHALL keep the existing dashboard layout and remain read-only (no editable prop). The `/people` route SHALL pass `editable={true}` to `PeopleAvailability.vue` and wire the emitted `update-cell` and `delete-person` events to the existing mutations.

#### Scenario: HomeView shows backend-derived availability in the existing dashboard slot

- **WHEN** the homepage renders and the availability query succeeds
- **THEN** the People Availability panel shows backend-derived people rows and day columns in the existing home dashboard layout
- **AND** the homepage renders the component without `editable` (defaults to `false`)
- **AND** the homepage does not render create, update, or delete controls for availability management

#### Scenario: PeopleView passes editable to the component

- **WHEN** the `/people` route renders around the shared availability query with data available
- **THEN** it renders `PeopleAvailability` with `editable={true}`
- **AND** clicking a status cell opens the inline picker for status selection
- **AND** `HomeView` continues to pass no `editable` prop, showing read-only pills

### Requirement: The `/people` route SHALL own the minimal availability management interactions via the editable matrix

The `/people` route SHALL be the surface for all availability management: adding a person, updating a single person's status for a single date, and removing a person. Status update interactions SHALL happen through the editable matrix's inline picker rather than a separate management card. Create-person SHALL remain as a dedicated form above the matrix. Delete-person SHALL be triggered from a button in the matrix person row. All interactions SHALL use the typed frontend client generated from the backend OpenAPI contract, while the homepage remains a read-only overview.

#### Scenario: Creating a person adds them to the editable matrix

- **WHEN** a person is created from the form above the matrix
- **THEN** the route refreshes to show the new person in the backend-backed matrix
- **AND** the new person's status cells can be clicked to set their availability

#### Scenario: Updating one status cell via inline picker refreshes the matrix

- **WHEN** a user clicks a status cell in editable mode and selects a new status
- **THEN** the component emits `update-cell` with `{ personId, dayIndex, status }`
- **AND** PeopleView maps the event to the upsert availability mutation
- **AND** on success the route refreshes to show the backend-confirmed status

#### Scenario: Clearing a cell status via inline picker refreshes the matrix

- **WHEN** a user clicks a status cell in editable mode and selects "Clear"
- **THEN** the component emits `update-cell` with `{ personId, dayIndex, status: null }`
- **AND** PeopleView maps the event to the delete availability mutation
- **AND** on success the route refreshes to show the default status for that cell

#### Scenario: Deleting a person from the matrix row removes them

- **WHEN** a user clicks the delete trigger on a person row in the matrix
- **THEN** the component emits `delete-person` with the person's `id`
- **AND** PeopleView maps the event to the delete person mutation
- **AND** on success the route refreshes to remove the person from the matrix

### Requirement: Frontend verification SHALL cover editable People Availability rendering

Frontend tests SHALL cover the editable mode component rendering, cell event emission, Popover presence in editable mode, and the updated `/people` route layout. The existing read-only SSR tests SHALL remain unchanged and passing.

#### Scenario: Frontend tests assert editable People Availability behavior

- **WHEN** the frontend component tests run
- **THEN** they assert that the component renders status cells as clickable triggers when `editable={true}`
- **AND** they assert that clicking a cell emits an `update-cell` event with the correct payload shape
- **AND** they assert that the read-only mode (`editable={false}` or omitted) renders plain badges without Popover triggers

#### Scenario: PeopleView tests cover the updated layout

- **WHEN** the `/people` route tests run
- **THEN** they assert that the create-person form renders above the matrix
- **AND** they assert that the matrix receives `editable={true}`
- **AND** they do not assert presence of the removed "Manage people" card content
- **AND** the route-level loading, error, and empty states still render correctly
