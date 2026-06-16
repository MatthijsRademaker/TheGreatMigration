## 1. Model the KPI card sub-elements

- [ ] 1.1 Add four `component` sub-elements inside the `dashboard` feature area in `.devagent/architecture/workspace.c4`: `peopleKpi`, `highPriorityKpi`, `unassignedKpi`, and `roomsKpi`, each with a description that matches the spec
- [ ] 1.2 Optionally tag each card element with `#feature` for visual consistency with other feature-area sub-elements

## 2. Add summarizes relationships for each data-backed card

- [ ] 2.1 Add `dashboard.peopleKpi -[summarizes]-> peopleData "Confirmed helpers count"` relationship
- [ ] 2.2 Add `dashboard.highPriorityKpi -[summarizes]-> taskData "High priority task count"` relationship
- [ ] 2.3 Add `dashboard.unassignedKpi -[summarizes]-> taskData "Unassigned job count"` relationship
- [ ] 2.4 Verify `dashboard.roomsKpi` has NO `summarizes` relationship (placeholder)

## 3. Refine existing dashboard-level summarizes relationships

- [ ] 3.1 Remove the generic `dashboard -[summarizes]-> taskData "High priority and under-staffed work"` relationship (now expressed at card level)
- [ ] 3.2 Remove the generic `dashboard -[summarizes]-> peopleData "Available helpers"` relationship (now expressed at card level)
- [ ] 3.3 Keep the `dashboard -[summarizes]-> scheduleData "Move days and next practical steps"` relationship (no card-level replacement yet)

## 4. Add the scoped dashboard KPIs view

- [ ] 4.1 Add a `view dashboard-kpis of migration.spa.dashboard` block in the `views { }` section with title "Dashboard KPI Cards" and a description explaining the data flow
- [ ] 4.2 Add `include *` to include the dashboard and its four card children
- [ ] 4.3 Add `include * -> migration.spa.dashboard.*` to show incoming data flow to the cards
- [ ] 4.4 Add `include migration.spa.taskData`, `include migration.spa.peopleData`, `include migration.spa.scheduleData` to show the data store dependencies
- [ ] 4.5 Optionally apply a `style` rule to color the KPI cards distinctly from the data stores (e.g., green for dashboard children, slate for data stores)
- [ ] 4.6 Add `autoLayout LeftRight` for a clear data-left-to-cards-right flow

## 5. Validate the updated LikeC4 model

- [ ] 5.1 Run `likec4 validate --json --no-layout --file .devagent/architecture/workspace.c4 .devagent/architecture/` and verify `filteredErrors` is 0
- [ ] 5.2 Verify all four existing views (`system-context`, `container-overview`, `feature-map`, `planning-flow`) have no validation errors
- [ ] 5.3 Run `scripts/precommit-run` from the repo root
