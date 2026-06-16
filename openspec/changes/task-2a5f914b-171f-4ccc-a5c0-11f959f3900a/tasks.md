## 1. Database Migration

- [ ] 1.1 Create migration `010_add_schedule_card_task_id.sql` that adds a nullable `task_id TEXT REFERENCES backlog_tasks(id)` column to `schedule_task_cards`
- [ ] 1.2 Update the seed data migration to set `task_id` for existing seeded schedule cards that match backlog task titles (optional — may be deferred to a follow-up)
- [ ] 1.3 Run `make generate` or equivalent to regenerate sqlc query artifacts

## 2. Backend — Schedule Card API Changes

- [ ] 2.1 Add `TaskId` field to `CreateScheduleCardInput` domain struct in `backend/api/daily_schedule.go`
- [ ] 2.2 Add `taskId` optional string field to `CreateScheduleCardRequestBody` (Huma request body struct) — keep all existing fields
- [ ] 2.3 Add `TaskId` field to `TaskCard` response struct (nullable `*string`, serialized as `"taskId"`)
- [ ] 2.4 Add task existence validation in `validateScheduleCardInput`: when `taskId` is non-empty, call `store.PersonExists`-style check (or a new `store.TaskExists()` method) and reject with 400 if the task doesn't exist
- [ ] 2.5 Modify the store's `CreateScheduleCard` implementation to persist `task_id` when provided, and to copy task fields (title, priority, room_area, people_needed) from the referenced backlog task when no explicit value is supplied for each field
- [ ] 2.6 Modify the store's `UpdateScheduleCard` implementation to support updating `task_id` (when provided, update the reference; when omitted, leave as-is)
- [ ] 2.7 Add `TaskExists(ctx, taskId string) (bool, error)` to the Store interface and implement it via a sqlc query on `backlog_tasks`
- [ ] 2.8 Update the daily-schedule read query (`GetDailyScheduleTaskCards`) or the Go adapter to include the `task_id` column from `schedule_task_cards` in the response models
- [ ] 2.9 Update the `TaskCard` adapter in `store.go` (the `GetDailySchedule` method) to populate the `TaskId` field from the persisted `task_id` column

## 3. Backend — Task Delete with Reference Check

- [ ] 3.1 Add `TaskHasScheduleCards(ctx, taskId string) (bool, error)` to the Store interface
- [ ] 3.2 Add sqlc query to check for referencing `schedule_task_cards` rows by `task_id`
- [ ] 3.3 Modify the `DeleteTask` handler to check for schedule card references before deleting; if references exist, return a 400 error with a clear message (e.g., "Cannot delete task: it has N scheduled cards. Remove them first.")

## 4. Backend — Artifacts and Tests

- [ ] 4.1 Regenerate the OpenAPI snapshot (`backend/openapi-snapshot.json`) to include the new `taskId` field on schedule card request/response bodies
- [ ] 4.2 Regenerate the frontend typed client (`frontend/src/client/`) from the updated OpenAPI snapshot
- [ ] 4.3 Add backend unit tests for `CreateScheduleCard` with `taskId`: create with valid taskId, create with invalid taskId (400 rejected), create with taskId + explicit field overrides
- [ ] 4.4 Add backend unit tests for delete with task references: verify that deleting a referenced task returns a clear error
- [ ] 4.5 Add backend integration tests for the taskId field in daily-schedule reads

## 5. Frontend — Calendar Add-Task Modal

- [ ] 5.1 Replace the free-form form body in `CalendarView.vue`'s add-task modal with a task-selector UI: import and use `useTaskBacklog()` data, render a searchable Select/Combobox populated with backlog task titles
- [ ] 5.2 Add state for the selected task ID (`formTaskId`) and a computed read-only display showing the selected task's priority, room, and people needed
- [ ] 5.3 Wire the submit handler to pass `taskId` (and `scheduledDate`, `assignedTo`) to `createScheduleCardMutation` — omit title/priority/roomArea/peopleNeeded from the body when taskId is set (backend inherits them)
- [ ] 5.4 Handle the empty backlog case: show a message with a link/guidance to create tasks in the task panel
- [ ] 5.5 Add search/filter capability to the task selector so users can find tasks by typing

## 6. Frontend — Calendar Edit Modal

- [ ] 6.1 In the edit path, detect whether the selected schedule card has a non-null `taskId`: if yes, show read-only task info (title, priority, room, people needed) and editable date/assignment fields; if no, fall back to the current full free-form edit modal
- [ ] 6.2 Update the edit submit handler: for taskId-linked cards, send `taskId`, `scheduledDate`, and `assignedTo` (omit title, priority, roomArea, peopleNeeded) to `updateScheduleCardMutation`

## 7. Frontend — Schedule Card Display

- [ ] 7.1 Update `DailySchedule.vue` or the schedule card template to show a small badge or indicator when a card has a `taskId` (e.g., "From backlog" with a link-style affordance)
- [ ] 7.2 Update the `ScheduleTaskCard` local interface in `useDailySchedule.ts` to include `taskId: string | null`
- [ ] 7.3 Update the adapter in `useDailySchedule.ts` to pass through the `taskId` field

## 8. Frontend — Tests

- [ ] 8.1 Add/update `CalendarView` tests: verify the add-task modal renders with backlog task options, verify create mutation sends `taskId`, verify empty backlog shows guidance message
- [ ] 8.2 Add/update edit modal tests: verify taskId-linked cards show read-only info, verify unreferenced cards show full form
- [ ] 8.3 Update route-render tests to cover the new modal behavior
