# Proposal: Date range planning configurable

## Why

The app's move timeline is currently scattered across hard-coded frontend literals with no shared planning-window source of truth. `CalendarView.vue` hard-codes five Mon–Fri day columns with narrative `focus` and `helpers` fields, `HomeView.vue` uses a static `"5"` for the Move days summary card, and `PeopleView.vue` stores helper availability as weekday strings unrelated to any calendar range. Changing the move period would require hunting through multiple files and create inconsistent behavior across views — exactly the coupling this task aims to eliminate.

## What Changes

- **Create `frontend/src/shared/lib/planWindow.ts`** as the single canonical source of truth for the move-plan date range. The module exports `PLAN_WINDOW_START` and `PLAN_WINDOW_END` ISO date string constants (defaulting to `2026-07-05` and `2026-08-13`, matching the task example "5 juli until 13 aug") plus derived helpers: `planWindowDays` (array of `{date: Date, label: string, dateString: string}`), `planWindowDayCount` (inclusive day count), and `formatPlanDayLabel(date: Date): string` (formats as e.g. "Sat 5 Jul").
- **Refactor `CalendarView.vue`** to derive its rendered day columns from `planWindowDays` instead of the hard-coded five-element Mon–Fri array. The per-day `focus` and `helpers` narrative fields are removed; each column displays the date label and a generic placeholder description. The grid layout adjusts from `md:grid-cols-5` to a responsive column count appropriate for the window size.
- **Refactor `HomeView.vue`** to derive its "Move days" summary card value from `planWindowDayCount` instead of the hard-coded string `"5"`. The card's description ("Scheduled working days in the plan") and icon remain unchanged.
- **Amend `openspec/specs/hello-world-integration/spec.md`** to narrow the static-content requirement: the "Other HomeView sections are preserved" scenario is updated so that the "Move days" card is no longer required to render with static content, while "Available today" and "Under-staffed" remain unchanged.
- **Add unit tests** in `frontend/tests/planWindow.test.ts` covering inclusive-range derivation, day count accuracy, label formatting, and propagation when constants change. Extend the existing route render test to assert that CalendarView renders the correct number of day columns matching `planWindowDayCount`.
- **Explicitly defer `PeopleView.vue`** — its weekday availability strings (`Mon, Wed, Fri`) represent availability patterns, not date-derived data. Converting them would expand scope into building real availability data structures that the non-goals forbid. The `planWindow` module documents its exports for future PeopleView integration when real availability CRUD exists.

## Impact

- **Affected specs**: `openspec/specs/hello-world-integration/spec.md` (amended — "Move days" card removed from static-content requirement; "Available today" and "Under-staffed" cards remain unchanged).
- **Affected code**: `frontend/src/calendar/CalendarView.vue` (derive day columns from planWindow), `frontend/src/home/HomeView.vue` (derive Move days count from planWindow), new `frontend/src/shared/lib/planWindow.ts` (canonical source of truth), new `frontend/tests/planWindow.test.ts` (unit tests), `frontend/tests/app-routes-render.test.ts` (cross-view render assertion).
- **Unaffected**: `frontend/src/people/PeopleView.vue` (explicitly deferred — its weekday availability strings are not date-derived and will integrate with planWindow when real availability CRUD is built).
- **Non-goals preserved**: No backend endpoints, no runtime settings UI, no new component directories, no PeopleView availability CRUD, no date-library dependency.