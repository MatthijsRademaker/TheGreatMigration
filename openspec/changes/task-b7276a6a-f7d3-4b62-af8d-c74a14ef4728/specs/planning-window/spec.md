## MODIFIED Requirements

### Requirement: CalendarView SHALL render the Daily Schedule board while deferring planning-window integration

`frontend/src/calendar/CalendarView.vue` SHALL replace its planning-window placeholder grid with the presentational `DailySchedule` component. In this change it SHALL not import `usePlanningWindow()`, SHALL not render planning-window loading or error states, and SHALL not render `data-testid="plan-day-column"` elements. The `/calendar` route SHALL remain inside the existing app shell and route metadata, but its visible content SHALL be driven by local Daily Schedule demo data or explicit props until a later wiring change is scheduled.

#### Scenario: CalendarView shows the Daily Schedule board
- **WHEN** the `/calendar` route renders after this change
- **THEN** it displays the `Daily Schedule` panel
- **AND** it no longer shows `Schedule board foundation`

#### Scenario: CalendarView is independent from planning-window data in this slice
- **WHEN** the `CalendarView.vue` implementation is inspected
- **THEN** it does not import `usePlanningWindow`
- **AND** it does not render planning-window loading or error branches
- **AND** it renders deterministically without backend planning-window data

### Requirement: Frontend tests SHALL keep planning-window coverage without using `/calendar` as the day-count assertion surface

Planning-window frontend tests SHALL continue to validate the shared composable and day-derivation behavior. `frontend/tests/app-routes-render.test.ts` SHALL stop asserting that `/calendar` renders the mocked planning-window day count and SHALL instead assert Daily Schedule content such as the panel title, representative day labels, task content, priority labels, staffing text, and `Add task` affordances.

#### Scenario: Planning-window logic stays covered in dedicated tests
- **WHEN** the planning-window frontend tests run
- **THEN** they continue to verify composable success, loading, and error states plus day derivation behavior

#### Scenario: The `/calendar` route test matches the presentational board
- **WHEN** `frontend/tests/app-routes-render.test.ts` renders `/calendar`
- **THEN** it asserts the `Daily Schedule` panel and representative board content
- **AND** it does not count `data-testid="plan-day-column"` occurrences
