## MODIFIED Requirements

### Requirement: New placeholder views SHALL render through the shared AppShell

This requirement replaces the original "New placeholder views SHALL render through the shared AppShell" requirement from the canonical `openspec/specs/sidebar-navigation/spec.md`. The restriction that SettingsView must contain only placeholder content with no interactive controls is removed. SettingsView is now permitted to contain an interactive planning-window card; its content contract is defined by the planning-window specification.

The `RoomsView.vue` and `SettingsView.vue` components SHALL render their content within the existing shared `AppShell` layout. They SHALL NOT introduce their own shell chrome, sidebar, or header. `RoomsView` SHALL contain only a minimal placeholder heading and body text. `SettingsView` SHALL contain an interactive planning-window card as defined in the planning-window specification; data fetching and interactive controls SHALL be permitted for the planning-window card.

#### Scenario: RoomsView renders through AppShell with minimal content

- **WHEN** the RoomsView component is rendered
- **THEN** the content is wrapped by the `AppShell` layout component
- **AND** no additional sidebar, header, or navigation chrome is rendered inside the view
- **AND** the view contains at most a heading and one paragraph of placeholder text

#### Scenario: SettingsView renders through AppShell with interactive planning-window card

- **WHEN** the SettingsView component is rendered
- **THEN** the content is wrapped by the `AppShell` layout component
- **AND** no additional sidebar, header, or navigation chrome is rendered inside the view
- **AND** the view contains an interactive "Planning window" card with date picker controls
- **AND** the card loads the current planning window range from the backend composable

### Requirement: Shell render tests SHALL assert the expanded settings content contract

This requirement replaces the original "Shell render tests SHALL assert the expanded sidebar contract" requirement from the canonical `openspec/specs/sidebar-navigation/spec.md`, updating the `/settings` route render test expectations to reflect the interactive planning-window form.

The route-render tests in `frontend/tests/app-routes-render.test.ts` SHALL assert the presence of the six navigation labels (Dashboard, Tasks, Schedule, People, Rooms / Areas, Settings) in the sidebar chrome. Additional route-render cases SHALL cover `/rooms` and `/settings`, asserting their route metadata and content strings. For `/settings`, the content assertion SHALL reflect the planning-window form elements rather than the previous "Feature coming soon" placeholder text. The project card content ("The Great Migration", "House move planner") and utility action labels ("Add note", "Help & Support") SHALL be asserted in the sidebar chrome test.

#### Scenario: Home route test asserts all six nav labels and footer content

- **WHEN** the home route (`/`) is rendered in tests
- **THEN** the rendered HTML contains "Rooms / Areas" and "Settings" as nav labels
- **AND** the rendered HTML contains "Add note" and "Help & Support" as footer labels

#### Scenario: New route cases exist for Rooms / Areas and Settings

- **WHEN** the test suite is executed
- **THEN** a test case exists for `/rooms` asserting its route metadata title, description, and placeholder content
- **AND** a test case exists for `/settings` asserting its route metadata title, description, and planning-window form content