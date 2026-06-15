## ADDED Requirements

### Requirement: Showcase route is registered

The application SHALL provide a `/showcase` route accessible through vue-router that renders within the existing AppShell layout.

#### Scenario: Direct navigation to showcase

- **WHEN** a user navigates to `/showcase`
- **THEN** the AppShell header displays the route meta title `Component showcase`
- **AND** the AppShell header displays the route meta description `Developer preview of reusable UI components and patterns.`
- **AND** the showcase page content renders inside the RouterView within AppShell

#### Scenario: Showcase route is lazy-loaded

- **WHEN** the application bundle is built
- **THEN** the showcase view module is loaded only when the `/showcase` route is first visited

### Requirement: Showcase page renders a curated component catalog

The showcase page SHALL present the shared UI component inventory in visually grouped sections with live rendered examples, not text-only links.

#### Scenario: Catalog renders grouped sections

- **WHEN** the showcase page is loaded
- **THEN** the page displays component groups with section headers (Actions, Inputs, Feedback, Layout, Overlays)
- **AND** each group renders one or more live examples of the relevant shared/ui components

#### Scenario: Button variant and size matrix

- **WHEN** the Actions group is displayed
- **THEN** Button examples show all six variants (default, outline, secondary, ghost, destructive, link) in at least three sizes (default, sm, lg)
- **AND** each example is labeled with its variant and size combination

#### Scenario: Badge variant display

- **WHEN** the Feedback group is displayed
- **THEN** Badge examples show all six variants (default, secondary, destructive, outline, ghost, link)
- **AND** each example is labeled with its variant name

#### Scenario: Skeleton placeholder

- **WHEN** the Feedback group is displayed
- **THEN** a Skeleton example renders as a pulsing row
- **AND** a usage note explains the skeleton is a loading placeholder

#### Scenario: Input example

- **WHEN** the Inputs group is displayed
- **THEN** an Input example renders in its default state
- **AND** a usage note describes the input as a text entry field

#### Scenario: Card subcomponent composition

- **WHEN** the Layout group is displayed
- **THEN** Card subcomponents (Card, CardHeader, CardTitle, CardDescription, CardContent) render as a composed example card with title and description text

#### Scenario: Separator orientations

- **WHEN** the Layout group is displayed
- **THEN** Separator examples show both horizontal and vertical orientations
- **AND** each orientation is labeled

#### Scenario: Sheet interactive demo

- **WHEN** the Overlays group is displayed
- **THEN** a Sheet demo renders with a trigger button
- **AND** clicking the trigger opens a sheet panel with title, description, and close action

#### Scenario: Tooltip hover demo

- **WHEN** the Overlays group is displayed
- **THEN** a Tooltip demo renders with a hover trigger element
- **AND** hovering the trigger displays tooltip content

### Requirement: Sidebar component family is documented as layout-coupled

The showcase SHALL document the Sidebar component family as a layout-level component that cannot be rendered standalone in a catalog card, with a reference to the live AppSidebar usage.

#### Scenario: Sidebar family prose entry

- **WHEN** the showcase page is loaded
- **THEN** a section acknowledges the Sidebar component family exists in shared/ui
- **AND** the section explains the Sidebar requires a SidebarProvider context and is used at the layout composition level (AppSidebar, AppShell)
- **AND** the section references the live AppSidebar as the canonical usage example
- **AND** no standalone Sidebar is rendered inside the showcase

### Requirement: Showcase has discoverable sidebar navigation

The application sidebar SHALL include a secondary navigation section that provides access to the showcase route while maintaining visual separation from the primary planning navigation.

#### Scenario: Developer section in sidebar

- **WHEN** the sidebar is rendered
- **THEN** a secondary "Developer" navigation section appears below a visual separator after the primary "Plan" section
- **AND** the Developer section contains a single entry linking to `/showcase`
- **AND** the showcase entry uses an icon to distinguish it from primary navigation items

#### Scenario: Showcase entry active state

- **WHEN** the user is on the `/showcase` route
- **THEN** the showcase sidebar entry renders in its active state

### Requirement: Showcase follows existing project conventions
The showcase implementation SHALL follow the established project conventions for feature-folder structure, shared/ui imports, and theme token usage.

#### Scenario: Feature-folder convention

- **WHEN** the showcase source file is inspected
- **THEN** it is located at `src/showcase/ShowcaseView.vue`
- **AND** it follows the same `<script setup lang="ts">` + `<template>` structure as `src/home/HomeView.vue` and other route views

#### Scenario: Theme token usage

- **WHEN** the showcase page is rendered
- **THEN** all styling uses the existing Tailwind theme tokens (e.g., `bg-card`, `text-muted-foreground`, `border`) and shared/ui components (Card, Badge, etc.)
- **AND** no bespoke CSS or inline styles are introduced

#### Scenario: No Storybook dependency

- **WHEN** the project dependencies are inspected
- **THEN** Storybook is not present in `package.json`
- **AND** no Storybook configuration files exist
- **AND** no separate build step or documentation runtime is required
