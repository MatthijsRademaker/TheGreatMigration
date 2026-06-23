## Why

The home screen dashboard is unusable on mobile: the Move-Day Readiness Journey card, five stacked KPI cards, and an oversized decorative image in the sidebar drawer all consume excessive vertical space and obstruct navigation. These three elements need targeted mobile treatment to make the core dashboard accessible on small screens.

## What Changes

- **Hide MigrationJourney on mobile** — the progress card is informational/decorative and not actionable on mobile; hidden below `sm` breakpoint via Tailwind responsive class.
- **KPI cards collapse to compact summary row on mobile** — all five metrics displayed inline as a single scrollable pill/chip row (`🚩3  👥4/6  💼7  ✓—  🔧3/5`) with a toggle to expand to full card grid; on `sm:` and above the full grid is always shown with no toggle.
- **Sidebar trail-portrait image hidden on mobile** — the decorative `trail-portrait.png` in the sidebar footer overlaps navigation items on mobile (Sheet) because its height is unconstrained and the `group-data-[collapsible=icon]:hidden` trick doesn't fire inside the Sheet context; guarded with `v-if="!isMobile"` using the already-available `isMobile` ref from `useSidebar()`.

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities
- `migration-journey`: now hidden below `sm` breakpoint; spec adds mobile visibility rule
- `kpi-summary-cards`: adds mobile-collapsed compact row behaviour with expand toggle; spec adds responsive display contract
- `sidebar-navigation`: sidebar footer decorative image is suppressed on mobile; spec adds the mobile guard rule

## Impact

- `frontend/src/home/HomeView.vue` — add responsive class to `<MigrationJourney>`
- `frontend/src/home/components/KpiCards.vue` — add collapse state, compact row render on mobile, expand toggle
- `frontend/src/shared/layout/app-sidebar/AppSidebar.vue` — add `v-if="!isMobile"` to trail-portrait img
- No API changes, no new dependencies, no breaking changes
