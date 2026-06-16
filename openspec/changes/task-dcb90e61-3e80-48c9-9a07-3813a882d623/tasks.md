## 1. Sidebar Composition

- [ ] 1.1 Expand `primaryNavigation` array in `AppSidebar.vue` to include Rooms / Areas (`/rooms`, `Building2Icon`) and Settings (`/settings`, `SettingsIcon`)
- [ ] 1.2 Split navigation into two `SidebarGroup` sections: `Plan` (Dashboard, Tasks, Schedule, People) and `Organization` (Rooms / Areas, Settings) separated by `SidebarSeparator`
- [ ] 1.3 Add imports for `Building2Icon` and `SettingsIcon` from `@lucide/vue`
- [ ] 1.4 Remove hardcoded `badge` properties from Tasks (`'12'`) and People (`'6'`) entries; add comment documenting badges should be re-added when real data subscriptions exist
- [ ] 1.5 Replace the `SidebarFooter` `GM` chip with a project card (`SidebarMenuButton` showing project name and subtitle) and two display-only utility actions (Add note with `PlusIcon`, Help & Support with `CircleHelpIcon`); add comment noting interactivity is deferred
- [ ] 1.6 Add import for `SidebarSeparator`, `PlusIcon`, and `CircleHelpIcon`

## 2. Placeholder Routes and Views

- [ ] 2.1 Create `frontend/src/rooms/RoomsView.vue` with minimal placeholder content rendering through the shared `AppShell` with route metadata
- [ ] 2.2 Create `frontend/src/settings/SettingsView.vue` with minimal placeholder content rendering through the shared `AppShell` with route metadata
- [ ] 2.3 Add `/rooms` and `/settings` routes to `frontend/src/app/routes.ts` following the existing lazy-loaded pattern with `title` and `description` metadata

## 3. Test Coverage

- [ ] 3.1 Update the home route test in `frontend/tests/app-routes-render.test.ts` to assert `Rooms / Areas` and `Settings` nav labels in sidebar chrome
- [ ] 3.2 Add route-render cases for `/rooms` and `/settings` asserting route metadata and view content
- [ ] 3.3 Assert project card content (project name, subtitle) and utility action labels (Add note, Help & Support) in the sidebar chrome assertion
- [ ] 3.4 Remove or update any test assertion that checks for hardcoded badge counts if present

## 4. Verification

- [ ] 4.1 Run `scripts/precommit-run` and ensure all checks pass
- [ ] 4.2 Confirm `vue-tsc --noEmit` passes with zero type errors
- [ ] 4.3 Confirm `vitest run` passes all existing and new tests
- [ ] 4.4 Visually verify the sidebar renders with correct active/inactive states across the six nav items