## ADDED Requirements

### Requirement: The system SHALL persist tools as an organizer-curated list with an optional single bringer

The system SHALL store tools in a `tools` table with columns `id` (stable identifier prefixed
`tool-`), `name` (non-empty), `brought_by` (nullable reference to `person.id`), and `sort_order`
(integer for stable ordering). A tool SHALL have at most one bringer. When a referenced person is
deleted, the person's claimed tools SHALL revert to having no bringer rather than blocking the
deletion (`ON DELETE SET NULL`).

#### Scenario: A tool with no bringer is open

- **WHEN** a tool has `brought_by` set to null
- **THEN** the tool is reported as open (not crossed off)

#### Scenario: A tool with a bringer is crossed off

- **WHEN** a tool has `brought_by` set to an existing person ID
- **THEN** the tool is reported as crossed off, attributed to that person

#### Scenario: Deleting a person reverts their claimed tools to open

- **WHEN** a person who is the bringer of one or more tools is deleted
- **THEN** the deletion succeeds
- **AND** each of those tools has its `brought_by` set to null

### Requirement: The system SHALL provide a read endpoint returning all tools and a coverage summary

The system SHALL expose `GET /api/tools` returning a body with a `summary` object
(`total`, `claimed`, `open`) and a `tools` array ordered by `sort_order`. Each tool entry SHALL
include `id`, `name`, and `broughtBy` (the bringer's person ID, or null when open). `claimed`
SHALL equal the count of tools with a non-null bringer, `open` the count with a null bringer, and
`total` their sum.

#### Scenario: Read returns tools with a derived coverage summary

- **WHEN** a client requests `GET /api/tools` and three tools exist, one of them claimed
- **THEN** the response lists all three tools ordered by `sort_order`
- **AND** `summary.total` is 3, `summary.claimed` is 1, and `summary.open` is 2

#### Scenario: Each tool reports its bringer or null

- **WHEN** the read endpoint returns a claimed tool
- **THEN** that tool's `broughtBy` is the bringer's person ID
- **AND** an unclaimed tool's `broughtBy` is null

### Requirement: Organizers SHALL be able to create, update, and delete tools

The system SHALL expose `POST /api/tools` to create a tool, `PUT /api/tools/{id}` to update its
name and sort order, and `DELETE /api/tools/{id}` to remove it. Create SHALL require a non-empty
`name`, assign a server-generated `tool-` ID and an append sort order, and create the tool with no
bringer. Update and delete SHALL return 404 when the tool ID is unknown.

#### Scenario: Create a tool

- **WHEN** a client POSTs `{ "name": "Ladder" }` to `/api/tools`
- **THEN** a new tool is created with a server-assigned `tool-` ID, an append sort order, and a
  null bringer
- **AND** the created tool is returned

#### Scenario: Create rejects an empty name

- **WHEN** a client POSTs a tool with an empty `name`
- **THEN** the request is rejected with 400 Bad Request

#### Scenario: Update an unknown tool returns 404

- **WHEN** a client PUTs to `/api/tools/{id}` for an ID that does not exist
- **THEN** the request is rejected with 404 Not Found

#### Scenario: Delete a tool

- **WHEN** a client DELETEs an existing tool
- **THEN** the tool is removed
- **AND** a subsequent read no longer lists it

### Requirement: Any person SHALL be able to claim or unclaim a tool as its bringer

The system SHALL expose `PUT /api/tools/{id}/bringer` with body `{ personId }` to set the tool's
bringer, and `DELETE /api/tools/{id}/bringer` to clear it. PUT SHALL validate that both the tool
and the person exist, returning 404 for an unknown tool and 400 for an unknown person. Claiming
SHALL NOT depend on the person's availability. Setting a bringer on an already-claimed tool SHALL
replace the existing bringer. DELETE SHALL be idempotent, succeeding even when the tool already
has no bringer.

#### Scenario: Claim an open tool

- **WHEN** a client PUTs `{ "personId": "<existing person>" }` to `/api/tools/{id}/bringer` for an
  open tool
- **THEN** the tool's bringer is set to that person
- **AND** the tool is reported as crossed off

#### Scenario: Claiming ignores availability

- **WHEN** a person claims a tool regardless of their availability status (including `off`)
- **THEN** the claim succeeds

#### Scenario: Claim with an unknown person is rejected

- **WHEN** a client PUTs a bringer with a `personId` that does not exist
- **THEN** the request is rejected with 400 Bad Request

#### Scenario: Re-claiming replaces the bringer

- **WHEN** a client claims a tool that is already claimed by someone else
- **THEN** the tool's bringer is replaced with the new person

#### Scenario: Unclaim is idempotent

- **WHEN** a client DELETEs the bringer of a tool that has no bringer
- **THEN** the request succeeds and the tool remains open
