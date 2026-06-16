# sidebar-navigation Specification

## Purpose
TBD - created by archiving change task-dcb90e61-3e80-48c9-9a07-3813a882d623. Update Purpose after archive.
## Requirements
### Requirement: Sidebar SHALL render six navigation items across two groups

The `AppSidebar` component SHALL render two `SidebarGroup` sections separated by a `SidebarSeparator`. The first group, labeled "Plan", SHALL contain navigation items for Dashboard, Tasks, Schedule, and People. The second group, labeled "Organization", SHALL contain navigation items for Rooms / Areas and Settings. Each navigation item SHALL use a `SidebarMenuButton` wrapping a `RouterLink` to its target route and displaying a lucide-vue icon alongside the item title.

#### Scenario: Plan group renders four primary navigation items

- **WHEN** the sidebar is rendered
- **THEN** a `SidebarGroup` with label "Plan" is present
- **AND** it contains `SidebarMenuButton` elements for Dashboard, Tasks, Schedule, and People
- **AND** each button wraps a `RouterLink` to its respective route (`/`, `/tasks`, `/calendar`, `/people`)

#### Scenario: Organization group renders two secondary navigation items

- **WHEN** the sidebar is rendered
- **THEN** a `SidebarSeparator` element follows the Plan group
- **AND** a `SidebarGroup` with label "Organization" is present
- **AND** it contains `SidebarMenuButton` elements for Rooms / Areas and Settings
- **AND** each button wraps a `RouterLink` to its respective route (`/rooms`, `/settings`)

#### Scenario: Navigation items use appropriate lucide icons

- **WHEN** the sidebar is rendered
- **THEN** the Dashboard item displays `HomeIcon`
- **AND** the Tasks item displays `ClipboardListIcon`
- **AND** the Schedule item displays `CalendarDaysIcon`
- **AND** the People item displays `UsersRoundIcon`
- **AND** the Rooms / Areas item displays `Building2Icon`
- **AND** the Settings item displays `SettingsIcon`

### Requirement: Active navigation item SHALL use soft-green fill state

The active navigation item (matching the current route path) SHALL render with `is-active` state applied to its `SidebarMenuButton`. The shadcn-vue sidebar primitives SHALL apply the active styling using semantic sidebar tokens: `--sidebar-accent` (soft green `#E6F3EA`) background, `--sidebar-accent-foreground` (brand green `#1E6B3E`) icon and text, and `font-medium` weight. Inactive items SHALL remain neutral with subtle hover emphasis via `--sidebar-accent` on hover.

#### Scenario: Active nav item receives accent styling

- **WHEN** the current route path is `/tasks`
- **THEN** the Tasks `SidebarMenuButton` has `is-active` state
- **AND** the button applies `data-active:bg-sidebar-accent` styling
- **AND** the Dashboard, Schedule, People, Rooms / Areas, and Settings buttons do not have `is-active` state

#### Scenario: Inactive nav items have neutral appearance

- **WHEN** the current route path is `/`
- **THEN** the Dashboard `SidebarMenuButton` has `is-active` state
- **AND** the Tasks, Schedule, People, Rooms / Areas, and Settings buttons render without `data-active` styling
- **AND** the inactive buttons apply `hover:bg-sidebar-accent` styling on hover

### Requirement: Sidebar SHALL NOT display hardcoded badge counts

The `AppSidebar` component SHALL NOT render `SidebarMenuBadge` components with hardcoded numeric values. The `badge` property SHALL be absent from all navigation item definitions. A code comment SHALL document that badge counts should be re-added when real data subscriptions become available.

#### Scenario: No hardcoded badges are rendered

- **WHEN** the sidebar is rendered
- **THEN** no `SidebarMenuBadge` element contains the text "12" or "6"
- **AND** no `SidebarMenuBadge` element is present in the rendered output

### Requirement: Sidebar footer SHALL display project card and utility actions

The `SidebarFooter` section SHALL replace the current `GM` monogram chip with a project card `SidebarMenuButton` displaying the project name ("The Great Migration") and subtitle ("House move planner"). The footer SHALL also include two display-only utility action items: "Add note" (with `PlusIcon`) and "Help & Support" (with `CircleHelpIcon`). These utility actions SHALL be rendered as `SidebarMenuButton` elements without `RouterLink` wrappers — they are non-interactive placeholders for this change. A code comment SHALL document that interactivity for these items is deferred to a follow-up change.

#### Scenario: Footer contains project card with branding text

- **WHEN** the sidebar is rendered
- **THEN** the `SidebarFooter` contains a `SidebarMenuButton` displaying "The Great Migration" and "House move planner"
- **AND** no `GM` monogram chip is present

#### Scenario: Footer contains display-only utility actions

- **WHEN** the sidebar is rendered
- **THEN** the `SidebarFooter` contains a `SidebarMenuButton` with "Add note" label and `PlusIcon`
- **AND** the `SidebarFooter` contains a `SidebarMenuButton` with "Help & Support" label and `CircleHelpIcon`
- **AND** neither utility action button wraps a `RouterLink`

### Requirement: Every visible navigation item SHALL be backed by a route

Every navigation item rendered in the sidebar SHALL have a corresponding route entry in `frontend/src/app/routes.ts`. The routes for Rooms / Areas (`/rooms`) and Settings (`/settings`) SHALL use the existing lazy-loaded pattern. Each new route SHALL define `meta.title` and `meta.description` properties for the shared shell.

#### Scenario: Rooms / Areas route exists and renders placeholder content

- **WHEN** the application routes to `/rooms`
- **THEN** the route resolves to a `RoomsView` component
- **AND** the rendered HTML contains the route's `meta.title` value
- **AND** the rendered HTML contains minimal placeholder content indicating the feature is forthcoming

#### Scenario: Settings route exists and renders planning-window form

- **WHEN** the application routes to `/settings`
- **THEN** the route resolves to a `SettingsView` component
- **AND** the rendered HTML contains the route's `meta.title` value
- **AND** the rendered HTML contains the interactive planning-window card with date picker controls and Save/Reset buttons

### Requirement: New placeholder views SHALL render through the shared AppShell

The `RoomsView.vue` and `SettingsView.vue` components SHALL render their content within the existing shared `AppShell` layout. They SHALL NOT introduce their own shell chrome, sidebar, or header. `RoomsView` SHALL contain only a minimal placeholder heading and body text. `SettingsView` SHALL contain an interactive planning-window card as defined in the planning-window specification; data fetching and interactive controls SHALL be permitted for the planning-window card.

#### Scenario: RoomsView renders through AppShell with minimal content

- **WHEN** the RoomsView component is rendered
- **THEN** the content is wrapped by the `AppShell` layout component
- **AND** no additional sidebar, header, or navigation chrome is rendered inside the view
- **AND** the view contains at most a heading and one paragraph of placeholder text

#### Scenario: SettingsView renders through AppShell with interactive planning-window card

- **WHEN** the SettingsView component is rendered
- **THEN** the content is wrapped by the `AppShell` layout component
- **AND** no additional sidebar, header, or navigation chrome is rendered inside the view
- **AND** the view contains an interactive "Planning window" card with date picker controls
- **AND** the card loads the current planning window range from the backend composable

### Requirement: Shell render tests SHALL assert the expanded settings content contract

The route-render tests in `frontend/tests/app-routes-render.test.ts` SHALL assert the presence of the six navigation labels (Dashboard, Tasks, Schedule, People, Rooms / Areas, Settings) in the sidebar chrome. Additional route-render cases SHALL cover `/rooms` and `/settings`, asserting their route metadata and content strings. For `/settings`, the content assertion SHALL reflect the planning-window form elements rather than the previous "Feature coming soon" placeholder text. The project card content ("The Great Migration", "House move planner") and utility action labels ("Add note", "Help & Support") SHALL be asserted in the sidebar chrome test.

#### Scenario: Home route test asserts all six nav labels and footer content

- **WHEN** the home route (`/`) is rendered in tests
- **THEN** the rendered HTML contains "Rooms / Areas" and "Settings" as nav labels
- **AND** the rendered HTML contains "Add note" and "Help & Support" as footer labels

#### Scenario: New route cases exist for Rooms / Areas and Settings

- **WHEN** the test suite is executed
- **THEN** a test case exists for `/rooms` asserting its route metadata title, description, and placeholder content
- **AND** a test case exists for `/settings` asserting its route metadata title, description, and planning-window form content

### Requirement: Sidebar SHALL preserve existing variant and rail configuration

The `Sidebar` root component SHALL retain its existing `collapsible="icon"` and `variant="inset"` configuration. The `SidebarRail` SHALL remain present. The `SidebarHeader` brand block (`NotebookTabsIcon`, "The Great Migration", "House move planner") SHALL be unchanged.

#### Scenario: Sidebar configuration is preserved

- **WHEN** the sidebar is rendered
- **THEN** the `Sidebar` root element has `data-collapsible="icon"` and `data-variant="inset"` attributes
- **AND** the `SidebarRail` element is present
- **AND** the `SidebarHeader` contains the `NotebookTabsIcon` and the text "The Great Migration"

### Requirement: Sidebar SHALL reuse existing sidebar primitives and tokens exclusively

The implementation SHALL use only the existing shadcn-vue sidebar primitives exported from `@/shared/ui/sidebar`. No new shared-ui components or component directories SHALL be created. All styling SHALL use the existing semantic sidebar CSS tokens (`--sidebar-*`) from `frontend/src/app/styles.css`. No raw hex color values, inline styles with colors, or component-specific `<style>` blocks SHALL be introduced.

#### Scenario: No new components or raw styling introduced

- **WHEN** the implementation files are inspected
- **THEN** no new `.vue` files exist in `frontend/src/shared/ui/` beyond those already present before this change
- **AND** `AppSidebar.vue` uses only Tailwind token classes (e.g., `text-sidebar-foreground`, `bg-sidebar-accent`) for styling
- **AND** no raw hex values appear in the component template or style blocks

### Requirement: Mobile trigger SHALL be visible in the shared AppShell header on viewports ≤ 768px

On mobile-width viewports, the shared `AppShell.vue` header SHALL render a `SidebarTrigger` button visible via the `md:hidden` utility class. The trigger SHALL be placed left-aligned in the header, before the date range and title content. The trigger SHALL use the existing `SidebarTrigger` component which calls `toggleSidebar()` and correctly dispatches to `setOpenMobile()` when `isMobile` is true.

#### Scenario: Mobile trigger is visible and opens the sidebar Sheet

- **GIVEN** the viewport width is ≤ 768px
- **WHEN** the app shell renders
- **THEN** a `SidebarTrigger` element is present in the `AppShell` header
- **AND** the trigger is visually displayed (not hidden by `md:hidden`)
- **AND** clicking the trigger opens the mobile sidebar `Sheet`
- **AND** the Sheet contains the same nav groups (Plan, Organization) and footer content as the desktop sidebar

#### Scenario: Mobile trigger is hidden on desktop viewports

- **GIVEN** the viewport width is > 768px
- **WHEN** the app shell renders
- **THEN** the trigger element is present in the DOM but visually hidden via the `md:hidden` CSS class
- **AND** the desktop sidebar renders normally with its in-sidebar `SidebarTrigger` visible

### Requirement: In-sidebar SidebarTrigger SHALL be hidden on mobile viewports

The `SidebarTrigger` rendered inside `AppSidebar.vue`'s `SidebarHeader` SHALL be conditionally hidden when `isMobile` is true. This prevents a redundant close-toggle control from appearing inside the already-open mobile Sheet. On desktop viewports, the in-sidebar trigger SHALL remain visible to provide the collapse/expand control.

#### Scenario: In-sidebar trigger hidden on mobile

- **GIVEN** the viewport width is ≤ 768px
- **WHEN** the app sidebar renders
- **THEN** the `SidebarTrigger` inside `AppSidebar.vue`'s header is not rendered in the DOM

#### Scenario: In-sidebar trigger visible on desktop

- **GIVEN** the viewport width is > 768px
- **WHEN** the app sidebar renders
- **THEN** the `SidebarTrigger` inside `AppSidebar.vue`'s header is rendered and visible

### Requirement: Mobile sidebar SHALL close automatically on route change

When the mobile sidebar Sheet is open and the route changes (any navigation), the sidebar SHALL close by calling `setOpenMobile(false)`. The close-on-navigation behavior SHALL be implemented via a `watch` on `route.path` in `AppSidebar.vue`. The watcher SHALL guard against the initial mount by only closing when `openMobile` is already true.

#### Scenario: Mobile Sheet closes when navigating to a new route

- **GIVEN** the viewport width is ≤ 768px
- **AND** the mobile sidebar Sheet is open (openMobile is true)
- **WHEN** the route changes to a different path
- **THEN** the mobile sidebar Sheet closes (openMobile becomes false)
- **AND** the destination content is fully visible without the drawer overlay

#### Scenario: Mobile Sheet stays closed on initial page load

- **GIVEN** the viewport width is ≤ 768px
- **WHEN** the application first loads at a route
- **THEN** the mobile sidebar Sheet remains closed (openMobile is false)
- **AND** the route watcher does not fire `setOpenMobile(false)` on the initial mount

#### Scenario: Mobile Sheet closes via any navigation path

- **GIVEN** the viewport width is ≤ 768px
- **AND** the mobile sidebar Sheet is open
- **WHEN** the user navigates by clicking a RouterLink in the sidebar
- **THEN** the mobile sidebar Sheet closes
- **AND** the destination route content is visible

### Requirement: Desktop sidebar configuration SHALL be preserved after mobile support changes

All existing desktop sidebar behaviors SHALL remain unchanged after adding mobile support. This includes: `collapsible="icon"` and `variant="inset"` on the `Sidebar` root, the `SidebarRail` element, the `SidebarHeader` brand block, the two-group navigation structure, the footer content with project card and utility actions, the cookie-persisted open/closed state, and the keyboard shortcut (Ctrl+B / Cmd+B).

#### Scenario: Desktop sidebar configuration unchanged

- **GIVEN** the viewport width is > 768px
- **WHEN** the sidebar is rendered
- **THEN** the `Sidebar` root element has `data-collapsible="icon"` and `data-variant="inset"` attributes
- **AND** the `SidebarRail` element is present
- **AND** the in-sidebar `SidebarTrigger` is visible in the `SidebarHeader`
- **AND** the two-group navigation (Plan, Organization) renders with all six items
- **AND** the footer content (project card, utility actions) renders unchanged

### Requirement: Mobile sidebar behavior SHALL be covered by automated jsdom component tests

A new test file SHALL be created at `frontend/tests/sidebar-mobile.test.ts` using the jsdom environment. The tests SHALL mock the mobile breakpoint via `window.matchMedia`, mount sidebar components using `@vue/test-utils` mount, and assert: mobile trigger visibility, Sheet opening on trigger click, and Sheet closing on route change. Existing SSR tests in `frontend/tests/app-routes-render.test.ts` SHALL remain unchanged.

#### Scenario: Test asserts mobile trigger exists when matchMedia returns matches: true

- **GIVEN** `window.matchMedia` is mocked to return `{ matches: true }` for `(max-width: 768px)`
- **WHEN** the mobile sidebar test mounts the sidebar
- **THEN** the test asserts a `SidebarTrigger` element with `data-sidebar="trigger"` is present in the rendered output

#### Scenario: Test asserts Sheet closes on route navigation

- **GIVEN** `window.matchMedia` is mocked to return `{ matches: true }`
- **AND** the mobile Sheet is open
- **WHEN** a route navigation is triggered
- **THEN** the test asserts `openMobile` becomes false after navigation

### Requirement: Mobile sidebar SHALL reuse existing sidebar primitives and tokens exclusively

The mobile sidebar implementation SHALL use only the existing shadcn-vue sidebar primitives exported from `@/shared/ui/sidebar`. No new shared-ui components or component directories SHALL be created for mobile support. All mobile Sheet styling SHALL use the existing semantic sidebar CSS tokens (`--sidebar-*`) from `frontend/src/app/styles.css`. The `SidebarProvider.vue`, `Sidebar.vue`, and `SidebarTrigger.vue` primitives SHALL NOT be modified.

#### Scenario: No new components or primitives modified

- **WHEN** the mobile sidebar implementation is inspected
- **THEN** no new `.vue` files exist in `frontend/src/shared/ui/sidebar/` beyond those already present
- **AND** `SidebarProvider.vue` has zero changes
- **AND** `Sidebar.vue` has zero changes
- **AND** `SidebarTrigger.vue` has zero changes
