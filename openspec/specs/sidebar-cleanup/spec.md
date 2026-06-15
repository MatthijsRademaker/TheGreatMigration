# sidebar-cleanup Specification

## Purpose
TBD - created by archiving change task-94fecd0e-d08d-4e09-93ef-a558237c09f7. Update Purpose after archive.
## Requirements
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
