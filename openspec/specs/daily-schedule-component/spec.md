# daily-schedule-component Specification

## Purpose
TBD - created by archiving change task-ad30ae03-429d-4556-93db-b3f37727e49c. Update Purpose after archive.
## Requirements
### Requirement: Daily Schedule add-task affordances SHALL emit future-ready add intents without scheduling persistence

If `frontend/src/calendar/DailySchedule.vue` is updated in this change, the component SHALL keep its header `Add task` control and per-day `+ Add task` affordances visible while emitting typed add intents instead of owning planner persistence. The header-level add control SHALL emit a generic add-task intent. Per-day add affordances SHALL emit the selected day's `date` so future planner handlers can open the reusable add-operation modal with day context. This change SHALL NOT introduce scheduling mutations, automatic modal opening, or any persisted daily planner workflow.

#### Scenario: Header add control emits a generic add intent

- **WHEN** a user activates the Daily Schedule header `Add task` control
- **THEN** the component emits a typed add-task intent for the parent to handle
- **AND** the component does not perform a scheduling mutation itself

#### Scenario: Per-day add control emits day context

- **WHEN** a user activates a per-day `+ Add task` affordance
- **THEN** the component emits a typed add-task intent that includes that column's `date`
- **AND** the parent can decide whether to open the shared add-operation modal or ignore the event

#### Scenario: Daily Schedule remains persistence-free in this change

- **WHEN** the component source and tests are inspected after this change
- **THEN** the component still does not create scheduled tasks or call planner write endpoints
- **AND** any new behavior is limited to event emission for future reuse
