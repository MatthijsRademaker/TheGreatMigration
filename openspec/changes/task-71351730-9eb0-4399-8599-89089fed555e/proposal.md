## Why

The "People" navigation entry in the sidebar is unreliably clickable — users report it either does not respond to clicks or requires very precise targeting. Since every other nav item uses the same template pattern and works correctly, the issue stems from a subtle interaction in the tooltip-wrapped `SidebarMenuButton` rendering chain that degrades pointer-event forwarding for certain entries.

## What Changes

- **Fix the tooltip-wrapped router-link clickability**: Eliminate the `RouterLink` event-forwarding fragility by restructuring how `SidebarMenuButton` wraps navigation links when a tooltip is present.
- **Remove the duplicate branding block from `SidebarFooter`**: The footer contains an identical copy of the header's "The Great Migration / House move planner" brand card. This is a copy-paste artifact that adds vertical clutter (48px extra height per item).
- **Ensure all six nav items have consistent, full-width click targets in both expanded and collapsed sidebar states**.
- Add a targeted unit test that verifies each nav `RouterLink` in the sidebar is clickable and navigates to its expected route.

## Capabilities

### Modified Capabilities

- `sidebar-navigation`: The sidebar navigation items' click targets need guaranteed full-width pointer-event coverage, and the footer duplicate branding must be removed.

## Impact

- **frontend/src/shared/layout/app-sidebar/AppSidebar.vue**: Template restructuring for nav items; removal of duplicate footer branding block.
- **frontend/src/shared/ui/sidebar/SidebarMenuButton.vue**: May need structural changes to the tooltip/navigation wrapping pattern.
- **frontend/src/shared/ui/sidebar/SidebarMenuButtonChild.vue**: Potentially simplify the `as-child` forwarding chain.
- **frontend/tests/sidebar-mobile.test.ts**: New test for nav-item clickability.
- No backend, API, or dependency changes.
