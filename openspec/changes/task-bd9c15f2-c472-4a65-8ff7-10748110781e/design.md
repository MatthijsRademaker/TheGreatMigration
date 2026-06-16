## Context

The current home route already renders `KpiCards`, but the lower half of `frontend/src/home/HomeView.vue` is still a placeholder composition built from a `Today’s plan` card and an inline `Move notes` card. The repository now contains reusable `TaskManagementPanel`, `PeopleAvailability`, and `DailySchedule` components that match the dashboard regions described by the archived ASCII homepage blueprint. The blueprint is still useful for row hierarchy, but parts of it are stale: it predates the current `DailySchedule.vue` component and uses KPI labels that no longer match the current `KpiCards.vue` contract.

## Goals

- Compose the `/` route into the dashboard hierarchy shown by the ASCII blueprint: KPI row, Tasks/People row, and Daily Schedule/Move Notes row.
- Reuse the existing `KpiCards`, `TaskManagementPanel`, `PeopleAvailability`, and `DailySchedule` components from their current feature locations.
- Keep the change frontend-presentational and composition-only.
- Preserve the existing design-system direction and responsive layout patterns already used by the dashboard surfaces.
- Update route-render coverage so `/` asserts the new composed sections instead of `Today’s plan`.

## Non-Goals

- Adding the aspirational top toolbar from the archived ASCII blueprint.
- Changing `AppShell.vue`, `AppSidebar.vue`, routes, route metadata, or navigation behavior.
- Reworking the internals or data contracts of `KpiCards`, `TaskManagementPanel`, `PeopleAvailability`, or `DailySchedule` beyond minimal composition needs.
- Wiring `TaskManagementPanel`, `PeopleAvailability`, or `DailySchedule` to live backend data in this slice.
- Creating a reusable Move Notes component or adding notes CRUD.
- Regenerating clients, editing backend code, or publishing canonical OpenSpec artifacts.

## Decisions

### 1. Keep `HomeView.vue` as the composition root

Implement the change entirely in `frontend/src/home/HomeView.vue`. It already owns the `/` route composition, and the accepted scope does not require route, shell, or sidebar changes.

### 2. Prefer current implemented components and current specs over stale ASCII details

Use the archived ASCII artifact as the spatial blueprint, but keep the current repository vocabulary where the archive is outdated. That means retaining `KpiCards.vue` and its current KPI labels (`People available today`, `High priority tasks`, `Unassigned jobs`, `Rooms completed`) and composing the existing `DailySchedule.vue` component instead of treating Daily Schedule as aspirational.

### 3. Compose two independent dashboard rows beneath `KpiCards`

Render one responsive grid for `TaskManagementPanel` beside `PeopleAvailability`, and a second responsive grid for `DailySchedule` beside the static `Move notes` card. Reuse the existing `xl:grid-cols-[1.4fr_0.9fr]` proportion so the wider left-column panels mirror the dashboard blueprint without introducing new layout primitives.

### 4. Remove placeholder-only home content instead of adapting it

Delete the `upcomingWork` array, the `Today’s plan` card, and the now-unused `Badge` import from `HomeView.vue`. Keep the existing `Move notes` reminder copy inside a `Card` so the route gains the intended structure without introducing a new notes feature.

### 5. Limit verification to the existing route-render surface and precommit script

Update `frontend/tests/app-routes-render.test.ts` so the home route asserts `Task Management`, `People availability`, `Daily Schedule`, and `Move notes`, uses `Move notes` as the unique `/` route content string, and removes the old `Today’s plan` assertions while preserving the current KPI mock assertions. Run `scripts/precommit-run` before claiming completion.

## Risks

- Reusing the existing `xl:grid-cols-[1.4fr_0.9fr]` split for two stacked rows could expose spacing or overflow issues at intermediate breakpoints.
- `TaskManagementPanel` emits placeholder actions that will remain unhandled on the home route in this slice, though that is a low-risk cosmetic issue.
- `DailySchedule` demo content will now appear on both `/calendar` and `/`, which is acceptable for this presentational slice but may need future data reconciliation.

## Conflict Resolution

- **KPI vocabulary:** Keep the current `KpiCards.vue` labels and canonical KPI spec rather than reviving the archived ASCII card labels.
- **Daily Schedule availability:** Compose the existing `frontend/src/calendar/DailySchedule.vue` component and treat the archived “no Vue component exists” note as stale.
- **Toolbar scope:** Leave the current `AppShell.vue` header unchanged; the richer top-toolbar concept is explicitly deferred.
- **Home-route test anchor:** Use `Move notes` as the unique `/` route content string in `routeCases` because `Task Management` and `Daily Schedule` also appear on other routes.

## Traceability

- Task: `bd9c15f2-c472-4a65-8ff7-10748110781e`
- Dossier: `2026-06-16T06:43:58.786Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial`
