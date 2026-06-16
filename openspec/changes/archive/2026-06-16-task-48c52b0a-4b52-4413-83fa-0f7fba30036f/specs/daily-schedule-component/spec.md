# daily-schedule-component Specification (Delta)

## Purpose

Update the component capability so Daily Schedule remains a presentational board while the app wires it to the dashboard BFF and lets `/calendar` own persisted schedule management.

## MODIFIED Requirements

### Requirement: Daily Schedule add-task affordances SHALL emit future-ready add intents without scheduling persistence

`frontend/src/calendar/DailySchedule.vue` SHALL keep its header `Add task` control and per-day `+ Add task` affordances visible and continue to emit typed add intents. The component SHALL stay presentational: it SHALL accept typed props for days and task cards, and it SHALL NOT fetch `GET /api/dashboard/daily-schedule` or call schedule write endpoints directly. Route-level containers MAY use the emitted intents to launch persisted create, update, and delete flows. The home dashboard SHALL remain read-only, while `/calendar` SHALL own any editable schedule workflow introduced by this change.

#### Scenario: Header add control emits an intent without a direct mutation

- **WHEN** a user activates the Daily Schedule header `Add task` control
- **THEN** the component emits a typed add-task intent for the parent to handle
- **AND** the component does not perform a scheduling mutation itself

#### Scenario: Per-day add control emits the selected date without a direct mutation

- **WHEN** a user activates a per-day `+ Add task` affordance
- **THEN** the component emits a typed add-task intent that includes that column's `date`
- **AND** the component does not call a schedule write endpoint itself

#### Scenario: The home dashboard stays read-only

- **WHEN** the home route renders Daily Schedule from backend-derived props
- **THEN** it shows the shared schedule board data
- **AND** it does not expose schedule create, update, or delete controls

#### Scenario: The calendar route can bind emitted intents to a persisted workflow

- **WHEN** the `/calendar` route renders Daily Schedule with editable route-owned state
- **THEN** the route can attach handlers that persist schedule changes through the BFF
- **AND** `DailySchedule.vue` remains a prop-driven presentation component

## ADDED Requirements

### Requirement: The frontend SHALL adapt dashboard daily schedule data into the Daily Schedule component contract

The frontend SHALL provide a shared daily-schedule read path that wraps the generated `GET /api/dashboard/daily-schedule` query and adapts the generated payload into the local contract rendered by `DailySchedule.vue`. The adapter SHALL treat nullable `days` and `tasks` arrays as empty arrays, narrow `priority` to `high`, `medium`, or `low`, narrow `staffingStatus` to `fullyStaffed` or `underStaffed`, and ignore or sanitize invalid rows instead of letting them crash rendering.

#### Scenario: Successful query data adapts into component props without template changes

- **WHEN** the daily-schedule query succeeds
- **THEN** the adapted output provides the day and task props that `DailySchedule.vue` already renders
- **AND** the component template does not require generated client types

#### Scenario: Null arrays become empty collections

- **WHEN** the generated response omits or nulls `days` or `tasks`
- **THEN** the adapter returns empty arrays for those collections
- **AND** rendering still succeeds

#### Scenario: Invalid enum values or malformed rows do not crash rendering

- **WHEN** the generated response contains non-canonical `priority` or `staffingStatus` values, or malformed task rows
- **THEN** the adapter filters or sanitizes those values before rendering
- **AND** the route shows a stable board state instead of throwing

### Requirement: The homepage and calendar route SHALL share the BFF-backed schedule read path

`frontend/src/home/HomeView.vue` and `frontend/src/calendar/CalendarView.vue` SHALL use the same shared daily-schedule read path so both routes render backend-derived schedule data. Both routes SHALL handle loading, backend-unavailable, empty, and success states explicitly. The home dashboard SHALL remain a read-only overview. The `/calendar` route SHALL layer schedule create, update, and delete controls on top of the shared read path.

#### Scenario: The home route renders backend-derived schedule data

- **WHEN** the homepage query succeeds
- **THEN** the Daily Schedule panel renders backend-derived days and task cards
- **AND** hard-coded demo fixtures are not shown as the success-state data source

#### Scenario: The calendar route renders the same read model and owns management controls

- **WHEN** the `/calendar` route query succeeds
- **THEN** it renders the same shared schedule data contract as the homepage
- **AND** the route may expose the editable schedule workflow for persisted changes

#### Scenario: Explicit query states replace silent demo fallback

- **WHEN** the shared daily-schedule query is pending, fails, returns no days, or succeeds
- **THEN** the route shows loading, backend-unavailable, empty, or success UI explicitly
- **AND** it does not silently fall back to demo task data for those states

### Requirement: Frontend verification SHALL cover BFF-backed schedule rendering and route state handling

Frontend tests SHALL cover the shared adapter or composable, component rendering with backend-shaped props, and SSR route rendering for `/` and `/calendar` with mocked daily-schedule responses. Existing demo-data assertions SHALL be replaced with BFF-shaped fixture assertions.

#### Scenario: Component and composable tests cover explicit query states

- **WHEN** frontend tests run for Daily Schedule integration
- **THEN** they cover success, loading, backend-unavailable, and empty states using backend-shaped fixtures

#### Scenario: Route render tests mock the daily-schedule endpoint

- **WHEN** the SSR route render suite runs for `/` and `/calendar`
- **THEN** it provides mocked responses for `GET /api/dashboard/daily-schedule`
- **AND** rendering remains stable without relying on demo defaults
