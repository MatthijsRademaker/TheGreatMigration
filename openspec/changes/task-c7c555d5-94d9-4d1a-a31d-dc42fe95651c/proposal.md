## Why

The project has a growing library of reusable UI components under `src/shared/ui/` (badge, button, card, input, separator, sheet, skeleton, sidebar, tooltip) with no discovery surface for developers. There is no single place to review what components exist, what variants they support, or when to use them — developers must inspect source files or cross-reference existing views to understand the design-system inventory. The task backlog requests a Storybook-like showcase route to fill this gap without introducing Storybook's tooling, infrastructure, or build dependencies.

## What Changes

**[New route — /showcase]**
- From: No showcase route exists. The router defines only `/`, `/tasks`, `/calendar`, and `/people`.
- To: A new `/showcase` route registered in `src/app/router.ts` with lazy-loaded `src/showcase/ShowcaseView.vue`, route meta title `Component showcase`, and description `Developer preview of reusable UI components and patterns.`
- Reason: Provide a dedicated discovery surface for the shared UI component library.
- Impact: Non-breaking. Adds one route alongside existing planning routes.

**[New page — ShowcaseView]**
- From: No showcase page exists.
- To: `src/showcase/ShowcaseView.vue` following the feature-folder convention (parallel to `src/home/HomeView.vue`, `src/tasks/TasksView.vue`). The page renders a manually curated component catalog grouped into sections: Actions (Button with all variant × size states), Inputs (Input), Feedback (Badge, Skeleton), Layout (Card subcomponents, Separator), and Overlays (Sheet with trigger demo, Tooltip with hover demo). Each group shows live rendered examples with state labels and brief usage guidance. The Sidebar component family is documented as layout-coupled with a prose note (not rendered standalone). No code snippets in v1.
- Reason: Centralize reusable UI inventory in one navigable page.
- Impact: Non-breaking. New file in new directory.

**[Modified sidebar — AppSidebar.vue]**
- From: Single primary navigation group under "Plan" label, followed by a "Move focus" secondary section.
- To: Appends a new "Developer" secondary navigation section below the existing separator, containing one entry for the showcase route with `ComponentIcon` from `lucide-vue`.
- Reason: Provide discoverable access to the showcase while visually separating it from the primary move-planning workflow (satisfies acceptance criteria #5 and #7).
- Impact: Non-breaking. One additional entry in the sidebar, visually separated from the "Plan" group.

## Capabilities

### New Capabilities
- `component-showcase`: A dedicated route and page that catalogs and renders the shared UI component library for developer discovery.

### Modified Capabilities
- (none — this is a new capability with no prior spec)

## Impact

**Affected code**: `src/app/router.ts` (new route), `src/showcase/ShowcaseView.vue` (new file), `src/shared/layout/app-sidebar/AppSidebar.vue` (new secondary navigation group).

**Dependencies**: No new external dependencies. Uses existing `vue-router`, `lucide-vue` (already in project), and existing `shared/ui` components.

**Breaking changes**: None. All existing routes and views remain untouched.

**Risk**: Manual catalog may drift if new components are added to `shared/ui` without a corresponding showcase entry. Mitigation: the catalog is a static object within ShowcaseView — adding a component requires a one-line entry alongside the component's export, documented in the change as an explicit team convention.