## 1. Backend Implementation

- [ ] 1.1 Create `backend/tasks.go` with request/response types (`TaskBacklogInput`, `TaskBacklogOutput`, `TaskBacklogBody`, `TaskSummary`, `TaskRow`, `PriorityLegend`, `TaskStatusLegend`), in-memory seed data (at least 10 tasks exercising all priority/status values and 4+ rooms), and `registerTasksBacklog(api huma.API)` handler following the `dashboard.go` pattern.
- [ ] 1.2 Register `registerTasksBacklog(api)` in `backend/main.go` after `registerPlanningWindow(api)`.
- [ ] 1.3 Update `newTestAPI()` in `backend/main_test.go` to call `registerTasksBacklog(api)`.

## 2. Backend Testing

- [ ] 2.1 Add `TestTaskBacklog` in `backend/main_test.go` verifying: 200 JSON response, all top-level fields present (`summary`, `tasks`, `priorities`, `statuses`), summary count consistency (highPriorityTasks matches priority=high count, unassignedTasks matches empty assignedTo count, understaffedTasks matches 0 < len(assignedTo) < peopleNeeded count), canonical priority/status values, OpenAPI publication of new path.
- [ ] 2.2 Verify existing `TestHelloEndpoint`, `TestDashboardPeopleAvailability`, and `TestPlanningWindowEndpoint` still pass (non-regression).
- [ ] 2.3 Run `scripts/precommit-run` and confirm all checks pass.

## 3. Documentation

- [ ] 3.1 Ensure the OpenSpec spec (`task-backlog-api/spec.md`) is published with this change.
- [ ] 3.2 Confirm the endpoint appears in `/openapi.json` when the backend is running.

## 4. Deferred (Explicitly Out of Scope)

- [ ] 4.1 Regenerate `frontend/openapi-snapshot.json` to include the new endpoint (follow-up task).
- [ ] 4.2 Wire `HomeView.vue` summary cards to consume `GET /api/tasks/backlog` (follow-up task).
- [ ] 4.3 Wire `TasksView.vue` backlog table to consume `GET /api/tasks/backlog` (follow-up task).