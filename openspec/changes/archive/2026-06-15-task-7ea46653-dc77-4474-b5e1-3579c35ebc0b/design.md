## Context

The Great Migration is a lightweight house-move planner. The homepage dashboard (`/`) shows summary cards and a planning overview. The design (`designs/home-page.png`) includes a "People available today 6 / 8" summary card and a "People Availability" multi-day grid (2 Jul–5 Jul). Currently `HomeView.vue` hardcodes these as static values, and the backend serves only `GET /api/hello`.

The architecture model (`.devagent/architecture/workspace.c4`) positions the Dashboard as a summarizer of `peopleData`, confirming a BFF-style read endpoint is architecturally correct. The design system (`docs/design-system-v2.md`) defines four canonical availability statuses: `available` (success), `busy` (destructive), `partial` (warning), `off` (muted/neutral). No persistence model exists yet, making in-memory seeded data the appropriate first step.

## Goals

- Expose a single `GET /api/dashboard/people-availability` Huma-registered endpoint for homepage people-availability data.
- Return a combined payload — range metadata, summary counts, per-person daily statuses, and status legend — in one fetch.
- Constrain availability statuses to the four design-backed values.
- Seed in-memory data with at least 8 people covering all four status states.
- Add backend tests proving response shape and range/status correctness.
- Preserve existing `/api/hello`, CORS, and `/openapi.json` behavior.

## Non-Goals

- Building or wiring the homepage `/people` frontend UI.
- Implementing unrelated dashboard APIs (tasks, schedule, move notes).
- Adding authentication, persistence, or write/edit workflows.
- Creating a generic people CRUD API under `/api/people/`.

## Decisions

### 1. Single combined BFF endpoint

**Decision:** Expose `GET /api/dashboard/people-availability` as a single endpoint returning range metadata, summary counts, per-person daily availability, and a status legend in one combined payload.

**Rationale:** The homepage design shows both a summary card and a multi-day grid. Serving both from one endpoint avoids sequential fetches, reduces latency, and matches the BFF read-model pattern confirmed by the architecture model (Dashboard summarizes peopleData).

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`

### 2. Dashboard route namespace

**Decision:** Place the endpoint under `/api/dashboard/` rather than `/api/people/`.

**Rationale:** This signals the endpoint is a dashboard-optimized BFF aggregate, not a generic people CRUD API. The architecture C4 model explicitly shows the Dashboard summarizes peopleData. A future generic people API can live independently under `/api/people/` without conflict.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, workspace.c4

### 3. In-memory seeded data for first slice

**Decision:** Use in-memory data seeded with at least 8 people exercising all four canonical statuses (`available`, `busy`, `partial`, `off`) across the returned date range.

**Rationale:** No persistence model exists in the repository. The acceptance criteria explicitly allow seeded/in-memory data. The seed must include stable person keys (IDs, names, initials) and enough variety to drive the `6 / 8` summary card.

**Sources:** `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, acceptance criteria

### 4. Handler in separate `backend/dashboard.go` file

**Decision:** Place the new handler, data types, and seed data in a new `backend/dashboard.go` file under package `main`. Register the route in `backend/main.go`.

**Rationale:** Separating the handler keeps `main.go` readable and lets the seed data and handler evolve independently when persistence is added later.

**Sources:** `round:1:agent:swarm-architect`

### 5. Query parameter defaults and semantics

**Decision:** Accept optional query parameters `start` (ISO 8601 date string, default server-local today) and `days` (positive integer, default 4, inclusive of `start`). `selectedDate` in the response defaults to `startDate`; no separate query parameter for it in this slice.

**Rationale:** `days` inclusive of `start` makes `days=4` produce exactly 4 dates, matching the 2 Jul–5 Jul window in the design. Making `start` optional with a server-local default keeps the API simple for initial usage; the endpoint description documents that clients should pass `start` explicitly for timezone-correct results. `selectedDate` as a computed field (equal to `startDate`) means the summary card always references the first date in the range; a separate query param can be added later if the frontend needs a distinct reference date.

**Sources:** Round 1 architect/lead-dev questions, `round:1:agent:swarm-architect` risks

### 6. `availableToday` counting rule

**Decision:** `availableToday` counts only people whose status on `selectedDate` equals `available`. `partial` status does not count.

**Rationale:** The design-system description says "Helpers with confirmed availability" which implies full-day `available` status. The counting rule is documented in the endpoint description to avoid frontend ambiguity.

**Sources:** `round:1:agent:swarm-architect` suggested requirements, design-system-v2.md

### 7. Status legend in response

**Decision:** The top-level `statuses` array contains the full canonical status definitions as a legend: `{"id", "label", "colorIntent"}` for each of the four statuses.

**Rationale:** Providing the full legend lets the frontend render status chips without hardcoding display metadata, and supports future status additions without API changes.

**Sources:** `round:1:agent:swarm-reviewer` questions, design-system-v2.md

## Risks

| Risk | Mitigation |
|------|------------|
| Route naming: `/api/dashboard/` prefix may be unexpected if the team views this as a domain API. | Document in the endpoint description that this is a dashboard-optimized BFF aggregate, and a generic people API can live under `/api/people/` independently. The workspace.c4 architecture confirms this orientation. |
| `availableToday` semantics: frontend may expect `partial` to count toward "available today." | Document the counting rule in the endpoint description. Return the status legend in the response so the frontend can interpret statuses explicitly. |
| Default date behavior: server-local default may differ from client timezone. | Document that the `start` parameter defaults to server-local date and clients should pass `start` explicitly for timezone-correct results. |
| Seed data shape stability: when persistence is added, the response shape may change. | The acceptance criteria define a stable contract (`range`, `summary`, `people[].availability`) that should survive a backend swap. |

## Traceability

- **Task**: `7ea46653-dc77-4474-b5e1-3579c35ebc0b`
- **Dossier**: `2026-06-15T17:12:21.851Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Validated round outputs**: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- **Artifact base**: `initial` snapshot for `task-7ea46653-dc77-4474-b5e1-3579c35ebc0b`
