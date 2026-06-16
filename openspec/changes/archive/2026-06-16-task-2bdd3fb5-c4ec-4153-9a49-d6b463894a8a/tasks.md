## 1. Create ASCII Homepage Mockup Artifact

- [x] 1.1 Create `openspec/changes/task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a/home-page-ascii.md` with a preamble stating it is a non-canonical structural reference, not a visual spec.
- [x] 1.2 Add a legend section mapping each dashboard region to its corresponding repo source:
  - Sidebar → `frontend/src/shared/layout/app-sidebar/AppSidebar.vue` (implemented, `collapsible="icon"`)
  - Top toolbar → `frontend/src/shared/layout/app-shell/AppShell.vue` (GAP: current is page-title + Planning mode badge; design target has date-range/Today/arrows/alerts/profile)
  - KPI cards → `GET /api/tasks/backlog` summary (`openspec/specs/task-backlog-api/spec.md`) and `GET /api/dashboard/people-availability` summary (`openspec/specs/dashboard-people-availability/spec.md`)
  - Tasks Backlog → `frontend/src/tasks/components/TaskManagementPanel.vue` (reusable component exists)
  - People Availability → `frontend/src/people/PeopleAvailability.vue` (reusable component exists)
  - Daily Schedule → `openspec/specs/dashboard-daily-schedule/spec.md` (API contract only; no Vue component exists)
  - Move Notes → Static design-only panel; no spec or component exists
- [x] 1.3 Add a desktop-width monospaced ASCII wireframe (expanded sidebar) with:
  - A persistent left rail reflecting `AppSidebar.vue` structure (brand header, Plan group: Dashboard/Tasks/Schedule/People, Organization group: Rooms/Settings, footer with project card + Add note/Help & Support)
  - A top toolbar annotated as aspirational (date-range, Today button, arrows, alerts, profile chip) with gap note referencing current `AppShell.vue`
  - Row 1: Four equal-width KPI cards labeled Total Tasks, Available People, High Priority Tasks, Move Days
  - Row 2: Tasks Backlog (wider, ~60%) beside People Availability matrix (~40%)
  - Row 3: Daily Schedule multi-column board (~60%) beside Move Notes panel (~40%) with aspirational annotations
- [x] 1.4 Add a collapsed-sidebar ASCII variant showing how the content area reflows when `AppSidebar` is collapsed to icon-only.
- [x] 1.5 Add a gap-analysis section cataloging deviations between current `HomeView.vue` and the target composition:
  - Hello world card → replaced by Total Tasks KPI
  - Available today (6) / Under-staffed (3) → replaced by Available People and High Priority Tasks KPIs
  - Move days card → retained as Move Days KPI (already wired to `usePlanningWindow()`)
  - Today's plan static card → replaced by Tasks Backlog + People Availability row
  - Move notes card → replaced by Daily Schedule + Move Notes row
  - AppShell header toolbar gap (simple title vs. design toolbar)

## 2. Validation

- [x] 2.1 Verify the artifact does not modify any file under `openspec/specs/`, `docs/`, `frontend/src/`, or `backend/`.
- [x] 2.2 Confirm all section labels use existing repo/spec vocabulary (no invented terminology).
- [x] 2.3 Confirm the legend includes exact file paths and spec references (not directory names).
- [x] 2.4 Confirm the four KPI cards are enumerated by name and mapped to their API contracts.
- [x] 2.5 Confirm Daily Schedule and Move Notes are annotated as aspirational/not-yet-built.
