## ADDED Requirements

### Requirement: Backend SHALL run the demo seed dataset only when DB_SEED is enabled

After the schema `goose.Up()` completes, the backend SHALL conditionally apply the demo seed
dataset embedded from `backend/seed/`. Seed application SHALL be gated by the `DB_SEED`
environment variable, defaulting to `false` (opt-in). When `DB_SEED` is truthy and auto-migrate
is enabled, the backend SHALL run a second `goose.Up()` over the seed dataset using a distinct
version table (`goose_seed_version`) and SHALL restore the schema version table name afterward.
If the seed pass fails, the backend SHALL log the error and exit. The seed pass SHALL NOT run when
`DB_AUTO_MIGRATE` is disabled.

#### Scenario: Seed runs when enabled
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=true` and `DB_SEED=true` against an empty database
- **THEN** the schema migrations apply, the seed dataset applies under `goose_seed_version`, and the backend serves HTTP with seeded demo data

#### Scenario: Seed is skipped by default
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=true` and `DB_SEED` unset against an empty database
- **THEN** the schema migrations apply, no seed dataset is applied, and the backend serves HTTP with empty domain tables

#### Scenario: Seed does not run when migrations are disabled
- **WHEN** the backend starts with `DB_AUTO_MIGRATE=false` and `DB_SEED=true`
- **THEN** neither the schema migrations nor the seed dataset are applied

#### Scenario: Seed pass is idempotent across restarts
- **WHEN** the backend restarts with `DB_SEED=true` against a database where the seed dataset is already applied
- **THEN** the seed `goose.Up()` completes without re-inserting rows and the backend proceeds normally
