## Context

The app already models **people** (with per-date availability) and a **task backlog** (rows
with a status and a list of assigned person IDs). The tool list is a new, smaller domain that
sits beside them: an organizer-curated list of required physical tools where the same people
can claim what they will bring. The backend is Go + Huma (OpenAPI) over pgx/Postgres with a
generated query layer under `backend/db/`; the frontend is Vue 3 + Pinia Colada consuming a
client generated from the OpenAPI spec. Two existing patterns are reused wholesale, so this
change is mostly assembly, not invention:

- **`task-backlog-api`**: table + read endpoint returning rows plus a derived summary + CRUD.
- **People availability `PUT/DELETE /api/people/{id}/availability/{date}`**: an idempotent
  upsert/delete pair that validates person existence — the exact shape of claim/unclaim.

## Goals / Non-Goals

**Goals:**

- A single source of truth for required tools and who is bringing each.
- Claiming a tool is one action and equals crossing it off; unclaiming reverts it.
- Reuse the existing people list for the bringer picker — no parallel people model.
- Match the existing `tasks` / `people` code layout, generated-client flow, and spec style.

**Non-Goals:**

- No quantities or multiple bringers per tool (each tool is one item, one bringer).
- No separate "confirmed present" state distinct from "claimed".
- No coupling to availability, the planning window, scheduling, or rooms.
- No per-helper add permission — only organizers curate the list (helpers claim only).

## Decisions

### A single nullable `brought_by` FK is the entire state machine

`tools(id, name, brought_by NULL → person.id, sort_order)`. `brought_by IS NULL` = open;
non-null = crossed off by that person.

- **Why**: The user confirmed claim = done and single-item-each. A boolean `done` column would
  be redundant (it would always equal `brought_by IS NOT NULL`) and could drift out of sync.
- **Alternative considered**: a `tool_bringers` join table (mirroring `task` assignments).
  Rejected — that models *many* bringers; here exactly one person brings one tool, so a scalar
  FK is simpler and makes "who brings it" a single column read.
- **FK behavior**: `brought_by` references `person(id)` with `ON DELETE SET NULL` so deleting a
  person reverts their claimed tools to open rather than blocking the delete or orphaning rows.
  This is intentionally looser than tasks (which block person deletion on references), because a
  claim is a soft commitment, not an assignment that schedule cards depend on.

### Claim / unclaim as an upsert/delete pair on a sub-resource

`PUT /api/tools/{id}/bringer` with `{ personId }` claims; `DELETE /api/tools/{id}/bringer`
unclaims. Both validate the tool exists; PUT also validates the person exists (reusing
`store.PersonExists`). DELETE is idempotent.

- **Why**: Structurally identical to the availability upsert/delete already in `people.go`, so
  it reuses the validation idiom and generates a clean mutation pair in the client.
- **Alternative considered**: folding the bringer into `PUT /api/tools/{id}` (full-row update).
  Rejected — claiming is the high-frequency helper action and should not require sending the
  whole row or carry organizer-edit semantics.

### Read endpoint returns rows + a coverage summary

`GET /api/tools` → `{ summary: { total, claimed, open }, tools: [...] }`, mirroring
`TaskBacklogBody`'s summary+rows shape. The home KPI card and the `/tools` screen both read this
one endpoint; the KPI shows `claimed / total`.

- **Why**: One read source of truth, derived counts computed server-side (consistent with
  `TaskSummary`), avoids the client recomputing coverage.

### Frontend mirrors `frontend/src/tasks/`

New `frontend/src/tools/` with `types.ts`, `ToolsView.vue`, a `useTools` composable wrapping the
generated query, and row/picker components. The bringer picker is fed by the existing people
query (the same one `usePeopleAvailability` uses) so it is literally "the same people from the
availability screen". A new `/tools` route is registered in `frontend/src/app/routes.ts`.

### KPI and sidebar are enumerated specs → treated as modified capabilities

`kpi-summary-cards` enumerates exactly four cards in order, and `sidebar-navigation` enumerates
exactly six items. Adding a tools KPI card and a Tools nav item changes those enumerations, so
both are MODIFIED-requirement deltas, not incidental edits.

## Risks / Trade-offs

- **ON DELETE SET NULL silently drops claims** → Acceptable and intended: a deleted person's
  commitments should reopen, not block deletion. Documented in the spec scenario so it is a
  deliberate behavior, not a surprise.
- **No concurrency guard on claiming** → Two helpers could claim the same open tool near-
  simultaneously; last write wins. Acceptable for a small trusted group; the read refresh after
  mutation surfaces the final state. No optimistic-locking column added (would be speculative).
- **Modifying the prescriptive KPI/sidebar specs** → These specs assert exact counts/orders, so
  the deltas must restate the full updated requirement (five cards, seven nav items) to avoid
  losing detail at archive time.
- **Generated-client regen is a required step** → If proto/OpenAPI generation is skipped, the
  frontend won't see the new endpoints. Captured as explicit tasks gated behind the Docker
  verification scripts.

## Migration Plan

1. Add the `tools` table via a new SQL migration (additive; no existing table touched).
2. Add queries + regenerate the `backend/db/` layer; implement `backend/api/tools.go`; register
   endpoints.
3. Regenerate OpenAPI + the frontend client.
4. Build the `/tools` UI, KPI card, and nav entry.
5. Verify via `scripts/precommit-run`.

Rollback is clean: drop the `tools` table and remove the route/card/nav additions; nothing else
depends on tools.

## Open Questions

None — all scoping decisions (claim=done, single-item, organizer-curated, own route + KPI,
anyone-can-claim) were resolved during exploration.
