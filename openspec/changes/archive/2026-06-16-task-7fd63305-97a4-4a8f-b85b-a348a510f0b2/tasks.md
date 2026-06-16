## 1. Sidebar Trigger Relocation

- [x] 1.1 Remove `SidebarTrigger` import from `AppShell.vue`
- [x] 1.2 Remove `<SidebarTrigger class="-ml-1" />` and adjacent `<Separator>` from the AppShell header template
- [x] 1.3 Add `SidebarTrigger` import to `AppSidebar.vue` from `@/shared/ui/sidebar`
- [x] 1.4 Add a `SidebarMenuButton` wrapping `SidebarTrigger` as the last item in `AppSidebar.vue`'s `SidebarHeader`, using `variant="ghost"` for proper icon-collapsed mode sizing and alignment
- [x] 1.5 Verify sidebar collapse/expand works from the new trigger location in both expanded and icon-collapsed modes
- [x] 1.6 Verify existing `SidebarRail` hover behavior and keyboard shortcuts continue to work as secondary collapse affordances

## 2. Planning Window Range Formatter

- [x] 2.1 Add `formatPlanWindowRange(startDate: string, endDate: string, days: number): string` to `planWindow.ts` using UTC `Intl.DateTimeFormat` with `en-US` locale, producing format like `"5 Jul – 13 Aug 2026 · 40 days"`
- [x] 2.2 Export the new function from `planWindow.ts`
- [x] 2.3 Follow existing UTC pattern from `formatPlanDayLabel` for deterministic output regardless of runtime timezone

## 3. Composable Extension

- [x] 3.1 Import `formatPlanWindowRange` into `usePlanningWindow.ts`
- [x] 3.2 Add a `formattedRange` computed that derives `startDate`, `endDate`, and `days` from `query.data.value` (raw API response) and delegates to `formatPlanWindowRange()`
- [x] 3.3 Return `null` from `formattedRange` when `query.data.value` is `undefined` (loading/error states)
- [x] 3.4 Export `formattedRange` from the composable return

## 4. Header Visualizer in AppShell

- [x] 4.1 Import `usePlanningWindow` into `AppShell.vue`
- [x] 4.2 Destructure `formattedRange`, `isLoading`, `isError` from `usePlanningWindow()`
- [x] 4.3 Remove the `<Badge>Planning mode</Badge>` element from the header template
- [x] 4.4 Add an inline visualizer `<span>` or `<div>` that renders `formattedRange` when data is available, a compact loading skeleton during loading, and nothing (or a muted "—") on error
- [x] 4.5 Use a fixed-height container for the visualizer area to prevent layout shift between loading/error/data states

## 5. Test Updates

- [x] 5.1 Replace `expect(html).toContain('Planning mode')` assertion with assertions for visualizer content derived from the mocked planning-window data (e.g., `"Jul"`, `"2026"`, `"40 days"`)
- [x] 5.2 Verify all existing route-specific assertions continue to pass across all six routes
- [x] 5.3 If needed, add assertion for the absence of the old `SidebarTrigger` in the header area (implicitly verified by the visualizer assertion replacing the badge assertion)

## 6. Verification

- [x] 6.1 Run `scripts/precommit-run` and ensure all checks pass
- [x] 6.2 Confirm `vue-tsc --noEmit` passes with zero type errors
- [x] 6.3 Confirm `vitest run` passes all existing and new tests
- [ ] 6.4 Visually verify the sidebar trigger works from its new location in both expanded and icon-collapsed desktop modes
- [ ] 6.5 Visually verify the date-range visualizer renders correctly in the header across all routes
