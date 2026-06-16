## 1. Wire the shared read path into the existing surfaces

- [x] Add a shared `usePeopleAvailability` composable/adapter that queries `GET /api/dashboard/people-availability`, converts nullable generated fields into the component contract, formats ISO dates for display, and exposes loading, error, and empty states.
- [x] Update `frontend/src/home/HomeView.vue` to render `PeopleAvailability` from the shared composable while preserving the existing dashboard layout and keeping the homepage read-only.
- [x] Update `frontend/src/people/PeopleView.vue` to use the same composable and show explicit loading, backend-unavailable, and empty states instead of route-level demo defaults.

## 2. Add the minimal write surface required for a functional matrix

- [x] Extend `backend/queries/people_availability.sql`, generated sqlc output, and `backend/store.go` with person create/update/delete and availability upsert/delete operations backed by the existing `people` and `availability` tables.
- [x] Register Huma endpoints for `POST /api/people`, `PUT /api/people/{id}`, `DELETE /api/people/{id}`, `PUT /api/people/{id}/availability/{date}`, and `DELETE /api/people/{id}/availability/{date}`.
- [x] Enforce client-supplied stable person IDs, canonical availability statuses, planning-window date validation, `404` for missing people, and `409` for deletes blocked by backlog or schedule assignment references.
- [x] Preserve the existing `GET /api/dashboard/people-availability` contract and default missing-availability behavior while adding the write surface.

## 3. Regenerate client artifacts and wire `/people` mutations

- [x] Refresh `frontend/openapi-snapshot.json` and regenerate committed `frontend/src/client/` artifacts through the existing workflow after backend contract changes.
- [x] Add `/people` route mutation flows for creating a person, updating a single date status, and deleting an unreferenced person using the typed generated client.
- [x] Keep mutation controls on `/people` only; do not add CRUD controls to the homepage panel.

## 4. Expand verification and run repository checks

- [x] Update frontend tests to cover the adapter/composable, component rendering with backend-shaped props, and `/` plus `/people` route rendering with mocked BFF data.
- [x] Add backend tests for CRUD success paths, `400`/`404`/`409` failures, OpenAPI path inclusion, persistence behavior, and preservation of existing hello/planning-window/dashboard availability tests.
- [x] Run `scripts/precommit-run`.
