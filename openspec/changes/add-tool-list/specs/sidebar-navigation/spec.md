## MODIFIED Requirements

### Requirement: Sidebar SHALL render seven navigation items across two groups

The `AppSidebar` component SHALL render two `SidebarGroup` sections separated by a `SidebarSeparator`. The first group, labeled "Plan", SHALL contain navigation items for Dashboard, Tasks, Schedule, People, and Tools. The second group, labeled "Organization", SHALL contain navigation items for Rooms / Areas and Settings. Each navigation item SHALL use a `SidebarMenuButton` wrapping a `RouterLink` to its target route and displaying a lucide-vue icon alongside the item title.

#### Scenario: Plan group renders five primary navigation items

- **WHEN** the sidebar is rendered
- **THEN** a `SidebarGroup` with label "Plan" is present
- **AND** it contains `SidebarMenuButton` elements for Dashboard, Tasks, Schedule, People, and Tools
- **AND** each button wraps a `RouterLink` to its respective route (`/`, `/tasks`, `/calendar`, `/people`, `/tools`)

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
- **AND** the Tools item displays `WrenchIcon`
- **AND** the Rooms / Areas item displays `Building2Icon`
- **AND** the Settings item displays `SettingsIcon`
