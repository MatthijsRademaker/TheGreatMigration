# shared-shell-header Specification

## Purpose
TBD - created by archiving change task-7fd63305-97a4-4a8f-b85b-a348a510f0b2. Update Purpose after archive.
## Requirements
### Requirement: SidebarTrigger SHALL be rendered inside AppSidebar SidebarHeader, not in AppShell header

The primary visible sidebar toggle SHALL be moved from `AppShell.vue`'s content header into `AppSidebar.vue`'s `SidebarHeader` section. The `SidebarTrigger` SHALL be rendered as the last `SidebarMenuButton` item in the `SidebarHeader` after the existing branding `RouterLink`, using `variant="ghost"` and the existing shadcn-vue sidebar primitives. The `AppShell.vue` header SHALL NOT import, render, or reference `SidebarTrigger`.

#### Scenario: SidebarTrigger is not rendered in AppShell header

- **WHEN** the application shell is rendered
- **THEN** the `AppShell.vue` header does NOT contain a `data-sidebar="trigger"` element
- **AND** the `AppShell.vue` component does NOT import `SidebarTrigger`

#### Scenario: SidebarTrigger is rendered in AppSidebar SidebarHeader

- **WHEN** the sidebar is rendered
- **THEN** the `SidebarHeader` contains a `SidebarMenuButton` wrapping a `SidebarTrigger` component
- **AND** the trigger button is positioned as the last element in the `SidebarHeader` after the branding row
- **AND** the trigger button uses `data-sidebar="trigger"` attribute

#### Scenario: SidebarTrigger works in expanded and icon-collapsed modes

- **WHEN** the sidebar is in expanded mode
- **THEN** clicking the `SidebarTrigger` in the `SidebarHeader` collapses the sidebar to icon mode
- **AND** the `SidebarMenuButton` wrapping the trigger renders with full width styling
- **WHEN** the sidebar is in icon-collapsed mode
- **THEN** clicking the `SidebarTrigger` in the `SidebarHeader` expands the sidebar
- **AND** the `SidebarMenuButton` wrapping the trigger renders with icon-only styling (`size-8`)

#### Scenario: Secondary collapse affordances remain functional

- **WHEN** the sidebar is rendered
- **THEN** the `SidebarRail` element is present and provides hover-based collapse/expand
- **AND** the keyboard shortcut for sidebar toggle continues to work

### Requirement: SidebarTrigger in AppSidebar SHALL use SidebarMenuButton with ghost variant

The `SidebarTrigger` inside `AppSidebar.vue`'s `SidebarHeader` SHALL be rendered inside a `SidebarMenuButton` component using `variant="ghost"` for visual consistency with the sidebar's design language. The trigger button SHALL use the existing `PanelLeftIcon` from lucide-vue and a `sr-only` label reading "Toggle Sidebar".

#### Scenario: Trigger button uses ghost styling

- **WHEN** the sidebar is rendered with the trigger in `SidebarHeader`
- **THEN** the wrapping `SidebarMenuButton` has `data-variant="ghost"` or equivalent ghost styling
- **AND** the button displays the `PanelLeftIcon` icon

### Requirement: planWindow.ts SHALL export a pure formatPlanWindowRange helper

The `planWindow.ts` module SHALL export a new function `formatPlanWindowRange(startDate: string, endDate: string, days: number): string`. The function SHALL accept ISO 8601 date strings (`startDate`, `endDate`) and an inclusive day count (`days`). It SHALL use `Intl.DateTimeFormat` with `en-US` locale and `UTC` timezone to produce deterministic labels regardless of runtime timezone, following the existing `formatPlanDayLabel` pattern. The function SHALL return a compact compound string in the format `"5 Jul – 13 Aug 2026 · 40 days"` using an en-dash between dates and a middle-dot before the day count.

#### Scenario: Formatter produces compact range string from valid inputs

- **GIVEN** `startDate` is `"2026-07-05"`, `endDate` is `"2026-08-13"`, and `days` is `40`
- **WHEN** `formatPlanWindowRange("2026-07-05", "2026-08-13", 40)` is called
- **THEN** the returned string is `"5 Jul – 13 Aug 2026 · 40 days"`

#### Scenario: Formatter uses UTC for deterministic output

- **GIVEN** any valid ISO date strings
- **WHEN** `formatPlanWindowRange()` is called from different runtime timezones
- **THEN** the returned string is identical regardless of the server or browser timezone

#### Scenario: Formatter accepts only string date inputs

- **WHEN** `formatPlanWindowRange()` is inspected
- **THEN** the `startDate` parameter type is `string`, not `Date`
- **AND** the `endDate` parameter type is `string`, not `Date`

### Requirement: usePlanningWindow SHALL expose a formattedRange computed

The `usePlanningWindow()` composable SHALL expose a new `formattedRange` computed property. This computed SHALL derive `startDate`, `endDate`, and `days` from `query.data.value` (the raw API response body) and delegate to `formatPlanWindowRange()`. It SHALL return `null` when `query.data.value` is `undefined` (loading or error state). It SHALL NOT construct `Date` objects or reference `planWindowDays` entries.

#### Scenario: formattedRange returns formatted string when data is available

- **GIVEN** `query.data.value` is `{ startDate: "2026-07-05", endDate: "2026-08-13", days: 40 }`
- **WHEN** `formattedRange.value` is accessed
- **THEN** the value is `"5 Jul – 13 Aug 2026 · 40 days"`

#### Scenario: formattedRange returns null when data is not available

- **GIVEN** `query.data.value` is `undefined`
- **WHEN** `formattedRange.value` is accessed
- **THEN** the value is `null`

#### Scenario: formattedRange is SSR-safe

- **WHEN** `formattedRange` is evaluated during SSR `renderToString`
- **THEN** no `Date` serialization errors occur
- **AND** the returned value is a plain string or `null`

### Requirement: AppShell header SHALL render a planning-window date-range visualizer

The `AppShell.vue` header template SHALL replace the `<Badge>Planning mode</Badge>` with an inline visualizer that consumes `usePlanningWindow().formattedRange`. The visualizer SHALL render the formatted range string when data is available, a compact loading skeleton placeholder during loading, and nothing (or a muted fallback) on error. The visualizer area SHALL use a fixed-height container to prevent layout shift between states. The route title and description SHALL remain unchanged in the header.

#### Scenario: Visualizer renders formatted range when data is loaded

- **GIVEN** the planning-window query has resolved successfully
- **WHEN** the application shell is rendered
- **THEN** the header contains the formatted range string (e.g., `"5 Jul – 13 Aug 2026 · 40 days"`)
- **AND** the header does NOT contain the text "Planning mode"

#### Scenario: Visualizer shows loading skeleton when query is pending

- **GIVEN** the planning-window query is in a loading state
- **WHEN** the application shell is rendered
- **THEN** the header visualizer area shows a compact loading skeleton or placeholder
- **AND** the header layout does not shift or break

#### Scenario: Visualizer shows fallback on error

- **GIVEN** the planning-window query has failed
- **WHEN** the application shell is rendered
- **THEN** the header visualizer area shows either nothing or a muted fallback indicator
- **AND** the header layout does not shift or break

#### Scenario: Route title and description are preserved

- **WHEN** the application shell is rendered for any route
- **THEN** the header contains the route's `meta.title` value as an `<h1>`
- **AND** the header contains the route's `meta.description` value as a `<p>`

### Requirement: Shell SHALL preserve existing sidebar variant, rail, and branding configuration

The `Sidebar` root component SHALL retain its existing `collapsible="icon"` and `variant="inset"` configuration. The `SidebarRail` SHALL remain present. The `SidebarHeader` brand block (`NotebookTabsIcon`, "The Great Migration", "House move planner") SHALL be unchanged aside from the addition of the `SidebarTrigger` menu item after the branding row.

#### Scenario: Sidebar configuration is preserved

- **WHEN** the sidebar is rendered
- **THEN** the `Sidebar` root element has `data-collapsible="icon"` attribute
- **AND** the `SidebarRail` element is present
- **AND** the `SidebarHeader` contains the `NotebookTabsIcon` and the text "The Great Migration"
- **AND** the `SidebarHeader` contains the text "House move planner"

### Requirement: Implementation SHALL reuse existing shell primitives and tokens exclusively

The implementation SHALL use only the existing shadcn-vue sidebar primitives exported from `@/shared/ui/sidebar`. No new shared-ui components or component directories SHALL be created. All styling SHALL use the existing Tailwind utility classes and semantic CSS tokens. No raw hex color values, inline styles with colors, or component-specific `<style>` blocks SHALL be introduced for the new visualizer or trigger placement.

#### Scenario: No new components or raw styling introduced

- **WHEN** the implementation files are inspected
- **THEN** no new `.vue` files exist in `frontend/src/shared/ui/` beyond those already present before this change
- **AND** `AppShell.vue` uses only Tailwind utility classes for the visualizer
- **AND** `AppSidebar.vue` uses only existing sidebar primitives for the relocated trigger
- **AND** no raw hex color values appear in the new or modified template code

### Requirement: SSR shell tests SHALL assert the new header and sidebar chrome contract

The SSR route-render tests in `frontend/tests/app-routes-render.test.ts` SHALL remove the `expect(html).toContain('Planning mode')` assertion. The tests SHALL instead assert that the rendered output contains date-related text from the visualizer output derived from the mocked planning-window data (`startDate: "2026-07-05"`, `endDate: "2026-08-13"`, `days: 40`). Existing assertions for route titles, descriptions, navigation labels, and footer content SHALL continue to pass.

#### Scenario: Home route test asserts visualizer content instead of Planning mode badge

- **WHEN** the home route (`/`) is rendered in the SSR test
- **THEN** the rendered HTML does NOT contain the text "Planning mode"
- **AND** the rendered HTML contains date-related visualizer text (e.g., "Jul", "2026", "40 days")

#### Scenario: All route-specific assertions continue to pass

- **WHEN** the full test suite is executed
- **THEN** all route-render assertions for metadata titles, descriptions, view content, and sidebar nav labels pass
- **AND** the `/calendar` route assertions (tasks, dates, priority variants) continue to pass
- **AND** the `/` route KPI assertions continue to pass

### Requirement: Change SHALL be frontend-only with no backend modifications

This change SHALL NOT introduce, modify, or remove any backend API endpoints, types, or server logic. It SHALL consume only the existing `GET /api/planning-window` endpoint and the existing `PlanningWindowBody` response type. The header visualizer SHALL be read-only display chrome with no mutation capability.

#### Scenario: No backend files are modified

- **WHEN** the implementation is complete
- **THEN** no files outside the `frontend/` directory are modified
- **AND** the `GET /api/planning-window` endpoint implementation is unchanged
