# Homepage Header Timeline Design

## Layout

The homepage header is a **sticky top bar** that spans the full viewport width above the scrollable content area. It uses a **left-to-right horizontal flex layout** with `align-items: center` and a fixed height of `h-10` (40px).

The header has a subtle bottom border (`border-b`) for visual separation, a semi-transparent background (`bg-background/90`), and a backdrop blur effect (`backdrop-blur`) that makes the header translucent when content scrolls behind it. Horizontal padding is `px-4` (16px on each side), and gaps between sibling elements use `gap-2` (8px spacing).

## Control Grouping and Order

From left to right, the header contains these controls in sequence:

1. **Mobile sidebar trigger** (hidden on `md` breakpoint and above): a button with `data-sidebar="trigger"` that toggles the side navigation.
2. **Brand mark**: the text "The Great Migration" in a `text-sm font-semibold tracking-tight` span, followed by a decorative birds image (`/images/birds.png`) at `h-5 w-auto`.
3. **Flexible spacer** (implicit, via remaining flex space): pushes the right-side controls to the far end.
4. **Date range label**: a `text-sm text-muted-foreground` span showing the current 4-day page window, e.g. "2 Jul – 5 Jul, 2024". The format is: numeric day, abbreviated month (`short`) for the start date, then an en-dash, then numeric day, abbreviated month for the end date, a comma, and the year. Uses `en-US` locale and UTC timezone for deterministic rendering.
5. **Today button**: a `Button` with `variant="outline"` and `size="sm"`, labeled "Today". Always enabled. Navigates to the page containing the real current date within the planning window. If today falls outside the planning window entirely, navigates to page 1 (the start of the planning window).
6. **Previous chevron button**: a `Button` with `variant="ghost"` and `size="sm"`, containing a `ChevronLeft` icon from `lucide-vue` at `size-4`. Disabled when on page 1. Navigates to the previous 4-day page window.
7. **Next chevron button**: a `Button` with `variant="ghost"` and `size="sm"`, containing a `ChevronRight` icon from `lucide-vue` at `size-4`. Disabled when on the last page. Navigates to the next 4-day page window.

## Visibility

This timeline toolbar **only appears on the homepage route** (`/`). On all other routes (`/calendar`, `/people`, `/tasks`, `/rooms`, `/settings`), the header shows the existing static planning-window range ("5 Jul – 13 Aug 2026 · 40 days") instead.

## Interactive Behavior

- **Previous/Next**: Clicking the chevrons moves the visible 4-day window backward or forward by one page (4 days). Both the People availability grid and Daily schedule card update to show the same 4-day window.
- **Today**: Clicking "Today" navigates to the page whose 4-day window contains the current calendar date. If the current date is not within the planning window, navigates to page 1.
- **Disabled states**: The Previous chevron is disabled on page 1. The Next chevron is disabled on the last page. The Today button is never disabled.

## Loading State

While the planning window is loading, the toolbar area shows a pulsing skeleton placeholder (`h-3 w-32 animate-pulse rounded bg-muted`) instead of the range label and buttons.

## What Was Removed

The in-card pagination bar (Previous button, Next button, "Page X of Y" text, and date range label) that was previously rendered inside the Daily Schedule card header on the homepage is suppressed. The day columns and task cards remain, as does the "View by: Day" button.
