# people-availability-design-alignment Specification

## Purpose

This specification defines the text-only design-alignment requirements for the shared `PeopleAvailability` Vue component. It translates visual intent from `designs/home-page.png` and `designs/components.png` into concrete, testable frontend requirements, governed by the `docs/design-system-v2.md` token and semantics contract. The spec follows the established pattern from `task-management-design-alignment` and `daily-schedule-design-alignment`.

## ADDED Requirements

### Requirement: PeopleAvailability card SHALL render a title-only CardHeader

The `CardHeader` in `PeopleAvailability.vue` SHALL contain only a `CardTitle` element. The `CardDescription` element SHALL be removed. The panel SHALL render the title text "People availability" (lowercase 'a'). This heading SHALL render identically in both read-only (home dashboard) and editable (`/people`) variants.

#### Scenario: Title renders on home dashboard

- **WHEN** `HomeView.vue` renders `<PeopleAvailability v-bind="availabilityData" />`
- **THEN** the rendered output contains "People availability" inside a CardTitle element
- **AND** the rendered output does NOT contain "Track who is available and where each person can help."

#### Scenario: Title renders on /people route

- **WHEN** `PeopleView.vue` renders `<PeopleAvailability v-bind="availabilityData" :editable="true" />`
- **THEN** the rendered output contains "People availability" inside a CardTitle element
- **AND** the rendered output does NOT contain "Track who is available and where each person can help."

### Requirement: PeopleAvailability card SHALL NOT render a summary row

The shared component template SHALL NOT render the "{{ availableToday }} of {{ totalPeople }} available today" summary row. The KpiCards component on the home dashboard already surfaces the available-today count. The matrix itself is the primary data surface on `/people`.

#### Scenario: Summary row is absent on home dashboard

- **WHEN** `HomeView.vue` renders `<PeopleAvailability v-bind="availabilityData" />` with `availableToday=1` and `totalPeople=4`
- **THEN** the rendered output does NOT contain "1 of 4 available today"

#### Scenario: Summary row is absent on /people route

- **WHEN** `PeopleView.vue` renders `<PeopleAvailability v-bind="availabilityData" :editable="true" />`
- **THEN** the rendered output does NOT contain "of" and "available today" in a summary context

#### Scenario: Summary row absent with zero people

- **WHEN** `PeopleAvailability` renders with `availableToday=0` and `totalPeople=0`
- **THEN** the rendered output does NOT contain "0 of 0 available today"

### Requirement: Matrix SHALL omit the Person column header

The `<thead>` row SHALL NOT render a `<th>` element with the text "Person". The first data column (avatar + name) SHALL render without a column header label. The remaining day-column headers and the optional Actions column header SHALL render as before.

#### Scenario: Person column header is absent

- **WHEN** `PeopleAvailability` renders with default people data
- **THEN** the table `<thead>` does NOT contain a `<th>` with text "Person"
- **AND** the first `<td>` in each `<tr>` still renders the avatar and name content

### Requirement: Day column headers SHALL use date+weekday format

The day column `<th>` elements SHALL render labels in date+weekday format (e.g., "Sun 5 Jul", "Mon 6 Jul") matching the output of `formatPlanDayLabel`. The component's default `days` prop value SHALL be updated from `['Mon', 'Tue', 'Wed', 'Thu']` to an empty array, as the adapter always provides date+weekday labels. Demo-only abbreviated labels SHALL no longer be the default.

#### Scenario: Date+weekday headers render from adapter data

- **WHEN** `PeopleAvailability` renders with `days` prop `["Sun 5 Jul", "Mon 6 Jul"]`
- **THEN** the column headers render "Sun 5 Jul" and "Mon 6 Jul"
- **AND** the headers do NOT render abbreviated labels like "Mon", "Tue"

#### Scenario: Default days prop is empty

- **WHEN** `PeopleAvailability` renders with no `days` prop override (using defaults)
- **THEN** no day column headers are rendered (empty `days` array)

### Requirement: Status cells SHALL render as compact capitalized-text Badges

Each status cell SHALL render a `Badge` component with its `variant` set to the canonical availability status (`available`, `busy`, `partial`, `off`) and visible text showing the capitalized status label (e.g., "Available", "Busy", "Partial", "Off"). The Badge SHALL use soft semantic surfaces per `design-system-v2.md`. Status cells SHALL remain text-visible for screen-reader accessibility.

#### Scenario: All four status labels render

- **WHEN** `PeopleAvailability` renders with people data exercising all four statuses
- **THEN** the rendered output contains "Available", "Busy", "Partial", and "Off" text
- **AND** each badge has the correct `data-variant` attribute matching its status

#### Scenario: Status badges render in read-only mode

- **WHEN** `PeopleAvailability` renders without `editable` prop (read-only)
- **THEN** status cells render `Badge` components with capitalized status text
- **AND** no Popover triggers wrap the badges

#### Scenario: Status badges render in editable mode

- **WHEN** `PeopleAvailability` renders with `editable=true`
- **THEN** status cells render `Badge` components with capitalized status text
- **AND** each badge is wrapped in a Popover trigger button

### Requirement: Legend SHALL render as a footer row below the matrix

The status legend SHALL move from inside the `overflow-x-auto` `div` to a footer position below the `<table>` element, outside the horizontal scroll container. The legend SHALL render as a single `flex flex-wrap items-center gap-3` row of compact `Badge` elements ordered: Available, Busy, Partial, Off. Each legend badge SHALL use its canonical `variant` and display its label text.

#### Scenario: Legend renders as footer row

- **WHEN** `PeopleAvailability` renders with default legend
- **THEN** the legend appears after the closing `</table>` tag, outside any `overflow-x-auto` container
- **AND** the legend contains exactly four Badge elements with `data-variant` attributes `available`, `busy`, `partial`, `off`
- **AND** the legend badges display labels "Available", "Busy", "Partial", "Off" in that order

#### Scenario: Legend renders with empty people list

- **WHEN** `PeopleAvailability` renders with `people=[]` and an empty legend array
- **THEN** the legend area renders no Badge elements (empty footer)

### Requirement: All editable-mode behaviour SHALL be preserved

The `PeopleAvailability` component SHALL continue to render the Actions column header (when `editable=true`), Popover triggers per status cell, status-selection popover content with all four status options plus a Clear action, Delete buttons per person row, and the delete confirmation `Dialog` with Cancel and confirm-Delete buttons. The `deletingPersonId` and `updating` props SHALL continue to control button disabled state and text.

#### Scenario: Actions column and Delete buttons render in editable mode

- **WHEN** `PeopleAvailability` renders with `editable=true` and 4 people
- **THEN** the table `<thead>` contains an "Actions" column header
- **AND** each person row contains a Delete button (4 total)

#### Scenario: Popover triggers render per cell in editable mode

- **WHEN** `PeopleAvailability` renders with `editable=true`, 4 people, and 4 days
- **THEN** exactly 16 Popover trigger elements are present

#### Scenario: Status selection emits update-cell event

- **WHEN** a user opens a Popover and selects "Busy" from the status options
- **THEN** the component emits `update-cell` with `{ personId, dayIndex, status: "busy" }`

#### Scenario: Clear action emits update-cell with null status

- **WHEN** a user opens a Popover and selects "Clear (reset to off)"
- **THEN** the component emits `update-cell` with `{ personId, dayIndex, status: null }`

#### Scenario: Delete confirmation dialog flow works

- **WHEN** a user clicks a Delete button
- **THEN** a confirmation Dialog opens with the person's name in the title
- **AND** clicking the confirm Delete button emits `delete-person` with the person id
- **AND** clicking Cancel closes the dialog without emitting `delete-person`

#### Scenario: DeletingPersonId prop disables Delete button

- **WHEN** `PeopleAvailability` renders with `deletingPersonId="p1"` and person "p1" in the list
- **THEN** the Delete button for person "p1" is disabled and shows "Deleting…" text
- **AND** the Delete button for other persons is enabled and shows "Delete"

#### Scenario: Updating prop disables Popover triggers

- **WHEN** `PeopleAvailability` renders with `updating=true` and `editable=true`
- **THEN** all Popover trigger buttons have the `disabled` attribute

### Requirement: No actions or editable controls SHALL render in read-only mode

When `editable` is false (or omitted with its default `false`), the component SHALL NOT render the Actions column header, Popover triggers, Delete buttons, or delete confirmation Dialog content.

#### Scenario: Read-only mode omits all editable controls

- **WHEN** `PeopleAvailability` renders without `editable` prop (default `false`)
- **THEN** the rendered output does NOT contain "Actions" column header
- **AND** the rendered output does NOT contain Popover trigger elements
- **AND** the rendered output does NOT contain Delete buttons
- **AND** the rendered output does NOT contain Dialog content elements

### Requirement: Shared adapter SHALL stop producing removed props

The `adaptToComponentProps` function in `usePeopleAvailability.ts` SHALL no longer produce the `description`, `availableToday`, or `totalPeople` fields in its adapted output. The `rawData` ref SHALL continue to expose the full `DashboardBody` for future use.

#### Scenario: Adapted output excludes removed props

- **WHEN** `adaptToComponentProps` is called with a valid `DashboardBody`
- **THEN** the returned `props` object does NOT contain a `description` field
- **AND** the returned `props` object does NOT contain `availableToday` or `totalPeople` fields
- **AND** the returned `props` object still contains `title`, `days`, `people`, and `legend`

#### Scenario: Adapted output handles null data gracefully

- **WHEN** `adaptToComponentProps` is called with `undefined`
- **THEN** the returned `props` object contains empty `people`, `days`, and `legend` arrays
- **AND** does NOT throw or produce undefined property access errors

### Requirement: PeopleAvailabilityProps type SHALL remove description

The `PeopleAvailabilityProps` interface in `frontend/src/people/types.ts` SHALL remove the `description` field. The `availableToday` and `totalPeople` fields SHALL remain in the interface for backward compatibility with the `usePeopleAvailability` composable return type, even though the component template no longer renders them.

#### Scenario: Type interface is updated

- **WHEN** the change is implemented
- **THEN** `PeopleAvailabilityProps` does NOT declare a `description` field
- **AND** `PeopleAvailabilityProps` still declares `availableToday` and `totalPeople` as optional `number` fields

### Requirement: Frontend tests SHALL assert the refined design structure

Component-level and route-level tests SHALL be updated to assert the new design anatomy and to remove assertions against removed elements.

#### Scenario: PeopleAvailability component tests assert new anatomy

- **WHEN** `PeopleAvailability.test.ts` runs
- **THEN** the test asserting summary text "1 of 4 available today" is removed or updated to assert its absence
- **AND** the test asserting "0 of 0 available today" for empty data is updated to assert absence of summary
- **AND** the test asserting `data-slot="card-title"` still passes
- **AND** a test asserts that "Person" column header text is NOT present
- **AND** the test asserting date headers expects date+weekday format (e.g., "Sun 5 Jul")
- **AND** all editable-mode tests (Popover triggers, status selection emits, Delete dialog flow, `deletingPersonId`, `updating`) still pass unchanged

#### Scenario: PeopleView route tests assert editable matrix without removed elements

- **WHEN** `PeopleView.test.ts` runs
- **THEN** tests assert person names ("Sophia Chen", "Marcus Rivera") are still present
- **AND** tests assert the Actions column is still present in editable mode
- **AND** tests assert Delete buttons are still present (2 for 2 people)
- **AND** tests assert status badges with `data-variant` attributes are still present
- **AND** tests do NOT assert on removed summary text or description

#### Scenario: App-routes-render test updates /people assertions

- **WHEN** `app-routes-render.test.ts` renders the `/people` route
- **THEN** the test asserts `expect(html).toContain('People availability')` (heading)
- **AND** the test does NOT assert `expect(html).toContain('6 of 8 available today')`
- **AND** the test asserts legend footer badges with `data-variant="available"` and `data-variant="busy"` are present
- **AND** the test asserts person names "Sophia Chen" and "Marcus Rivera" are present in the matrix

#### Scenario: App-routes-render test preserves home route assertions

- **WHEN** `app-routes-render.test.ts` renders the home route `/`
- **THEN** the test still asserts `expect(html).toContain('People availability')`
- **AND** the test still asserts KPI card content ("People available today", "6 / 8", "available on Jul 5")
- **AND** the test still asserts Tasks Backlog and Daily Schedule presence

### Requirement: Existing canonical vocabulary and backend contract SHALL remain unchanged

The canonical availability status vocabulary (`available`, `busy`, `partial`, `off`) SHALL remain exactly as defined in `dashboard-people-availability/spec.md`. The `GET /api/dashboard/people-availability` endpoint, its response shape, and the OpenAPI contract SHALL not be modified.

#### Scenario: Backend is untouched

- **WHEN** the change is implemented
- **THEN** no backend Go files have diffs
- **AND** the OpenAPI specification is unchanged
- **AND** the canonical four-status vocabulary is unchanged
