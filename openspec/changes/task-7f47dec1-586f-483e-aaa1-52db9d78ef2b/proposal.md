## Why

The KPI summary cards on the home dashboard currently use a generic layout — icons sit in small circular badges in the top-right corner, and content (title, KPI value, subtitle) flows in a single column inside CardHeader/CardContent sub-components. The daily schedule task cards use a distinctive left-accent treatment (`border-l-4` with semantic colors) that creates stronger visual hierarchy and makes priority signals immediately scannable. Applying a similar left-accent pattern to KPI cards — combined with a restructured layout that elevates icons from small badges to a full-height visual column — improves dashboard consistency, makes KPI values more scannable, and gives each card a clearer visual identity tied to its semantic meaning (people, tasks, rooms).

## What Changes

- **Restructure each KPI card's inner layout** from `CardHeader` + `CardContent` with a top-right icon circle to a two-column grid: a left accent column (full-height icon area) and a right content stack (title row, KPI value row, subtitle row).
- **Add a colored left accent bar** to each KPI card, using the same semantic token mapping as the current icon backgrounds: info for people availability, destructive for high-priority tasks, warning for unassigned jobs, success for the rooms-completed placeholder.
- **Elevate icons from `size-10` (40px) circular badges** to a full-height left panel with a larger icon, using the accent color as the panel background and white/foreground for the icon.
- **Remove the `CardHeader`/`CardContent` split** in favor of a custom layout inside the base `Card` shell. The Card component itself remains unchanged — only the per-card inner template is restructured.
- **Keep the leaf decoration** (`/images/leaf.png`) in the bottom-right of the content column, not overlapping the accent column.
- **Preserve all data contracts, loading/error/empty states, query wiring, and the `rooms-completed` placeholder treatment** — no backend or data-layer changes.
- **Update the SSR route-render test** if needed to account for structural DOM changes (e.g., if the old `data-slot="card-header"` assertions break).

## Capabilities

### New Capabilities

*(none — this is a visual refresh of existing capability surfaces)*

### Modified Capabilities

- `kpi-summary-cards`: The card layout requirement changes from `CardHeader`+`CardContent` primitives with top-right icon circles to a two-column left-accent layout with full-height icon panel. The styling requirement is amended to use a colored left accent bar per semantic card type instead of icon circles with `bg-secondary`/`bg-destructive-soft`/`bg-warning-soft` backgrounds.

## Impact

- **Frontend only**: `frontend/src/home/components/KpiCards.vue` — complete internal template restructure
- **No backend changes**: Data contracts, queries, API endpoints unchanged
- **No shared component changes**: `Card.vue`, `CardHeader.vue`, `CardContent.vue`, etc. remain untouched — only the consumer (`KpiCards.vue`) changes
- **No new dependencies**: Uses existing lucide-vue icons, semantic CSS tokens, and design-system primitives
- **Test update**: `frontend/tests/app-routes-render.test.ts` may need assertion updates if `data-slot="card-header"` patterns are removed from SSR output
