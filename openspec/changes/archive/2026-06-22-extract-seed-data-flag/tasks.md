## 1. Extract seed into its own dataset

- [x] 1.1 Create `backend/seed/` and move the three seed migrations into it, renumbered as a fresh
      goose sequence: `002_seed_demo_data.sql` → `001_seed_demo_data.sql`,
      `004_seed_backlog_schedule_data.sql` → `002_seed_backlog_schedule.sql`,
      `006_seed_rooms_areas.sql` → `003_seed_rooms_areas.sql`. Keep the goose `Up`/`Down` markers.
- [x] 1.2 Delete `002`, `004`, `006` from `backend/migrations/`. Do **not** renumber the remaining
      schema migrations — leave the version gaps (goose tolerates them on a fresh DB).
- [x] 1.3 Rewrite the schedule-card insert in `002_seed_backlog_schedule.sql` to the post-`009`
      schema: replace the `day_group` column with `scheduled_date` using the mapping in design.md
      (day_group 0→2026-07-05 … 3→2026-07-08). Leave `task_id` (added in `010`) unset.
- [x] 1.4 Append sequence-advance statements to the seed dataset (in the relevant seed files or a
      trailing `004_advance_sequences.sql`): `setval('people_id_seq', 8)`,
      `setval('rooms_areas_id_seq', 8)`, and the backlog-task sequence to 11 — matching the
      `setval(MAX...)` logic in migrations `007`/`011`/`008`. Confirm `schedule_task_cards`'
      `SERIAL` sequence needs no manual advance.

## 2. Backend startup wiring

- [x] 2.1 Add `//go:embed seed/*.sql` and a `seedFS embed.FS` alongside the existing
      `migrationsFS` in `backend/main.go`.
- [x] 2.2 Add `shouldSeed()` reading `DB_SEED`, defaulting to **false** (opt-in), mirroring the
      `shouldAutoMigrate()` parsing style.
- [x] 2.3 After the schema `goose.Up`, if `shouldSeed()` is true: `goose.SetTableName("goose_seed_version")`,
      run `goose.Up(sqlDB, "seed")`, then restore `goose.SetTableName("goose_db_version")`. Fail
      loudly (same pattern as the schema pass) on error.
- [x] 2.4 Verify the backend compiles and vets via `scripts/check`.

## 3. Compose

- [x] 3.1 Add `DB_SEED: "true"` to the `backend` service environment in `compose.yml` so
      `docker compose up` keeps its current seeded behavior.

## 4. Tests

- [x] 4.1 Update `backend/main_integration_test.go` (and any test asserting 11 tasks / 8 rooms /
      seeded people) to run the seed dataset explicitly after the schema `goose.Up` — using the
      `goose_seed_version` table — so seeded-count assertions still hold.
- [x] 4.2 Add a test that applies **schema only** (no seed pass) against an empty database and
      asserts the domain tables (`people`, `backlog_tasks`, `rooms_areas`, `schedule_task_cards`)
      are empty and the server starts.
- [x] 4.3 Add a test that, after the seed pass, creating a new person/room/task generates the next
      sequential ID (`p9`, `room-9`, `task-12`) — proving the sequence-advance fix prevents
      collisions.
- [x] 4.4 Verify with `scripts/check` and `scripts/test` (Docker-backed go vet/build/test).

## 5. Verification

- [x] 5.1 Run `scripts/precommit-run` and resolve failures.
- [x] 5.2 Manually confirm: `docker compose up` (dev) shows seeded data; the same backend image
      started with `DB_SEED` unset against a fresh DB serves with empty domain tables.
- [x] 5.3 Confirm all spec scenarios in `backend-persistence`, `backend-service`, and
      `compose-orchestration` are satisfied.
