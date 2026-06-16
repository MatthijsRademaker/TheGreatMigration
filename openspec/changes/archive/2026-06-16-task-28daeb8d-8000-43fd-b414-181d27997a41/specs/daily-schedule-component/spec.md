## ADDED Requirements

### Requirement: The schedule-card create/edit SHALL select room or area from the backend-managed room list

The `/calendar` route SHALL replace the free-form text Input for room/area in the schedule-card create and edit modal with a Select component populated from the `listRoomsQuery` (`GET /api/rooms`). The Select SHALL display the room name and set the `roomArea` form value to the selected room's name string. While the rooms query is pending, the Select SHALL show a loading placeholder. If the rooms query fails, the Select SHALL display an error message with a retry control. The room/area field SHALL remain required in the create and update request body.

#### Scenario: Room select shows managed rooms on success

- **WHEN** a user opens the schedule-card create or edit modal and the rooms query resolves
- **THEN** the Room / Area field renders as a Select populated with room names from the backend response
- **AND** the user can select a room from the dropdown instead of typing free-form text

#### Scenario: Room select shows loading state while rooms fetch

- **WHEN** a user opens the schedule-card modal and the rooms query is pending
- **THEN** the Room / Area field shows a loading placeholder
- **AND** the rest of the form remains interactive

#### Scenario: Room select shows error state on failure

- **WHEN** the rooms query fails
- **THEN** the Room / Area field displays an error message indicating rooms could not be loaded
- **AND** the field provides a retry control for the rooms query

#### Scenario: Selected room name is submitted as roomArea

- **WHEN** a user selects a room from the dropdown and submits the form
- **THEN** the submitted `roomArea` value is the selected room's name string

### Requirement: The schedule-card create/edit SHALL use the shared DatePicker for scheduled date

The `/calendar` route SHALL replace the free-form text Input for scheduled date in the schedule-card modal with the existing shared `DatePicker` component from `@/shared/ui/date-picker`. The DatePicker SHALL display the current `formScheduledDate` value (from edit prefill or default) and emit date changes in `YYYY-MM-DD` format matching the backend contract.

#### Scenario: DatePicker renders with prefill value in edit mode

- **WHEN** a user opens the edit modal for an existing schedule card
- **THEN** the DatePicker displays the card's existing scheduled date
- **AND** the date can be changed through the DatePicker calendar popover

#### Scenario: DatePicker emits YYYY-MM-DD on selection

- **WHEN** a user selects a date through the DatePicker
- **THEN** the underlying `formScheduledDate` value is set to the selected date in `YYYY-MM-DD` format
