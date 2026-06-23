## Context

The home dashboard was built desktop-first. Three elements cause critical mobile UX problems:

1. **MigrationJourney** renders as a full-width card that consumes ~160px of vertical space. It's purely informational — no tappable actions — so it wastes prime screen real estate on mobile.
2. **KpiCards** renders five stacked full-height cards on a single-column mobile grid (~600px of scroll before reaching any actionable content).
3. **Sidebar trail-portrait image** (`trail-portrait.png`) in `SidebarFooter` has no height cap. The `group-data-[collapsible=icon]:hidden` guard that hides it in icon-collapsed desktop mode does NOT fire inside the mobile `Sheet` context (the Sheet renders outside the `group.peer` div), so the image always renders at its full natural dimensions on mobile, overflowing into and overlapping the navigation items above it.

## Goals / Non-Goals

**Goals:**
- MigrationJourney is invisible on mobile; desktop experience unchanged
- KpiCards shows a compact single-row summary on mobile by default; tapping a toggle reveals the full card grid; on tablet/desktop the full grid always shows with no toggle
- Sidebar footer image is completely suppressed on mobile; desktop experience unchanged
- No new dependencies, no new components, no API changes

**Non-Goals:**
- Redesigning the KPI cards themselves (layout, colors, values)
- Making the MigrationJourney interactive or collapsible on any viewport
- Responsive work on any screen beyond the home route and sidebar

## Decisions

### Decision 1 — MigrationJourney: `hidden sm:block` wrapper in HomeView

Add `class="hidden sm:block"` to the `<MigrationJourney>` element in `HomeView.vue`. The component renders zero DOM on mobile. Simplest possible approach; no changes to `MigrationJourney.vue` itself.

**Alternatives considered:**
- Collapsible card header on mobile — adds interactive complexity to a purely decorative element, not worth it.
- Compact chip summary line — still consumes space; doesn't solve the "too much screen" problem.

### Decision 2 — KpiCards: compact summary row on mobile with expand toggle

On `< sm` viewports `KpiCards.vue` renders a single horizontal pill row:

```
┌────────────────────────────────────────────────────────────────┐
│  🚩 3   👥 4/6   💼 7   ✓ —   🔧 3/5            [▾ Show KPIs] │
└────────────────────────────────────────────────────────────────┘
```

Each metric is an icon + value pair. The row scrolls horizontally if it overflows. A `ChevronDown/Up` toggle on the right expands/collapses the full card grid beneath. Default state is collapsed. On `sm:` and above, the component always renders the full card grid with no toggle in the DOM.

Implementation: add an `expanded` ref (default `false`) to `KpiCards.vue`. Render the compact row and toggle conditionally with `sm:hidden` / `hidden sm:flex`. The full card grid uses `v-show` (not `v-if`) so the queries don't re-fetch on expand.

**Alternatives considered:**
- Show 2 priority cards + "Show 3 more" (Option A) — still uses two full card heights on mobile, defeats the "too much space" goal.
- Accordion/disclosure card — heavier DOM structure for the same outcome.

### Decision 3 — Sidebar image: `v-if="!isMobile"` in AppSidebar

`isMobile` is already imported and destructured from `useSidebar()` in `AppSidebar.vue`. Add `v-if="!isMobile"` to the `<img>` element in `<SidebarFooter>`. Zero new code — one attribute addition. The image continues to render exactly as before on desktop (both expanded and icon-collapsed states).

**Alternatives considered:**
- `max-h-32 object-cover` — preserves a sliver of image but still wastes sidebar space and requires tuning.
- `hidden md:block` CSS class — works for the initial render but the Sheet context means the `md:` breakpoint class may not match perfectly with the JS `isMobile` boundary used by the sidebar primitive. Using the same `isMobile` ref keeps it consistent.

## Risks / Trade-offs

- **KPI expanded state resets on navigation** — `expanded` is local component state; navigating away and back will reset it to collapsed. This is acceptable (mobile users expect fresh page state) and avoids complexity of persisting toggle state.
- **Compact row truncation** — if values are long (e.g. `99 / 99`) the row may overflow on very narrow viewports. Horizontal scroll via `overflow-x-auto` handles this gracefully.
- **`v-show` vs `v-if` for card grid** — `v-show` keeps all five queries live even in collapsed state. This is intentional: queries are already in-flight from the page render; suppressing them with `v-if` would cause a re-fetch on expand.
