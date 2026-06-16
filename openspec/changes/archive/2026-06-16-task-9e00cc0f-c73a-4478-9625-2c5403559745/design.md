## Context

The LikeC4 architecture model at `.devagent/architecture/workspace.c4` currently models the Dashboard as a single `featureArea` element with a generic description and three high-level `summarizes` relationships to task, people, and schedule data stores. The frontend implementation now renders four distinct KPI summary cards — each with a specific label, value source, icon, and data dependency:

| Card                         | Data Source       | Status        |
| ---------------------------- | ----------------- | ------------- |
| People available today       | `peopleData`      | Live (backend) |
| High priority tasks          | `taskData`        | Live (backend) |
| Unassigned jobs              | `taskData`        | Live (backend) |
| Rooms completed              | (static)          | Placeholder   |

The architecture diagram does not reflect this decomposition, making it harder for readers to understand which data flows serve which card — and creating a blind spot when the fourth placeholder card inevitably gains a real data source.

## Goals / Non-Goals

**Goals:**
- Decompose the Dashboard feature area in the C4 model into four sub-elements, one per KPI card, each with a targeted `summarizes` relationship to its data source.
- Add a scoped LikeC4 component view (`dashboard-kpis`) that renders only the Dashboard, its four KPI card children, and their data-flow relationships.
- Update the Dashboard's existing relationship to `taskData` from a single generic summary to two explicit summaries (high-priority tasks, unassigned jobs).
- Keep all existing views (`system-context`, `container-overview`, `feature-map`, `planning-flow`) rendering correctly after the change.

**Non-Goals:**
- Changing any frontend or backend source files — no Vue components, no Go handlers, no SQL queries.
- Regenerating the OpenAPI snapshot or client artifacts.
- Adding or modifying application tests.
- Changing the C4 specification (element/relationship kinds remain the same).
- Adding any new LikeC4 configuration, dependencies, or build tooling.

## Decisions

### 1. KPI cards as sub-elements of the Dashboard feature area

Model the four KPI cards as direct child elements of `migration.spa.dashboard` using the `component` kind with a `#feature` tag. Sub-elements inside a `featureArea` are naturally grouped under the parent in scoped views, which matches the visual hierarchy (cards inside the dashboard page) and prevents polluting the top-level feature-area namespace.

**Alternatives considered:**
- **Separate `featureArea` per card**: Too heavyweight. Cards are part of the dashboard, not independent page-level features.
- **Sibling `component` elements**: Would require explicit grouping in views and lose the parent-child hierarchy that LikeC4 scoping provides.
- **No sub-elements, just refine descriptions**: Would lose the structural clarity that makes the four-card contract visible in the diagram.

Chosen: Sub-elements inside the dashboard `featureArea { }` block.

### 2. One shared `summarizes` relationship for both task-backlog cards

Both "High priority tasks" and "Unassigned jobs" pull from the same `taskData` store. Model one `summarizes` relationship for each card to its data source — `taskBacklog` → `taskData` and `unassignedJobs` → `taskData` — making the shared data source explicit while keeping each card's semantic distinct.

### 3. Rooms completed as a zero-data placeholder

The "Rooms completed" card currently has no backend data source. Model it with a `component` element and a `notes` property explaining it is a placeholder, but give it no `summarizes` relationship. This makes the gap visible in the architecture diagram, clearly signaling the future integration point.

### 4. Scoped component view for the dashboard

Add a new scoped view `dashboard-kpis of migration.spa.dashboard` that includes `*` (the dashboard + its four card children) plus the three `planningData` elements that the cards summarize. This keeps the view focused and avoids showing unrelated features (tasks, calendar, people feature areas).

## Risks / Trade-offs

- **[LikeC4 syntax gap]** KPI card sub-elements must use identifiers that are valid LikeC4 identifiers (letters, digits, hyphens, underscores; no dots). Discovery solution: keep identifiers short and flat (e.g., `peopleKpi`, `highPriorityKpi`, `unassignedKpi`, `roomsKpi`).
- **[View layout drift]** Adding new elements and relationships may change the auto-layout of existing views (`feature-map`, `planning-flow`). Mitigation: validate with `likec4 validate --json --no-layout` after editing, and if views need layout tuning, re-run `likec4 export` to regenerate rendered diagrams.
- **[Cross-file FQN requirement]** Sub-elements defined inside a nested block still use relative FQN resolution within the file. No cross-file references are needed since this is a single-file project.
