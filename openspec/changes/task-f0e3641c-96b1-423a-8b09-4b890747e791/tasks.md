## 1. Refine TaskManagementPanel

- [ ] 1.1 Change `CardTitle` from "Task Management" to "Tasks Backlog".
- [ ] 1.2 Add a total-task count chip in the header area using `data.summary.totalTasks` (e.g., `<Badge variant="secondary">{{ data.summary.totalTasks }} tasks</Badge>`), positioned adjacent to the heading.
- [ ] 1.3 Add `edit-task` and `delete-task` emits to the component's `defineEmits` alongside existing `filter` and `add-task`.
- [ ] 1.4 In non-readOnly mode, render per-row Edit/Delete action buttons inside `TaskRow` (via a named slot or conditional template).
- [ ] 1.5 Wire the per-row Edit button to emit `edit-task` with the task row object.
- [ ] 1.6 Wire the per-row Delete button to emit `delete-task` with the task ID.
- [ ] 1.7 Gate the Filter and Add Task toolbar buttons behind `v-if="!readOnly"` (this is already the case — verify and keep intact).
- [ ] 1.8 Keep the priority legend footer unchanged.
- [ ] 1.9 Keep the loading, error, and empty state markup unchanged.

## 2. Update TaskRow for editable actions

- [ ] 2.1 Add an optional named slot (e.g., `#actions`) to `TaskRow.vue` for per-row action buttons.
- [ ] 2.2 Render the slot content only when provided (no action cells in read-only mode).
- [ ] 2.3 Do NOT add a grip/drag handle.
- [ ] 2.4 Preserve existing column anatomy: title, priority badge, people-needed icon+count, room, status badge.

## 3. Unify TasksView under TaskManagementPanel

- [ ] 3.1 Remove the standalone `<Card>` wrapper and inline `<ul>` task list from `TasksView.vue` template.
- [ ] 3.2 Render `<TaskManagementPanel>` inside the existing `<section>` wrapper.
- [ ] 3.3 Wire `@add-task` to `startNewTask()`.
- [ ] 3.4 Wire `@edit-task` to `startEdit(task)`.
- [ ] 3.5 Wire `@delete-task` to `handleDelete(id)`.
- [ ] 3.6 Provide Edit/Delete buttons via the `TaskRow` actions slot.
- [ ] 3.7 Keep the `AddOperationModal` and all CRUD form state, mutation handlers, room-select loading/error/retry states, and query invalidation intact.
- [ ] 3.8 Remove the unused `CardDescription` text and inline `+ Add Task` button (now handled by the panel).
- [ ] 3.9 Verify modal open/close behavior works correctly with emit-wired triggers.

## 4. Verify HomeView read-only contract

- [ ] 4.1 Confirm `HomeView.vue` continues to render `<TaskManagementPanel read-only />`.
- [ ] 4.2 Confirm that heading change and count chip render on the home dashboard.
- [ ] 4.3 Confirm no Filter, Add Task, Edit, or Delete controls appear on the home route.

## 5. Update route-render tests

- [ ] 5.1 Update home route (`/`) assertions: replace `expect(html).toContain('Task Management')` with `expect(html).toContain('Tasks Backlog')`.
- [ ] 5.2 Add home route assertion for count chip text (`expect(html).toContain('3 tasks')`).
- [ ] 5.3 Add home route assertions for column headers: "Task", "Priority", "People Needed", "Room / Area", "Status".
- [ ] 5.4 Preserve home route assertions for `expect(html).not.toContain('Filter')`, `expect(html).not.toContain('Add Task')`, `expect(html).not.toContain('>Edit<')`, `expect(html).not.toContain('>Delete<')`.
- [ ] 5.5 Update `/tasks` route assertions: replace `expect(html).toContain('Task backlog')` with `expect(html).toContain('Tasks Backlog')`.
- [ ] 5.6 Add `/tasks` route assertion for count chip text.
- [ ] 5.7 Add `/tasks` route assertions for column headers.
- [ ] 5.8 Preserve `/tasks` route assertions for task row content ("Pack kitchen boxes", "Disassemble bed frames", "Move living room furniture") and management controls ("Edit", "Delete").
- [ ] 5.9 Preserve `/tasks` route assertion for `expect(html).toContain('Add Task')`.
- [ ] 5.10 Run `scripts/precommit-run` and verify all tests pass.

## 6. Add TaskManagementPanel component tests

- [ ] 6.1 Create `frontend/src/tasks/components/__tests__/TaskManagementPanel.spec.ts`.
- [ ] 6.2 Test heading renders "Tasks Backlog".
- [ ] 6.3 Test count chip renders with mock `summary.totalTasks` value.
- [ ] 6.4 Test readOnly mode: Filter and Add Task toolbar buttons are absent.
- [ ] 6.5 Test non-readOnly mode: Filter and Add Task toolbar buttons are present.
- [ ] 6.6 Test per-row Edit/Delete buttons absent in readOnly mode.
- [ ] 6.7 Test per-row Edit/Delete buttons present in non-readOnly mode.
- [ ] 6.8 Test loading, error, and empty states render correct text.
- [ ] 6.9 Test column headers render: Task, Priority, People Needed, Room / Area, Status.
- [ ] 6.10 Test priority legend footer renders with all three priority badges.