## 1. Feature-local task data and helpers

- [ ] 1.1 Add a tasks-local type or types file that mirrors the future backend `TaskRow` shape (`id`, `title`, `priority`, `peopleNeeded`, `room`, `status`, `assignedTo`)
- [ ] 1.2 Add dedicated static fixture data under `frontend/src/tasks/` including the design example row (`Painting hall`) and enough rows to exercise High/Medium/Low priorities, people-needed counts, rooms/areas, and empty `assignedTo`
- [ ] 1.3 Add a centralized helper or mapping for row display state so empty `assignedTo` renders the `Unassigned` pill label

## 2. Task management components

- [ ] 2.1 Create `TaskManagementPanel` under `frontend/src/tasks/components/` using existing `Card`, `Button`, and `Badge` primitives
- [ ] 2.2 Create a `TaskRow` sub-component under `frontend/src/tasks/components/` for the five-column row rendering
- [ ] 2.3 Implement the panel header with `Task Management` title, Filter outline button, and `+ Add Task` primary button using verified `@lucide/vue` icons
- [ ] 2.4 Emit placeholder `filter` and `add-task` events from the toolbar controls without wiring backend or mutation behavior
- [ ] 2.5 Render the task rows with Task, Priority, People Needed, Room / Area, and Status columns using an overflow-capable table/grid layout
- [ ] 2.6 Render People Needed as an icon-plus-count inline group and use the existing `secondary` badge variant for the `Unassigned` status pill
- [ ] 2.7 Render a separate priority legend showing High, Medium, and Low with the existing priority badge variants

## 3. Route integration

- [ ] 3.1 Replace the `Task foundation` placeholder in `frontend/src/tasks/TasksView.vue` with the new Task Management panel composition
- [ ] 3.2 Keep the route frontend-only and avoid any new API client, query, backend, or routing changes

## 4. Test coverage

- [ ] 4.1 Update `frontend/tests/app-routes-render.test.ts` so `/tasks` asserts the new task-management content instead of `Task foundation`
- [ ] 4.2 Cover representative strings such as `Task Management`, `Painting hall`, `People Needed`, `Room / Area`, `Unassigned`, `Filter`, and `Add Task`

## 5. Verification

- [ ] 5.1 Verify the selected `@lucide/vue` icon names resolve before finalizing the task components
- [ ] 5.2 Run `scripts/precommit-run` and confirm the frontend changes and updated tests pass
