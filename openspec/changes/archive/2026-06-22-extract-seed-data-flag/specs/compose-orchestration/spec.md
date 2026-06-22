## MODIFIED Requirements

### Requirement: Backend SHALL depend on the database service

The backend service in `compose.yml` SHALL receive a `DATABASE_URL` environment variable set to `postgres://app:app@db:5432/the_great_migration?sslmode=disable`. The backend service SHALL declare `depends_on` with `db: service_healthy` so that the backend waits for Postgres to be ready before starting. The backend service SHALL also receive `DB_AUTO_MIGRATE=true` and `DB_SEED=true`, so that the dev stack applies both the schema and the demo seed data on startup.

#### Scenario: Backend waits for database before starting
- **WHEN** `docker compose up` is executed
- **THEN** the backend service does not start until the `db` service passes its healthcheck

#### Scenario: Backend can connect to database over compose network
- **WHEN** the compose stack is running
- **THEN** the backend service can resolve and connect to `db:5432`

#### Scenario: Dev stack applies schema and seed data
- **WHEN** `docker compose up` is executed against an empty database volume
- **THEN** the backend applies the schema migrations and, because `DB_SEED=true`, the demo seed dataset
- **AND** the seeded demo data (people, tasks, schedule cards, rooms/areas) is present
