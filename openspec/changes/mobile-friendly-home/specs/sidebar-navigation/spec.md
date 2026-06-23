## ADDED Requirements

### Requirement: Sidebar footer decorative image SHALL be hidden on mobile viewports
The decorative `trail-portrait.png` image rendered in `AppSidebar.vue`'s `SidebarFooter` SHALL NOT be rendered when `isMobile` is `true`. The guard SHALL use `v-if="!isMobile"` with the `isMobile` ref sourced from `useSidebar()` — the same composable already used in the component. On desktop viewports (`isMobile` is `false`) the image SHALL render exactly as it does today, including being hidden when the sidebar is collapsed to icon mode via the existing `group-data-[collapsible=icon]:hidden` class.

#### Scenario: Footer image is absent on mobile sidebar
- **WHEN** the mobile sidebar Sheet is open on a viewport where `isMobile` is `true`
- **THEN** the `trail-portrait.png` image element is not present in the rendered DOM
- **AND** no oversized image overlaps the sidebar navigation items

#### Scenario: Footer image renders on desktop when sidebar is expanded
- **WHEN** the sidebar renders on a desktop viewport (`isMobile` is `false`) and the sidebar is in expanded state
- **THEN** the `trail-portrait.png` image is present in the rendered DOM
- **AND** the image fills the footer area as it did before this change

#### Scenario: Footer image remains hidden on desktop when sidebar is icon-collapsed
- **WHEN** the sidebar renders on a desktop viewport and the sidebar is collapsed to icon mode
- **THEN** the `trail-portrait.png` image is not visible (hidden via the existing `group-data-[collapsible=icon]:hidden` class)
