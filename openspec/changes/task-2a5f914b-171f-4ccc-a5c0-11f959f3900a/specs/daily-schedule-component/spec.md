## ADDED Requirements

### Requirement: The calendar add-task modal SHALL let users select from existing backlog tasks

`CalendarView.vue` SHALL replace the free-form add-task modal with a modal that lets users search and select from existing backlog tasks. The modal SHALL query the backlog via the existing `useTaskBacklog()` composable and render a searchable selector populated with task titles. When a task is selected, the modal SHALL display the task's priority, room, and people needed as read-only info. The only editable fields SHALL be the scheduled date (via the existing DatePicker) and the assigned people (via the existing Checkbox grid).

#### Scenario: Calendar add modal shows backlog tasks in a selector

- **WHEN** a user opens the add-task modal on `/calendar`
- **THEN** the modal shows a searchable selector populated with backlog task titles from `useTaskBacklog()`
- **AND** the modal does not show free-form text inputs for title, priority, room, or people needed

#### Scenario: Selected task displays read-only inherited values

- **WHEN** a user selects a task from the backlog selector
- **THEN** the modal displays the task's priority, room/area, and people needed as read-only information
- **AND** the user can still pick a scheduled date and optionally assign people for that day

#### Scenario: Empty backlog shows a helpful message

- **WHEN** a user opens the add-task modal and the backlog is empty
- **THEN** the modal displays a message indicating no tasks exist yet
- **AND** the modal provides a link or guidance to create tasks in the task panel

### Requirement: Calendar task creation SHALL submit a `taskId` reference to the backend

When a user selects a task and submits the calendar add-task form, `CalendarView.vue` SHALL call the existing `createScheduleCardMutation` with the selected task's `id` passed as `taskId` in the request body. The backend SHALL resolve the referenced task's title, priority, room, and people needed at creation time. The `assignedTo` and `scheduledDate` fields SHALL come from the form as before.

#### Scenario: Calendar creates schedule card with taskId reference

- **WHEN** a user selects a backlog task, picks a date, assigns people, and submits
- **THEN** the frontend sends `POST /api/schedule/cards` with `taskId` set to the selected task's ID
- **AND** the created schedule card inherits the task's title, priority, room, and people needed

### Requirement: Calendar edit modal SHALL show the linked task and allow date/assignment changes

When editing an existing schedule card that has a `taskId` reference, the edit modal SHALL display the referenced task's title, priority, room, and people needed as read-only info. The user SHALL be able to change the scheduled date and adjust assigned people. If the schedule card has no `taskId` reference (legacy card), the edit modal SHALL fall back to the current free-form behavior.

#### Scenario: Edit modal shows linked task info for referenced cards

- **WHEN** a user edits a schedule card that has a non-null `taskId`
- **THEN** the modal shows the referenced task's details as read-only
- **AND** the date and assignment fields remain editable

#### Scenario: Edit modal shows full form for unreferenced cards

- **WHEN** a user edits a schedule card with a null `taskId`
- **THEN** the modal shows the existing free-form fields (title, priority, room, people needed, date, assignments)
- **AND** all fields are editable

## MODIFIED Requirements

### Requirement: The schedule-card create/edit SHALL use the shared DatePicker for scheduled date

**Change**: The DatePicker is now used in the task-selection modal context instead of the free-form modal. The behavior and integration remain the same.

The `/calendar` route SHALL use the existing shared `DatePicker` component from `@/shared/ui/date-picker` for the scheduled date field in both create and edit modes. The DatePicker SHALL display the current `formScheduledDate` value (from edit prefill or default) and emit date changes in `YYYY-MM-DD` format matching the backend contract.

#### Scenario: DatePicker renders with prefill value in edit mode

- **WHEN** a user opens the edit modal for an existing schedule card
- **THEN** the DatePicker displays the card's existing scheduled date
- **AND** the date can be changed through the DatePicker calendar popover

#### Scenario: DatePicker emits YYYY-MM-DD on selection

- **WHEN** a user selects a date through the DatePicker
- **THEN** the underlying `formScheduledDate` value is set to the selected date in `YYYY-MM-DD` format

## REMOVED Requirements

### Requirement: The schedule-card create/edit SHALL select room or area from the backend-managed room list

**Reason**: The calendar modal no longer has a free-form room/area field. Room/area is inherited from the selected backlog task at creation time. The room select behavior still exists in the task panel's create/edit modal, which remains unchanged.

**Migration**: Remove the room Select from `CalendarView.vue`'s add-task modal. The room Select component and its loading/error states remain in `TasksView.vue` for task CRUD.
