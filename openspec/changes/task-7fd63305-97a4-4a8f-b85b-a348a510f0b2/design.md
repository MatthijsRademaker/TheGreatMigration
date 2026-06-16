## Context

The app shell (`AppShell.vue`) wraps every route and currently renders two problematic elements: a `SidebarTrigger` in the content header (outside the sidebar boundary it controls) and a static `<Badge>Planning mode</Badge>`. The canonical move-planning timeline is available via `GET /api/planning-window` and consumed by `usePlanningWindow()`, but the header ignores it. The design system (`docs/design-system-v2.md`) specifies date-range controls in the top bar, and visual artifacts (`designs/home-page.png`, `designs/design-system.png`) show a sidebar-first composition without a freestanding header nav toggle.

The sidebar infrastructure is complete: shadcn-vue primitives (`SidebarTrigger`, `SidebarMenuButton`, `SidebarHeader`) are in `frontend/src/shared/ui/sidebar/`, semantic sidebar tokens are defined in `styles.css`, and `SidebarTrigger` is placement-agnostic (it uses `useSidebar()` context). The planning-window composable already queries the API but lacks a formatted range output. The SSR test already mocks `/api/planning-window`, making the visualizer testable.

A previous home-dashboard change (`task-bd9c15f2`) explicitly deferred `AppShell`/`AppSidebar` changes, and an archived blueprint (`home-page-ascii.md`) documents the richer toolbar as a current gap — this task closes both.

## Goals / Non-Goals

### Goals
- Relocate the primary visible `SidebarTrigger` from `AppShell.vue` header into `AppSidebar.vue`'s `SidebarHeader` as the last element after the branding row, preserving collapse/expand and icon-collapsed mode behavior.
- Replace the `Planning mode` badge with a compact read-only date-range visualizer sourced from the canonical `usePlanningWindow()` composable.
- Extend `planWindow.ts` with an SSR-safe `formatPlanWindowRange()` pure formatter using ISO string inputs for deterministic UTC labels.
- Extend `usePlanningWindow()` with a `formattedRange` computed that delegates to the new formatter.
- Update the SSR route-render test to assert the new shell chrome and drop the `Planning mode` assertion.
- Reuse existing sidebar primitives, semantic CSS tokens, and lucide icons exclusively — no new component directories or raw styling.

### Non-Goals
- Implementing the full aspirational top toolbar (`Today`, arrows, notifications, profile chip).
- Making the planning window editable, persisting new dates, or adding a settings UI for date management.
- Introducing new backend endpoints, new route behavior, or cross-page filtering tied to the header visualizer.
- Reworking dashboard panels, notes/help actions, or sidebar information architecture beyond the trigger relocation.
- Adding an interactive range picker or previous/next date navigation controls.
- Retaining a separate `Planning mode` status chip alongside the visualizer.

## Decisions

### D1: SidebarTrigger placed in SidebarHeader as last element after branding
**Decision**: Move `SidebarTrigger` from `AppShell.vue` into `AppSidebar.vue`'s `SidebarHeader` as the last `SidebarMenuButton` item after the existing branding `RouterLink`. Use `SidebarMenuButton` with `variant="ghost"` for proper alignment in both expanded and icon-collapsed modes.
**Rationale**: All three refinement participants agreed on `SidebarHeader`. The lead-dev specifically identified this placement keeps the trigger visible at the top of the mobile sheet (addressing the suppressed `SheetContent` close button concern). `SidebarFooter` is scrolled-to on mobile and would make the trigger undiscoverable. `SidebarTrigger` is placement-agnostic (`useSidebar()` context), so this is purely a template relocation.
**Evidence**: Architect R1 suggested requirement #2; Lead-dev R1 suggested requirement #1; Reviewer R1 blocker #2 about undefined placement resolution; dossier open question about trigger location.

### D2: formatPlanWindowRange accepts ISO strings, not Date objects
**Decision**: The new `formatPlanWindowRange(startDate: string, endDate: string, days: number)` pure helper SHALL accept ISO 8601 date strings (from the raw API `PlanningWindowBody`) and use `Intl.DateTimeFormat` with `en-US` locale and `UTC` timezone to produce deterministic labels. It SHALL NOT accept `Date` objects.
**Rationale**: The reviewer identified that `planWindowDays` entries contain raw `Date` objects (`new Date(cursor)`) which Vue's SSR runtime cannot serialize. The raw API response fields (`startDate`, `endDate`) are ISO strings that are SSR-safe. Using UTC `Intl.DateTimeFormat` mirrors the existing `formatPlanDayLabel` pattern and guarantees deterministic output regardless of server timezone. This approach also avoids coupling the header visualizer to the day-by-day `planWindowDays` array.
**Evidence**: Reviewer R1 blocker #1 (SSR Date serialization); Architect R1 suggested requirement #3; Lead-dev R1 suggested requirement #2; `planWindow.ts` existing `formatPlanDayLabel` UTC pattern; `PlanningWindowBody` type exposes `startDate: string`, `endDate: string`, `days: number`.

### D3: Visualizer format: compound range string with day count
**Decision**: `formatPlanWindowRange()` SHALL return a single compound string in the format `"5 Jul – 13 Aug 2026 · 40 days"` using an en-dash separator and middle-dot before the day count. The day count is included as a companion label.
**Rationale**: The architect proposed this compound format. The compact single-string format simplifies the template (no need for multiple `<span>` elements or complex conditional composition) while surfacing all three canonical pieces of planning-window information (start, end, inclusive days). The en-dash and middle-dot conventions produce a professional, design-system-aligned appearance.
**Evidence**: Architect R1 suggested requirement #3 (proposed compound format); Lead-dev R1 question #2 about day count inclusion; Acceptance criteria #3 from dossier; `PlanWindowBody` type includes `days`.

### D4: Planning mode badge completely removed
**Decision**: The `<Badge>Planning mode</Badge>` in `AppShell.vue` SHALL be removed entirely. No smaller status chip or indicator SHALL remain alongside the range visualizer. The visualizer is a complete replacement.
**Rationale**: All three refinement participants treat the visualizer as a direct replacement for the badge. A separate status chip would add visual noise and serve no informational purpose once the actual planning-window dates are displayed. The acceptance criteria call the visualizer a replacement, not a companion.
**Evidence**: Dossier acceptance criteria #1 (visualizer replaces badge); Architect recommendation (replace the badge); Lead-dev recommendation (replace the badge); Reviewer question #3 about badge coexistence.

### D5: usePlanningWindow gains formattedRange computed
**Decision**: Add a `formattedRange` computed to `usePlanningWindow()` that derives `startDate`, `endDate`, and `days` from `query.data.value` (the raw API response) and delegates to `formatPlanWindowRange()`. It SHALL return `null` when data is not yet available.
**Rationale**: Keeping formatting logic out of `AppShell.vue`'s template follows separation of concerns. The composable already owns planning-window state; adding a derived formatted output is a natural extension that the template can consume with a single `v-if`/`v-else` chain.
**Evidence**: Lead-dev R1 suggested requirement #3; Architect R1 suggested requirement #4; dossier identifies composable formatting gap.

### D6: Existing sidebar configuration preserved
**Decision**: The `Sidebar` root SHALL retain its `collapsible="icon"` and `variant="inset"` configuration. The `SidebarRail` and its keyboard shortcut behavior SHALL remain unchanged. The `SidebarHeader` brand block (`NotebookTabsIcon`, "The Great Migration", "House move planner") SHALL be unchanged.
**Rationale**: The dossier explicitly states the sidebar rail and keyboard shortcuts remain acceptable secondary collapse affordances. The scope is header cleanup, not sidebar redesign. The footer (with its redundant brand block and utility actions) is also preserved as-is — cleaning that up is out of scope.
**Evidence**: Dossier assumptions (rail/keyboard preserved); Non-goals (no sidebar information architecture changes); archived sidebar-navigation spec requirement about preserving configuration.

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Planning-window query returns error or loading state during SSR, causing visualizer slot to render empty or break layout | Medium | Medium | The `formattedRange` computed returns `null` on loading/error. AppShell template uses `v-if` to render the visualizer only when data is available, with a fixed-height container to prevent layout shift. |
| SidebarTrigger in SidebarHeader creates visual misalignment in icon-collapsed mode with the branding icon square | Low | Low | Use `SidebarMenuButton` with the existing `size="lg"` pattern. The sidebar primitives handle icon-collapsed sizing via `group-data-[collapsible=icon]:size-8!`. Visual check confirms alignment. |
| Mobile sheet close affordance: Sidebar.vue suppresses the default `SheetContent` close button (`[&>button]:hidden`), so moving the trigger inside means mobile users must find it within the sidebar | Medium | Medium | Placing the trigger in `SidebarHeader` keeps it at the very top of the mobile sheet where it is immediately visible. `SidebarFooter` placement was rejected specifically to avoid this issue. |
| Test false-negative: if `formattedRange` output changes format during implementation, the SSR test assertion will fail | Medium | Low | Fix the formatter output first, then align test assertions. The test's mock data is fixed (`2026-07-05`/`2026-08-13`/40), making the expected format deterministic. |
| Visualizer loading state causes layout shift when data arrives after client hydration | Low | Low | Use a fixed-height container with `inline-block` sizing for the visualizer area. The loading skeleton placeholder and data state occupy the same vertical space. |

## Traceability

- **Task**: `7fd63305-97a4-4a8f-b85b-a348a510f0b2`
- **Design contract**: `docs/design-system-v2.md` — navigation states, top bar date-range controls
- **Visual artifacts**: `designs/home-page.png`, `designs/design-system.png`
- **Source files**: `AppShell.vue`, `AppSidebar.vue`, `planWindow.ts`, `usePlanningWindow.ts`
- **Sidebar primitives**: `frontend/src/shared/ui/sidebar/` (SidebarTrigger, SidebarMenuButton, useSidebar)
- **SSR test**: `frontend/tests/app-routes-render.test.ts`
- **API contract**: `GET /api/planning-window` → `PlanningWindowBody` (startDate, endDate, days)
- **Dossier**: `2026-06-16T07:22:29.526Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Rounds**: 1 (architect, lead-dev, reviewer)
- **Preceding deferral**: `task-bd9c15f2-c472-4a65-8ff7-10748110781e` (deferred AppShell/AppSidebar changes)
- **Archived gap**: `home-page-ascii.md` in archive `2026-06-16-task-2bdd3fb5`