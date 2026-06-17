## Why

The task backlog UI is split across two different surfaces that have drifted from each other and from the approved design artifacts. `TaskManagementPanel.vue` (the reusable dashboard panel) renders heading "Task Management" while `TasksView.vue` (the `/tasks` route) renders heading "Task backlog" and its own inline `<ul>` task list with Edit/Delete buttons — completely separate from the panel's `<table>` layout. The design artifacts (section 4 of `designs/components.png` and the backlog panel in `designs/home-page.png`) show a unified backlog surface titled "Tasks Backlog" with a count indicator.

The backend already provides `summary.totalTasks` via `useTaskBacklog()`, but no UI surface renders it. The priority legend footer duplicates information already present in each row's priority badge. These disconnects make the backlog UI inconsistent and harder to maintain.

This change consolidates both surfaces into one design-aligned shell, making the task description fully implementable from text alone without inspecting PNG files.

## What Changes

1. **Heading alignment**: Change the task panel heading from "Task Management" to "Tasks Backlog" in `TaskManagementPanel.vue`. Adopt the same heading in `TasksView.vue` (which now reuses the panel directly).

2. **Total-task count chip**: Surface `summary.totalTasks` as a compact count chip in the `TaskManagementPanel` header area, positioned adjacent to the heading (e.g., "12 tasks"). This uses data already available from `useTaskBacklog()` without changing the backend contract.

3. **Priority legend footer preserved**: Keep the existing footer priority legend in `TaskManagementPanel`. It remains a useful reference that matches the original design decision D4 from the archived change.

4. **Unify TasksView under TaskManagementPanel**: `TasksView.vue` drops its standalone `<Card>` wrapper and inline `<ul>` task list, and instead renders `<TaskManagementPanel>` (non-readOnly) directly. The panel's `add-task` event wires to `startNewTask()`, opening the existing CRUD modal. In editable mode, the panel exposes per-row Edit/Delete action buttons that emit `edit-task` and `delete-task` events, wired to `startEdit(task)` and `handleDelete(id)`. All CRUD form state, room-select loading/error states, mutation handlers, and query invalidation remain owned by `TasksView.vue`.

5. **Home dashboard stays read-only**: `HomeView.vue` continues to render `<TaskManagementPanel read-only />`. No Filter or Add Task controls appear on the home route. The heading and count chip changes are automatic since they live inside the panel component.

6. **No grip/drag handle**: `TaskRow.vue` does not add a leading grip/drag affordance. The drag handle shown in the design mockup is deferred until real reordering behavior is implemented (explicit non-goal).

7. **Route-render test updates**: `frontend/tests/app-routes-render.test.ts` updates assertions: home route asserts "Tasks Backlog", the count chip text, and absence of Filter/Add Task/Edit/Delete; `/tasks` route asserts "Tasks Backlog", the count chip, presence of Add Task, per-row Edit/Delete, and column headers matching the panel's 5-column table. Home route assertions for "Task Management" are replaced.

8. **No backend or data-contract changes**: All canonical task vocabulary (`priority`, `status`, `assignedTo`), backlog query semantics (`GET /api/tasks/backlog`), types in `frontend/src/tasks/types.ts`, helpers in `frontend/src/tasks/helpers.ts`, and the `useTaskBacklog` composable remain unchanged.

## Impact

- **Affected specs**: `task-management-ui` (existing spec extended by this change's new design-alignment spec)
- **Affected files**:
  - `frontend/src/tasks/components/TaskManagementPanel.vue` — heading, count chip, editable-row action buttons, edit/delete emits
  - `frontend/src/tasks/components/TaskRow.vue` — optional per-row action slot for Edit/Delete
  - `frontend/src/tasks/TasksView.vue` — replace custom card+list with `<TaskManagementPanel>`, wire emits
  - `frontend/src/home/HomeView.vue` — no code changes needed (panel updates propagate automatically)
  - `frontend/tests/app-routes-render.test.ts` — update heading/control/row assertions
- **No impact** on: backend code, OpenAPI contract, generated client, database schema, other dashboard panels, calendar, people, rooms, or settings views.
