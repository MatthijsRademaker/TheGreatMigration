# people-availability-integration Specification

## ADDED Requirements

### Requirement: The composable SHALL accept page and daysPerPage options for day pagination

The `usePeopleAvailability()` composable SHALL accept optional `page` (positive integer, defaults to 1) and `daysPerPage` (positive integer, defaults to 7) parameters in its options object. When `page` is provided, the composable SHALL compute the effective `start` date from the planning window start date plus `(page - 1) * daysPerPage` days, and SHALL pass the computed `start` and `daysPerPage` as the `start` and `days` query parameters to the dashboard API. The composable SHALL expose `currentPage` (number), `totalPages` (number computed from planning window total days / daysPerPage), `daysPerPage` (number), and `totalDays` (number from planning window) as reactive properties alongside the existing `data`, `daysISO`, `isLoading`, `isError`, and `isEmpty`.

#### Scenario: Composable passes computed start and days to the API

- **WHEN** `usePeopleAvailability({ page: 2, daysPerPage: 7 })` is called and the planning window starts on `2026-07-05`
- **THEN** the dashboard query is made with `start="2026-07-12"` (7 days after planning window start) and `days=7`

#### Scenario: Default page value is 1

- **WHEN** `usePeopleAvailability()` is called with no page option
- **THEN** `currentPage` is 1

#### Scenario: Days per page defaults to 7

- **WHEN** `usePeopleAvailability()` is called with no daysPerPage option
- **THEN** `daysPerPage` is 7

#### Scenario: totalPages correctly reflects planning window length

- **WHEN** the planning window spans 40 days and `daysPerPage` is 7
- **THEN** `totalPages` is 6 (ceil(40 / 7))

### Requirement: The composable SHALL pass offset and limit to the dashboard API

The `usePeopleAvailability()` composable SHALL accept optional `offset` (non-negative integer, defaults to 0) and `limit` (non-negative integer, defaults to 0, meaning all people) parameters. When `offset` and `limit` are provided, they SHALL be passed as query parameters to the dashboard API. The composable SHALL expose `totalPeople` (from the response's `summary.totalPeople`) as a reactive property.

#### Scenario: Default offset and limit pass through to the API

- **WHEN** `usePeopleAvailability()` is called without offset or limit
- **THEN** the dashboard query is made with no `offset` or `limit` query parameters (backward compatible)

#### Scenario: Explicit offset and limit pass through to the API

- **WHEN** `usePeopleAvailability({ offset: 0, limit: 10 })` is called
- **THEN** the dashboard query includes `offset=0` and `limit=10`

### Requirement: The PeopleView SHALL render day-navigation controls

The `PeopleView.vue` component SHALL render navigation controls above the People Availability matrix when day pagination is active. The controls SHALL include:

- A date range label showing the first and last date of the visible page (e.g., "Sun 5 Jul – Sat 11 Jul")
- A page indicator showing the current page and total pages (e.g., "Page 1 of 6")
- A "Previous" button disabled when `currentPage` is 1
- A "Next" button disabled when `currentPage` equals `totalPages`
- Clicking Previous or Next SHALL decrement or increment the page and re-fetch data

#### Scenario: Navigation controls render with correct labels

- **WHEN** the People Availability panel renders with page 2 of 6 and daysPerPage of 7
- **THEN** the date range label shows the computed date range for page 2
- **AND** the page indicator shows "Page 2 of 6"
- **AND** the Previous button is enabled
- **AND** the Next button is enabled

#### Scenario: Previous button is disabled on first page

- **WHEN** `currentPage` is 1
- **THEN** the Previous button is disabled

#### Scenario: Next button is disabled on last page

- **WHEN** `currentPage` equals `totalPages`
- **THEN** the Next button is disabled

#### Scenario: Clicking Next increments the page and refreshes data

- **WHEN** the user clicks Next on page 1
- **THEN** `currentPage` becomes 2 and the dashboard query is re-fetched with the new page's date range

### Requirement: Frontend verification SHALL cover paginated composable and component behavior

Frontend tests SHALL cover the composable's pagination parameter derivation and the component's navigation control rendering.

#### Scenario: Page navigation recomputes query parameters

- **WHEN** the composable's `page` changes from 1 to 2 with `daysPerPage=7` and a planning window starting `2026-07-05`
- **THEN** the effective `start` date changes from `2026-07-05` to `2026-07-12`
- **AND** `days` remains 7

#### Scenario: Route render test verifies navigation controls

- **WHEN** the `/people` route renders with backend-shaped data
- **THEN** the pagination controls (prev, next, page indicator, date range label) are present in the DOM
