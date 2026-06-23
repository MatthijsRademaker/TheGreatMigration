## Context

The room/area is modeled as free text in two tables:

```
rooms_areas               backlog_tasks            schedule_task_cards
  id    TEXT PK             room  TEXT  ──┐          room_area  TEXT  ──┐
  name  TEXT                (no FK)       │          (no FK)            │
  type  CHECK(room|area)                  └── string that *happens* ────┘
                                              to match rooms_areas.name
```

The forms already populate their room Select from `GET /api/rooms` but bind `:value="room.name"` (`TasksView.vue:206`, and the schedule-card modal), so values are catalog-constrained at write time yet stored and compared by name. `rooms_areas.name` has **no UNIQUE constraint**, so name→id is not guaranteed 1:1.

On the read side, `area`/`roomArea` reaches `DailySchedule.vue`'s local `TaskCard` interface but is dropped by `TaskBoardCard.vue`, which neither declares nor renders it.

## Goals / Non-Goals

**Goals:**
- `backlog_tasks` and `schedule_task_cards` reference `rooms_areas(id)` via a real FK.
- Existing free-text values migrate without data loss; the catalog self-heals for orphans.
- The daily schedule and task board show the area name with a deterministic, stable color.
- Forms persist the area id.

**Non-Goals:**
- Filtering or grouping the schedule by area (display-only this change).
- A user-editable color per area (color is derived, not stored).
- Editing the area from the schedule card itself (still via the existing edit modal).
- Reworking room/area CRUD beyond the value the forms bind.

## Decisions

### Decision 1 — Replace the free-text column with `area_id` FK (no compat shim)

Per project rule 7 (no backwards compatibility), the old columns are dropped, not kept alongside. Migration order:

```sql
-- 1. add nullable area_id to both tables
ALTER TABLE backlog_tasks       ADD COLUMN area_id TEXT;
ALTER TABLE schedule_task_cards  ADD COLUMN area_id TEXT;

-- 2. backfill: match by name, dedup deterministically (see Decision 2)
-- 3. auto-create catalog rows for orphan strings (see Decision 2)
-- 4. enforce
ALTER TABLE backlog_tasks      ALTER COLUMN area_id SET NOT NULL;
ALTER TABLE backlog_tasks      ADD CONSTRAINT fk_backlog_area
  FOREIGN KEY (area_id) REFERENCES rooms_areas(id);
ALTER TABLE schedule_task_cards ALTER COLUMN area_id SET NOT NULL;
ALTER TABLE schedule_task_cards ADD CONSTRAINT fk_card_area
  FOREIGN KEY (area_id) REFERENCES rooms_areas(id);

-- 5. drop the old columns
ALTER TABLE backlog_tasks      DROP COLUMN room;
ALTER TABLE schedule_task_cards DROP COLUMN room_area;
```

### Decision 2 — Self-healing backfill (match → dedup → auto-create)

For each distinct legacy string `s` in `backlog_tasks.room` ∪ `schedule_task_cards.room_area`:

```
matches = rooms_areas WHERE name = s
  if matches has ≥1 row → area_id = MIN(id) over matches   (deterministic on dupes)
  if matches is empty   → INSERT rooms_areas(id, name, type)
                            VALUES (new 'room-<seq>', s, 'area') → use that id
```

`MIN(id)` makes the duplicate-name resolution deterministic and repeatable. Auto-create chooses `type='area'` (the safe generic) — an orphan string was never a validated room, so it becomes an area. No row is left without an `area_id`, which is what lets step 4 enforce `NOT NULL` + FK without failure.

**Alternatives considered:**
- *Fail the migration on any orphan* (fail-fast) — rejected: the existing data legitimately contains free-text strings the catalog never had; blocking deploy on manual cleanup is friction for no integrity gain once auto-create guarantees a valid target.
- *Map orphans to a single "Unassigned" sentinel* — rejected: collapses distinct real areas ("Garage", "Shed") into one, losing the very distinction this feature exists to surface.

### Decision 3 — Color derived by hashing `area.id`, computed in the frontend

No `color` column. A pure util maps the stable area id to a fixed palette:

```ts
// frontend/src/shared/lib/areaColor.ts
const PALETTE = ['#3b82f6', '#f97316', '#10b981', '#a855f7', '#ef4444', '#eab308', '#06b6d4', '#ec4899']
export function areaColor(id: string): string {
  let h = 5381
  for (let i = 0; i < id.length; i++) h = ((h << 5) + h) ^ id.charCodeAt(i) // djb2
  return PALETTE[Math.abs(h) % PALETTE.length]
}
```

Hashing the **id** (not the name) means renaming an area keeps its color. Color is a design concern, so it lives in the frontend; the API stays clean and returns only `{ id, name }`.

**Alternatives considered:**
- *`color` column on `rooms_areas`* — the user explicitly chose hash-derived; avoids a schema change and a CRUD-form field. Trade-off: you can't pin "Kitchen = orange"; the hash decides. Accepted.

### Decision 4 — Area chip placement avoids colliding with the priority accent

`TaskBoardCard.vue` already uses the left border (`border-l-4`) for the **priority** accent. The area color must not compete with it, so the area renders as an inline chip (a colored dot/pill + name) on its own line under the title:

```
┌─────────────────────────┐
│ Paint            ⚑ high │   ← priority badge (top-right), priority border (left edge)
│ ● Living Room           │   ← area chip: dot tinted via areaColor(area.id) + name
│ Anna, Ben               │
│ 1 / 2 — needs help      │
└─────────────────────────┘
```

The same chip component renders in the task-board row for consistency.

### Decision 5 — Read model returns a nested `area` object

Both read models replace the scalar with `area: { id, name }`. Frontend adapters map it; `TaskBoardCard` reads `area.name` for the label and `areaColor(area.id)` for the tint. Write inputs (`areaId`) replace the old `roomArea`/`room` string; the field stays required.

## Risks / Trade-offs

- **Backfill correctness is one-shot** — once old columns are dropped, a wrong name→id mapping (from pre-existing duplicate names) can't be recovered from the dropped column. Mitigation: `MIN(id)` is deterministic; duplicate names are rare in a small managed catalog; the `UNIQUE(name)` open question prevents recurrence.
- **Auto-created areas inherit `type='area'`** — a string that was conceptually a room becomes an area. Acceptable: it was never validated as a room, and type only affects the room/area distinction in the catalog, not the chip.
- **Palette collisions** — with 8 colors, >8 areas guarantees repeats. Acceptable for the current scale; the palette can grow without a contract change.
- **Color not stable across a palette edit** — changing `PALETTE` reshuffles all colors. Acceptable; it's a presentational constant, not persisted state.
