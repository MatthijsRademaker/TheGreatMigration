## Context

The home dashboard (`/` route) renders four static KPI cards via `frontend/src/home/components/KpiCards.vue` — People available today, High priority tasks, Unassigned jobs, and Rooms completed (placeholder). Each card currently uses the project's shared `Card`, `CardHeader`, `CardContent`, `CardTitle`, and `CardDescription` sub-components with a layout that places a small icon circle (`size-10`) in the top-right corner of the header and stacks title, KPI value, and subtitle vertically within the CardHeader/CardContent flow.

The daily schedule component (`frontend/src/calendar/DailySchedule.vue`) uses a different visual language: each task card has a `border-l-4` left accent stripe colored by priority (`border-l-destructive`, `border-l-warning`, `border-l-success`), creating a strong left-side semantic signal.

This change brings KPI cards into visual alignment with that pattern while also restructuring the layout to make the icon a prominent full-height feature rather than a small corner badge.

**Key constraints:**

- All data contracts, queries, loading/error/empty state handling, and the rooms-completed placeholder treatment remain unchanged
- The shared `Card.vue` and its sub-components (`CardHeader`, `CardContent`, etc.) must not be modified
- No new global theme tokens, shared component directories, or raw color values may be introduced
- All CSS must use existing design system v2 semantic tokens from `frontend/src/app/styles.css`

## Goals / Non-Goals

**Goals:**

- Restructure each KPI card's inner template from `CardHeader`+`CardContent` sub-components to a custom two-column layout inside the base `Card`
- Add a full-height colored left accent column per card, using the same semantic token mapping as the current icon circle backgrounds
- Elevate icons from `size-10` circles to a larger size centered in the left accent column
- Stack title, KPI value, and subtitle vertically in the right content column
- Match the card's `rounded-lg` border-radius so the accent column's left corners are rounded
- Preserve the leaf decoration in the bottom-right of the content area
- Keep the `data-testid="kpi-placeholder-rooms-completed"` and code comment on the placeholder card
- Ensure loading, error, and empty states render within the same two-column structure
- Pass `scripts/precommit-run`

**Non-Goals:**

- Changing data contracts, API queries, or backend Go code
- Modifying the shared `Card.vue`, `CardHeader.vue`, `CardContent.vue`, `CardTitle.vue`, or `CardDescription.vue` components
- Changing the responsive grid (`md:grid-cols-2 xl:grid-cols-4`) or the component's placement in `HomeView.vue`
- Adding new icons, icon libraries, or design tokens
- Implementing real rooms-completed logic (remains a placeholder)
- Changing the DailySchedule component or any other dashboard surface

## Decisions

### 1. Override Card padding per-instance rather than modifying the shared Card component

The base `Card.vue` sets `py-panel` and `gap-panel` as part of its default classes. To achieve a flush left accent column, each KPI card instance will override these with `!py-0 !gap-0`. This avoids any changes to the shared primitive while keeping each card's custom layout self-contained.

- **Alternatives considered:**
  - *Negative margins on the icon column*: More fragile and could clip at different viewport sizes
  - *Modify Card.vue to accept a `flush` variant*: Cross-cutting change to a shared primitive with no other consumers
  - *Absolute positioning for the icon column*: Creates z-index and spacing complications with the leaf image

### 2. Two-column flex layout inside each Card instance

Each KPI card uses a single flex row as the Card's direct child:

```
Card (relative, overflow-hidden, !py-0 !gap-0)
├── flex-row
│   ├── Left accent column (w-[72px], self-stretch, semantic bg + text)
│   │   └── Icon (size-8 / 32px, centered)
│   └── Right content column (flex-1, p-panel)
│       ├── Title row (text-caption, text-muted-foreground)
│       ├── KPI value row (text-3xl, font-semibold)
│       └── Subtitle row (text-sm, text-muted-foreground)
└── img.leaf (absolute, bottom-0 right-0, pointer-events-none)
```

- **Width `w-[72px]` chosen** because it's wide enough to hold a `size-8` icon comfortably with padding, but not so wide that it dominates the card at small viewport sizes. On a 4-column grid (~250px per card), 72px is ~29% — acceptable for the accent role.
- **`self-stretch`** ensures the accent column spans the full card height, including when content has a tall loading/error state.
- The accent column uses `rounded-l-lg` to match the card's `rounded-lg`, though `overflow-hidden` on the Card already clips it.

### 3. Semantic color mapping per card type

| Card | Left column class | Icon color class |
|---|---|---|
| People available today | `bg-info-soft` | `text-info` |
| High priority tasks | `bg-destructive-soft` | `text-destructive` |
| Unassigned jobs | `bg-warning-soft` | `text-warning` |
| Rooms completed | `bg-success-soft` | `text-success` |

- Uses the same `*-soft` token pattern already used for the current icon circle backgrounds (`bg-destructive-soft`, `bg-warning-soft`), extended to `info-soft` and `success-soft` which are already defined in `styles.css`
- Soft variants provide a subtle background that doesn't compete with the KPI value for attention
- Full-strength `text-*` on the icon provides contrast against the soft background

### 4. Per-card config object drives colors and icons

The existing `backlogCards` pattern (a `computed<BacklogCardConfig[]>`) will be extended to cover all four cards, unifying the color/icon mapping in one place. Each card config includes:

- `label`, `value`, `description`, `icon component`, `accentClass` (for the left column `bg-* text-*`)

This eliminates the current duplication between the inline People card template and the `v-for` backlog cards.

### 5. Loading/error/empty states stay in the content column

The two-column layout wraps both the loading indicator ("Loading…") and error message ("Backend unavailable") in the content column's title/KPI-value area, using the same structure. This ensures the icon column always renders (providing consistent visual framing) while the content column handles state-specific text.

### 6. Leaf decoration position

The leaf image stays `absolute bottom-0 right-0` so it sits in the bottom-right corner of the content area, not overlapping the accent column. Since the accent column is on the left and is part of the same `relative` Card, the leaf is automatically positioned within the Card's bounds.

## Risks / Trade-offs

- **[Override churn] Overriding Card's `py-panel` with `!py-0`**: If the Card component's internal padding is ever changed, the `!` override will still win — but a future developer might not expect the override. Mitigation: the override is scoped to `KpiCards.vue` only, with a code comment explaining why.
- **[Responsive squeeze] 72px icon column at small viewports**: On very narrow screens (e.g., single-column grid on mobile, ~300px card width), a 72px column is ~24%. Acceptable. The icon shrinks via `size-8` (32px) which is large enough to remain recognizable.
- **[SSR test fragility] Assertions on `data-slot="card-header"` may break**: The current SSR test checks for `data-slot="card-title"` patterns indirectly through label text matching. Since we're removing `CardHeader`/`CardTitle` sub-components, label text assertions (`toContain("People available today")`) should still pass — they're pure text matching. No test change expected, but verify during implementation.
- **[Placeholder consistency] Rooms completed gets the same treatment**: Giving the placeholder card a `bg-success-soft` accent with the full layout could falsely imply it's data-driven. Mitigation: keep `data-testid="kpi-placeholder-rooms-completed"` and the existing code comment, plus keep `—` as the KPI value so it's visually distinct from live cards.
