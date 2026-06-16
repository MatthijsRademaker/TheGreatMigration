## Why

The current `AppShell.vue` header architecture is inverted: `SidebarTrigger` — a navigation control — is rendered outside the sidebar component it controls. The header also displays a static `Planning mode` badge that conveys no real information while the canonical planning-window data (`GET /api/planning-window`) is already available via `usePlanningWindow()`. The design system contract (`docs/design-system-v2.md`) specifies that the top bar should surface date-range controls, and the visual artifacts (`designs/home-page.png`, `designs/design-system.png`) show sidebar-first navigation composition without a freestanding content-header nav toggle.

This task fixes both architectural problems: move the primary sidebar toggle into sidebar-owned chrome, and replace the stale placeholder badge with a read-only planning-window visualizer that surfaces the move timeline context globally across all routes.

## What Changes

- **`AppShell.vue`** — Remove the `SidebarTrigger` import and its template usage from the content header. Remove the `<Badge>Planning mode</Badge>`. Add an inline date-range visualizer driven by `usePlanningWindow().formattedRange` with compact loading/error fallback states. Keep route title and description unchanged.
- **`AppSidebar.vue`** — Add a `SidebarTrigger` menu item as the last element in `SidebarHeader` after the branding row, using `SidebarMenuButton` with `variant="ghost"` for icon-collapsed and expanded mode compatibility. All existing sidebar chrome (branding, nav groups, footer, `SidebarRail`) is preserved.
- **`planWindow.ts`** — Export a new `formatPlanWindowRange(startDate: string, endDate: string, days: number): string` pure helper that formats the planning window as a compact range string (e.g., `"5 Jul – 13 Aug 2026 · 40 days"`) using UTC-based `Intl.DateTimeFormat`. Uses ISO date strings (SSR-safe), not `Date` objects.
- **`usePlanningWindow.ts`** — Add a `formattedRange` computed that delegates to `formatPlanWindowRange()`, deriving its inputs from the raw API response (`query.data.value.startDate`, `endDate`, `days`) for SSR safety.
- **`app-routes-render.test.ts`** — Replace the `expect(html).toContain('Planning mode')` assertion with assertions for the visualizer's formatted output. The test already mocks `/api/planning-window` with `startDate: '2026-07-05'`, `endDate: '2026-08-13'`, `days: 40`.

## Impact

- **Affected code**: `frontend/src/shared/layout/app-shell/AppShell.vue`, `frontend/src/shared/layout/app-sidebar/AppSidebar.vue`, `frontend/src/shared/lib/planWindow.ts`, `frontend/src/shared/composables/usePlanningWindow.ts`, `frontend/tests/app-routes-render.test.ts`
- **Dependencies**: No new packages. Uses existing shadcn-vue sidebar primitives (`SidebarTrigger`, `SidebarMenuButton`), existing `usePlanningWindow()` composable, and semantic CSS tokens.
- **Non-breaking**: Route metadata (title, description), sidebar navigation items, sidebar variant (`collapsible="icon"`, `variant="inset"`), `SidebarRail`, keyboard shortcut behavior, and all existing routes remain unchanged.
- **Out of scope**: Full aspirational toolbar (Today/arrows/notifications/profile), interactive range picker or date navigation, editable planning window, new backend endpoints, sidebar information architecture changes, dashboard panel or notes/help rewiring.
