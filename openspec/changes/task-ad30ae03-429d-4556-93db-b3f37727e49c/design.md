## Context

`frontend/src/tasks/TasksView.vue` currently renders task create/edit inside a right-side `Sheet` and uses two different open paths: the Add flow relies on `SheetTrigger`, while the Edit flow sets `sheetOpen` directly. Mutation errors are rendered outside the sheet, and the task form footer uses the existing sheet footer layout. The project already has shared overlay conventions built on Reka Dialog primitives and design-system guidance for calm surfaces, semantic tokens, restrained shadows, and 12px/16px radius hierarchy.

## Goals

- Replace the `/tasks` add/edit sheet with a reusable modal shell while preserving create, edit, cancel, delete, error, and backlog-refresh behavior.
- Keep the modal domain-agnostic so callers provide form content and handlers while the shell owns accessible dialog structure, sizing, and default actions.
- Ensure submit, cancel, and open-state changes propagate cleanly to caller-owned handlers.
- Meet the responsive sizing requirement: approximately 60% viewport width/height on desktop/tablet and full-screen on mobile.
- Add focused verification for modal event propagation and `/tasks` regression coverage.

## Non-Goals

- Implementing daily planner persistence or a scheduling mutation flow.
- Redesigning task fields, task row presentation, or the backlog API contract beyond moving the form into a reusable modal shell.
- Refactoring `TaskManagementPanel.vue`, `HomeView.vue`, People, or Rooms to adopt the modal in this change.
- Routing delete confirmation through the add-operation modal.
- Changing backend endpoints or generated clients unless existing frontend tests force compatibility fixes.

## Decisions

### 1. Use a two-layer shared modal architecture

Create a low-level dialog primitive family under `frontend/src/shared/ui/dialog/` that mirrors the existing `sheet/` structure but uses centered modal positioning. Compose those primitives into `frontend/src/shared/components/AddOperationModal.vue`, which provides the reusable add-operation contract without embedding domain logic.

### 2. Unify add and edit open-state handling under explicit `v-model:open`

Remove `SheetTrigger` from the `/tasks` Add flow. `TasksView.vue` owns a single modal open state, and both `startNewTask()` and `startEdit()` explicitly open the modal. This removes the existing trigger race and makes add/edit lifecycle behavior consistent.

### 3. Keep task-domain state, mutations, and errors in `TasksView.vue`

`TasksView.vue` continues to own form refs, assignment toggling, create/update/delete mutations, backlog query invalidation, and reset behavior. The shared modal stays domain-agnostic. Task-specific mutation errors are rendered by the caller inside the modal body slot so errors stay visible without teaching the shell about task-domain error types.

### 4. Provide a default footer with override flexibility

`AddOperationModal.vue` provides default submit/cancel actions controlled by props for labels, disabled state, and submitting state, and also exposes a named footer override so future adopters can replace the footer when they need different actions. Header and footer remain stable while the form body scrolls.

### 5. Size the modal for reuse and responsive behavior

On `sm` and above, modal content is centered and sized to approximately `60vw` by `60vh` with sensible max/min constraints and a scrollable content body. Below `sm`, the modal fills the viewport. Styling should use existing semantic tokens and calm floating-surface treatments from the design system rather than ad hoc colors.

### 6. Keep route architecture stable in this slice

`TasksView.vue` remains the management route owner, and `TaskManagementPanel.vue` remains the existing read-only/home-oriented panel contract. This change only replaces the `/tasks` form shell; it does not unify the two task-panel implementations.

### 7. Future-proof `DailySchedule.vue` through emitted intents only

If `DailySchedule.vue` is updated as part of this work, its header `Add task` action emits a typed add intent and its per-day `+ Add task` affordances emit typed add intents with day `date` context. No scheduling persistence, modal opening logic, or planner mutation work is introduced here.

### 8. Verify the shared contract and the route regression path

Add focused tests for modal `submit`, `cancel`, and `update:open` behavior, plus `/tasks` regression coverage proving create/update payload handling, backlog invalidation, modal reset/close behavior after success, and error preservation on failure.

## Risks

- If `startNewTask()` is migrated without explicitly opening the modal, the Add Task button will stop working.
- A 60vw modal can become too large on wide screens unless the implementation applies sensible caps while still meeting the “about 60%” requirement.
- The modal must keep header/footer stable and scroll only the body region so the longer task form remains usable on shorter desktop heights.
- Adding typed `DailySchedule` emits changes that component’s public contract, so tests must keep the existing visible button labels and SSR output stable.

## Conflict Resolution

- **Component layering**: Resolved to a two-layer design (`shared/ui/dialog` primitives plus `shared/components/AddOperationModal.vue`) instead of a single monolithic modal wrapper, matching the accepted architecture decision and existing `sheet/` pattern.
- **Open-state race**: Resolved by removing `SheetTrigger` from `/tasks` and requiring explicit modal opening in `startNewTask()` and `startEdit()`.
- **Error rendering**: Resolved by keeping mutation error ownership in `TasksView.vue` and rendering errors inside the caller-provided modal body content rather than adding task-aware shell logic.
- **Task panel fragmentation**: Resolved by keeping `TasksView.vue` and `TaskManagementPanel.vue` separate in this slice; only the `/tasks` form shell changes.
- **Footer contract**: Resolved by requiring shell-provided default submit/cancel actions plus a named footer override for callers that need different controls.

## Traceability

- Task: `ad30ae03-429d-4556-93db-b3f37727e49c`
- Dossier: `2026-06-16T13:13:44.229Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial`