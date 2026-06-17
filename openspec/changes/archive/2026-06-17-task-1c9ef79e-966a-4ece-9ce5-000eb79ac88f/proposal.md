## Why

On mobile, the app header — which contains the page title, date-range visualizer, and sidebar toggle — scrolls away with content, forcing users to scroll back up to see the current page context or open the mobile sidebar. Making the header sticky on mobile keeps navigation context always visible and matches a well-established mobile UX pattern.

## What Changes

- The `AppShell.vue` header SHALL become sticky (`position: sticky; top: 0`) so it remains at the top of the viewport when content overflows on mobile.
- The `SidebarInset` scroll behavior SHALL change from `overflow-hidden` to allowing vertical scrolling (`overflow-y: auto`) so content taller than the viewport can actually scroll.
- These changes are layout-only — no new components, no backend changes, no new state management.

## Capabilities

### New Capabilities

*(None — this modifies the existing shell header behavior.)*

### Modified Capabilities

- `shared-shell-header`: Adding a requirement that the header SHALL be sticky at the top of the `SidebarInset` scroll container on all breakpoints (mobile and desktop), with content scrolling beneath it.

## Impact

- **`frontend/src/shared/layout/app-shell/AppShell.vue`** — the only file modified
- No new components, no new composables, no backend changes
- No view components (`HomeView`, `TasksView`, `CalendarView`, `PeopleView`, etc.) are touched
- Desktop sidebar inset layout is unaffected — the header sticks within its own scroll area
- Horizontal scroll containers inside views (e.g., `DailySchedule`) continue to work independently
