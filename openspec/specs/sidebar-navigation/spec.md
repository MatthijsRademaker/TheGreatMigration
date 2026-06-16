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

#### Scenario: Settings route exists and renders placeholder content

- **WHEN** the application routes to `/settings`
- **THEN** the route resolves to a `SettingsView` component
- **AND** the rendered HTML contains the route's `meta.title` value
- **AND** the rendered HTML contains minimal placeholder content indicating the feature is forthcoming

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
