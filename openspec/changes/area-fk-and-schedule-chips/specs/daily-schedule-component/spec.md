## MODIFIED Requirements

### Requirement: The frontend SHALL adapt dashboard daily schedule data into the Daily Schedule component contract

The frontend SHALL provide a shared daily-schedule read path that wraps the generated `GET /api/dashboard/daily-schedule` query and adapts the generated payload into the local contract rendered by `DailySchedule.vue`. The adapter SHALL treat nullable `days` and `tasks` arrays as empty arrays, narrow `priority` to `high`, `medium`, or `low`, narrow `staffingStatus` to `fullyStaffed` or `underStaffed`, map the `area { id, name }` object onto the local task contract, and ignore or sanitize invalid rows instead of letting them crash rendering. A task whose `area` is missing or malformed SHALL be rendered with a stable fallback rather than crashing.

#### Scenario: Successful query data adapts into component props without template changes

- **WHEN** the daily-schedule query succeeds
- **THEN** the adapted output provides the day and task props that `DailySchedule.vue` renders
- **AND** each task carries its `area` object through to the card
- **AND** the component template does not require generated client types

#### Scenario: Null arrays become empty collections

- **WHEN** the generated response omits or nulls `days` or `tasks`
- **THEN** the adapter returns empty arrays for those collections
- **AND** rendering still succeeds

#### Scenario: Invalid enum values or malformed rows do not crash rendering

- **WHEN** the generated response contains non-canonical `priority` or `staffingStatus` values, a missing `area`, or malformed task rows
- **THEN** the adapter filters or sanitizes those values before rendering
- **AND** the route shows a stable board state instead of throwing

### Requirement: The schedule-card create/edit SHALL select room or area from the backend-managed room list

The `/calendar` route SHALL render the room/area field in the schedule-card create and edit modal as a Select populated from the `listRoomsQuery` (`GET /api/rooms`). The Select SHALL display the room name and set the form value to the selected room's **id**, which is submitted as `areaId`. While the rooms query is pending, the Select SHALL show a loading placeholder. If the rooms query fails, the Select SHALL display an error message with a retry control. The area field SHALL remain required in the create and update request body.

#### Scenario: Room select shows managed rooms on success

- **WHEN** a user opens the schedule-card create or edit modal and the rooms query resolves
- **THEN** the Room / Area field renders as a Select populated with room names from the backend response
- **AND** the user can select a room from the dropdown

#### Scenario: Room select shows loading state while rooms fetch

- **WHEN** a user opens the schedule-card modal and the rooms query is pending
- **THEN** the Room / Area field shows a loading placeholder
- **AND** the rest of the form remains interactive

#### Scenario: Room select shows error state on failure

- **WHEN** the rooms query fails
- **THEN** the Room / Area field displays an error message indicating rooms could not be loaded
- **AND** the field provides a retry control for the rooms query

#### Scenario: Selected room id is submitted as areaId

- **WHEN** a user selects a room from the dropdown and submits the form
- **THEN** the submitted `areaId` value is the selected room's id
