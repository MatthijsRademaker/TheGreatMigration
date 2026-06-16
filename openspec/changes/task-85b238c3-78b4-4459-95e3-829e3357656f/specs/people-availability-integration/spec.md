## ADDED Requirements

### Requirement: The frontend SHALL derive the dashboard query start date from the planning window

`usePeopleAvailability()` SHALL accept an optional `start` parameter. When `start` is not provided, the composable SHALL call `usePlanningWindow()` to obtain the planning window's `startDate` and SHALL pass that value as the `start` query parameter to `GET /api/dashboard/people-availability`. The composable SHALL defer the dashboard query until the planning window query has resolved (via an `enabled` guard or `dependsOn`), so that the request URL always carries an explicit start date. When `start` is explicitly provided, it SHALL be used as-is without consulting the planning window.

#### Scenario: Composable passes planning window startDate as query parameter

- **WHEN** `usePeopleAvailability()` is called without a `start` argument and the planning window resolves with `startDate=2026-07-05`
- **THEN** the dashboard query is issued with `?start=2026-07-05` in the URL

#### Scenario: Composable waits for planning window before firing dashboard query

- **WHEN** `usePeopleAvailability()` is called without a `start` argument
- **THEN** the dashboard query does not fire until the planning window query has resolved

#### Scenario: Explicit start parameter bypasses planning window

- **WHEN** `usePeopleAvailability({ start: '2024-06-01' })` is called
- **THEN** the dashboard query is issued with `?start=2024-06-01` in the URL without consulting the planning window

### Requirement: The composable SHALL expose a `daysISO` computed ref aligned with the adapted `days` array

The `adaptToComponentProps` function SHALL produce a `daysISO` array — a `string[]` of ISO date strings (`YYYY-MM-DD`) parallel to the adapted `days` label array. The `usePeopleAvailability()` composable SHALL expose `daysISO` as a computed ref. Both arrays SHALL be built in the same UTC-based iteration loop so they remain in lockstep.

#### Scenario: daysISO array matches the adapted days array length and order

- **WHEN** `usePeopleAvailability()` returns data for a range spanning 5 days
- **THEN** `daysISO.value.length` is 5
- **AND** `daysISO.value[i]` is the ISO date string for the date displayed at `data.value.days[i]`

#### Scenario: daysISO values are valid ISO 8601 date strings

- **WHEN** userPeopleAvailability() returns resolved data
- **THEN** every entry in `daysISO.value` matches the pattern `YYYY-MM-DD` and is a valid date

### Requirement: PeopleView date derivation SHALL use the adapted `daysISO` array

`frontend/src/people/PeopleView.vue` SHALL remove the existing `getISODate` function. The `handleCellUpdate` handler SHALL derive the date string for a cell update by looking up `daysISO.value[dayIndex]`. If `dayIndex` is out of range (should not happen in normal operation, but a defensive guard), the handler SHALL set `updateError` and return without making an API call.

#### Scenario: handleCellUpdate derives date from daysISO by index

- **WHEN** a cell emits `{ personId, dayIndex: 2, status: 'busy' }` and `daysISO` is `['2026-07-05', '2026-07-06', '2026-07-07', '2026-07-08']`
- **THEN** the upsert call is made with `path.date: '2026-07-07'`

#### Scenario: Out-of-range dayIndex produces a clear error message

- **WHEN** a cell emits `{ personId, dayIndex: 99, status: 'busy' }` and `daysISO.length` is 4
- **THEN** no API call is made
- **AND** `updateError` is set to a clear message indicating the selected cell cannot be mapped to a date

### Requirement: PeopleView error handling SHALL differentiate 400 subtypes

`frontend/src/people/PeopleView.vue` SHALL inspect the backend error body for distinct 400-response subtypes instead of using the single catch-all "Invalid status or date." The handler SHALL parse the `detail` field from the Huma error model and produce a distinct user-facing message for each known subtype: date outside planning window, invalid date format, invalid status value. For unrecognized 400 errors, a generic fallback message including the `detail` text SHALL be displayed.

#### Scenario: Date outside planning window shows a distinct message

- **WHEN** the backend returns `400` with `detail: "date is outside the planning window"`
- **THEN** `updateError` is set to a message indicating the date is outside the planning window

#### Scenario: Invalid status shows a distinct message

- **WHEN** the backend returns `400` with `detail: "status must be one of..."` (containing the canonical status list)
- **THEN** `updateError` is set to a message indicating the status value is invalid

#### Scenario: Invalid date format shows a distinct message

- **WHEN** the backend returns `400` with `detail: "date must be a valid ISO 8601 date"`
- **THEN** `updateError` is set to a message indicating the date format is invalid

#### Scenario: Unrecognized 400 shows a fallback with the backend detail text

- **WHEN** the backend returns `400` with `detail: "some other error"`
- **THEN** `updateError` is set to a fallback message that includes the backend's `"some other error"` text
