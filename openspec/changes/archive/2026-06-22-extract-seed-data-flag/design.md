# Design — Extract seed data behind a flag

## Context

`backend/main.go` does `//go:embed migrations/*.sql` and `goose.Up(sqlDB, "migrations")`,
running the entire folder unconditionally. Seed data lives *inside* that folder as goose
migrations `002`, `004`, `006`, interleaved with schema migrations. There is no runtime toggle:
`DB_AUTO_MIGRATE=false` skips **all** migrations including schema, which is useless for a fresh DB.

## Goal

Schema always applies; seed applies only when `DB_SEED=true`. Seed must be physically extracted
from the schema chain (not just gated in-place) so the coupling is removed for good.

## Two non-obvious constraints discovered

### 1. Seed is order-dependent on schema (the `day_group` drop)

`004` seeds `schedule_task_cards` with a `day_group` column. Migration `009` later **drops
`day_group`** and replaces it with a `NOT NULL scheduled_date` backfilled *from* `day_group`.

Consequence: the naive "run all seed as one blob after all schema" fails — by then `day_group`
no longer exists. Because we run seed after the full schema, the extracted schedule-card seed
**must be rewritten to the final schema**: insert `scheduled_date` directly.

Mapping (planning window starts `2026-07-05`, `day_group` → `start + N days`):

| day_group | scheduled_date |
|-----------|----------------|
| 0         | 2026-07-05     |
| 1         | 2026-07-06     |
| 2         | 2026-07-07     |
| 3         | 2026-07-08     |

### 2. ID sequences are set from existing rows (the sequence trap)

Migrations `007` (`rooms_areas_id_seq`), `011` (`people_id_seq`), and the backlog-task sequence
do `setval(seq, COALESCE(MAX(existing rows), 0))`. Today seed runs *before* them, so the sequence
advances past seeded IDs.

After extraction, schema (including those `setval`s) runs first against **empty** tables → each
sequence sits at 0. Then seed inserts hardcoded `p1..p8`, `room-1..room-8`, `task-1..task-11`
without advancing the sequence → the next `CreatePerson` regenerates `p1` and **collides**.

Fix: the seed dataset SHALL re-run `setval` for `people_id_seq`, `rooms_areas_id_seq`, and the
backlog-task sequence after inserting, advancing them past the seeded maximum. `schedule_task_cards`
uses a `SERIAL` whose sequence auto-advances on insert, so it needs no manual `setval`.

#### 2a. Schema `setval(seq, 0)` crashes on empty tables (discovered during implementation)

Constraint #2 was incomplete. Migrations `007` and `011` do
`setval(seq, COALESCE(MAX(...), 0))`. They previously ran *after* seed (table
populated, `MAX = 8`). After extraction they run *first* against **empty**
tables → `setval(seq, 0)`, which Postgres rejects (`value 0 is out of bounds`,
sequence min is 1). The seed dataset cannot fix this — the schema `goose.Up`
crashes before the seed pass runs.

Fix: `007` and `011` were rewritten to tolerate empty tables —
`setval(seq, GREATEST(COALESCE(MAX,0),1), EXISTS(SELECT 1 FROM table))`. On an
empty table this leaves the sequence at its start (`is_called=false`, next
`nextval` → 1); on a populated table it is identical to the old 2-arg form
(`is_called=true`, next → `MAX+1`). This is required by the new
"schema-only leaves domain tables empty AND the backend starts" scenario.

#### 2b. Migration `012` (tools) had the same seed coupling (discovered during implementation)

`012_create_tools_table.sql` (added after this change was authored) baked demo
tool rows referencing seeded people `p1`/`p3` into a schema migration — the
exact coupling this change removes. Schema-only mode crashed on its FK to the
(empty) `people` table. The tool `INSERT` + `setval('tools_id_seq', 5)` were
extracted into `seed/005_seed_tools.sql`; `012` is now structure-only
(`CREATE TABLE` + `CREATE SEQUENCE`).

## Decision: separate goose dataset with its own version table

```
backend/migrations/   schema only — 001,003,005,007,008,009,010,011,012 (gaps kept)
                      embed migrations/*.sql → goose.Up(table: goose_db_version)   ALWAYS

backend/seed/         001_seed_demo_data.sql        (was 002, unchanged)
                      002_seed_backlog_schedule.sql (was 004, rewritten to scheduled_date)
                      003_seed_rooms_areas.sql      (was 006, unchanged)
                      + sequence-advance (in the relevant seed files or a trailing step)
                      embed seed/*.sql → goose.Up(table: goose_seed_version)   IF DB_SEED
```

Why a **separate version table**: goose tracks applied versions in one table by version id. Two
datasets sharing `goose_db_version` would collide (schema `001` vs seed `001`). `goose.SetTableName`
(reset afterward) or the provider API gives the seed pass its own `goose_seed_version` table, so
the two datasets are tracked independently and each is idempotent on restart.

Why **keep gaps** in `migrations/` instead of renumbering `003→002`, etc.: renumbering rewrites
recorded version ids and breaks every already-migrated dev database. Goose applies pending
migrations in version order and tolerates gaps on a fresh DB. Surgical > tidy.

### startup flow in main.go

```
if shouldAutoMigrate() {
    goose.Up(schemaFS)                       // always; table goose_db_version
    if shouldSeed() {                        // DB_SEED truthy
        goose.SetTableName("goose_seed_version")
        goose.Up(seedFS)                     // demo data, idempotent
        goose.SetTableName("goose_db_version")  // restore
    }
}
```

`shouldSeed()` mirrors `shouldAutoMigrate()` but **defaults to false** (seed is opt-in), whereas
auto-migrate defaults to true.

## Alternatives rejected

- **Path X — in-place SQL guard**: keep seed at `002/004/006`, set an `app.seed` GUC on the
  connection, wrap each seed file's inserts in `IF current_setting('app.seed')`. Preserves order
  with no rewrite, but leaves seed permanently welded into the schema chain and makes the flag
  invisible plpgsql magic. Rejected: does not achieve extraction, which is the point.
- **Plain DML seed (no goose)**: seed inserts are not idempotent (no `ON CONFLICT`), so re-running
  on container restart would error. goose's version tracking is exactly what prevents that.
  Rejected in favor of a goose dataset with its own table.

## Risks

- **Existing dev volumes**: a `pgdata` volume that already recorded `002/004/006` in
  `goose_db_version` will now find those versions missing from the schema folder. On a fresh
  `docker compose down -v && up` this is moot; for persisted volumes, document that the seed pass
  re-applies under `goose_seed_version` and the orphaned schema entries are harmless (goose only
  applies *pending* versions). Acceptable for a dev-only convenience path.
- **Test fixtures**: any test calling `goose.Up(db, "migrations")` and asserting seeded counts must
  add the seed pass; covered in tasks.
