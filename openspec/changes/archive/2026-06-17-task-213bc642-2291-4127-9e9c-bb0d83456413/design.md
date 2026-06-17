## Context

The Daily Schedule presents a board of day columns with task cards showing who is assigned to what. The current flow:

```
SettingsView                    Daily Schedule
    │                                 │
    │── PUT /api/planning-window      │
    │   { startDate, endDate }        │
    │                                 │
    │                                 │── GET /api/dashboard/daily-schedule
    │                                 │   (no start/days → backend defaults to days=4)
    │                                 │── Returns 4 days of schedule cards
    │                                 │
    │   Planning window = 40 days     │
    │   But only first 4 days visible │←── problem
```

The architecture has three clear layers that must be respected:

1. **Backend**: `GET /api/dashboard/daily-schedule` already accepts optional `start` and `days` query params with correct defaults. No changes needed.
2. **Composable (useDailySchedule)**: Must gain planning-window-aware page derivation, navigation helpers, and a planning-window-change page-reset watcher — exactly mirroring `usePeopleAvailability`.
3. **Presentational component (DailySchedule.vue)**: Must stay prop-driven (`days`, `readOnly`) and NOT acquire data-fetching or pagination logic.

## Goals / Non-Goals

**Goals:**
- Users can navigate through all days of the planning window in the Daily Schedule board (prev/next controls, page size of 4 days by default)
- The composable derives explicit `start` and `days` query parameters from the planning window and current page state
- Pagination controls render at the route layer (HomeView, CalendarView), preserving DailySchedule.vue as purely presentational
- Both HomeView (read-only) and CalendarView (editable) display pagination navigation
- Page state is in-memory only — no URL sync or localStorage persistence
- Existing tests are preserved and extended with pagination-specific coverage

**Non-Goals:**
- Redesigning schedule cards, staffing indicators, or schedule CRUD workflows
- Adding new backend day-range semantics beyond the existing `start` and `days` query contract
- Adding person/task pagination, filtering, or search to the Daily Schedule board
- Changing planning-window settings behavior or unrelated dashboard panels
- URL/state persistence for the selected schedule page
- Changing the `DailySchedule.vue` component contract

## Decisions

### 1. Pagination controls live at the route layer, not in DailySchedule.vue

Both `HomeView.vue` and `CalendarView.vue` render prev/next navigation controls outside the `<DailySchedule>` component, exactly mirroring how `PeopleView.vue` renders navigation above `<PeopleAvailability>`. DailySchedule.vue remains presentational with only `days` and `readOnly` props.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Put pagination props/emits in DailySchedule.vue | Violates its OpenSpec contract (`openspec/specs/daily-schedule-component/spec.md`) that mandates it stay presentational and route-agnostic |
| Extract a shared `<SchedulePagination>` sub-component | Premature abstraction — if the two routes diverge later, extraction is straightforward |

**Why this wins:** Controlled duplication of pagination markup across two route views is cheaper than violating the component boundary. The PeopleView pattern is already proven in production.

### 2. daysPerPage defaults to 4 (not 7)

Each schedule day column is `min-w-[280px]`, so 7 columns would need ~1960px. Defaulting to 4 preserves current board readability. Routes can override via the composable's `daysPerPage` option if product later decides on a different page size.

**Why this wins:** Preserves existing visible density, matches backend's historical default, and the prop-based design makes tuning trivial.

### 3. In-memory page state only

Page position is composable-local (a Vue `ref<number>`). Neither URL query params nor localStorage are involved. This matches `usePeopleAvailability` behavior and is scoped by task non-goals.

### 4. No backend changes

The endpoint already accepts `start` and `days` query params with `days` defaulting to 4. The frontend now derives and passes these explicitly instead of relying on backend defaults. Backward compatibility is preserved — omitting `start`/`days` still returns the first 4 days.

### 5. Planning-window-change page reset (mirroring usePeopleAvailability)

When the planning window changes (e.g., admin updates the date range in Settings), the composable resets `page` to 1 if `page > 1`, using the same `oldDates.length > 0` guard as `usePeopleAvailability` to avoid resetting on initial load or SSR hydration.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|---|---|---|
| Board width vs page size: 4 columns needs ~1120px | Low — current layout already accommodates this | `daysPerPage` is configurable; default 4 matches current behavior |
| Duplicate `usePlanningWindow()` calls | Low — both `usePeopleAvailability` and the paginated `useDailySchedule` each call `usePlanningWindow()` independently | Pinia Colada deduplicates identical query keys automatically |
| Page reset timing during SSR hydration | Low — the `oldDates.length > 0` guard prevents spurious resets | Exact replication of the tested `usePeopleAvailability` watch pattern |
| `DailyScheduleState.range` reflects paginated window, not full planning window | Low — no existing consumer relies on `range` for full-window display | Documented in implementation notes; no code changes needed |
| Route-level SSR tests need mock updates | Medium — `app-routes-render.test.ts` currently mocks a 4-day schedule | Extend mocks with planning-window data and add pagination control assertions |

## Traceability

- **Task**: `213bc642-2291-4127-9e9c-bb0d83456413`
- **Dossier**: `2026-06-17T09:19:37.138Z` — problem framing, goals, non-goals, open questions
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation` — all three refinement agents independently converged on route-level pagination, daysPerPage=4, in-memory state, no backend changes
- **Reference**: People Availability pagination change `task-2abd2a16-d49c-4fb5-ab32-ca036cfa1010`
