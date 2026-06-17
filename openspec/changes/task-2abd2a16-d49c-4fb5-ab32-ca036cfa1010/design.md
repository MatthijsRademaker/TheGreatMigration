## Context

The People Availability panel (`/people`) renders a matrix of people × days where users assign availability statuses (`available`, `busy`, `partial`, `off`). The current flow:

```
SettingsView                    People Availability Panel
   │                                      │
   │── PUT /api/planning-window           │
   │   { startDate, endDate }             │
   │                                      │
   │                                      │── GET /api/dashboard/people-availability
   │                                      │   (start=planWindow.startDate, days=4)
   │                                      │── Returns ALL people × 4 days
   │                                      │
   │   Planning window = 40 days          │
   │   But only first 4 days visible      │←── problem
```

Two gaps need closing:

1. **Day pagination**: The matrix is hardcoded to 4 days. Users need to navigate through all days in the planning window.
2. **Person scalability**: `GetAllPeople` fetches every person unconditionally. As the organization grows, the payload grows unboundedly.

This design addresses both with a consistent pagination model.

## Goals / Non-Goals

**Goals:**

- Users can navigate through all days of the planning window in the People Availability matrix (prev/next controls, configurable page size)
- The backend API supports person-level pagination (`offset`, `limit`) so the frontend can request a subset of people
- The backend response includes pagination metadata (`totalPeople`, `page`, `perPage`)  for pagination-aware UI rendering
- API backward compatibility: omitting `offset`/`limit` returns all people (unchanged behavior)
- Frontend composable accepts `days` and `page` parameters, computes the effective date range from the planning window
- All existing tests pass without modification

**Non-Goals:**

- Adding person-level pagination controls to the UI (the frontend will initially request all people; pagination metadata enables future UI pagination without backend changes)
- Changing the upsert/delete availability endpoints — they already work with any date in the planning window
- Adding bulk editing, drag-and-drop, or search/filter to the matrix
- Changing the homepage (`HomeView.vue`) layout or composition
- Changing the `PeopleAvailability.vue` component's cell-rendering contract (badges, popovers)

## Decisions

### 1. Backend: Person pagination via `offset` and `limit` query parameters

`GET /api/dashboard/people-availability` adds optional `offset` (default 0) and `limit` (default 0, meaning no limit / return all) query parameters. The `DashboardBody` response adds `totalPeople` (total count regardless of pagination), `page` (computed from offset/limit), and `perPage` (the limit value, or total count if no limit).

**When `limit` is 0 (default):** All people are returned, `page=1`, `perPage=totalPeople`. This preserves backward compatibility — existing callers see no behavioral change.

**When `limit` > 0:** Only the requested slice of people is returned. The store uses a paginated SQL query (`LIMIT $1 OFFSET $2`). A separate `SELECT COUNT(*) FROM people` provides the total.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Cursor-based pagination | Simpler to implement offset-based for this use case; no ordering concerns beyond `ORDER BY id` |
| Keep all-people always | Works for today's ~10-person seed data, but doesn't scale to 50+ people × 40-day windows |
| Only add `days` as a query override, skip person pagination | Correct but incomplete — person pagination is cheap to add now and avoids a breaking change later |

**Why this wins:** Adding `offset`/`limit` to the existing endpoint is a backward-compatible additive change. It future-proofs the API without increasing complexity for today's callers. The pagination metadata lets the frontend render sensible UI ("Showing 1-10 of 34 people") without a second API call.

### 2. Frontend: Day pagination driven by planning window + page state

The `usePeopleAvailability()` composable accepts a `page` (1-indexed) and `daysPerPage` (default 7) in its options. It computes:

```
currentStartDate = planningWindow.startDate + (page - 1) * daysPerPage
totalDays = planningWindow.totalDays
totalPages = ceil(totalDays / daysPerPage)
```

This replaces the hardcoded `start` (planning window start) and `days` (4) with dynamic values computed from the page state and planning window.

**Page controls** render in `PeopleView.vue` above the matrix, showing:

- "Mon 5 Jul – Sun 11 Jul" (date range label)
- "Page 1 of 6" (page indicator)
- [< Prev] [Next >] buttons (disabled at boundaries)

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Horizontal scroll with all days visible | 40+ columns is unusable on most screens; requires both scroll and zoom interaction |
| Infinite scroll (load more days on demand) | Over-engineered for the day dimension; users think in pages/weeks |
| Dropdown date range selector | More interaction steps; pagination is more familiar for linear navigation |

**Why this wins:** Pagination is the most recognizable navigation pattern. A weekly cadence (7 days) aligns with the planning mental model. The implementation is entirely frontend-local — no backend changes needed for day navigation.

### 3. Frontend: People pagination metadata enables future UI

The composable already passes `offset` and `limit` to the API based on the page state. For the initial implementation, `offset=0` and `limit=0` (all people) are used — meaning the backend returns all people as before. The pagination metadata from the response (`totalPeople`, `page`, `perPage`) is exposed on the composable return value so future UI work can add person-page controls without backend changes.

This means: **the backend adds person pagination support now, the frontend consumes the metadata, but person-level pagination controls are not part of this change.** The day pagination controls ARE part of this change.

### 4. `availableToday` summary scoped to visible page

The `availableToday` count currently reflects all people on the selected date (the first day of the visible range). With pagination, this count is scoped to the **visible people** on the **first visible date**. This is consistent behavior — it already counts from the returned data, and people pagination changes which people are returned.

### 5. OpenAPI regeneration

After adding `offset`/`limit` to the `DashboardInput` struct and `totalPeople`/`page`/`perPage` to the `DashboardBody` struct, the OpenAPI snapshot and generated frontend client must be regenerated.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|---|---|---|
| `limit=0` semantics (0 = no limit) conflicts with SQL `LIMIT` | SQL `LIMIT 0` returns zero rows, but we want all rows | Store code checks `if limit <= 0 { /* no LIMIT clause */ }` before building the query |
| Day pagination computes `start` from planning window; if planning window changes mid-session, stale page | Shows wrong dates or breaks page count | `usePlanningWindow()` is reactively computed; when planning window data changes, the composable recomputes and the page state resets to page 1 |
| `SELECT COUNT(*)` on every paginated request adds DB overhead | Slight latency on each request | The people table is expected to be small (<1000 rows); count is fast on a PK index. If it becomes a bottleneck, cache the count with a short TTL. |
| Frontend pagination controls increase component complexity | More code to maintain | Keep pagination state in the composable, not the component. PeopleView only calls the composable and renders controls via props — no direct pagination logic. |
| OpenAPI client regeneration may break if types shift | Generated types mismatch expectations | `DashboardBody` adds new optional fields; existing code that destructures only known fields is unaffected. Run `scripts/precommit-run` to verify. |
