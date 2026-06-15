## Context

Tech-onboarding scaffolding (archived as change `task-c7c555d5-94d9-4d1a-a31d-dc42fe95651c`) introduced a developer showcase route, extra sidebar sections, placeholder screen content, and a broad set of generated `shared/ui` primitives. The current design boards — `designs/design-system.png` (navigation, controls, badges, chips, spacing tokens) and `designs/components.png` (main-screen reusable UI pieces) — document the intended product UI but show no developer catalog, no "Developer" sidebar section, no "Move focus" prose block, and no onboarding footer copy.

The repo now mixes useful foundation pieces (Sheet, Tooltip, Button, Badge, Card, Input, Separator, active sidebar subcomponents) with onboarding-only artifacts that must be removed. The sidebar has dependency constraints that prevent naive deletion: `Sidebar.vue` uses Sheet for mobile behavior, and `SidebarMenuButton.vue` uses Tooltip for collapsed-state labels.

The `openspec/specs/component-showcase/spec.md` is a canonical spec that mandates the very artifacts being removed, so it must be superseded for internal consistency. `docs/design-system-v2.md` includes an out-of-scope clause excluding `designs/components.png` compositions, which creates a circular authority gap when the cleanup references that design board for navigation/layout decisions.

## Goals / Non-Goals

### Goals
- Remove all clearly onboarding-only developer artifacts: `/showcase` route, `ShowcaseView.vue`, Developer sidebar group, Move focus prose block, placeholder footer copy.
- Prune unused `shared/ui/sidebar` subcomponent files and their barrel exports.
- Remove `shared/ui/skeleton/` (zero product consumers after cleanup).
- Remove the local `TooltipProvider.vue` wrapper (zero consumers after cleanup).
- Update sidebar nav label from "Calendar" to "Schedule" to match `designs/design-system.png`.
- Supersede `openspec/specs/component-showcase/spec.md` with a retirement spec.
- Narrow the `docs/design-system-v2.md` out-of-scope clause to resolve the design-contract contradiction.

### Non-Goals
- Fully implementing the designed dashboard, backlog table, people matrix, schedule board, rooms/areas, or settings screens.
- Removing or replacing placeholder product views (Home, Tasks, Calendar, People) — these are the current shell awaiting follow-up work.
- Removing shared primitives with active dependencies (Sheet, Tooltip) or active product usage (Button, Badge, Card, Input, Separator).
- Adding Storybook or another developer-doc surface.
- Changing route paths (`/calendar` stays `/calendar`).
- Introducing new domain component directories.

## Decisions

### 1. Remove /showcase and src/showcase/ entirely
**Decision**: Delete the route from `routes.ts`, delete `ShowcaseView.vue`, and remove the `src/showcase/` directory.
**Rationale**: The showcase was created during tech onboarding and serves as a developer catalog. Neither `designs/design-system.png` nor `designs/components.png` shows a developer catalog surface. The non-goals explicitly exclude replacing the showcase with Storybook or another doc surface.
**Evidence**: `src/app/routes.ts` lines defining `/showcase`; `src/showcase/ShowcaseView.vue` imports and demos 7+ primitives; archived change `task-c7c555d5-94d9-4d1a-a31d-dc42fe95651c` confirms onboarding provenance.

### 2. Remove Developer sidebar group, Move focus block, and footer copy from AppSidebar
**Decision**: Remove the "Developer" `SidebarGroup`, the "Move focus" prose block, their `SidebarSeparator` delimiters, and replace the footer text.
**Rationale**: None of these elements appear in `designs/design-system.png` which lists Dashboard, Tasks, Schedule, People, Rooms/Areas, and Settings as the intended navigation items.
**Evidence**: `AppSidebar.vue` lines 76-105 (Move focus and Developer groups), line 113 (footer copy).

### 3. Remove 7 unused sidebar subcomponents
**Decision**: Delete the `.vue` files and remove their barrel exports for: `SidebarGroupAction`, `SidebarInput`, `SidebarMenuAction`, `SidebarMenuSkeleton`, `SidebarMenuSub`, `SidebarMenuSubButton`, `SidebarMenuSubItem`.
**Rationale**: Grep confirms zero consumers of these exports outside their barrel export in `index.ts`.
**Evidence**: `src/shared/ui/sidebar/index.ts` export list; grep result across all `src/**/*.{ts,vue}`.

### 4. Remove Skeleton primitive
**Decision**: Delete `src/shared/ui/skeleton/Skeleton.vue`, its barrel export, and the directory. Remove `SidebarMenuSkeleton.vue` first (it imports Skeleton).
**Rationale**: Skeleton's only consumers are `ShowcaseView.vue` (being removed) and `SidebarMenuSkeleton.vue` (unused and being removed). Zero product consumers remain.
**Evidence**: `Skeleton.vue` import usage verified via grep; `SidebarMenuSkeleton.vue` confirmed unused outside barrel.

### 5. Remove local TooltipProvider wrapper
**Decision**: Delete `src/shared/ui/tooltip/TooltipProvider.vue` and remove its barrel export.
**Rationale**: The only consumer (`ShowcaseView.vue`) is being removed. `SidebarProvider.vue` imports `reka-ui`'s `TooltipProvider` directly with `:delay-duration="0"`. No other code imports the local wrapper.
**Note**: `SidebarProvider` uses `delay-duration=0` vs the wrapper's default `700`. If a future consumer needs the 700ms default, it can be re-added — this is not a concern for current usage.

### 6. Retain Sheet and Tooltip
**Decision**: Do NOT remove `src/shared/ui/sheet/` or `src/shared/ui/tooltip/` (except the local `TooltipProvider.vue` wrapper).
**Rationale**: `Sidebar.vue` imports Sheet/SheetContent for mobile sidebar behavior. `SidebarMenuButton.vue` imports Tooltip/TooltipContent/TooltipTrigger for collapsed-sidebar label popups. Neither can be removed without breaking sidebar functionality.
**Evidence**: `src/shared/ui/sidebar/Sidebar.vue` lines importing from `@/shared/ui/sheet`; `src/shared/ui/sidebar/SidebarMenuButton.vue` lines importing from `@/shared/ui/tooltip`.

### 7. Retain placeholder product views
**Decision**: Do NOT remove or modify `HomeView.vue`, `TasksView.vue`, `CalendarView.vue`, or `PeopleView.vue` beyond what the sidebar/navigation label changes indirectly affect.
**Rationale**: These are the current product UI shell. Removing them would leave empty routes. Rebuilding them to match the design comps is explicitly out of scope per acceptance criteria non-goals.
**Evidence**: Placeholder view contents ("Task foundation", "Schedule board foundation", "People availability foundation"); acceptance criteria: "keep the implementation scope framed as cleanup/alignment work, not a full rebuild of the dashboard."

### 8. Update sidebar nav label "Calendar" → "Schedule"
**Decision**: Change the display label in `AppSidebar.vue` `primaryNavigation` array. Do NOT change the route path or component name.
**Rationale**: `designs/design-system.png` lists "Schedule" in the navigation, not "Calendar". Changing only the display label avoids cascading route/path changes while aligning the visible UI with the design board.
**Evidence**: `AppSidebar.vue` line 25: `{ title: 'Calendar', to: '/calendar', ... }`; `designs/design-system.png` navigation list.

### 9. Supersede component-showcase OpenSpec spec
**Decision**: Create a retirement spec at `specs/component-showcase-retirement/spec.md` documenting that the `/showcase` route, Developer sidebar section, and component catalog have been removed. The original spec at `openspec/specs/component-showcase/spec.md` is superseded.
**Rationale**: The canonical spec currently mandates the artifacts being removed. Leaving it active would create a spec-implementation contradiction.
**Evidence**: `openspec/specs/component-showcase/spec.md` requirements; archived copy at `changes/archive/2026-06-15-task-c7c555d5-94d9-4d1a-a31d-dc42fe95651c/specs/component-showcase/spec.md` preserves the original.

### 10. Resolve design-contract contradiction
**Decision**: Amend `docs/design-system-v2.md` to narrow the out-of-scope clause: `designs/design-system.png` governs navigation and layout decisions; the composition exclusion from `designs/components.png` applies only to full dashboard compositions, not to layout/navigation source-of-truth decisions.
**Rationale**: The task prompt requires alignment with `designs/design-system.png` for navigation, but `docs/design-system-v2.md` broadly excludes compositions from `designs/components.png`. Clarifying the scope prevents future drift.
**Evidence**: `docs/design-system-v2.md` out-of-scope clause; task prompt requirement; reviewer condition.

## Risks / Trade-offs

| Risk | Severity | Mitigation |
|------|----------|------------|
| Sidebar feels sparse after removing Move focus and Developer sections | Low | Acceptable. Rooms/Areas and Settings navigation items will be added in future work. |
| Nav label "Schedule" diverges from route path `/calendar` | Low | Acceptable. Route path is an internal identifier; the display label is user-facing. |
| Sheet/Tooltip accidentally removed by over-eager cleanup | Medium | Excluded from removal tasks with explicit DO NOT REMOVE annotation. |
| Skeleton removal misses stale barrel export, causing import failures | Medium | Remove the entire `src/shared/ui/skeleton/` directory, not just the `.vue` file. |
| Placeholder views mistaken for design-approved content | Low | Documented in the non-goals and retained-items section. |
| Future developer needs the local TooltipProvider wrapper's 700ms delay defaults | Low | Can be re-added trivially. No current consumer needs it. |

## Traceability

| Source | Artifact | Decision |
|--------|----------|----------|
| Task prompt 94fecd0e | Task description | "Remove scaffolded components … get the repo back in line with designs" |
| Dossier `affectedAreas[0]` | Routing/navigation shell | Remove /showcase route, Developer sidebar group |
| Dossier `affectedAreas[2]` | Shared UI primitive inventory | Remove unused sidebar exports, Skeleton, local TooltipProvider |
| Architect R1 | Decision `swarm-architect-recommendation` | Core removal targets + retention of Sheet/Tooltip |
| Lead-Dev R1 | Decision `swarm-lead-dev-recommendation` | Precise cleanup list + TooltipProvider + nav label update |
| Reviewer R1 | Decision `swarm-reviewer-recommendation` | APPROVE with 4 conditions → all resolved |
| Architect blocker-1 | Non-blocking claim | Supersede component-showcase spec → Decision 9 |
| Lead-Dev blocker-1 | Non-blocking claim | Supersede component-showcase spec → Decision 9 |
| Lead-Dev blocker-2 | Non-blocking claim | Design contract contradiction → Decision 10 |
| Reviewer blocker-1 | Non-blocking claim | Design contract contradiction → Decision 10 |
| Architect blocker-2 | Non-blocking claim | Skeleton removal order → Decision 4 |
| `src/shared/ui/sidebar/Sidebar.vue` | Evidence | Sheet dependency confirmed → Decision 6 |
| `src/shared/ui/sidebar/SidebarMenuButton.vue` | Evidence | Tooltip dependency confirmed → Decision 6 |
| `src/app/routes.ts` | Evidence | /showcase route confirmed → Decision 1 |
| `src/shared/layout/app-sidebar/AppSidebar.vue` | Evidence | Developer/Move focus/footer confirmed → Decision 2 |
| `openspec/specs/component-showcase/spec.md` | Evidence | Canonical spec mandates showcase → Decision 9 |
| `docs/design-system-v2.md` | Evidence | Out-of-scope contradiction → Decision 10 |
