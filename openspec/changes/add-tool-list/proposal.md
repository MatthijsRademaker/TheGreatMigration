## Why

Migration days need physical tools (ladder, drill, tarps), but today there is no shared
record of what is required or who is bringing each item — coordination happens out-of-band
and items get forgotten or doubled up. An organizer-curated tool list where the existing
helpers can claim "I'll bring this" closes that gap with a single source of truth.

## What Changes

- Add a **tools** domain: an organizer-curated list of required tools, each with a name and
  an optional bringer (a reference to an existing person).
- A tool is **claimed / crossed off** when it has a bringer, and **open** when it has none —
  a single nullable person reference is the entire state. Claiming is the same act as
  crossing it off (no separate "done" step).
- Each tool is a **single item** brought by exactly one person; a person may bring many tools.
- **Anyone** in the people list can claim or unclaim a tool — no availability check.
- New backend endpoints: read all tools (with coverage summary), organizer CRUD on the list,
  and claim/unclaim a bringer.
- New `/tools` route with a tool-list screen, reusing the existing people query for the
  bringer picker.
- A new home-dashboard KPI card showing tool coverage (claimed / total).
- A new **Tools** sidebar navigation item.

## Capabilities

### New Capabilities

- `tool-list-api`: Backend persistence and HTTP endpoints for the tool list — a read endpoint
  returning all tools plus a coverage summary, organizer CRUD (create / rename-reorder /
  delete), and claim/unclaim of a bringer via an upsert/delete pair on a person reference.
- `tool-list-ui`: The `/tools` route, its view and composable, the bringer picker (sourced from
  the existing people query), and the claim/unclaim and organizer-edit interactions.

### Modified Capabilities

- `kpi-summary-cards`: The home dashboard adds a fifth KPI card for tool coverage; the
  enumerated card set and order change.
- `sidebar-navigation`: The Plan group adds a **Tools** navigation item linking to `/tools`;
  the enumerated item count and icon set change.

## Impact

- **Database**: new `tools` table (`id`, `name`, `brought_by` nullable FK → person, `sort_order`)
  plus a migration and generated query layer under `backend/db/`.
- **Backend API**: new `backend/api/tools.go` (read + CRUD + bringer endpoints) registered with
  Huma; OpenAPI regenerated.
- **Generated client**: regenerated `frontend/src/client/` (types + Pinia Colada query/mutation
  helpers) from the updated OpenAPI spec.
- **Frontend**: new `frontend/src/tools/` (view, composable, types, components); a new KPI card
  in `frontend/src/home/components/KpiCards.vue`; a new nav entry in the sidebar and a route in
  `frontend/src/app/routes.ts`.
- **Decoupling**: tools are independent of availability and scheduling — no changes to the
  planning window, calendar, or assignment logic.
