## Why

Demo seed data (people, availability, backlog tasks, schedule cards, rooms/areas) is currently
welded into the goose migration chain as migrations `002`, `004`, and `006`. The backend embeds
`migrations/*.sql` and runs the whole folder unconditionally on startup, so **there is no way to
stand up a clean database without the demo data**. A production deployment needs the schema but
not the eight fictional people and eleven made-up tasks.

This change extracts the seed data out of the schema chain into a separate, flag-gated goose
dataset, so the schema always applies and the seed applies only when explicitly requested
(`DB_SEED=true`, the dev default). It is the prerequisite for a production compose variant.

## What Changes

- **Extract** the three seed migrations (`002_seed_demo_data.sql`, `004_seed_backlog_schedule_data.sql`,
  `006_seed_rooms_areas.sql`) out of `backend/migrations/` into a new `backend/seed/` directory as
  their own goose dataset with its **own version table** (`goose_seed_version`), embedded via a
  second `//go:embed`.
- **Gate** seed application behind a new `DB_SEED` environment variable. When `DB_SEED=true` (and
  auto-migrate is on), the backend runs a second `goose.Up` over the seed dataset after the schema
  migrations. When unset/false, no seed runs and the database contains schema only.
- **Rewrite** the extracted schedule-card seed to the *final* schema: because seed now runs after
  all schema migrations, the `schedule_task_cards` insert SHALL use the post-`009` `scheduled_date`
  column instead of the dropped `day_group` column.
- **Re-advance ID sequences** inside the seed dataset: after inserting seeded rows the seed SHALL
  `setval` the `people_id_seq`, `rooms_areas_id_seq`, and backlog-task sequence past the seeded
  maximum, because the schema's `setval` migrations now run against empty tables.
- **Wire the dev compose** `backend` service with `DB_SEED: "true"` so `docker compose up` keeps its
  current seeded behavior unchanged.
- **Update tests** that assume seed data is always present to apply the seed dataset explicitly, and
  add coverage proving a schema-only (no-seed) startup yields empty domain tables.

This change does **not** add the production compose file or the optimized frontend build — those are
follow-up changes that depend on this one.

## Capabilities

### Modified Capabilities

- `backend-persistence`: Seed data is no longer part of the schema migration set; it is a separate
  flag-gated dataset. The "seed migrations are applied" guarantee becomes conditional on `DB_SEED`.
- `backend-service`: Migration startup gains a second, conditional seed-dataset pass gated by
  `DB_SEED`, tracked in a distinct goose version table.
- `compose-orchestration`: The dev `backend` service explicitly sets `DB_SEED=true`; seeded demo
  data on `docker compose up` is now an explicit opt-in rather than an implicit side effect of
  migrating.

## Impact

- **Migrations**: `backend/migrations/` loses `002`, `004`, `006` (gaps at those version numbers are
  retained — goose tolerates non-contiguous versions on a fresh DB; existing schema migrations are
  **not** renumbered to avoid rewriting recorded version history). New `backend/seed/` dataset with
  re-numbered `001`/`002`/`003` seed files plus a trailing sequence-advance step.
- **Backend**: `backend/main.go` gains a second `//go:embed seed/*.sql`, a `DB_SEED` check, a
  `goose.SetTableName("goose_seed_version")`-scoped second `goose.Up`, and restores the schema table
  name afterward.
- **Tests**: `backend/main_integration_test.go` (and any test asserting 11 tasks / 8 rooms) applies
  the seed dataset explicitly; a new test asserts schema-only startup leaves domain tables empty.
- **Compose**: `compose.yml` `backend` environment gains `DB_SEED: "true"`.
- **No frontend, API, or schema-shape changes** beyond the seed-file rewrite to the final schema.
