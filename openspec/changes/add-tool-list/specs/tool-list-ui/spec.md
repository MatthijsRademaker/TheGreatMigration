## ADDED Requirements

### Requirement: The application SHALL provide a /tools route rendering the tool list

The application SHALL register a `/tools` route that renders a tool-list view backed by a
composable wrapping the `GET /api/tools` query. The view SHALL list each tool with its name and
its claim state, and SHALL distinguish open tools from crossed-off tools (which display the
bringer). The view SHALL render loading, error, and empty states.

#### Scenario: Tools route renders the list

- **WHEN** a user navigates to `/tools` and tools exist
- **THEN** each tool is rendered with its name
- **AND** crossed-off tools are visually distinguished from open tools and show their bringer

#### Scenario: Empty tool list

- **WHEN** the `/tools` route renders and no tools exist
- **THEN** an empty state is shown

### Requirement: The view SHALL let a user claim a tool via a person picker sourced from the people list

The tool-list view SHALL provide, for each open tool, a way to claim it by selecting a person from
a picker populated by the existing people query (the same people shown on the availability screen).
Selecting a person SHALL call `PUT /api/tools/{id}/bringer` and refresh the list so the tool shows
as crossed off.

#### Scenario: Claim a tool from the picker

- **WHEN** a user selects a person from the picker for an open tool
- **THEN** the tool is claimed via `PUT /api/tools/{id}/bringer`
- **AND** after the list refreshes the tool shows as crossed off by that person

#### Scenario: Picker offers all people

- **WHEN** the bringer picker opens
- **THEN** it lists people from the existing people query without filtering by availability

### Requirement: The view SHALL let a user unclaim a crossed-off tool

The tool-list view SHALL provide, for each crossed-off tool, a way to remove its bringer. Doing so
SHALL call `DELETE /api/tools/{id}/bringer` and refresh the list so the tool returns to open.

#### Scenario: Unclaim a tool

- **WHEN** a user removes the bringer from a crossed-off tool
- **THEN** the tool is unclaimed via `DELETE /api/tools/{id}/bringer`
- **AND** after the list refreshes the tool shows as open

### Requirement: The view SHALL let an organizer add and remove tools

The tool-list view SHALL provide controls to add a tool by name (calling `POST /api/tools`) and to
remove a tool (calling `DELETE /api/tools/{id}`), refreshing the list after each.

#### Scenario: Add a tool

- **WHEN** a user submits a new tool name
- **THEN** the tool is created via `POST /api/tools`
- **AND** the new tool appears in the refreshed list as open

#### Scenario: Remove a tool

- **WHEN** a user removes a tool
- **THEN** the tool is deleted via `DELETE /api/tools/{id}`
- **AND** it no longer appears in the refreshed list
