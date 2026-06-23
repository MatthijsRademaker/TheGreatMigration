## Why

Schedule cards and backlog tasks store their room/area as **free text** (`schedule_task_cards.room_area`, `backlog_tasks.room`), even though a managed catalog (`rooms_areas`) already exists. The create/edit forms already constrain users to catalog entries but persist the **name string** rather than the catalog `id` — a soft foreign key pretending to be free text. Two problems follow:

1. **The daily schedule never shows the area.** `roomArea` flows all the way to `<TaskBoardCard :task="task">`, but `TaskBoardCard.vue` drops it — its `BoardTaskCard` interface doesn't declare the field and the template never renders it. With multiple "Paint" cards scoped to different areas, the schedule is ambiguous: every card just says "Paint".
2. **Referential integrity is fiction.** Renaming or deleting a catalog room silently orphans every card/task that stored the old name; duplicate names are indistinguishable; typos are unvalidated.

This change makes the room/area a real foreign key and surfaces it visually on the schedule with a color-coded chip.

## What Changes

- **Real FK to `rooms_areas`** — replace `backlog_tasks.room TEXT` and `schedule_task_cards.room_area TEXT` with `area_id TEXT REFERENCES rooms_areas(id)`. A backfill migration matches each existing string to a catalog row by name (deterministic `MIN(id)` on duplicates) and **auto-creates** a `rooms_areas` row (`type='area'`) for any string with no match, so no data is dropped and the catalog self-heals.
- **Read models expose a nested area object** — the daily-schedule and task-backlog responses return `area { id, name }` instead of a bare `roomArea` / `room` string. Write endpoints accept `areaId`.
- **Color-coded area chip on schedule cards** — `TaskBoardCard.vue` renders the area name with a deterministic color derived by hashing `area.id` into a fixed palette (no new DB column, no CRUD change). The same chip is applied to task-board rows for cross-surface consistency.
- **Forms bind the area id** — the schedule-card and task create/edit Selects bind `:value="room.id"` instead of `room.name`.

## Capabilities

### New Capabilities
- `area-reference-integrity`: schedule cards and backlog tasks reference `rooms_areas` by FK; a self-healing backfill migrates existing free-text values without data loss.
- `schedule-area-chips`: the daily schedule and task board render a color-coded area chip, with color derived deterministically from the area id.

### Modified Capabilities
- `dashboard-daily-schedule`: task cards expose `area { id, name }` instead of a `roomArea` string; writes accept `areaId`.
- `task-backlog-api`: task rows expose `area { id, name }` instead of a `room` string; writes accept `areaId`.
- `daily-schedule-component`: the read-path adapter maps the area object; the schedule-card form binds the area id.
- `task-management-ui`: the task create/edit form binds the area id and the row displays the area chip.

## Impact

- `backend/migrations/` — new migration: add `area_id`, backfill (match/dedup/auto-create), drop old columns, add FK.
- `backend/db/queries/` + regenerated `*.sql.go` — joins to resolve `area.id`/`area.name`; writes take `area_id`.
- `backend/api/daily_schedule.go`, `backend/api/tasks.go` — read models return the area object; create/update inputs take `areaId`.
- `frontend/src/shared/lib/` — new `areaColor(id)` palette util.
- `frontend/src/calendar/components/TaskBoardCard.vue` — render area chip.
- `frontend/src/calendar/DailySchedule.vue`, `frontend/src/tasks/components/TaskRow.vue` — pass area through / render chip.
- `frontend/src/tasks/TasksView.vue` + the schedule-card modal — bind `:value="room.id"`.
- Tests asserting `roomArea` / `room` strings updated to the area object.
- **Breaking**: the `roomArea`/`room` string fields are removed from the API (no backwards-compat shim, per project rules).

## Open Questions

- Add `UNIQUE(name)` to `rooms_areas` after dedup to prevent the soft-FK ambiguity from reappearing? Only safe if there are no intentional duplicate names today.
