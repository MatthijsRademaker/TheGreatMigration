## MODIFIED Requirements

### Requirement: Sidebar SHALL render six navigation items across two groups

The `AppSidebar` component SHALL render two `SidebarGroup` sections separated by a `SidebarSeparator`. The first group, labeled "Plan", SHALL contain navigation items for Dashboard, Tasks, Schedule, and People. The second group, labeled "Organization", SHALL contain navigation items for Rooms / Areas and Settings. Each navigation item SHALL use a `SidebarMenuButton` wrapping a native `<a>` element with an `href` attribute matching the target route. A `title` attribute SHALL provide a native browser tooltip for collapsed-mode identification. The `RouterLink` component SHALL NOT be used for sidebar navigation items. The `<a>` element SHALL handle single-click navigation via a `@click.prevent` handler that calls `router.push()`, while preserving right-click → open in new tab via the native `href`.

#### Scenario: Plan group renders four primary navigation items

- **WHEN** the sidebar is rendered
- **THEN** a `SidebarGroup` with label "Plan" is present
- **AND** it contains `SidebarMenuButton` elements for Dashboard, Tasks, Schedule, and People
- **AND** each button wraps a native `<a>` element with `href` set to the target route (`/`, `/tasks`, `/calendar`, `/people`)
- **AND** each `<a>` element has a `title` attribute matching the item label

#### Scenario: Organization group renders two secondary navigation items

- **WHEN** the sidebar is rendered
- **THEN** a `SidebarSeparator` element follows the Plan group
- **AND** a `SidebarGroup` with label "Organization" is present
- **AND** it contains `SidebarMenuButton` elements for Rooms / Areas and Settings
- **AND** each button wraps a native `<a>` element with `href` set to the target route (`/rooms`, `/settings`)

#### Scenario: Navigation items use appropriate lucide icons

- **WHEN** the sidebar is rendered
- **THEN** the Dashboard item displays `HomeIcon`
- **AND** the Tasks item displays `ClipboardListIcon`
- **AND** the Schedule item displays `CalendarDaysIcon`
- **AND** the People item displays `UsersRoundIcon`
- **AND** the Rooms / Areas item displays `Building2Icon`
- **AND** the Settings item displays `SettingsIcon`

#### Scenario: Navigation items are clickable via standard and alternative interactions

- **WHEN** a user left-clicks a sidebar navigation item
- **THEN** the browser navigates to the route via SPA routing without a full page reload
- **AND** when a user right-clicks the same item and selects "Open in new tab"
- **THEN** the browser opens the target route in a new tab
- **AND** when a user Ctrl+clicks the item
- **THEN** the browser opens the target route in a new tab

### Requirement: Sidebar footer SHALL display utility actions without duplicate branding

The `SidebarFooter` section SHALL contain two display-only utility action items: "Add note" (with `PlusIcon`) and "Help & Support" (with `CircleHelpIcon`). These utility actions SHALL be rendered as `SidebarMenuButton` elements without `<a>` wrappers — they are non-interactive placeholders for this change. A code comment SHALL document that interactivity for these items is deferred to a follow-up change. The `SidebarFooter` SHALL NOT contain the project brand card ("The Great Migration / House move planner") — that content SHALL appear only in the `SidebarHeader`.

#### Scenario: Footer contains display-only utility actions without branding

- **WHEN** the sidebar is rendered
- **THEN** the `SidebarFooter` contains a `SidebarMenuButton` with "Add note" label and `PlusIcon`
- **AND** the `SidebarFooter` contains a `SidebarMenuButton` with "Help & Support" label and `CircleHelpIcon`
- **AND** neither utility action button wraps an `<a>` or `RouterLink`
- **AND** the `SidebarFooter` does NOT contain the text "The Great Migration"
- **AND** the `SidebarFooter` does NOT contain the text "House move planner"

## ADDED Requirements

### Requirement: Collapsed-mode tooltips SHALL use native HTML `title` attributes

When the sidebar is in collapsed (icon-only) state, each navigation item SHALL display its label via the native HTML `title` attribute on the `<a>` element. The Reka UI `Tooltip` component SHALL NOT be used for nav item tooltips. This eliminates the event-forwarding chain that degrades click reliability. The `title` attribute SHALL be set to the navigation item's `title` property value.

#### Scenario: Native title attribute is present on all nav links

- **WHEN** the sidebar is rendered
- **THEN** the Dashboard `<a>` element has `title="Dashboard"`
- **AND** the Tasks `<a>` element has `title="Tasks"`
- **AND** the Schedule `<a>` element has `title="Schedule"`
- **AND** the People `<a>` element has `title="People"`
- **AND** the Rooms / Areas `<a>` element has `title="Rooms / Areas"`
- **AND** the Settings `<a>` element has `title="Settings"`

#### Scenario: No Reka UI Tooltip wraps nav items

- **WHEN** the sidebar is rendered
- **THEN** no `Tooltip` or `TooltipTrigger` component wraps navigation `SidebarMenuButton` elements
