# Home Page ASCII Blueprint

> **Non-canonical structural reference** — not a visual spec.
>
> This file lives under `openspec/changes/task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a/`. It does NOT modify any file under `openspec/specs/`, `docs/`, `frontend/src/`, or `backend/`. Its purpose is to make the target homepage dashboard composition legible to a text-only LLM without requiring PNG inspection. Proportions are inferred from the 12-column grid defined in `docs/design-system-v2.md`; exact pixel alignment with `designs/home-page.png` is not guaranteed.

---

## Legend

Each dashboard region is mapped to its corresponding repo source. File paths and spec references are exact — not directory names.

| Dashboard Region | Source | Status | API Contract |
|---|---|---|---|
| **Sidebar** | `frontend/src/shared/layout/app-sidebar/AppSidebar.vue` | Implemented; `collapsible="icon"` variant `inset` | N/A (client-side navigation) |
| **Top Toolbar** | `frontend/src/shared/layout/app-shell/AppShell.vue` | **GAP**: current renders page-title header + Planning mode badge; design target toolbar (date-range, Today button, arrows, alerts, profile chip) does not exist | N/A |
| **KPI: Total Tasks** | N/A (to be composed in `HomeView.vue`) | **Not yet implemented** | `GET /api/tasks/backlog` → `summary.totalTasks` (`openspec/specs/task-backlog-api/spec.md`) |
| **KPI: Available People** | N/A (to be composed in `HomeView.vue`) | **Not yet implemented** | `GET /api/dashboard/people-availability` → `summary.availableToday` / `summary.totalPeople` (`openspec/specs/dashboard-people-availability/spec.md`) |
| **KPI: High Priority Tasks** | N/A (to be composed in `HomeView.vue`) | **Not yet implemented** | `GET /api/tasks/backlog` → `summary.highPriorityTasks` (`openspec/specs/task-backlog-api/spec.md`) |
| **KPI: Move Days** | Partially wired in `HomeView.vue` via `usePlanningWindow()` | **Existing**: currently rendered as a summary card with `planWindowDayCount` | `GET /api/planning-window` → `planWindowDayCount` (`openspec/specs/planning-window/spec.md`) |
| **Tasks Backlog** | `frontend/src/tasks/components/TaskManagementPanel.vue` | Implemented (reusable component with `TaskRow`, table columns: Task, Priority, People Needed, Room/Area, Status) | `GET /api/tasks/backlog` (`openspec/specs/task-backlog-api/spec.md`) |
| **People Availability** | `frontend/src/people/PeopleAvailability.vue` | Implemented (reusable component with availability matrix: Person × Day columns, status badges: Available/Busy/Partial/Off) | `GET /api/dashboard/people-availability` (`openspec/specs/dashboard-people-availability/spec.md`) |
| **Daily Schedule** | None | **Aspirational** — API contract exists but no Vue component | `GET /api/dashboard/daily-schedule` (`openspec/specs/dashboard-daily-schedule/spec.md`) |
| **Move Notes** | None | **Aspirational** — static design-only panel; no spec or component exists | N/A |

---

## Desktop Wireframe — Expanded Sidebar

```
┌──────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────────────────┐
│                                      │                                                                                          │
│  ┌────────────────────────────────┐  │  ┌────────────────────────────────────────────────────────────────────────────────────┐  │
│  │  📋 The Great Migration        │  │  │  [☰]  │  <  Jul 5–8, 2026  >  [Today]  │  ← →  │  🔔  │  👤 Alex       ▲  ← GAP  │  │
│  │     House move planner         │  │  │        ──────────────────────────────────────────────────────────────────────────  │  │
│  └────────────────────────────────┘  │  │  ⚠ Aspirational toolbar. AppShell.vue currently renders only a page-title header    │  │
│                                      │  │    with Planning mode badge — not the date-range / arrows / alerts / profile chip.  │  │
│  ─── Plan ────────────────────────   │  └────────────────────────────────────────────────────────────────────────────────────┘  │
│                                      │                                                                                          │
│  ● Dashboard              (active)  │  ┌──────────────────────────┬──────────────────────────┬──────────────────────────┬───────┐
│  ● Tasks                            │  │   Total Tasks    42      │  Available People  8/12  │  High Priority     7     │ Move  │
│  ● Schedule                         │  │   ──────────────────     │  ───────────────────────  │  ──────────────────      │ Days  │
│  ● People                           │  │   All tasks in backlog   │  Available today          │  Needs attention         │  40   │
│                                      │  │   GET /api/tasks/        │  GET /api/dashboard/      │  GET /api/tasks/         │       │
│  ─── Organization ───────────────   │  │   backlog                 │  people-availability      │  backlog                 │       │
│                                      │  └──────────────────────────┴──────────────────────────┴──────────────────────────┴───────┘
│  ● Rooms / Areas                    │
│  ● Settings                         │  ┌───────────────────────────────────────────────────────┬──────────────────────────────┐
│                                      │  │  Tasks Backlog                    [Filter] [+ Add]   │  People availability         │
│                                      │  │  ──────────────────────────────────────────────────  │  ────────────────────────────│
│  ─────────────────────────────────   │  │  Task          Priority  Needed  Room        Status  │  5 of 12 available today     │
│                                      │  │  Pack kitchen   High      3      Kitchen     Ready   │                              │
│  ┌────────────────────────────────┐  │  │  Disassemble    Medium    2      Bedroom     Backlog │  Person   Mon  Tue  Wed  Thu │
│  │  📋 The Great Migration        │  │  │  bed frame                                        │  Alex     Avail Avail Busy Avail│
│  │     House move planner         │  │  │  Label boxes    Low       1      Garage      Assigned│  Morgan   Busy  Part  Avail Off │
│  └────────────────────────────────┘  │  │  ...                                             │  Sam      Off   Avail Part  Busy │
│                                      │  │                                                   │  Riley    Part  Busy Avail Avail│
│  + Add note                          │  │  Priority: [High] [Medium] [Low]                  │                              │
│  ? Help & Support                    │  │                                                   │  [Available] [Busy]           │
│                                      │  │  ← uses TaskManagementPanel.vue (reusable)        │  [Partial] [Off]              │
│                                      │  │     GET /api/tasks/backlog                        │                              │
│                                      │  │                                                   │  ← uses PeopleAvailability.vue│
│                                      │  │                                                   │     GET /api/dashboard/       │
│                                      │  │                                                   │     people-availability       │
│                                      │  └───────────────────────────────────────────────────┴──────────────────────────────┘
│                                      │
│                                      │  ┌───────────────────────────────────────────────────────┬──────────────────────────────┐
│                                      │  │  Daily Schedule                                       │  Move Notes                  │
│                                      │  │  ──────────────────────────────────────────────────── │  ────────────────────────────│
│                                      │  │  ⚠ Aspirational — no Vue component exists.           │  ⚠ Aspirational — static     │
│                                      │  │    API contract: GET /api/dashboard/daily-schedule    │    design-only panel. No     │
│                                      │  │    spec: openspec/specs/dashboard-daily-schedule/     │    spec or component exists. │
│                                      │  │                                                       │                              │
│                                      │  │  ┌──────────┬──────────┬──────────┬──────────┐        │  ┌────────────────────────┐  │
│                                      │  │  │ Sun 7/5  │ Mon 7/6  │ Tue 7/7  │ Wed 7/8  │        │  │ Check building access  │  │
│                                      │  │  │ 8 avail  │ 6 avail  │ 7 avail  │ 5 avail  │        │  │ times before assigning │  │
│                                      │  │  ├──────────┼──────────┼──────────┼──────────┤        │  │ early morning jobs.    │  │
│                                      │  │  │ Pack     │ Move     │ Clean    │ Label    │        │  │                        │  │
│                                      │  │  │ kitchen  │ furniture│ garage   │ boxes    │        │  │ Keep tea, chargers,    │  │
│                                      │  │  │ High 2/3 │ Med  2/2 │ Low  1/1 │ Low  1/2 │        │  │ tape, markers, and bin │  │
│                                      │  │  │ ⚠ under  │ ✓ full   │ ✓ full   │ ⚠ under  │        │  │ bags in the first-day  │  │
│                                      │  │  └──────────┴──────────┴──────────┴──────────┘        │  │ essentials box.        │  │
│                                      │  │                                                       │  └────────────────────────┘  │
│                                      │  └───────────────────────────────────────────────────────┴──────────────────────────────┘
│                                      │
└──────────────────────────────────────┴──────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Desktop Wireframe — Collapsed Sidebar (Icon-Only)

When `AppSidebar` is toggled to `collapsible="icon"`, the sidebar shrinks to an icon rail and the content area expands to fill the reclaimed horizontal space.

```
┌─────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│         │                                                                                                                       │
│  ┌────┐  │  ┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐  │
│  │ 📋 │  │  │  [☰]  │  <  Jul 5–8, 2026  >  [Today]  │  ← →  │  🔔  │  👤          ▲  ← GAP (aspirational toolbar)            │  │
│  └────┘  │  └─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘  │
│         │                                                                                                                       │
│  ─────   │  ┌───────────────────────────┬───────────────────────────┬───────────────────────────┬──────────────────────────────┐  │
│         │  │   Total Tasks     42      │  Available People  8/12   │  High Priority      7     │       Move Days     40      │  │
│  🏠 D   │  │   ───────────────────     │  ──────────────────────── │  ───────────────────      │                              │  │
│  📋 T   │  │   All tasks in backlog    │  Available today           │  Needs attention          │  Scheduled working days      │  │
│  📅 S   │  └───────────────────────────┴───────────────────────────┴───────────────────────────┴──────────────────────────────┘  │
│  👥 P   │                                                                                                                       │
│         │  ┌────────────────────────────────────────────────────────────────┬──────────────────────────────────────────────────┐  │
│  ─────   │  │  Tasks Backlog                           [Filter]  [+ Add]    │  People availability                             │  │
│         │  │  ──────────────────────────────────────────────────────────── │  ──────────────────────────────────────────────── │  │
│  🏢 R   │  │  Task           Priority  Needed  Room         Status         │  5 of 12 available today                         │  │
│  ⚙️ S   │  │  Pack kitchen    High      3      Kitchen      Ready          │                                                  │  │
│         │  │  Disassemble     Medium    2      Bedroom      Backlog        │  Person    Mon    Tue    Wed    Thu               │  │
│         │  │  bed frame                                                   │  Alex      Avail  Avail  Busy   Avail             │  │
│         │  │  Label boxes     Low       1      Garage       Assigned       │  Morgan    Busy   Part   Avail  Off               │  │
│         │  │  ...                                                          │  Sam       Off    Avail  Part   Busy              │  │
│         │  │                                                                │  Riley     Part   Busy   Avail  Avail             │  │
│  ─────   │  │  Priority: [High] [Medium] [Low]                             │                                                  │  │
│         │  │                                                                │  [Available] [Busy] [Partial] [Off]               │  │
│  ┌────┐  │  │  ← TaskManagementPanel.vue (reusable)                        │  ← PeopleAvailability.vue (reusable)              │  │
│  │ 📋 │  │  └────────────────────────────────────────────────────────────────┴──────────────────────────────────────────────────┘  │
│  └────┘  │                                                                                                                       │
│         │  ┌────────────────────────────────────────────────────────────────┬──────────────────────────────────────────────────┐  │
│  + 📝    │  │  Daily Schedule                                                │  Move Notes                                      │  │
│  ? ❓    │  │  ──────────────────────────────────────────────────────────── │  ──────────────────────────────────────────────── │  │
│         │  │  ⚠ Aspirational — no Vue component exists.                    │  ⚠ Aspirational — static design-only panel.      │  │
│         │  │  API: GET /api/dashboard/daily-schedule                        │  No spec or component exists.                    │  │
│         │  │                                                                │                                                  │  │
│         │  │  ┌───────────┬───────────┬───────────┬───────────┐             │  ┌────────────────────────────────────────────┐  │  │
│         │  │  │ Sun 7/5   │ Mon 7/6   │ Tue 7/7   │ Wed 7/8   │             │  │ Check building access times before         │  │  │
│         │  │  │ 8 avail   │ 6 avail   │ 7 avail   │ 5 avail   │             │  │ assigning early morning jobs.              │  │  │
│         │  │  ├───────────┼───────────┼───────────┼───────────┤             │  │                                            │  │  │
│         │  │  │ Pack      │ Move      │ Clean     │ Label     │             │  │ Keep tea, chargers, tape, markers, and     │  │  │
│         │  │  │ kitchen   │ furniture │ garage    │ boxes     │             │  │ bin bags in the first-day essentials box.  │  │  │
│         │  │  │ High 2/3  │ Med  2/2  │ Low  1/1  │ Low  1/2  │             │  └────────────────────────────────────────────┘  │  │
│         │  │  │ ⚠ under   │ ✓ full    │ ✓ full    │ ⚠ under   │             │                                                  │  │
│         │  │  └───────────┴───────────┴───────────┴───────────┘             │                                                  │  │
│         │  └────────────────────────────────────────────────────────────────┴──────────────────────────────────────────────────┘  │
│         │                                                                                                                       │
└─────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Key differences in collapsed state:**
- Sidebar reduces from ~260px to icon-only rail (~56px). Navigation items show as icon tooltips.
- Content area gains the reclaimed ~200px of horizontal space.
- All dashboard regions (KPI cards, Tasks Backlog, People Availability, Daily Schedule, Move Notes) remain present with the same proportions relative to the wider content area.

---

## Gap Analysis — Current `HomeView.vue` vs. Target Composition

The current placeholder at `frontend/src/home/HomeView.vue` renders four cards in two rows. The target composition requires a three-row dashboard with different data sources and components.

| Current `HomeView.vue` Element | Target Composition Element | Status |
|---|---|---|
| **Hello world** card (live backend message via `GET /api/hello`) | **Total Tasks** KPI card (from `GET /api/tasks/backlog` `summary.totalTasks`) | Replace. The Hello card is a connectivity probe, not a dashboard metric. |
| **Available today (6)** card (hardcoded, `UsersRoundIcon`) | **Available People** KPI card (from `GET /api/dashboard/people-availability` `summary.availableToday` / `summary.totalPeople`) | Replace. The hardcoded "6" must be sourced from the people-availability API. |
| **Under-staffed (3)** card (hardcoded, `MapPinnedIcon`) | **High Priority Tasks** KPI card (from `GET /api/tasks/backlog` `summary.highPriorityTasks`) | Replace. The under-staffed placeholder is replaced by the high-priority count from the task backlog API. |
| **Move days** card (`planWindowDayCount` via `usePlanningWindow()`) | **Move Days** KPI card (same data source) | **Retained**. Already wired to `usePlanningWindow()`. Re-render as one of four equal-width KPI cards in the first content row. |
| **Today's plan** card (static `upcomingWork` array, 3 hardcoded items) | **Tasks Backlog + People Availability** row (two reusable components side by side) | Replace. The static "Today's plan" is split into a live task backlog table (`TaskManagementPanel.vue`) and a live availability matrix (`PeopleAvailability.vue`). |
| **Move notes** card (static text snippets) | **Daily Schedule + Move Notes** row | Replace. The Move Notes panel moves to the third content row alongside Daily Schedule. Daily Schedule is aspirational (no Vue component). Move Notes remains aspirational (design-only, no spec). |
| **AppShell header** (page-title + Planning mode badge) | **Design toolbar** (date-range picker, Today button, ← → arrows, 🔔 alerts, 👤 profile chip) | **GAP**. `AppShell.vue` currently renders a simple `<header>` with `SidebarTrigger`, `Separator`, `Badge` (Planning mode), `<h1>` page title, and `<p>` description. The target toolbar from the design does not exist. This is a cross-cutting shell change, not scoped to `HomeView.vue` alone. |

### Additional Gaps

- **No API wiring**: `HomeView.vue` currently only calls `GET /api/hello` and `usePlanningWindow()`. The target composition requires wiring `GET /api/tasks/backlog`, `GET /api/dashboard/people-availability`, and `GET /api/dashboard/daily-schedule`.
- **No reusable component usage**: `TaskManagementPanel.vue` and `PeopleAvailability.vue` exist as standalone reusable components but are not imported or rendered in `HomeView.vue`.
- **Grid layout mismatch**: Current layout is a 2-row grid (4 KPI cards + 2 side-by-side cards). Target layout is a 3-row grid (4 KPI cards + Tasks/People row + Schedule/Notes row) with proportional column widths (Tasks Backlog ~60%, People Availability ~40%, etc.).

---

## Vocabulary Reference

All labels in the ASCII wireframe use existing repo/spec vocabulary:

| Context | Vocabulary | Source |
|---|---|---|
| **Task priorities** | High, Medium, Low | `openspec/specs/task-backlog-api/spec.md` — `priority` field; `priorities` legend |
| **Task statuses** | Backlog, Ready, Assigned | `openspec/specs/task-backlog-api/spec.md` — `status` field; `statuses` legend |
| **Availability statuses** | Available, Busy, Partial, Off | `openspec/specs/dashboard-people-availability/spec.md` — `status` field; `statuses` legend |
| **Staffing statuses** | fullyStaffed, underStaffed | `openspec/specs/dashboard-daily-schedule/spec.md` — `staffingStatus` field |
| **Navigation labels** | Dashboard, Tasks, Schedule, People, Rooms / Areas, Settings | `frontend/src/shared/layout/app-sidebar/AppSidebar.vue` — `planNavigation` and `organizationNavigation` arrays |
| **Section titles** | People availability, Task Management | `frontend/src/people/PeopleAvailability.vue` (default `title` prop), `frontend/src/tasks/components/TaskManagementPanel.vue` (`CardTitle`) |
| **Sidebar brand** | The Great Migration / House move planner | `AppSidebar.vue` header and footer project card text |
