# home-page-ascii-blueprint Specification

## Purpose
TBD - created by archiving change task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a. Update Purpose after archive.
## Requirements
### Requirement: Artifact SHALL be a non-canonical Markdown file under the change directory

The ASCII mockup SHALL be placed at `openspec/changes/task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a/home-page-ascii.md`. It SHALL NOT modify any file under `openspec/specs/`, `docs/`, `frontend/src/`, or `backend/`.

#### Scenario: Artifact exists at the correct location
- **WHEN** the change is complete
- **THEN** `openspec/changes/task-2bdd3fb5-c4ec-4153-9a49-d6b463894a8a/home-page-ascii.md` exists and contains markdown content

#### Scenario: No canonical files are modified
- **WHEN** the change is complete
- **THEN** no files under `openspec/specs/` have been modified
- **AND** no files under `docs/` have been modified
- **AND** no files under `frontend/src/` have been modified
- **AND** no files under `backend/` have been modified

### Requirement: Artifact SHALL include a legend mapping regions to repo sources

The artifact SHALL contain a legend section that maps each dashboard region to its corresponding repo component file path, OpenSpec spec reference, and data contract, using exact file paths rather than directory names.

#### Scenario: Legend maps sidebar region correctly
- **WHEN** the legend is read
- **THEN** the sidebar entry references `frontend/src/shared/layout/app-sidebar/AppSidebar.vue`
- **AND** the sidebar entry notes the component supports `collapsible="icon"`

#### Scenario: Legend maps top toolbar with gap annotation
- **WHEN** the legend is read
- **THEN** the toolbar entry references `frontend/src/shared/layout/app-shell/AppShell.vue`
- **AND** the toolbar entry documents the gap between the current page-title header and the aspirational design toolbar (date-range, Today, arrows, alerts, profile chip)

#### Scenario: Legend maps KPI cards to API contracts
- **WHEN** the legend is read
- **THEN** the KPI cards entry references `openspec/specs/task-backlog-api/spec.md` and `openspec/specs/dashboard-people-availability/spec.md` as data sources

#### Scenario: Legend maps Tasks Backlog to existing component
- **WHEN** the legend is read
- **THEN** the Tasks Backlog entry references `frontend/src/tasks/components/TaskManagementPanel.vue` as an existing reusable component

#### Scenario: Legend maps People Availability to existing component
- **WHEN** the legend is read
- **THEN** the People Availability entry references `frontend/src/people/PeopleAvailability.vue` as an existing reusable component

#### Scenario: Legend documents aspirational regions
- **WHEN** the legend is read
- **THEN** the Daily Schedule entry references `openspec/specs/dashboard-daily-schedule/spec.md` and SHALL note that no Vue component exists yet
- **AND** the Move Notes entry SHALL note that no spec or component exists — it is a static design-only panel

### Requirement: ASCII wireframe SHALL show all required dashboard regions in correct spatial hierarchy

The desktop-width ASCII wireframe (expanded sidebar) SHALL include the following regions in order from top to bottom, using monospaced character alignment:

1. Left sidebar (persistent, reflecting `AppSidebar.vue` navigation structure)
2. Top toolbar (with gap annotation noting aspirational elements)
3. Row 1: Four KPI cards (Total Tasks, Available People, High Priority Tasks, Move Days)
4. Row 2: Tasks Backlog (wider) beside People Availability
5. Row 3: Daily Schedule (wider) beside Move Notes

#### Scenario: Wireframe includes the sidebar region
- **WHEN** the wireframe is read
- **THEN** a left-rail sidebar region is visible with navigation items matching `AppSidebar.vue`: Dashboard, Tasks, Schedule, People, Rooms / Areas, Settings
- **AND** the sidebar shows the project card ("The Great Migration", "House move planner")

#### Scenario: Wireframe includes four KPI cards in the first content row
- **WHEN** the wireframe is read
- **THEN** four KPI cards are shown in a horizontal row
- **AND** the cards are labeled: Total Tasks, Available People, High Priority Tasks, Move Days
- **AND** each card label maps to a data contract noted in the legend

#### Scenario: Wireframe includes Tasks Backlog and People Availability in the second content row
- **WHEN** the wireframe is read
- **THEN** a wider Tasks Backlog section and a narrower People Availability section are shown side by side
- **AND** the Tasks Backlog section shows a table-like structure with task rows
- **AND** the People Availability section shows a matrix with person rows and day columns

#### Scenario: Wireframe includes Daily Schedule and Move Notes in the third content row
- **WHEN** the wireframe is read
- **THEN** a wider Daily Schedule section and a narrower Move Notes section are shown side by side
- **AND** both sections include aspirational annotations noting no Vue components exist

### Requirement: Artifact SHALL include a collapsed-sidebar ASCII variant

The artifact SHALL include a second narrower ASCII wireframe variant showing the homepage layout when the sidebar is collapsed to icon-only state, since `AppSidebar` already supports `collapsible="icon"`.

#### Scenario: Collapsed variant shows content-area reflow
- **WHEN** the collapsed variant is read
- **THEN** the sidebar region is shown as an icon-only rail (narrower)
- **AND** the content area width expands to fill the reclaimed space
- **AND** all dashboard regions (KPI cards, Tasks Backlog, People Availability, Daily Schedule, Move Notes) remain present

### Requirement: Artifact SHALL include a gap analysis section

The artifact SHALL include a gap-analysis section that catalogs specific differences between the current `frontend/src/home/HomeView.vue` placeholder composition and the target ASCII layout.

#### Scenario: Gap analysis documents placeholder cards to be replaced
- **WHEN** the gap analysis is read
- **THEN** it documents that the Hello world card SHALL be replaced by the Total Tasks KPI
- **AND** it documents that the Available today and Under-staffed cards SHALL be replaced by Available People and High Priority Tasks KPIs
- **AND** it documents that the Today's plan card SHALL be replaced by the Tasks Backlog + People Availability row
- **AND** it documents that the Move notes card SHALL be replaced by the Daily Schedule + Move Notes row

#### Scenario: Gap analysis documents the toolbar mismatch
- **WHEN** the gap analysis is read
- **THEN** it documents that `AppShell.vue` currently renders a page-title header with Planning mode badge
- **AND** it notes the design target toolbar (date-range, Today, arrows, alerts, profile chip) does not exist in the current shell

### Requirement: Section labels SHALL use existing repo and spec vocabulary

All section names, region labels, and sample content in the ASCII wireframe SHALL align with current repo/spec vocabulary rather than inventing new terminology.

#### Scenario: Labels match backend contract vocabulary
- **WHEN** the wireframe is read
- **THEN** priority labels use "High", "Medium", "Low" (matching `task-backlog-api/spec.md` priority vocabulary)
- **AND** availability statuses use "Available", "Busy", "Partial", "Off" (matching `dashboard-people-availability/spec.md` status vocabulary)
- **AND** task statuses use "Backlog", "Ready", "Assigned" (matching `task-backlog-api/spec.md` status vocabulary)

#### Scenario: Labels match existing component naming
- **WHEN** the wireframe is read
- **THEN** navigation labels match `AppSidebar.vue`: "Dashboard", "Tasks", "Schedule", "People", "Rooms / Areas", "Settings"
- **AND** the section title "People availability" matches `PeopleAvailability.vue` default title

### Requirement: Artifact SHALL be usable by a text-only LLM

The ASCII mockup SHALL communicate desktop layout relationships clearly enough for a text-only LLM to infer row/column structure and hierarchy without opening PNG files.

#### Scenario: Artifact is self-contained
- **WHEN** the artifact is read by a text-only LLM
- **THEN** the LLM can identify all major dashboard regions and their spatial relationships
- **AND** the LLM can identify which regions have existing Vue components vs. which are aspirational
- **AND** the LLM can identify which API contracts back each data-displaying region
