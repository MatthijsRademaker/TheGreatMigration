## 1. Shell Layout Changes

- [ ] 1.1 On `SidebarInset` in `AppShell.vue`, change `overflow-hidden` to `overflow-y-auto` so the main content area becomes a scroll container
- [ ] 1.2 Add `sticky top-0 z-10` to the `<header>` element's class list in `AppShell.vue` so it sticks at the top when content scrolls
- [ ] 1.3 Verify the `background`/`backdrop-blur` and `border-b` on the header already provide the correct visual separation — no additional styling needed

## 2. Verification

- [ ] 2.1 Run the project's precommit checks: `scripts/precommit-run`
- [ ] 2.2 Visually confirm on mobile viewport (≤768px) that the header stays at the top when content overflows and scrolls
- [ ] 2.3 Visually confirm on desktop viewport that the header stays at the top within the `SidebarInset` with `variant="inset"` margins preserved
- [ ] 2.4 Confirm the mobile `SidebarTrigger` button inside the header remains clickable during scroll
- [ ] 2.5 Confirm the Calendar view's horizontal scrolling (`overflow-x-auto` on `DailySchedule`) works independently of the new vertical scroll
