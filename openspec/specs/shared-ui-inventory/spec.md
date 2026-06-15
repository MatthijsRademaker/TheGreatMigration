# shared-ui-inventory Specification

## Purpose
TBD - created by archiving change task-94fecd0e-d08d-4e09-93ef-a558237c09f7. Update Purpose after archive.
## Requirements
### Requirement: Sheet and Tooltip primitives are retained

Sheet and Tooltip shared UI primitives SHALL remain intact with all their existing exports, as they are required dependencies of the sidebar component family.

#### Scenario: Sheet exports are preserved

- **WHEN** `src/shared/ui/sheet/` is inspected
- **THEN** all Sheet-related files and exports are unchanged from their pre-cleanup state

#### Scenario: Tooltip core exports are preserved

- **WHEN** `src/shared/ui/tooltip/` is inspected
- **THEN** `Tooltip.vue`, `TooltipContent.vue`, `TooltipTrigger.vue`, and their barrel exports are unchanged from their pre-cleanup state

### Requirement: No other shared UI primitives are affected

Primitives not explicitly targeted for removal (Badge, Button, Card, Input, Separator, and the retained sidebar/tooltip/sheet subsets) SHALL remain unchanged.

#### Scenario: Non-targeted primitives are untouched

- **WHEN** `src/shared/ui/badge/`, `src/shared/ui/button/`, `src/shared/ui/card/`, `src/shared/ui/input/`, and `src/shared/ui/separator/` are inspected
- **THEN** all files and exports are identical to their pre-cleanup state

### Requirement: Actively-used sidebar exports are preserved

When the unused subcomponents are removed, the actively-used sidebar exports SHALL remain in `src/shared/ui/sidebar/index.ts`.

#### Scenario: Actively-used sidebar exports remain

- **WHEN** `src/shared/ui/sidebar/index.ts` is inspected
- **THEN** exports for `Sidebar`, `SidebarContent`, `SidebarFooter`, `SidebarGroup`, `SidebarGroupContent`, `SidebarGroupLabel`, `SidebarHeader`, `SidebarInset`, `SidebarMenu`, `SidebarMenuBadge`, `SidebarMenuButton`, `SidebarMenuItem`, `SidebarProvider`, `SidebarRail`, `SidebarSeparator`, `SidebarTrigger`, and `useSidebar` are present
