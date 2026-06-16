## Context

The calendar route (`CalendarView.vue`) and task management route (`TasksView.vue`) both allow creating task-like entities, but they target completely separate database tables with zero relationship:

- **Task Panel** → `backlog_tasks` (via `POST /api/tasks`) — the canonical task definition
- **Calendar** → `schedule_task_cards` (via `POST /api/schedule/cards`) — independent copies with no link back

Users must define the same task twice if they want it in both the backlog and the schedule, and the two copies can drift (different priority, different room, different staffing needs). This creates confusion and extra data entry.

The change introduces a **lightweight reference** from schedule cards to backlog tasks. The calendar modal shifts from free-form creation to selecting existing tasks, while still supporting one-off schedule cards without a backlog reference.

## Goals / Non-Goals

**Goals:**

- Calendar "Add Task" modal selects from existing backlog tasks instead of free-form creation
- Schedule cards can optionally reference a backlog task via `taskId`
- When `taskId` is provided, title/priority/room/peopleNeeded inherit from the referenced task at creation time (snapshot, not live reference)
- Existing free-form schedule card creation remains supported (for ad-hoc scheduling)
- Task panel CRUD is untouched

**Non-Goals:**

- Live sync between task edits and existing schedule cards (they remain snapshots)
- Bulk scheduling (multiple days at once)
- Drag-and-drop scheduling
- Changing the task panel's CRUD workflow
- Removing the free-form schedule card creation endpoint

## Decisions

### Decision 1: Snapshot copy vs. live reference

**Chosen: Snapshot copy with optional reference**

When a schedule card is created with `taskId`, the backend copies the task's current title, priority, room/area, and peopleNeeded into the schedule card at creation time. The `task_id` FK is stored so the UI can show "linked to task-3" and optionally navigate back.

**Rationale:**

- Schedule cards represent scheduling intent for a specific day, which may differ from the general task definition (e.g., you need 2 people for "Pack kitchen" in general, but only 1 is available on Tuesday)
- Live references would require cascading updates or force schedule cards to change when the task definition changes — surprising behavior
- The FK still enables traceability and UI affordances ("View task details") without coupling card behavior to task changes

**Alternatives considered:**

- **Live reference**: Schedule card fields always reflect the latest task state. Rejected because it breaks per-day adjustments.
- **No reference at all**: Current state. Rejected because it creates data duplication with no traceability.

### Decision 2: `taskId` as optional field, not replacement

**Chosen: Existing create/update fields remain; `taskId` is an optional addition**

The `POST /api/schedule/cards` request body keeps all existing fields (title, priority, roomArea, peopleNeeded, scheduledDate, assignedTo) and adds `taskId` as an optional string. When `taskId` is provided and no conflicting fields are supplied, the server copies values from the referenced task. If the caller supplies explicit values alongside `taskId`, the explicit values win (enabling per-day overrides).

**Rationale:**

- Backward compatible — existing API clients continue to work unchanged
- Enables the simplest possible frontend implementation: the calendar modal can still let users override fields if needed
- No breaking changes to the existing endpoint

**Alternatives considered:**

- **Replace fields with taskId**: Makes the endpoint single-purpose. Rejected because existing integration tests and potential future ad-hoc usage would break.
- **New endpoint `POST /api/schedule/cards/from-task`**: Unnecessary complexity. One endpoint with an optional field is simpler.

### Decision 3: Frontend modal becomes a task selector

**Chosen: Replace the full form with a searchable task selector + date picker + optional people override**

The calendar's `AddOperationModal` body changes from a full creation form to:

1. A searchable Select/Combobox populated from `useTaskBacklog()` data
2. The existing DatePicker (unchanged)
3. A read-only info panel showing the selected task's priority, room, and people needed
4. An optional people assignments section (same Checkbox grid, same as today)

When a task is selected, the modal shows the inherited values but the backend receives them implicitly via `taskId`. The user can still override assigned people for that day.

**Rationale:**

- Eliminates the data duplication UX
- Reuses the existing backlog query — no new API calls needed
- The task selector surfaces the same tasks the user already manages in the task panel
- Optionally overriding people per day is a reasonable customization without reintroducing full free-form fields

**Decision against auto-fill approach:** A simpler approach would keep the existing form and just pre-fill it when a task is selected. This is rejected because it would still allow the user to change title, priority, and room — reproducing the divergence problem. The selector approach enforces the separation: tasks are defined in the panel, scheduled in the calendar.

### Decision 4: No new DB migration — nullable FK

**Chosen: `schedule_task_cards.task_id` as nullable TEXT FK**

A new migration adds `task_id TEXT REFERENCES backlog_tasks(id)` to `schedule_task_cards`. NULL means the card has no backlog reference (existing cards and future ad-hoc cards). The FK ensures referential integrity — deleting a backlog task that has referencing schedule cards will fail at the DB level (the existing `DeleteTask` handler should either cascade or return an error with a clear message).

**Rationale:**

- Simple, standard pattern
- No default values needed for existing rows
- FK protects against orphan references
- TEXT type matches `backlog_tasks.id` (which is already TEXT)

## Risks / Trade-offs

| Risk | Severity | Mitigation |
|------|----------|------------|
| Users delete a backlog task that has schedule cards referencing it | Medium | DB FK prevents deletion; need to update the `DeleteTask` handler to check for references and return a clear error message (e.g., "Cannot delete task: it has 3 scheduled cards") |
| Backlog has many tasks, making the selector overwhelming | Low | Implement search/filter on the selector. The existing backlog query returns all tasks — filtering is frontend-only |
| Existing schedule cards have NULL task_id — users can't tell which task they came from | Low | Acceptable trade-off for backward compatibility. New cards will have the reference. A future migration could attempt to match by title, but not needed now |
| Schedule card fields (copied at creation) get stale if the task is later updated | Low | This is by design (Decision 1). The card represents a scheduling decision, not a live pointer |
