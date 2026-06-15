## 1. Persistence schema and seed data

- [ ] Extend `backend/migrations/` with Postgres tables for backlog tasks, task assignments, schedule task cards, and schedule task assignments using separate read-model families for backlog and daily schedule.
- [ ] Seed the new tables from the current in-memory backlog and daily-schedule demo data so the existing deterministic responses, vocabulary coverage, assignment variety, and stable IDs are preserved.

## 2. Store, queries, and generated code

- [ ] Add sqlc query definitions and committed generated code for task backlog reads, schedule reads, and daily availability counts.
- [ ] Extend `backend/store.go` and `backend/store_mock_test.go` with `GetTaskBacklog` and `GetDailySchedule`, and implement the new methods in `PgStore` and `MockStore`.

## 3. Handler refactor

- [ ] Refactor `registerTasksBacklog` to accept `Store` and return the existing task backlog contract from Postgres-backed data.
- [ ] Refactor `registerDailySchedule` to accept `Store`, default omitted `start` from `GetPlanningWindow`, derive `availablePeopleCount` from availability data, and preserve the current schedule task-card contract.
- [ ] Update `backend/main.go` and test harness registration so all four read endpoints are wired through the Store-backed path.

## 4. Verification

- [ ] Expand unit tests to cover Store-backed success and failure paths for backlog and daily-schedule handlers without requiring Postgres.
- [ ] Expand integration tests to validate planning window, people availability, task backlog, daily schedule, and `/openapi.json` path coverage against a migrated Postgres database.
- [ ] Run `scripts/precommit-run` and ensure the Docker-backed verification flow passes.