# people-management-api Specification

## Purpose
TBD - created by archiving change task-9d7ba566-1e1d-4230-8f96-a108ed06490a. Update Purpose after archive.
## Requirements
### Requirement: Backend SHALL expose a minimal CRUD surface for people records

The backend SHALL register `POST /api/people`, `PUT /api/people/{id}`, and `DELETE /api/people/{id}` through Huma and include them in `/openapi.json`. The `POST /api/people` endpoint SHALL assign a server-generated sequential ID using the `'p' || nextval('people_id_seq')` pattern. The request body SHALL accept only `name` and `initials`. Successful create and update operations SHALL make the person visible through `GET /api/dashboard/people-availability`. Deleting a person SHALL succeed only when backlog and schedule assignment references do not block the delete.

#### Scenario: OpenAPI includes the people write endpoints

- **WHEN** the backend OpenAPI document is fetched
- **THEN** it includes write operations for `/api/people` and `/api/people/{id}` alongside the existing dashboard availability read path
- **AND** the `POST /api/people` request body schema contains only `name` and `initials` — no `id` field

#### Scenario: Created people receive a server-assigned sequential ID

- **WHEN** a person is created via `POST /api/people` with only `name` and `initials`
- **THEN** the response contains a full person record with `id` set to a value matching the pattern `p` followed by an integer (e.g. `p9`, `p10`)
- **AND** the `id` is one greater than the highest existing person ID at the time of creation

#### Scenario: Created or updated people appear in dashboard availability reads

- **WHEN** a person is created or updated successfully
- **THEN** a subsequent `GET /api/dashboard/people-availability` reflects that person in the `people` array

#### Scenario: Referenced people cannot be deleted

- **WHEN** `DELETE /api/people/{id}` is called for a person referenced by `backlog_task_assignments` or `schedule_task_assignments`
- **THEN** the endpoint returns `409 Conflict`
- **AND** the person remains present in subsequent availability reads

### Requirement: Backend SHALL expose single-date availability write endpoints within the planning window

The backend SHALL register `PUT /api/people/{id}/availability/{date}` to create or update a single availability entry and `DELETE /api/people/{id}/availability/{date}` to remove one entry. The write surface SHALL accept only the canonical statuses `available`, `busy`, `partial`, and `off`, SHALL reject malformed or out-of-window dates with `400`, and SHALL return `404` when the addressed person does not exist.

#### Scenario: Upserting one status changes the dashboard availability matrix
- **WHEN** a person’s status is upserted successfully for a planning-window date
- **THEN** a subsequent `GET /api/dashboard/people-availability` returns that status for the matching person and date
- **AND** the response summary remains consistent with the returned per-person rows

#### Scenario: Invalid status or date is rejected
- **WHEN** the availability write endpoint receives a status outside `available`, `busy`, `partial`, or `off`, or a malformed or out-of-window date
- **THEN** the endpoint returns `400`

#### Scenario: Missing people return 404 for availability writes
- **WHEN** an availability write or delete is requested for a person ID that does not exist
- **THEN** the endpoint returns `404`

#### Scenario: Deleting an availability entry restores the read-model default
- **WHEN** an availability entry is deleted successfully
- **THEN** subsequent dashboard availability reads no longer return the deleted stored row
- **AND** the read model falls back to its existing default behavior for missing availability on that date

### Requirement: Store-backed persistence SHALL support the new write surface without breaking existing reads

The Store interface, sqlc queries, and Postgres-backed implementation SHALL persist person and availability writes against the existing `people` and `availability` tables. Availability writes SHALL use the existing unique `(person_id, date)` constraint to support idempotent single-date upserts. The existing `GET /api/dashboard/people-availability` contract, status vocabulary, and summary invariants SHALL remain unchanged.

#### Scenario: Existing dashboard reads continue to satisfy the current contract
- **WHEN** the write surface is added
- **THEN** `GET /api/dashboard/people-availability` still returns `range`, `summary`, `people`, and `statuses`
- **AND** the status legend and summary invariants remain unchanged

#### Scenario: Persisted writes survive beyond a single request
- **WHEN** a person or availability record is written successfully
- **THEN** subsequent reads are served from Store-backed persisted data rather than per-request in-memory mutation

### Requirement: The committed frontend API artifacts SHALL reflect the people write surface

When the backend contract gains the people and availability write endpoints, the committed OpenAPI snapshot and generated client artifacts SHALL be refreshed so `/people` can call the write surface type-safely without depending on a live backend during normal verification.

#### Scenario: Snapshot-backed generation includes the write surface
- **WHEN** the frontend API artifacts are refreshed for this change
- **THEN** `frontend/openapi-snapshot.json` includes `/api/people`, `/api/people/{id}`, and `/api/people/{id}/availability/{date}` write paths
- **AND** the committed generated client under `frontend/src/client/` includes typed operations for those endpoints

### Requirement: Backend verification SHALL cover CRUD, validation, and contract preservation

Backend tests SHALL cover successful create/update/delete flows with the new server-assigned ID contract, single-date availability upsert/delete, `400` validation failures, `404` missing people, `409` blocked deletes, OpenAPI path inclusion, and preservation of the existing hello, planning-window, and dashboard-people-availability contracts. The create-person test SHALL NOT send an `id` field in the request body and SHALL verify that the response contains a server-assigned `id`.

#### Scenario: Create endpoint accepts only name and initials

- **WHEN** `POST /api/people` is called with a body containing `name`, `initials`, and no `id`
- **THEN** the endpoint returns `201` with a person record containing an auto-generated `id`
- **AND** a subsequent read confirms the person was persisted

#### Scenario: CRUD contract tests pass

- **WHEN** the backend test suite runs
- **THEN** the people and availability CRUD tests pass alongside the existing read-path tests

#### Scenario: OpenAPI includes the committed write surface

- **WHEN** backend tests inspect `/openapi.json`
- **THEN** the document includes the new people and availability write paths
