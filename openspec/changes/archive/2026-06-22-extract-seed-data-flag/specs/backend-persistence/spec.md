## MODIFIED Requirements

### Requirement: Goose migrations SHALL create the persistence schema

The existing migrations for `planning_windows`, `people`, and `availability` SHALL remain intact. Additional goose migrations SHALL extend the schema with separate read-model tables for:

- backlog tasks, including stable task identifiers, titles, canonical priority values, required people counts, room names, canonical backlog statuses, and deterministic ordering;
- backlog task assignments, including foreign keys to both backlog tasks and `people(id)`;
- daily-schedule task cards, including stable seeded identifiers, titles, canonical priority values, room-area labels, required people counts, a scheduled date, and deterministic ordering; and
- daily-schedule task assignments, including foreign keys to both schedule task cards and `people(id)`.

The schema migrations SHALL enforce the existing canonical vocabularies at the database level. The schema migration set SHALL contain **structure only** — it SHALL NOT insert demo rows. Demo seed data is provided by a separate, flag-gated dataset (see "Demo seed data SHALL be a separate, flag-gated goose dataset"). Schema migration version numbers SHALL NOT be renumbered when the seed migrations are removed; goose tolerates the resulting version gaps on a fresh database.

#### Scenario: Existing persistence tables are preserved while backlog and schedule tables are added
- **WHEN** the schema goose migrations are applied to an empty Postgres database
- **THEN** the `planning_windows`, `people`, and `availability` tables still exist with their current columns
- **AND** separate persisted read models exist for backlog tasks and daily-schedule task cards with assignment tables referencing `people(id)`

#### Scenario: Schema-only migration leaves domain tables empty
- **WHEN** the schema goose migrations are applied to an empty Postgres database and the seed dataset is NOT applied
- **THEN** the `people`, `backlog_tasks`, `rooms_areas`, and `schedule_task_cards` tables contain zero rows
- **AND** the backend starts and serves HTTP

## ADDED Requirements

### Requirement: Demo seed data SHALL be a separate, flag-gated goose dataset

Demo seed data SHALL live in `backend/seed/` as its own goose dataset, embedded via a dedicated `//go:embed` directive and tracked in a goose version table distinct from the schema table (`goose_seed_version`). The seed dataset SHALL seed the exact backlog IDs `task-1` through `task-11`, the eight demo people (`p1`–`p8`) and their availability, the eight rooms/areas (`room-1`–`room-8`), and the deterministic daily-schedule cards and assignments needed to reproduce the current dashboard contract. Because the seed dataset is applied after the full schema, schedule-card seed rows SHALL set the `scheduled_date` column directly (the `day_group` column no longer exists). After inserting, the seed dataset SHALL advance the `people_id_seq`, `rooms_areas_id_seq`, and backlog-task ID sequences past the seeded maximum so subsequently created entities receive non-colliding sequential IDs.

#### Scenario: Seed dataset reproduces the deterministic demo responses
- **WHEN** the schema migrations are applied and then the seed dataset is applied
- **THEN** the backlog task read model contains the same 11 seeded tasks and assignment variety
- **AND** the schedule read model contains the deterministic seeded task-card data and assignments needed to reproduce the current default schedule behavior
- **AND** eight demo people with their availability and eight rooms/areas exist

#### Scenario: Seed dataset advances ID sequences past seeded rows
- **WHEN** the seed dataset has been applied and a new person, room, and backlog task are created
- **THEN** they receive the next sequential identifiers (`p9`, `room-9`, `task-12`) with no collision against seeded IDs

#### Scenario: Seed dataset is tracked independently of schema
- **WHEN** the backend starts twice against a database where schema and seed have both been applied
- **THEN** neither the schema pass nor the seed pass re-applies any migration, and startup succeeds both times
