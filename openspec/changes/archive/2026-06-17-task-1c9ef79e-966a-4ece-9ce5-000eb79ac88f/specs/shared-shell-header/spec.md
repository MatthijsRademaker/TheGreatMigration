# shared-shell-header Specification

## ADDED Requirements

### Requirement: AppShell header SHALL be sticky at the top of the SidebarInset scroll container

The `<header>` element inside `AppShell.vue`'s `SidebarInset` SHALL use `position: sticky` with `top: 0` so it remains visible when the `SidebarInset` content overflows and scrolls vertically. The header SHALL have a `z-index` value sufficient to stay above scrolled content (at minimum `z-10`). The `SidebarInset` SHALL use `overflow-y: auto` (not `overflow-hidden`) to enable vertical scrolling of overflowing content.

#### Scenario: Header is sticky at the top on mobile viewport

- **GIVEN** the viewport is narrower than the `md` breakpoint (768px)
- **WHEN** the `SidebarInset` content is taller than the viewport height minus the header
- **AND** the user scrolls vertically within the `SidebarInset`
- **THEN** the `<header>` element remains fixed at the top of the `SidebarInset` scrollport
- **AND** the content below the header scrolls beneath it
- **AND** the mobile `SidebarTrigger` remains in the header and accessible during scroll

#### Scenario: Header is sticky at the top on desktop viewport

- **GIVEN** the viewport is at least the `md` breakpoint (768px) wide
- **WHEN** the `SidebarInset` content is taller than the available viewport height
- **AND** the user scrolls vertically within the `SidebarInset`
- **THEN** the `<header>` element remains fixed at the top of the `SidebarInset` scrollport
- **AND** the content below the header scrolls beneath it

#### Scenario: SidebarInset scroll container can scroll vertically

- **GIVEN** the `AppShell` is rendered
- **WHEN** the content inside `SidebarInset` exceeds the available height
- **THEN** the `SidebarInset` element has `overflow-y` set to a value that allows vertical scrolling (not `hidden` or `clip`)
- **AND** the user can scroll to reveal hidden content using standard touch or scroll-wheel gestures

#### Scenario: Horizontal scroll containers inside views remain functional

- **GIVEN** a view inside `RouterView` contains a horizontally scrollable element (e.g., `DailySchedule` day columns with `overflow-x: auto`)
- **WHEN** the user scrolls horizontally within that element
- **THEN** horizontal scrolling works independently of the `SidebarInset` vertical scroll
- **AND** vertical scrolling of the `SidebarInset` continues to work when scrolling outside the horizontally-scrollable element

#### Scenario: Sticky header has backdrop-blur effect on scroll

- **GIVEN** the header has `backdrop-blur` and `bg-background/90` classes
- **WHEN** content scrolls beneath the sticky header
- **THEN** the header displays the frosted-glass visual effect against scrolled content
- **AND** the header's `border-b` creates a visible separation line at the bottom of the header
