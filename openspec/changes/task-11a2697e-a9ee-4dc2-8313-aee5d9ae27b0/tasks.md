## 1. Add the backend task write surface

- [ ] Extend the task Store contract with `CreateTask`, `UpdateTask`, and `DeleteTask` methods alongside the existing backlog read method.
- [ ] Add sqlc task write queries and Postgres-backed Store implementations that create/update/delete `backlog_tasks` and replace `backlog_task_assignments` transactionally.
- [ ] Register `POST /api/tasks`, `PUT /api/tasks/{id}`, and `DELETE /api/tasks/{id}` through Huma, including validation for canonical priority/status values, non-empty title, `peopleNeeded >= 1`, required room, and existing assigned person IDs.
- [ ] Preserve `GET /api/tasks/backlog` as the post-write read contract, including stable task row fields, legends, and derived summary invariants.

## 2. Wire the task UI to the BFF and add route-owned management flows

- [ ] Introduce `frontend/src/tasks/composables/useTaskBacklog.ts` to wrap the generated backlog query, adapt the response into panel props, and expose loading, backend-unavailable, and empty states.
- [ ] Refactor `frontend/src/tasks/components/TaskManagementPanel.vue` to remove fixture-backed rendering, accept a documented `readOnly` variant, and keep `Unassigned` as a derived display label from empty `assignedTo` arrays.
- [ ] Update `frontend/src/tasks/TasksView.vue` to own Sheet-based create/edit/delete flows and assignment updates using the typed generated client plus backlog-query invalidation/refetch after successful mutations.
- [ ] Update `frontend/src/home/HomeView.vue` to render the same BFF-backed task surface in read-only mode without filter/add/edit/delete controls.
- [ ] Tighten task-row typing to the canonical priority/status vocabulary from the generated client or matching union types.

## 3. Refresh generated artifacts and expand verification

- [ ] Refresh `frontend/openapi-snapshot.json` and regenerate committed `frontend/src/client/` artifacts after the backend task write endpoints are added.
- [ ] Add backend tests for CRUD success paths, validation failures, missing-task handling, assignment persistence/replacement, delete behavior, OpenAPI path inclusion, and preservation of the existing backlog read invariants.
- [ ] Add frontend tests for the task backlog composable/adapter, loading/error/empty states, create/edit/delete mutation flows with query invalidation, and the `readOnly` panel behavior.
- [ ] Update `frontend/tests/app-routes-render.test.ts` to mock realistic task backlog payloads for `/api/tasks/backlog` and assert BFF-backed rendering on `/tasks` and `/` instead of relying on fixture leakage.
- [ ] Run `scripts/precommit-run`.
