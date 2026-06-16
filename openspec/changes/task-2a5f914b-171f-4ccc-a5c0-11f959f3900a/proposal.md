## Why

The calendar's "Add Task" modal is a duplicate of the task panel's creation form. Both allow free-form entry of the same data (title, priority, room, people needed, assignments) — the only difference is the calendar adds a date picker. This creates two independent data sources (`backlog_tasks` and `schedule_task_cards`) with no relationship between them, forcing duplicate data entry and allowing task definitions to diverge across surfaces.

The task panel (`TasksView.vue`) is the right place to define what work exists. The calendar should be a pure scheduling surface: pick from the defined backlog, choose a date, and optionally adjust per-day assignments.

## What Changes

- **Calendar "Add Task" modal**: Replace the free-form creation form with a task selector that lets users pick from existing backlog tasks
- **Schedule card creation with task reference**: `POST /api/schedule/cards` gains an optional `taskId` field; when provided, the card inherits title, priority, room, and people needed from the referenced task
- **Schedule card display in calendar**: Show the referenced task's status and allow navigating back to the task panel
- **Existing task panel remains unchanged**: Full CRUD continues as-is — this is the canonical task definition surface
- **Existing free-form schedule card creation remains supported**: For cases where a card needs to exist without a backlog task (e.g., ad-hoc or one-off scheduling)

## Capabilities

### New Capabilities

None — this is a modification of existing capabilities, not introduction of a new domain.

### Modified Capabilities

- `daily-schedule-component`: The calendar's add-task workflow changes from free-form creation to existing-task selection. The frontend calendar modal picks from the backlog query and passes a `taskId` to the create endpoint.
- `dashboard-daily-schedule`: The schedule-card create/update endpoint adds an optional `taskId` body field. When provided, title/priority/room/peopleNeeded are inherited from the referenced backlog task. Validation ensures the referenced task exists.
- `task-management-ui`: The task panel stays as-is — no changes. Its existing CRUD remains the canonical task definition surface.

## Impact

| Layer | Impact |
|-------|--------|
| **Backend** | `POST /api/schedule/cards` and `PUT /api/schedule/cards/{id}` gain an optional `taskId` field. `schedule_task_cards` table gains a nullable `task_id` FK to `backlog_tasks`. New validation: referenced task must exist. No existing fields removed. |
| **Frontend** | `CalendarView.vue`: modal body replaced with task selector (searchable, using backlog query data). `DailySchedule.vue`: schedule card display optionally shows linked task info. No changes to `TasksView.vue`. |
| **API artifacts** | OpenAPI snapshot and generated frontend client refreshed to include the new `taskId` field on schedule card bodies. |
| **DB** | New migration: `ALTER TABLE schedule_task_cards ADD COLUMN task_id TEXT REFERENCES backlog_tasks(id)`. Existing cards get `NULL`. |
