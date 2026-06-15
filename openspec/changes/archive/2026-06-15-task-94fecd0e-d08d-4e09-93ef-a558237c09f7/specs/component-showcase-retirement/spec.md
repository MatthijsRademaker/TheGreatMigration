# component-showcase-retirement Specification

## Purpose
Document the retirement of the developer-facing component showcase surface (`/showcase` route, `ShowcaseView.vue`, Developer sidebar section) that was introduced during tech onboarding (change `task-c7c555d5-94d9-4d1a-a31d-dc42fe95651c`) and is superseded by this change. The showcase served as a developer catalog but has no counterpart in the current product design boards (`designs/design-system.png`, `designs/components.png`) and is not part of the product UI roadmap.

## REMOVED Requirements

### Requirement: Showcase route is registered

The application SHALL NOT provide a `/showcase` route.

#### Scenario: Showcase route returns 404

- **WHEN** a user navigates to `/showcase`
- **THEN** the router does not match any route for `/showcase`
- **AND** vue-router handles the unmatched path according to its default behavior

#### Scenario: Routes array excludes showcase

- **WHEN** the routes array in `src/app/routes.ts` is inspected
- **THEN** no route entry exists with `path: "/showcase"`

### Requirement: Showcase page renders a curated component catalog

The `src/showcase/` directory and `ShowcaseView.vue` file SHALL NOT exist in the source tree.

#### Scenario: Showcase source directory is absent

- **WHEN** the `src/showcase/` path is checked
- **THEN** no `ShowcaseView.vue` file exists at that location
- **AND** the `src/showcase/` directory is either absent or empty

### Requirement: Showcase has discoverable sidebar navigation

The AppSidebar SHALL NOT include a "Developer" navigation group or any link to a showcase route.

#### Scenario: Sidebar excludes Developer group

- **WHEN** the AppSidebar is rendered
- **THEN** no `SidebarGroup` with label "Developer" exists
- **AND** no `SidebarMenuButton` links to `/showcase`

## ADDED Requirements

### Requirement: Retirement is documented in OpenSpec

The retirement of the component showcase SHALL be traceable through OpenSpec change records.

#### Scenario: Original spec is archived

- **WHEN** the original component-showcase spec is referenced
- **THEN** it is superseded by this retirement spec
- **AND** the original implementation is preserved in `openspec/changes/archive/2026-06-15-task-c7c555d5-94d9-4d1a-a31d-dc42fe95651c/`

#### Scenario: No Storybook or alternative catalog is introduced

- **WHEN** project dependencies and configuration are inspected after this change
- **THEN** Storybook is not present in `package.json`
- **AND** no new developer catalog page or route replaces `/showcase`
