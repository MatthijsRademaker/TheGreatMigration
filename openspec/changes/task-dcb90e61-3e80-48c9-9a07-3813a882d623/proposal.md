## Why

The current app shell sidebar (`AppSidebar.vue`) exposes only four navigation items (Dashboard, Tasks, Schedule, People) in a single group with a generic `GM` footer chip. The design artifacts — `designs/components.png` section 2 and `designs/home-page.png` — show a fuller left-rail navigation system: six nav items across two groups (Dashboard, Tasks, Schedule, People, Rooms / Areas, Settings), a project card, and lower utility actions (Add note, Help & Support). A previous cleanup (archived change `task-94fecd0e`) removed onboarding-era sidebar clutter but deliberately stopped at four nav items, deferring the design-derived expansion to this change.

The sidebar infrastructure is ready: shadcn-vue sidebar primitives are in `frontend/src/shared/ui/sidebar/`, semantic sidebar tokens (`--sidebar`, `--sidebar-accent`, `--sidebar-accent-foreground`) are defined in `styles.css` matching the design-system-v2 state contract (pale-green background, soft-green active fill, neutral inactive with subtle hover), and `lucide-vue` provides icons. This task completes the navigation alignment by updating the sidebar composition, adding placeholder routes for the two new destinations, replacing the footer chip with design-derived content, and updating shell render tests.

## What Changes

- **AppSidebar.vue** — Reorganize navigation into two groups (`Plan`: Dashboard, Tasks, Schedule, People; `Organization`: Rooms / Areas, Settings) separated by `SidebarSeparator`. Replace the `GM` footer chip with a project card section (`SidebarFooter` containing `SidebarMenuButton`) showing the project name and subtitle from existing branding. Add lower utility actions (Add note, Help & Support) as non-interactive display-only items. Remove hardcoded badge counts on Tasks and People.
- **Route inventory** — Add two lightweight placeholder routes: `/rooms` → `RoomsView.vue` and `/settings` → `SettingsView.vue`. Each renders through the shared `AppShell` with route metadata and minimal placeholder content only (no feature implementation).
- **Placeholder views** — Create `frontend/src/rooms/RoomsView.vue` and `frontend/src/settings/SettingsView.vue` with minimal shell content.
- **Shell render tests** — Extend `frontend/tests/app-routes-render.test.ts` to assert the expanded nav labels (Rooms / Areas, Settings), project card content, and route-render cases for the two new placeholder routes.

## Impact

- **Affected code**: `frontend/src/shared/layout/app-sidebar/AppSidebar.vue`, `frontend/src/app/routes.ts`, `frontend/src/rooms/RoomsView.vue` (new), `frontend/src/settings/SettingsView.vue` (new), `frontend/tests/app-routes-render.test.ts`
- **Dependencies**: No new packages required. Uses existing shadcn-vue sidebar primitives, `lucide-vue` icons (`Building2Icon`, `SettingsIcon`, `CircleHelpIcon`, `PlusIcon`), and semantic CSS tokens.
- **Non-breaking**: Existing routes (`/`, `/tasks`, `/calendar`, `/people`), views, and the `@/shared/ui/sidebar/` barrel are unchanged. The sidebar variant (`inset`, `collapsible="icon"`) and existing `SidebarRail` are preserved.
- **Out of scope**: Full Rooms / Areas or Settings feature implementation, wiring badges or utility actions to real backend data, new logo assets, rebranding the brand block, dark mode-specific sidebar changes.