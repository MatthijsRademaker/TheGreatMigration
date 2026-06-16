## ADDED Requirements

### Requirement: Mobile trigger SHALL be visible in the shared AppShell header on viewports ≤ 768px

On mobile-width viewports, the shared `AppShell.vue` header SHALL render a `SidebarTrigger` button visible via the `md:hidden` utility class. The trigger SHALL be placed left-aligned in the header, before the date range and title content. The trigger SHALL use the existing `SidebarTrigger` component which calls `toggleSidebar()` and correctly dispatches to `setOpenMobile()` when `isMobile` is true.

#### Scenario: Mobile trigger is visible and opens the sidebar Sheet

- **GIVEN** the viewport width is ≤ 768px
- **WHEN** the app shell renders
- **THEN** a `SidebarTrigger` element is present in the `AppShell` header
- **AND** the trigger is visually displayed (not hidden by `md:hidden`)
- **AND** clicking the trigger opens the mobile sidebar `Sheet`
- **AND** the Sheet contains the same nav groups (Plan, Organization) and footer content as the desktop sidebar

#### Scenario: Mobile trigger is hidden on desktop viewports

- **GIVEN** the viewport width is > 768px
- **WHEN** the app shell renders
- **THEN** the trigger element is present in the DOM but visually hidden via the `md:hidden` CSS class
- **AND** the desktop sidebar renders normally with its in-sidebar `SidebarTrigger` visible

### Requirement: In-sidebar SidebarTrigger SHALL be hidden on mobile viewports

The `SidebarTrigger` rendered inside `AppSidebar.vue`'s `SidebarHeader` SHALL be conditionally hidden when `isMobile` is true. This prevents a redundant close-toggle control from appearing inside the already-open mobile Sheet. On desktop viewports, the in-sidebar trigger SHALL remain visible to provide the collapse/expand control.

#### Scenario: In-sidebar trigger hidden on mobile

- **GIVEN** the viewport width is ≤ 768px
- **WHEN** the app sidebar renders
- **THEN** the `SidebarTrigger` inside `AppSidebar.vue`'s header is not rendered in the DOM

#### Scenario: In-sidebar trigger visible on desktop

- **GIVEN** the viewport width is > 768px
- **WHEN** the app sidebar renders
- **THEN** the `SidebarTrigger` inside `AppSidebar.vue`'s header is rendered and visible

### Requirement: Mobile sidebar SHALL close automatically on route change

When the mobile sidebar Sheet is open and the route changes (any navigation), the sidebar SHALL close by calling `setOpenMobile(false)`. The close-on-navigation behavior SHALL be implemented via a `router.afterEach` hook in `AppSidebar.vue`. The hook SHALL close the Sheet only when `isMobile` and `openMobile` are both true.

#### Scenario: Mobile Sheet closes when navigating to a new route

- **GIVEN** the viewport width is ≤ 768px
- **AND** the mobile sidebar Sheet is open (openMobile is true)
- **WHEN** the route changes to a different path
- **THEN** the mobile sidebar Sheet closes (openMobile becomes false)
- **AND** the destination content is fully visible without the drawer overlay

#### Scenario: Mobile Sheet stays closed on initial page load

- **GIVEN** the viewport width is ≤ 768px
- **WHEN** the application first loads at a route
- **THEN** the mobile sidebar Sheet remains closed (openMobile is false)
- **AND** the route hook does not fire `setOpenMobile(false)` on the initial mount

#### Scenario: Mobile Sheet closes via any navigation path

- **GIVEN** the viewport width is ≤ 768px
- **AND** the mobile sidebar Sheet is open
- **WHEN** the user navigates by clicking a RouterLink in the sidebar
- **THEN** the mobile sidebar Sheet closes
- **AND** the destination route content is visible

### Requirement: Desktop sidebar configuration SHALL be preserved after mobile support changes

All existing desktop sidebar behaviors SHALL remain unchanged after adding mobile support. This includes: `collapsible="icon"` and `variant="inset"` on the `Sidebar` root, the `SidebarRail` element, the `SidebarHeader` brand block, the two-group navigation structure, the footer content with project card and utility actions, the cookie-persisted open/closed state, and the keyboard shortcut (Ctrl+B / Cmd+B).

#### Scenario: Desktop sidebar configuration unchanged

- **GIVEN** the viewport width is > 768px
- **WHEN** the sidebar is rendered
- **THEN** the `Sidebar` root element has `data-collapsible="icon"` and `data-variant="inset"` attributes
- **AND** the `SidebarRail` element is present
- **AND** the in-sidebar `SidebarTrigger` is visible in the `SidebarHeader`
- **AND** the two-group navigation (Plan, Organization) renders with all six items
- **AND** the footer content (project card, utility actions) renders unchanged

### Requirement: Mobile sidebar behavior SHALL be covered by automated jsdom component tests

A new test file SHALL be created at `frontend/tests/sidebar-mobile.test.ts` using the jsdom environment. The tests SHALL mock the mobile breakpoint via `window.matchMedia`, mount sidebar components using `@vue/test-utils` mount, and assert: mobile trigger visibility, Sheet opening on trigger click, and Sheet closing on route change. Existing SSR tests in `frontend/tests/app-routes-render.test.ts` SHALL remain unchanged.

#### Scenario: Test asserts mobile trigger exists when matchMedia returns matches: true

- **GIVEN** `window.matchMedia` is mocked to return `{ matches: true }` for `(max-width: 768px)`
- **WHEN** the mobile sidebar test mounts the sidebar
- **THEN** the test asserts a `SidebarTrigger` element with `data-sidebar="trigger"` is present in the rendered output

#### Scenario: Test asserts Sheet closes on route navigation

- **GIVEN** `window.matchMedia` is mocked to return `{ matches: true }`
- **AND** the mobile Sheet is open
- **WHEN** a route navigation is triggered
- **THEN** the test asserts `openMobile` becomes false after navigation

### Requirement: Mobile sidebar SHALL reuse existing sidebar primitives and tokens exclusively

The mobile sidebar implementation SHALL use only the existing shadcn-vue sidebar primitives exported from `@/shared/ui/sidebar`. No new shared-ui components or component directories SHALL be created for mobile support. All mobile Sheet styling SHALL use the existing semantic sidebar CSS tokens (`--sidebar-*`) from `frontend/src/app/styles.css`. The `SidebarProvider.vue`, `Sidebar.vue`, and `SidebarTrigger.vue` primitives SHALL NOT be modified.

#### Scenario: No new components or primitives modified

- **WHEN** the mobile sidebar implementation is inspected
- **THEN** no new `.vue` files exist in `frontend/src/shared/ui/sidebar/` beyond those already present
- **AND** `SidebarProvider.vue` has zero changes
- **AND** `Sidebar.vue` has zero changes
- **AND** `SidebarTrigger.vue` has zero changes
