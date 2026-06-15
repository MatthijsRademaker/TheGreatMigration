# task-backlog-api Specification (Delta)

## Purpose

Update the canonical task backlog specification so `GET /api/tasks/backlog` remains contract-stable while moving from in-memory demo data to Postgres-backed storage.

## MODIFIED Requirements

### Requirement: Seed data SHALL be in-memory and exercise all vocabulary values

The endpoint handler SHALL read task backlog data from a `Store` interface backed by Postgres tables and sqlc-generated queries. Seed data inserted through migrations SHALL reproduce the current deterministic backlog payload exactly enough to preserve:

- the stable task IDs `task-1` through `task-11`;
- all three priority values (`high`, `medium`, `low`);
- all three backlog status values (`backlog`, `ready`, `assigned`);
- at least four distinct room names; and
- the current mix of empty, partially filled, and fully filled `assignedTo` arrays relative to `peopleNeeded`.

The response summary SHALL continue to be derived from the returned task rows rather than stored as separate persisted counts.

#### Scenario: Seed data preserves the current deterministic backlog rows
- **WHEN** `GET /api/tasks/backlog` is called against seeded Postgres data
- **THEN** `tasks` contains the same stable IDs and seeded assignment variety currently defined by the in-memory backlog seed

#### Scenario: Database-backed rows still exercise all canonical vocabulary values
- **WHEN** the endpoint is called
- **THEN** across all returned tasks, the priorities `high`, `medium`, and `low` each appear at least once
- **AND** the statuses `backlog`, `ready`, and `assigned` each appear at least once
- **AND** at least one task has an empty `assignedTo` array
- **AND** at least one task has a non-empty `assignedTo` with `len(assignedTo) < peopleNeeded`

### Requirement: Existing behavior SHALL remain intact

The task backlog endpoint SHALL remain read-only and additive. `GET /api/hello`, `GET /api/planning-window`, and `GET /api/dashboard/people-availability` SHALL continue to return their expected responses, and `GET /api/tasks/backlog` SHALL preserve its existing JSON contract while switching to Postgres-backed storage.

#### Scenario: Task backlog contract remains stable after persistence wiring
- **WHEN** `GET /api/tasks/backlog` is called after the refactor
- **THEN** the response still contains `summary`, `tasks`, `priorities`, and `statuses`
- **AND** the task rows and derived summary counts still satisfy the existing contract invariants

### Requirement: Backend tests SHALL cover the new endpoint

Backend tests in `backend/main_test.go` SHALL continue to verify the existing task backlog contract, and they SHALL also cover Store-backed success and failure paths for the backlog handler. Integration tests SHALL validate the backlog contract against the seeded Postgres database.

#### Scenario: MockStore-backed task backlog tests pass
- **WHEN** `go test ./...` runs in `backend/`
- **THEN** task backlog tests verify the existing contract invariants using a Store-backed handler registration

#### Scenario: Real-Postgres task backlog integration test passes
- **WHEN** `go test -tags=integration ./...` runs in `backend/`
- **THEN** the integration suite verifies the seeded backlog payload and summary invariants through `GET /api/tasks/backlog`

## ADDED Requirements

### Requirement: The Store interface SHALL be injected into the backlog handler

The `registerTasksBacklog` function SHALL accept a `Store` parameter. The handler closure SHALL call `store.GetTaskBacklog(ctx)` and SHALL return a `huma.Error500InternalServerError` if the Store call fails.

#### Scenario: Backlog handler delegates to Store on success
- **WHEN** the handler is called and the Store returns backlog data
- **THEN** the response body matches the Store-backed backlog payload

#### Scenario: Backlog handler returns 500 on Store failure
- **WHEN** the handler is called and the Store returns an error
- **THEN** the response is a 500 Internal Server Error
