## Why

The backend already uses Postgres for `GET /api/planning-window` and `GET /api/dashboard/people-availability`, but `GET /api/tasks/backlog` and `GET /api/dashboard/daily-schedule` still serve their payloads from in-memory Go seed data. This leaves the read surface split across two storage models and conflicts with the new requirement to wire the recently added backend endpoints to Postgres-backed storage.

## What Changes

- Extend the existing pgx/sqlc/goose persistence layer to cover the remaining read endpoints: `GET /api/tasks/backlog` and `GET /api/dashboard/daily-schedule`.
- Add Postgres schema and seed data for backlog tasks, task assignments, schedule task cards, and schedule task assignments while preserving the current JSON contracts, canonical vocabularies, deterministic demo data, and stable seed IDs.
- Extend `Store`, `PgStore`, sqlc queries, generated code, and handler registration so backlog and daily-schedule handlers read through Store injection instead of directly from `seedTasks`, `seedTasksForDay`, `seedPeople`, or `planWindowStart`.
- Make the daily-schedule default `start` come from the Store-backed planning window when the query parameter is omitted.
- Keep `GET /api/planning-window` and `GET /api/dashboard/people-availability` on their existing Postgres-backed path and expand regression coverage so all four read endpoints are validated through MockStore-based unit tests and real-Postgres integration tests.
- Update the affected OpenSpec requirements that still describe task backlog and daily schedule as in-memory endpoints.

## Impact

Implementation stays within the current backend architecture and reuses the existing Postgres foundation instead of introducing a new persistence stack. Frontend consumers keep their current response shapes and legends, but the two remaining in-memory read endpoints become database-backed. The change also expands integration coverage and OpenAPI assertions so `/api/hello`, `/api/planning-window`, `/api/dashboard/people-availability`, `/api/tasks/backlog`, and `/api/dashboard/daily-schedule` are all verified together.
