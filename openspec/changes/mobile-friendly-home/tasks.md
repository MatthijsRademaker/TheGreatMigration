## 1. Hide MigrationJourney on Mobile

- [x] 1.1 In `frontend/src/home/HomeView.vue`, add `class="hidden sm:block"` to the `<MigrationJourney>` element

## 2. KPI Cards Compact Mobile Row

- [x] 2.1 In `KpiCards.vue`, add `expanded` ref (default `false`) scoped to mobile state
- [x] 2.2 Add compact summary row template (`sm:hidden`): icon + value for each of the five metrics displayed inline with `overflow-x-auto`
- [x] 2.3 Add expand/collapse toggle button (ChevronDown/ChevronUp) at the end of the summary row
- [x] 2.4 Wrap the full card grid in `v-show="expanded"` + `sm:!grid` so it always shows on sm+ regardless of the ref
- [x] 2.5 Verify queries are not re-fetched on expand (grid uses `v-show`, not `v-if`)

## 3. Sidebar Footer Image Mobile Guard

- [x] 3.1 In `AppSidebar.vue`, add `v-if="!isMobile"` to the `<img src="/images/trail-portrait.png">` element in `SidebarFooter`

## 4. Verification

- [x] 4.1 Run `scripts/check-dashboard` to confirm no type errors
- [x] 4.2 Run `scripts/test-go` and frontend unit tests to confirm no regressions
