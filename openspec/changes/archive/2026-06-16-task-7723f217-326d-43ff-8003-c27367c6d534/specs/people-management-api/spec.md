## MODIFIED Requirements

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
