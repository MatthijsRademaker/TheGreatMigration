# daily-schedule-pagination Specification

## ADDED Requirements

### Requirement: The useDailySchedule composable SHALL accept page and daysPerPage options for day pagination

The `useDailySchedule()` composable SHALL accept optional `page` (positive integer, defaults to 1) and `daysPerPage` (positive integer, defaults to 4) parameters in its options object. The composable SHALL use `usePlanningWindow()` to resolve the canonical planning window. When `page` is provided, the composable SHALL compute the effective `start` date from the planning window start date plus `(page - 1) * daysPerPage` days, and SHALL pass the computed `start` and `daysPerPage` as the `start` and `days` query parameters to the dashboard API. The composable SHALL expose `page` (ref<number>), `totalPages` (computed<number> from planning window total days / daysPerPage, minimum 1), `daysPerPage` (ref<number>), and `totalDays` (computed<number> from planning window) as reactive properties alongside the existing `data`, `isLoading`, `isError`, `isEmpty`, `refresh`, and `queryKey`.

#### Scenario: Composable passes computed start and days to the API

- **WHEN** `useDailySchedule({ page: 2, daysPerPage: 4 })` is called and the planning window starts on `2026-07-05`
- **THEN** the dashboard query is made with `start="2026-07-09"` (4 days after planning window start) and `days=4`

#### Scenario: Default page value is 1

- **WHEN** `useDailySchedule()` is called with no page option
- **THEN** `page.value` is 1

#### Scenario: Days per page defaults to 4

- **WHEN** `useDailySchedule()` is called with no daysPerPage option
- **THEN** `daysPerPage.value` is 4

#### Scenario: totalPages correctly reflects planning window length

- **WHEN** the planning window spans 40 days and `daysPerPage` is 4
- **THEN** `totalPages.value` is 10 (ceil(40 / 4))

#### Scenario: totalPages is at least 1

- **WHEN** the planning window is empty or still loading
- **THEN** `totalPages.value` is 1

### Requirement: The composable SHALL expose goToPrevPage and goToNextPage navigation helpers

The `useDailySchedule()` composable SHALL expose `goToPrevPage()` and `goToNextPage()` functions. `goToPrevPage` SHALL decrement `page` by 1 when `page > 1`. `goToNextPage` SHALL increment `page` by 1 when `page < totalPages`. Neither function SHALL move the page beyond the valid range.

#### Scenario: goToPrevPage advances page by 1 when above page 1

- **WHEN** `page` is 2 and `goToPrevPage()` is called
- **THEN** `page` becomes 1

#### Scenario: goToPrevPage does not go below 1

- **WHEN** `page` is 1 and `goToPrevPage()` is called
- **THEN** `page` remains 1

#### Scenario: goToNextPage advances page by 1 when below totalPages

- **WHEN** `page` is 1 and `totalPages` is 10 and `goToNextPage()` is called
- **THEN** `page` becomes 2

#### Scenario: goToNextPage does not go beyond totalPages

- **WHEN** `page` equals `totalPages` and `goToNextPage()` is called
- **THEN** `page` remains unchanged

### Requirement: The composable SHALL reset page to 1 when the planning window changes

When the planning window data changes (e.g., the admin updates the date range in Settings), the composable SHALL reset `page` to 1 if and only if `page > 1` and the planning window actually changed (not on initial load or loading-to-loaded transition). This SHALL be implemented via a `watch` on the planning window days array with an `oldDates.length > 0` guard, exactly mirroring the `usePeopleAvailability` implementation at `frontend/src/shared/composables/usePeopleAvailability.ts`.

#### Scenario: Page resets to 1 when planning window changes

- **WHEN** `page` is 3 and the planning window date range is updated
- **THEN** `page` resets to 1

#### Scenario: Page does not reset on initial planning window load

- **WHEN** the composable is first created and the planning window transitions from loading (empty) to loaded
- **THEN** `page` remains unchanged (does not spuriously reset)

### Requirement: The composable SHALL defer the dashboard query until the planning window resolves

When no explicit `start` option is provided, the composable SHALL defer the dashboard query until the planning window is done loading and a `startParam` can be derived. This mirrors the `queryEnabled` pattern in `usePeopleAvailability`.

#### Scenario: Query is deferred while planning window is loading

- **WHEN** `useDailySchedule()` is called and the planning window is loading with no resolved days
- **THEN** the dashboard query is not executed until planning window data is available

### Requirement: HomeView and CalendarView SHALL render pagination navigation controls at the route layer

Both `HomeView.vue` and `CalendarView.vue` SHALL render pagination navigation controls outside the `<DailySchedule>` component. The controls SHALL include:

- A date range label showing the first and last day of the visible page (derived from `scheduleData.days[0]?.label` and the last day's label)
- A page indicator showing the current page and total pages (e.g., "Page 1 of 10")
- A "Previous" button disabled when `page <= 1`
- A "Next" button disabled when `page >= totalPages`
- Clicking Previous or Next SHALL call `goToPrevPage()` or `goToNextPage()` respectively, triggering a query refetch

Pagination controls SHALL only be rendered in the success state (not during loading, error, or empty states).

#### Scenario: HomeView renders pagination controls

- **WHEN** the home route (`/`) renders with successful schedule data and a 40-day planning window
- **THEN** a pagination bar is visible showing a date range label, page indicator, and Previous/Next buttons
- **AND** the DailySchedule component remains read-only

#### Scenario: CalendarView renders pagination controls

- **WHEN** the calendar route (`/calendar`) renders with successful schedule data
- **THEN** a pagination bar is visible showing a date range label, page indicator, and Previous/Next buttons
- **AND** the DailySchedule component retains its editable controls (Add task, Edit, Delete)

#### Scenario: Previous button is disabled on first page

- **WHEN** `page` is 1
- **THEN** the Previous button renders with a `disabled` attribute

#### Scenario: Next button is disabled on last page

- **WHEN** `page` equals `totalPages`
- **THEN** the Next button renders with a `disabled` attribute

#### Scenario: Pagination controls are not visible in loading state

- **WHEN** the schedule data is still loading
- **THEN** no pagination controls are rendered

### Requirement: DailySchedule.vue SHALL remain presentational and unchanged

The `DailySchedule.vue` component SHALL NOT acquire pagination-related props, emits, data-fetching, or navigation rendering. Its existing `days` and `readOnly` props SHALL continue as its only interface. The component SHALL continue to emit `add-task`, `edit-task`, and `delete-task` events for route-level handlers.

#### Scenario: DailySchedule.vue does not accept pagination props

- **WHEN** DailySchedule.vue is inspected
- **THEN** its `defineProps` only declares `days` and `readOnly`
- **AND** it does not declare `page`, `totalPages`, `prevDisabled`, `nextDisabled`, or any pagination-related props

#### Scenario: DailySchedule.vue does not render navigation controls

- **WHEN** DailySchedule.vue is rendered with valid `days` props
- **THEN** the rendered output does not contain "Previous", "Next", or "Page X of Y" text

#### Scenario: DailySchedule.vue continues to emit add-task/edit-task/delete-task

- **WHEN** DailySchedule.vue is rendered in non-readOnly mode with valid `days` props
- **THEN** the component still emits `add-task` (with optional date), `edit-task` (with task), and `delete-task` (with taskId) events

### Requirement: Existing composable behavior SHALL remain backward-compatible

The `useDailySchedule()` composable SHALL continue to work when called with zero arguments (the current usage in HomeView and CalendarView). In this case, it SHALL default to `page=1`, `daysPerPage=4`, and derive `start` from the planning window — producing behavior equivalent to the current 4-day default but now explicitly derived rather than relying on backend defaults.

#### Scenario: Zero-argument call still works

- **WHEN** `useDailySchedule()` is called with no options
- **THEN** the composable returns `data`, `isLoading`, `isError`, `isEmpty`, `refresh`, and `queryKey` as before
- **AND** additionally returns `page`, `totalPages`, `daysPerPage`, `totalDays`, `goToPrevPage`, `goToNextPage`

#### Scenario: Existing adaptation logic is preserved

- **WHEN** the dashboard API returns schedule data
- **THEN** the composable still adapts `priority` to `high|medium|low`, `staffingStatus` to `fullyStaffed|underStaffed`, treats null arrays as empty, and filters non-canonical values

### Requirement: The backend API SHALL NOT be modified

The `GET /api/dashboard/daily-schedule` endpoint SHALL remain unchanged. It already accepts optional `start` and `days` query parameters with correct defaults. The frontend SHALL derive and pass these parameters explicitly from the composable instead of relying on backend defaults.

#### Scenario: Backend endpoint remains unchanged

- **WHEN** the backend is running after this change
- **THEN** `GET /api/dashboard/daily-schedule` without query parameters still returns the default 4-day window starting from the planning window start date
- **AND** `GET /api/dashboard/daily-schedule?start=2026-07-10&days=3` still returns 3 days starting from 2026-07-10

### Requirement: Frontend tests SHALL cover page derivation, boundary navigation, and route-level rendering

Frontend tests SHALL be extended to cover the composable's pagination parameter derivation, navigation boundary behavior, planning-window-change page reset, and the presence of pagination controls on both route surfaces.

#### Scenario: Composable test covers page initialization and navigation

- **WHEN** the composable test suite runs
- **THEN** it verifies `page` initializes at 1, `goToNextPage`/`goToPrevPage` move the page correctly, and boundary guards prevent out-of-range navigation

#### Scenario: Composable test covers planning-window-change page reset

- **WHEN** the composable test suite runs
- **THEN** it verifies page resets to 1 when the planning window changes and page > 1, and does NOT reset on initial load

#### Scenario: Route render test verifies pagination controls on home and calendar

- **WHEN** the route render test suite runs for `/` and `/calendar`
- **THEN** it verifies the pagination bar (date range label, page indicator, Previous/Next buttons) is present in the rendered output
- **AND** it verifies Previous button is disabled when on page 1

## ADDED Requirements (continued)

### Requirement: The homepage and calendar route SHALL share the BFF-backed schedule read path with pagination

`HomeView.vue` and `CalendarView.vue` SHALL continue to use the same shared daily-schedule read path. Both routes SHALL handle loading, backend-unavailable, empty, and success states explicitly. The success state SHALL now include pagination navigation controls at the route layer. The home dashboard SHALL remain a read-only overview. The `/calendar` route SHALL continue to layer schedule create, update, and delete controls on top of the shared read path.

#### Scenario: The home route renders paginated schedule data

- **WHEN** the homepage query succeeds with a planning window spanning more than 4 days
- **THEN** the Daily Schedule panel renders backend-derived days and task cards for the current page
- **AND** pagination controls allow navigating to subsequent pages
- **AND** the board remains read-only

#### Scenario: The calendar route renders paginated schedule data with edit controls

- **WHEN** the `/calendar` route query succeeds
- **THEN** it renders the same shared schedule data contract as the homepage for the current page
- **AND** pagination controls allow navigating to subsequent pages
- **AND** the route still exposes the editable schedule workflow (Add task, Edit, Delete)
