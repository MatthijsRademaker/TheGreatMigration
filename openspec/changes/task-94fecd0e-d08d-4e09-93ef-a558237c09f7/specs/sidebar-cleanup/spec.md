# sidebar-cleanup Specification

## Purpose
Define the cleaned navigation structure of the application sidebar (AppSidebar) after removal of onboarding-era developer artifacts. The sidebar SHALL present only product-aligned navigation items matching `designs/design-system.png`: Dashboard, Tasks, Schedule, and People, with footer copy reflecting the product brand.

## MODIFIED Requirements

### Requirement: Sidebar contains only product navigation items

The AppSidebar SHALL render a single "Plan" navigation group containing Dashboard, Tasks, Schedule, and People entries.

#### Scenario: Primary navigation items are present

- **WHEN** the sidebar is rendered
- **THEN** the "Plan" `SidebarGroup` contains exactly four `SidebarMenuItem` entries
- **AND** the entries link to `/` (Dashboard), `/tasks` (Tasks), `/calendar` (Schedule), and `/people` (People)

#### Scenario: Schedule label matches design system

- **WHEN** the sidebar is rendered
- **THEN** the navigation entry for `/calendar` displays the label "Schedule"
- **AND** the label "Calendar" is not displayed in any sidebar navigation item

#### Scenario: Dashboard entry uses HomeIcon

- **WHEN** the sidebar is rendered
- **THEN** the Dashboard entry renders the `HomeIcon` component as its icon

### Requirement: Sidebar footer displays product-aligned copy

The AppSidebar footer SHALL display text reflecting "The Great Migration" product, not placeholder onboarding copy.

#### Scenario: Footer text is product-aligned

- **WHEN** the sidebar footer is rendered
- **THEN** the footer `SidebarMenuButton` displays "The Great Migration"
- **AND** the text "Local planning app" is not present anywhere in the sidebar

### Requirement: Route paths remain unchanged

Route paths SHALL NOT be changed as part of this cleanup. The `/calendar` path SHALL remain `/calendar`.

#### Scenario: Calendar route path is preserved

- **WHEN** the routes array in `src/app/routes.ts` is inspected
- **THEN** a route with `path: "/calendar"` and `name: "calendar"` exists
- **AND** the route meta and component import are unchanged

### Requirement: Active sidebar primitives are preserved

The sidebar SHALL continue to function with its existing dependencies on Sheet and Tooltip primitives.

#### Scenario: Sheet dependency is intact

- **WHEN** the sidebar is rendered on a mobile viewport
- **THEN** the sidebar opens as a Sheet overlay on mobile

#### Scenario: Tooltip dependency is intact

- **WHEN** the sidebar is collapsed to icon-only mode
- **THEN** hovering a collapsed sidebar navigation item displays a tooltip with the item label

## REMOVED Requirements

### Requirement: Developer section is absent

The AppSidebar SHALL NOT render a "Developer" navigation group or any showcase-related navigation items.

#### Scenario: No Developer group exists

- **WHEN** the sidebar is rendered
- **THEN** no `SidebarGroup` with label "Developer" is present
- **AND** no `SidebarMenuButton` links to `/showcase`
- **AND** the `ComponentIcon` is not imported in `AppSidebar.vue`

### Requirement: Move focus section is absent

The AppSidebar SHALL NOT render the onboarding-era "Move focus" prose content block.

#### Scenario: No Move focus prose block

- **WHEN** the sidebar is rendered
- **THEN** no `SidebarGroup` with label "Move focus" is present
- **AND** the onboarding prose text "Keep the plan practical: what needs doing, who is available, and which jobs need more hands" is not present in the rendered sidebar
