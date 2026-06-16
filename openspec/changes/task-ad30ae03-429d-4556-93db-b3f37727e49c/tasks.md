## 1. Build shared dialog primitives

- [ ] Add a centered dialog primitive family under `frontend/src/shared/ui/dialog/` that follows the existing shared-ui overlay conventions and accessible Reka Dialog composition.
- [ ] Provide overlay, content, header, title, description, close, and footer building blocks needed by a centered modal instead of the existing side-sheet layout.
- [ ] Style the dialog primitives with semantic tokens, restrained floating-surface treatments, and responsive classes that support full-screen mobile behavior.

## 2. Create the reusable add-operation modal shell

- [ ] Add `frontend/src/shared/components/AddOperationModal.vue` on top of the dialog primitives.
- [ ] Support `v-model:open`, accessible title/description, caller-provided body content, configurable submit/cancel labels, disabled/submitting states, `submit`/`cancel`/`update:open` propagation, and a default footer with a named override slot.
- [ ] Ensure the shell sizes to roughly 60% viewport width and height on `sm` and above, fills the viewport below `sm`, keeps header/footer stable, and scrolls the body region for long forms.

## 3. Migrate `/tasks` create/edit flows to the shared modal

- [ ] Replace `Sheet`, `SheetContent`, `SheetFooter`, `SheetHeader`, `SheetTitle`, `SheetDescription`, and `SheetTrigger` usage in `frontend/src/tasks/TasksView.vue` with `AddOperationModal.vue`.
- [ ] Rename and wire route-local modal state so both add and edit flows explicitly open the modal through `v-model:open`.
- [ ] Update `startNewTask()` to reset task form state and open the modal directly.
- [ ] Keep `startEdit()` prefill behavior, create/update mutation branching, assignment toggling, backlog invalidation, and success-driven reset/close behavior inside `TasksView.vue`.
- [ ] Move the task mutation error banner into the caller-rendered modal body so failures remain visible without clearing user-entered values.
- [ ] Keep delete behavior outside the modal and preserve its current invalidation path.

## 4. Preserve future reuse hooks without expanding planner scope

- [ ] Leave `TaskManagementPanel.vue` and the home dashboard task panel unchanged in this slice.
- [ ] If `frontend/src/calendar/DailySchedule.vue` is touched, add typed emitted add intents for the header and per-day add affordances, including day `date` context where applicable.
- [ ] Do not introduce scheduling persistence, planner mutations, or automatic modal opening from `DailySchedule.vue` in this change.

## 5. Add focused frontend verification

- [ ] Add or update tests for the shared modal event contract covering `submit`, `cancel`, and `update:open` propagation plus default-footer behavior.
- [ ] Add or update `/tasks` tests covering add/edit modal opening, create/update payload handling, backlog invalidation after success, modal reset/close after success, and state preservation on mutation failure.
- [ ] Keep existing SSR route-render coverage for `/tasks` management controls and daily-schedule visible content passing.

## 6. Run repository verification

- [ ] Run `scripts/precommit-run`.