## Why

The current `/tasks` create/edit UI is embedded directly in `TasksView.vue` as a right-side `Sheet`, so the form shell, open-state behavior, and mutation wiring are tightly coupled to one route. The requested change is to replace that side panel with a reusable modal that can serve task add/edit flows now and future add flows such as the daily planner later, while preserving the existing task CRUD behavior and handler wiring.

## What Changes

- Add a shared centered dialog primitive family under `frontend/src/shared/ui/dialog/` and a higher-level `frontend/src/shared/components/AddOperationModal.vue` shell for reusable add/edit flows.
- Define the modal contract around `v-model:open`, accessible title/description, caller-provided body content, configurable submit/cancel labels, disabled/submitting states, default footer actions with override flexibility, submit/cancel/open event propagation, ~60% viewport sizing on desktop/tablet, internal body scrolling, and full-screen behavior on mobile.
- Migrate `frontend/src/tasks/TasksView.vue` from `Sheet`/`SheetTrigger` to the shared add-operation modal. Keep all task-domain form state, create/update/delete mutations, query invalidation, assignment toggling, and error handling in `TasksView.vue`.
- Unify add and edit open-state handling so both flows explicitly control modal visibility through `v-model:open`; `startNewTask()` must open the modal directly instead of relying on `SheetTrigger` timing.
- Keep mutation errors visible inside the modal body without moving task-specific error logic into the shared shell.
- Leave delete behavior outside the add-operation modal and preserve the existing post-success backlog refresh behavior.
- Keep `TaskManagementPanel.vue` and the read-only home-dashboard task panel unchanged in this slice.
- If `frontend/src/calendar/DailySchedule.vue` is touched, add typed `add-task` emits for the header and per-day add affordances without introducing scheduling persistence.
- Add focused frontend tests for the modal event contract and the `/tasks` add/edit regression path, while keeping existing route-render and daily-schedule coverage passing.
- Publish a spec delta updating task-management UI behavior from `Sheet`-based create/edit interactions to the reusable modal contract.

## Impact

- `/tasks` create/edit interactions move from a right-side sheet to a centered reusable modal.
- The shared frontend layer gains reusable dialog/modal building blocks for future add-operation flows.
- Task create, edit, cancel, delete, query invalidation, and error behavior remain on the `/tasks` route.
- Desktop users get an approximately 60vw × 60vh modal with a scrollable body; mobile users get a full-screen modal.
- No backend endpoint, generated client contract, delete flow, or daily planner persistence behavior changes in this slice.