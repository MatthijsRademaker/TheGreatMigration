## Context

The People Availability matrix on `/people` displays status pills for each person across a date range. The range comes from `GET /api/dashboard/people-availability`, which defaults to `time.Now()` + 4 days when no explicit `start` query parameter is sent. But the availability upsert endpoint (`PUT /api/people/{id}/availability/{date}`) validates every date against the **planning window** — a fixed date range stored in the `planning_windows` database table (seeded as `2026-07-05` to `2026-08-13`).

This creates an inconsistency: the dashboard can display dates outside the planning window (e.g., today = June 16), making those status pills appear editable (clicking opens the popover), but every write attempt returns a 400 error "date is outside the planning window."

Additionally, the frontend error handler maps **all** 400 errors to the same generic message "Invalid status or date," making the real problem opaque. And the date-derivation in `handleCellUpdate` relies on `rawData.range.startDate` via a fragile `getISODate` helper that separates the date from the component's adapted `days` array.

**Current data flow:**

```
Frontend                     Backend
─────────────────────────────────────────
GET /dashboard/people-avail
  (no start param) ──────►  defaults to time.Now()
  ◄── range.startDate       returns data for June 16-19

User clicks cell, emits
  { personId, dayIndex, status }

getISODate(dayIndex) ──►   computes "2026-06-17"
PUT /people/{id}/avail/2026-06-17 ──►  validates against planning_window
  ◄── 400                   "date is outside the planning window"

"Invalid status or date."   ← generic catch-all
```

## Goals / Non-Goals

**Goals:**

- Make the dashboard default to planning window dates so that displayed cells are always writable
- Surface distinct error messages for each 400 subtype (invalid status, invalid date format, date outside planning window)
- Replace the fragile `getISODate(rawData.range.startDate + dayIndex)` date derivation with a direct lookup from the adapted `days` array
- Pass explicit `start` query parameter from the planning window in the frontend query so the dependency is explicit, not implicit
- Existing explicit `start` and `days` query parameters remain functional for testing and manual override
- All existing tests pass; new tests cover the changed behavior

**Non-Goals:**

- Removing or relaxing the planning-window validation on the upsert endpoint (the validation is correct — availability should only be editable within the move timeline)
- Adding drag-and-drop, bulk editing, or date navigation to the matrix
- Changing the `PeopleAvailability.vue` component contract (emits, props)
- Changing the planning-window CRUD or mutation surface
- Changing `HomeView.vue` layout or composition
- Changing the daily-schedule or task-backlog endpoints

## Decisions

### 1. Backend default: planning window start date

When `GET /api/dashboard/people-availability` receives no `start` query parameter, the handler reads the planning window via `store.GetPlanningWindow(ctx)` and uses `pw.StartDate` as the start date instead of `time.Now()`. If no planning window exists and no `start` is provided, the handler returns a 400 error.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Keep `time.Now()` default, add frontend guard | Two sources of truth — backend still serves data for non-editable dates to any client besides the web UI |
| Make the frontend always pass explicit `start` | Defensive but doesn't protect other API consumers (CLI, scripts, future clients) |
| Remove the planning-window validation from upsert | Undermines the planning window as the source of truth for the move timeline |

**Why backend change wins:** The planning window is the canonical date authority. The dashboard and the upsert should both reference it. Changing the default in the handler makes this invariant true for all consumers.

### 2. Frontend: explicit `start` passed from planning window

The `usePeopleAvailability()` composable accepts an optional `start` parameter. By default, it calls `usePlanningWindow()` to get the planning window's `startDate` and passes it as `start` to the dashboard query. This makes the dependency explicit in the URL and provides a clear override mechanism for testing.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Trust the backend default exclusively | Works for the common case but makes testing harder — tests would need to set up a planning window in the DB |
| Hard-code nothing | Same as above |
| Read `start` from a route query param | Over-engineered for this use case; planning-window start is the only override needed |

**Why explicit wins:** Explicit query params are visible in network logs, cache keys, and test assertions. The composable becomes self-documenting — it's clear what date range is being requested without knowing the backend default logic.

### 3. Frontend date derivation: from adapted `days` array index

The `PeopleAvailability.vue` template emits `dayIndex` — the `v-for` index within `person.availability`. Since the backend always produces exactly `days` entries per person, this index matches the index in the adapted `days` array. The `handleCellUpdate` handler now looks up the ISO date by index from a derived `daysByIndex` map instead of computing it from `rawData.range.startDate`.

```ts
// Before (fragile):
function getISODate(dayIndex: number): string {
  const start = new Date(rawData.value.range.startDate)
  start.setUTCDate(start.getUTCDate() + dayIndex)
  return start.toISOString().slice(0, 10)
}

// After (direct lookup):
// Adapted composable also returns a `daysISO` array parallel to `days` labels.
function getISODate(dayIndex: number): string {
  return daysISO.value[dayIndex] ?? ''
}
```

The `daysISO` array is computed once in the composable alongside the day labels, using the same UTC-based iteration logic already present in `adaptToComponentProps`. This eliminates the separate `getISODate` function and its dependency on raw API internals.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Keep `getISODate(rawData.range.startDate + dayIndex)` | Still works today, but couples date knowledge to raw response shape. If the composable adapter changes how it builds days (e.g., filtering, reordering), dates silently diverge. |
| Emit the ISO date string from the cell itself | Would require changing `PeopleAvailability.vue` prop shape to carry ISO dates through to each cell — more coupling |
| Emit the day label and have the parent reverse-map it | Fragile string matching |

**Why direct lookup wins:** The adapted `days` prop is the authoritative list of displayed day labels. Keeping its ISO-date equivalents alongside it in a parallel array makes the mapping explicit, testable, and independent of raw API internals.

### 4. Frontend error differentiation: parse backend 400 body

Instead of mapping all HTTP 400 errors to "Invalid status or date," the handler inspects the error body's `detail` field (from the Huma error model):

```ts
const detail = (err as any)?.cause?.body?.detail ?? ''
if (detail.includes('outside the planning window')) {
  updateError.value = 'This date is outside the planning window.'
} else if (detail.includes('must be a valid ISO 8601 date')) {
  updateError.value = 'Invalid date format.'
} else if (detail.includes('status must be one of')) {
  updateError.value = 'Invalid status value.'
} else {
  updateError.value = `Failed to update: ${detail}`
}
```

**Alternative considered:** Create a typed error-parsing utility that extracts the Huma error detail. That's a good follow-up for shared error handling but adds scope to this bug-fix change.

**Why inline parsing wins:** Minimal code, directly addresses the user-facing confusion, and doesn't introduce a shared abstraction before the error-handling pattern is proven.

### 5. The composable returns `daysISO`

`usePeopleAvailability()` adds a `daysISO` computed ref — a `string[]` of ISO date strings (YYYY-MM-DD) aligned with the adapted `days` labels. This is derived in the `adaptToComponentProps` function alongside the label array, using the same UTC-day iteration.

This means `PeopleView.vue` can destructure:

```ts
const { data: availabilityData, rawData, daysISO, isLoading, isError, isEmpty } = usePeopleAvailability()
```

And use `daysISO.value[dayIndex]` directly in `handleCellUpdate`.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|---|---|---|
| Backend `dashboard` handler now calls `store.GetPlanningWindow()`, adding a DB round-trip | Slight latency increase on default requests | Planning window is a small singleton query (single row, indexed PK). Sub-millisecond overhead. |
| No planning window in database (fresh migration) | Dashboard returns 400 for default requests | The seeding migration always inserts a planning window. If absent, the explicit `start` override still works. Test with `failingStore`. |
| Frontend `usePlanningWindow()` is a separate async query — dashboard may render before planning window resolves | Flash of wrong dates or loading state | The composable already has loading/error states. `usePeopleAvailability()` can wait for the planning window before firing the dashboard query (leveraging Pinia Colada's `dependsOn` or an `enabled` guard). |
| If `usePlanningWindow()` is used on pages that don't need it, the query is fetched unnecessarily | Wasted network request | Pinia Colada caches by query key — duplicate requests are deduplicated. The planning window is fetched once per app session. |
| `daysISO` array could drift from `days` labels if the adapter logic changes | Silent date mismatch | Both arrays are built in the same loop within `adaptToComponentProps`. Adding a unit test that asserts `daysISO[i]` correlates to `days[i]` catches drift. |
| Inline error detail parsing is brittle | Backend error message text changes → broken parsing | The Huma error model is stable. Document the dependency in a code comment. If the backend changes error text, tests fail with a clear signal. |
