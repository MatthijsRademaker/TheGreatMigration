## Context

The project has two surfaces for displaying the task backlog:

1. **`TaskManagementPanel.vue`** (`frontend/src/tasks/components/`): A reusable `<Card>`-based panel with a 5-column semantic `<table>` using `TaskRow.vue` sub-components. It accepts a `readOnly` prop and emits `filter` and `add-task` events. Currently heading reads "Task Management", has no count indicator, and renders a priority legend in its footer.

2. **`TasksView.vue`** (`frontend/src/tasks/`): The `/tasks` route view that renders its own `<Card>` with heading "Task backlog", a description, a "+ Add Task" button, and an inline `<ul><li>` task list with per-row Edit/Delete buttons. It owns all CRUD state: form fields, mutation handlers, the `AddOperationModal`, room-select loading/error/retry states, and query invalidation. It does **not** reuse `TaskManagementPanel.vue`.

3. **`HomeView.vue`** (`frontend/src/home/`): The `/` dashboard renders `<TaskManagementPanel read-only />` in a grid beside People Availability. Read-only is enforced by the `readOnly` prop, tested in `app-routes-render.test.ts`.

The divergence between these two surfaces means: (a) design-alignment work on the panel alone would not affect `/tasks`, (b) two separate code paths must be maintained for the same backlog data, and (c) the original archived design decision D1 ("render the panel from TasksView.vue") was never completed.

The backend backlog API (`GET /api/tasks/backlog`) is stable and provides `summary.totalTasks`, task rows, priorities, and statuses via `useTaskBacklog()`. The task vocabulary (`priority`, `status`, `assignedTo`) and data contracts are well-preserved in `types.ts`, `helpers.ts`, and the composable.

The design artifacts referenced by the dossier are section 4 of `designs/components.png` (the task-backlog component anatomy showing heading, columns, row badges) and the backlog panel in `designs/home-page.png` (dashboard grid placement).

## Goals / Non-Goals

### Goals

1. Unify the task backlog under a single authoritative shell (`TaskManagementPanel.vue`), used by both the home dashboard (read-only) and the `/tasks` route (editable).
2. Align the panel heading to "Tasks Backlog" matching the design artifacts.
3. Surface `summary.totalTasks` as a compact count chip in the panel header.
4. Preserve the priority legend footer (matches archived design decision D4).
5. Wire `TasksView.vue` CRUD modal to the panel's emits so that Add Task, Edit, and Delete all work through the unified shell.
6. Keep the home dashboard strictly read-only (no Filter, Add Task, Edit, or Delete controls).
7. Update route-render tests to assert the new heading, count chip, column structure, and route-specific control visibility.
8. Produce a self-contained, text-only implementation brief.

### Non-Goals

- Backend changes, OpenAPI regeneration, generated-client updates, or persistence work.
- New filtering, searching, drag-and-drop, scheduling, or sorting behavior.
- Redesign of People Availability, Daily Schedule, Move Notes, sidebar, or top toolbar.
- Changes to canonical task vocabulary (`priority`, `status`, `assignedTo`) or backlog query semantics.
- Adding a grip/drag handle to `TaskRow.vue` (deferred until reordering is implemented).
- Pixel-perfect ornamental work beyond what repository text and existing design artifacts justify.
- Changing `useTaskBacklog` composable logic, data adaptation, or query-key semantics.
- Replacing the `AddOperationModal` shared shell.

## Decisions

### D1: Canonical heading is "Tasks Backlog"

Change the `CardTitle` in `TaskManagementPanel.vue` from "Task Management" to "Tasks Backlog". Since `TasksView.vue` will reuse this panel, the heading is automatically consistent on both `/` and `/tasks` routes.

**Rationale**: All three refinement agents agree on "Tasks Backlog" as the canonical heading. The current headings ("Task Management" and "Task backlog") both diverge from the design artifacts. This is the only heading that matches the design reference.

**Evidence**: `TaskManagementPanel.vue` line with `CardTitle` renders "Task Management". `TasksView.vue` renders "Task backlog". Dossier affected areas confirm the drift.

### D2: Add total-task count chip from `summary.totalTasks`

Add a compact `Badge` in the panel header, adjacent to the heading, displaying `{{ data.summary.totalTasks }} tasks` (e.g., "12 tasks"). The count chip renders in both read-only and editable variants.

**Rationale**: All three refinement agents agree the count chip should be added. `summary.totalTasks` is available from `useTaskBacklog()` but never rendered. This surfaces useful information at zero backend cost. The chip format "{N} tasks" is simple, text-only, and unambiguous.

**Evidence**: `useTaskBacklog.ts` adapts `summary` from the API response. `TaskManagementPanel.vue` template does not reference `summary` anywhere. Dossier and all three agents flag this as a missed opportunity.

### D3: Preserve the priority legend footer

Keep the existing `<div class="mt-4 flex items-center gap-3 border-t border-border pt-4">` footer with priority legend badges in `TaskManagementPanel.vue`.

**Rationale**: The lead-dev and reviewer both recommend keeping the legend; 2 of 3 agents support it. The archived design decision D4 explicitly included the legend. While the architect argued it's redundant with per-row badges, the legend serves as a quick reference for users unfamiliar with the badge color semantics.

**Evidence**: `TaskManagementPanel.vue` currently renders the footer. Archived `design.md` D4: "Use `priorityHigh`, `priorityMedium`, and `priorityLow` for row badges and the priority legend." Consensus: 2-1 for keeping.

### D4: TasksView.vue reuses TaskManagementPanel directly

`TasksView.vue` drops its standalone `<Card>` wrapper and inline `<ul>` task list. Instead it renders `<TaskManagementPanel>` (non-readOnly) and wires its emits to the existing CRUD handlers. The `AddOperationModal` and all form/mutation/room-select state remain owned by `TasksView.vue`.

To support editable mode, `TaskManagementPanel` adds two emits (`edit-task`, `delete-task`) and `TaskRow` accepts an optional named slot for row actions. When `readOnly` is false, `TaskRow` renders the slot (or default Edit/Delete buttons). `TasksView` provides the action buttons via the slot, wired to `startEdit(task)` and `handleDelete(task.id)`.

The panel continues to emit `add-task` (wired to `startNewTask()`), `filter` (placeholder — no behavior wired), `edit-task(taskRow)`, and `delete-task(taskId)`.

**Rationale**: The lead-dev and reviewer both recommend direct reuse; 2 of 3 agents support it. The original archived design decision D1 explicitly intended for `TasksView` to render the panel. Two separate surfaces that display the same data inevitably drift. Unifying under one shell is the only way to achieve lasting design alignment. The slot-based approach keeps the panel generic (no knowledge of TasksView's form state) while letting TasksView provide the Edit/Delete UI.

**Evidence**: `TasksView.vue` lines 61-138 render a completely separate card and `<ul>` task list. Archived `design.md` D1: "Create a TaskManagementPanel component and a TaskRow sub-component inside the tasks feature area, then render the panel from TasksView.vue." The current codebase never completed this unification.

### D5: Home dashboard remains strictly read-only

The `readOnly` prop on `TaskManagementPanel` controls visibility of: Filter button, Add Task button, and per-row Edit/Delete action buttons. On the home route (`/`), `HomeView.vue` passes `read-only`, so none of these controls render. The heading and count chip render regardless of `readOnly`.

**Rationale**: All three refinement agents agree the home dashboard stays read-only. The existing spec (`task-management-ui/spec.md`) explicitly states: "HomeView SHALL render the panel in read-only mode with no filter, add, edit, or delete controls." The route test already asserts `expect(html).not.toContain('Filter')` and `expect(html).not.toContain('Add Task')`. The design mockup showing Filter/Add Task controls on the dashboard is overridden by the existing product decision.

**Evidence**: `HomeView.vue` renders `<TaskManagementPanel read-only />`. `task-management-ui/spec.md` home read-only requirement. `app-routes-render.test.ts` asserts absence of Filter/Add Task on `/`.

### D6: No grip/drag handle in TaskRow

`TaskRow.vue` does not add a leading grip/drag affordance. The row template continues to start directly with the task title cell.

**Rationale**: All three refinement agents agree. The grip affordance shown in the design mockup is decorative without real reordering behavior. Adding it would create a false affordance. Reordering is explicitly excluded by the non-goals ("No new drag-and-drop or sorting behavior").

**Evidence**: `TaskRow.vue` template starts with `<td>{{ task.title }}</td>`. No grip handle currently rendered. Non-goals explicitly exclude drag-and-drop.

### D7: Route-render tests update to reflect the unified shell

`frontend/tests/app-routes-render.test.ts` updates assertions:
- Home (`/`): Assert "Tasks Backlog", count chip "3 tasks" (matching mock `totalTasks: 3`), column headers ("Task", "Priority", "People Needed", "Room / Area", "Status"), task row content ("Pack kitchen boxes", "Disassemble bed frames"), priority legend, absence of "Filter", "Add Task", "Edit", "Delete", and absence of the old heading "Task Management".
- `/tasks`: Assert "Tasks Backlog", count chip "3 tasks", column headers, task row content ("Pack kitchen boxes", etc.), presence of "Add Task" (toolbar), "Edit" and "Delete" (per-row actions), and absence of the old heading "Task backlog".

**Rationale**: Current test assertions will break when headings change. The tests must be updated in lockstep with the component changes. Adding column-header and count-chip assertions locks in the design contract.

**Evidence**: `app-routes-render.test.ts` currently asserts `expect(html).toContain('Task Management')` on `/` and `expect(html).toContain('Task backlog')` on `/tasks`. These will break.

## Conflict Resolution

Two conflicts emerged between the architect and the lead-dev/reviewer:

### Priority legend footer (Keep vs. Remove)
- **Architect**: Remove — redundant with per-row badges, reduces vertical space.
- **Lead-dev + Reviewer**: Keep — matches archived design decision D4, serves as a reference legend for badge semantics.
- **Resolution**: Keep the footer. 2 of 3 agents recommend keeping it. The archived design decision explicitly included it. It provides a compact visual reference that benefits new users.

### TasksView reuse of TaskManagementPanel (Direct reuse vs. Separate visual mirror)
- **Architect**: Keep TasksView as separate card shell, visually align using same TaskRow and 5-column layout.
- **Lead-dev + Reviewer**: Have TasksView directly reuse TaskManagementPanel, wiring its emits to the CRUD modal.
- **Resolution**: Direct reuse. 2 of 3 agents recommend it. The original archived design decision D1 explicitly intended this. Two separate surfaces that display the same data are a maintenance liability and will inevitably drift again. The slot-based row-action mechanism preserves TasksView's CRUD ownership without coupling the panel to TasksView-specific state.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| TasksView adopting the panel introduces nested Card shells if TasksView keeps its wrapper | Medium — visual nesting looks wrong | TasksView removes its outer `<Card>` wrapper and renders `<TaskManagementPanel>` directly inside its `<section>` |
| Adding per-row Edit/Delete to the panel may break home read-only expectations | Medium — home dashboard must not show Edit/Delete | Gating per-row actions behind `!readOnly` (slot or v-if) ensures home stays clean |
| Route-render tests break from heading change | Low — expected, must update in lockstep | Update all affected `expect(html).toContain(...)` and `expect(html).not.toContain(...)` assertions |
| Count chip layout on narrow viewports | Low — header uses flex layout that wraps naturally | Use a compact `Badge variant="secondary"` and let the flex layout handle wrapping |
| Room-select loading/error states in TasksView must survive the refactor | Medium — these are owned by TasksView, not the panel | TaskManagementPanel does not own modal/room-select state; TasksView keeps all of that. The panel's only connection is via emit wiring |
| No component-level tests for TaskManagementPanel exist | Medium — refactoring without unit coverage increases regression risk | Add basic component tests for heading, count chip, readOnly control visibility, and loading/error/empty states as part of implementation |

## Traceability

- **Task**: `f0e3641c-96b1-423a-8b09-4b890747e791`
- **Dossier**: `2026-06-17T20:09:38.687Z`
- **Accepted decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Non-blocking findings**: `round-1-swarm-reviewer-blocker-1` (placeholder proposal), `round-1-swarm-reviewer-blocker-2` (read-only vs mockup conflict), `round-1-swarm-reviewer-blocker-3` (heading inconsistency), `round-1-swarm-reviewer-blocker-4` (count chip missing)
- **Validated round outputs**: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- **Reference specs**: `openspec/specs/task-management-ui/spec.md`, `openspec/specs/task-backlog-api/spec.md`
- **Reference design**: Archived `design.md` from `openspec/changes/archive/2026-06-15-task-f6503ce2-ff41-47e9-84fc-639a69c67a8a/`
- **Artifact base**: `taskId:f0e3641c-96b1-423a-8b09-4b890747e791`, snapshot `initial`