# shared-ui-inventory Specification

## Purpose
Define the cleaned inventory of `src/shared/ui/` primitives after removing unused onboarding-generated subcomponents, dead-code primitives, and orphaned wrapper files. The inventory must retain all primitives with active product dependencies while removing files with zero consumers.

## REMOVED Requirements

### Requirement: Unused sidebar subcomponents are removed

The following sidebar subcomponent files SHALL NOT exist in the source tree and their exports SHALL NOT appear in `src/shared/ui/sidebar/index.ts`:
- `SidebarGroupAction.vue`
- `SidebarInput.vue`
- `SidebarMenuAction.vue`
- `SidebarMenuSkeleton.vue`
- `SidebarMenuSub.vue`
- `SidebarMenuSubButton.vue`
- `SidebarMenuSubItem.vue`

#### Scenario: Unused sidebar files are deleted

- **WHEN** the `src/shared/ui/sidebar/` directory is inspected
- **THEN** none of the seven listed `.vue` files exist

#### Scenario: Unused exports are removed from barrel

- **WHEN** `src/shared/ui/sidebar/index.ts` is inspected
- **THEN** the file does not contain export statements for any of the seven removed subcomponents

### Requirement: Skeleton primitive is removed

The `src/shared/ui/skeleton/` directory and all its contents SHALL NOT exist in the source tree.

#### Scenario: Skeleton directory is absent

- **WHEN** the `src/shared/ui/skeleton/` path is checked
- **THEN** no `Skeleton.vue` file exists
- **AND** no `index.ts` file exports Skeleton
- **AND** the `skeleton/` directory does not exist

#### Scenario: No stale Skeleton imports remain

- **WHEN** a project-wide search for `from '@/shared/ui/skeleton'` is performed
- **THEN** zero results are found

### Requirement: Local TooltipProvider wrapper is removed

The `src/shared/ui/tooltip/TooltipProvider.vue` file SHALL NOT exist and its export SHALL NOT appear in `src/shared/ui/tooltip/index.ts`.

#### Scenario: TooltipProvider wrapper is deleted

- **WHEN** `src/shared/ui/tooltip/TooltipProvider.vue` is checked
- **THEN** the file does not exist

#### Scenario: TooltipProvider barrel export is removed

- **WHEN** `src/shared/ui/tooltip/index.ts` is inspected
- **THEN** no line exports `TooltipProvider` from `"./TooltipProvider.vue"`

## ADDED Requirements

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
