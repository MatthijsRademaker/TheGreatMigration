## 1. Database: FK + self-healing backfill

- [x] 1.1 New migration: add nullable `area_id TEXT` to `backlog_tasks` and `schedule_task_cards`
- [x] 1.2 Backfill: for each distinct `room` / `room_area` string, match `rooms_areas.name` and set `area_id = MIN(id)` over matches
- [x] 1.3 Backfill: for orphan strings (no match), `INSERT rooms_areas(id, name, type='area')` using the room-id sequence, then link `area_id`
- [x] 1.4 Enforce: `ALTER COLUMN area_id SET NOT NULL` and add `FOREIGN KEY (area_id) REFERENCES rooms_areas(id)` on both tables
- [x] 1.5 Drop `backlog_tasks.room` and `schedule_task_cards.room_area`
- [x] 1.6 Write the `-- +goose Down` path (re-add text columns, backfill name from join, drop FK/area_id)
- [x] 1.7 Reorder seeds so `rooms_areas` (now `002`) seeds before backlog/schedule (now `003`); rewrite backlog/schedule seed to use `area_id` (FK consequence not in original plan)

## 2. Backend: queries + read/write models

- [x] 2.1 Update SQL queries for backlog + daily-schedule reads to join `rooms_areas` and select `area_id` and `rooms_areas.name`
- [x] 2.2 Update create/update SQL for both tables to take `area_id`
- [x] 2.3 Regenerate sqlc via `scripts/generate-proto` (or the sqlc step) and confirm generated structs
- [x] 2.4 `backend/api/daily_schedule.go`: replace `RoomArea string` with `Area { ID, Name }` in `TaskCard`; create/update input takes `AreaId`
- [x] 2.5 `backend/api/tasks.go`: replace `Room string` with `Area { ID, Name }` in the task row; create/update input takes `AreaId`
- [x] 2.6 Validation: `areaId` required and SHALL reference an existing `rooms_areas` row (via new `AreaExists`); return 400 on missing/unknown id

## 3. Frontend: color util + chip

- [x] 3.1 Add `frontend/src/shared/lib/areaColor.ts` (djb2 hash of id → fixed palette) with a unit test
- [x] 3.2 `TaskBoardCard.vue`: add `area { id, name }` to the card interface; render a colored area chip on its own line, not on the priority left-border
- [x] 3.3 `DailySchedule.vue`: thread `area` through its local `TaskCard` interface (replace `roomArea`)
- [x] 3.4 `TaskRow.vue`: render the same area chip for cross-surface consistency

## 4. Frontend: read-path adapter + forms

- [x] 4.1 Regenerate the API client (`npm run regen:api`) so `area` / `areaId` appear in generated types
- [x] 4.2 Daily-schedule adapter: map the generated `area` object into the component contract; sanitize a missing/null area to a stable fallback
- [x] 4.3 Schedule-card modal: replace free-text Input with a rooms Select bound to `:value="room.id"`; submit `areaId`
- [x] 4.4 `TasksView.vue`: bind the room Select to `:value="room.id"`; submit `areaId`

## 5. Tests

- [x] 5.1 Backend: migration backfill integration test — name match, duplicate-name `MIN(id)`, orphan auto-create (shared across both tables), FK enforced
- [x] 5.2 Backend: read models return `area { id, name }`; writes accept and persist `areaId`; unknown id rejected (unit + integration)
- [x] 5.3 Frontend: `areaColor` is deterministic and stable per id
- [x] 5.4 Frontend: `TaskBoardCard`/`TaskRow` render the area chip; existing `roomArea`/`room` fixtures/assertions updated to the area object
- [x] 5.5 Frontend: forms submit `areaId` (not name)

## 6. Verification

- [x] 6.1 `scripts/precommit-run` passes — `project check` (go vet/build/test, openapi snapshot, dashboard check/build, proto) and `project test` (frontend + backend) green; also verified `scripts/test` and `scripts/test-integration`
