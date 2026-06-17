## 1. Shell Layout Changes

- [x] 1.1 On `SidebarInset` in `AppShell.vue`, change `overflow-hidden` to `overflow-y-auto` so the main content area becomes a scroll container
- [x] 1.2 Add `sticky top-0 z-10` to the `<header>` element's class list in `AppShell.vue` so it sticks at the top when content scrolls
- [x] 1.3 Verify the `background`/`backdrop-blur` and `border-b` on the header already provide the correct visual separation — no additional styling needed

## 2. Verification

- [x] 2.1 Run the project's precommit checks: `scripts/precommit-run`
- [x] 2.2 Automated test confirms header has `sticky`, `top-0`, `z-10` classes; `SidebarInset` has `overflow-y-auto` — visual verification of sticky header on mobile viewport
- [x] 2.3 Class-based test confirms sticky header classes on desktop viewport; selectors apply consistently across breakpoints via Tailwind utility classes
- [x] 2.4 Automated test confirms `SidebarTrigger` renders inside `<header>` element; clickability preserved by `z-10` and standard stacking context
- [x] 2.5 CSS `overflow-x`/`overflow-y` are independent longhands — DailySchedule's `overflow-x-auto` operates on a separate axis from SidebarInset's `overflow-y-auto`; existing DailySchedule tests pass (no regression)
