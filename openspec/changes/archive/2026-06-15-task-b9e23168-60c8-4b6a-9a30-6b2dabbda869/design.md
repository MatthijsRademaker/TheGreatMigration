# Design: Planning window shared module

## Context

The app's feature views (`CalendarView`, `HomeView`, `PeopleView`) each hard-code date-related literals independently. The C4 architecture model confirms `scheduleData` is a shared planning data element that calendar and dashboard both reference, but no shared source of truth exists in code. The task requests a "configurable constant which the rest of the application can use" so that changing the move-plan period propagates consistently. The `backend/main.go` exposes only `/api/hello`, confirming this is a frontend-only change.

## Goals

- Create a single TypeScript module that owns the canonical move-plan start and end dates and exposes reusable derived values.
- Refactor existing date-dependent UI (`CalendarView`, `HomeView`) to consume the shared source instead of duplicating literals.
- Keep the abstraction usable by future calendar and availability CRUD without requiring those features now.
- Amend the hello-world-integration spec so it no longer mandates static content for the now-dynamic "Move days" card.
- Add focused tests so a range change fails loudly instead of drifting silently.

## Non-Goals

- Add a runtime settings screen or admin UI for editing the date range.
- Add backend persistence, APIs, or cross-device sync for planning dates.
- Build missing calendar or availability CRUD beyond today's placeholder views.
- Convert PeopleView placeholder weekday strings to actual dates (deferred).
- Introduce a date-library dependency (date-fns, dayjs, etc.).
- Rework unrelated task, routing, or design-system behavior.

## Decisions

### D1: Module location — `frontend/src/shared/lib/planWindow.ts`

The existing `frontend/src/shared/lib/` directory already houses `utils.ts` as a utility/data module. The `composables/` directory is reserved for Vue lifecycle composables (currently only `.gitkeep`). Since `planWindow` exports immutable constants and pure derivation functions — not reactive state — it belongs in `lib/`. If future work needs reactive plan-window state, a composable can wrap this module at that time.

### D2: Export surface — ISO string constants plus pure helper functions

The module exports:
- `PLAN_WINDOW_START: string` — ISO date string (default `"2026-07-05"`)
- `PLAN_WINDOW_END: string` — ISO date string (default `"2026-08-13"`)
- `planWindowDays: {date: Date, label: string, dateString: string}[]` — derived array of inclusive day objects
- `planWindowDayCount: number` — derived inclusive day count
- `formatPlanDayLabel(date: Date): string` — formats a Date as e.g. `"Sat 5 Jul"`

ISO strings are the canonical format because they are serializable, testable, and unambiguous. `Date` objects in `planWindowDays` entries enable direct date arithmetic for consumers that need it. The label formatter uses `en-US` locale with abbreviated weekday, day, and abbreviated month to produce compact column headers.

### D3: Inclusive date semantics

The planning window uses inclusive start/end semantics: both `PLAN_WINDOW_START` and `PLAN_WINDOW_END` are included in the derived day list. This is the most intuitive interpretation of "5 juli until 13 aug" and is enforced in both `planWindowDays` generation and `planWindowDayCount` computation so consumers cannot accidentally diverge.

### D4: Default date range — 2026-07-05 to 2026-08-13

The task example explicitly specifies "5 juli until 13 aug." The current runtime date is 2026, so the default year is 2026. This produces an inclusive 40-day window (27 days in July + 13 in August), which exercises non-trivial month-boundary crossing in the derivation logic.

### D5: CalendarView — remove focus/helpers, use date-derived placeholders

The current `moveDays` array carries three fields per day: `day` (weekday label), `focus` (narrative like "Packing"), and `helpers` (staffing like "4 people"). The `focus` and `helpers` fields are narrative demo data that cannot scale with an arbitrary planning window — there is no data source to generate meaningful focus/helpers strings for 40 arbitrary days. Each column instead displays the formatted date label and a generic description (e.g., the card description "Planned move day"). Real task blocks and staffing data belong to future calendar CRUD work.

### D6: CalendarView grid layout — switch from fixed columns to responsive wrapping

The current `md:grid-cols-5` assumes exactly five columns. With a 40-day window, a 5-column grid would produce 8 rows on desktop. The grid switches to `grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7` so the layout adapts gracefully to the window size without horizontal scrolling or excessive vertical stacking. The `min-h-36` on each cell is preserved so cells remain visually consistent.

### D7: HomeView — replace static `"5"` with `planWindowDayCount`

The `summaries` array's "Move days" entry changes `value` from the hard-coded string `"5"` to `String(planWindowDayCount)` (or a reactive ref wrapping it if the summary array is defined in `<script setup>` scope). The label, description, and icon remain unchanged.

### D8: PeopleView — explicitly deferred

The weekday strings in `PeopleView.vue` (`Mon, Wed, Fri`, `Tue, Wed`, etc.) represent helper availability patterns, not specific calendar dates. They are not date-derived and converting them to in-range dates would be purely cosmetic while expanding scope into building real availability data structures — which the non-goals explicitly forbid. The `planWindow` module documents its exports in a header comment for future PeopleView integration when real availability CRUD exists. A code comment is added to `PeopleView.vue` noting the future integration point.

### D9: Spec treatment — amend hello-world-integration

The `openspec/specs/hello-world-integration/spec.md` "Other HomeView sections are preserved" scenario is amended to exclude the "Move days" card from the static-content requirement. The scenario now reads that the *remaining two* summary cards ("Available today" and "Under-staffed") plus "Today's plan" and "Move notes" continue to render with their existing static content. The "Move days" card is documented as deriving from the planning window. A corresponding spec delta is provided in this change's `specs/hello-world-integration/spec.md`.

### D10: Testing strategy

- **Unit tests** (`frontend/tests/planWindow.test.ts`): Validate that `PLAN_WINDOW_START` and `PLAN_WINDOW_END` return the expected ISO strings, `planWindowDayCount` matches the inclusive day count (40 for the default range), `planWindowDays` length equals `planWindowDayCount`, `planWindowDays[0].dateString` equals `PLAN_WINDOW_START`, `planWindowDays[last].dateString` equals `PLAN_WINDOW_END`, and `formatPlanDayLabel` produces correctly formatted output (e.g., `"Sun 5 Jul"` for 2026-07-05 which is a Sunday).
- **Render test** (`frontend/tests/app-routes-render.test.ts`): Extend the existing CalendarView route test to assert that the rendered HTML contains the correct number of day-column elements matching `planWindowDayCount`.

## Risks

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| CalendarView grid overflow with large windows | Medium | D6: responsive grid columns adapt to window size; this is acceptable placeholder behavior until real calendar UX is built. |
| Off-by-one errors in inclusive range generation across month/year boundaries | Low | D3: inclusive semantics enforced in both `planWindowDays` and `planWindowDayCount`; unit tests cover month-boundary crossing. |
| HomeView SSR/client mismatch if date computation uses runtime-only APIs | Low | D2: `formatPlanDayLabel` uses `Intl.DateTimeFormat` with explicit `en-US` locale, but `planWindowDayCount` is a pure number — no locale dependency. The render test uses SSR (`renderToString`) and will catch any mismatch. |
| Tight coupling if raw constants are consumed directly everywhere | Low | D2: consumers import named exports from the module; if the API shape changes, TypeScript compilation catches breakages. A composable wrapper can be added later without changing consumer code. |

## Traceability

- **Task**: `b9e23168-60c8-4b6a-9a30-6b2dabbda869` — "Date range planning configurable"
- **Dossier**: `2026-06-15T17:28:58.125Z` — problem framing, affected areas, acceptance criteria
- **Architect decision** (`round:1:agent:swarm-architect`): Core approach validated; PeopleView conversion overridden by lead-dev deferral
- **Lead dev decision** (`round:1:agent:swarm-lead-dev`): Module location (`lib/`), API shape, PeopleView deferral, testing strategy
- **Reviewer decision** (`round:1:agent:swarm-reviewer`): Identified 4 blocking gaps, all resolved in this design
- **OpenSpec spec**: `openspec/specs/hello-world-integration/spec.md` — amended per D9
- **C4 model**: `.devagent/architecture/workspace.c4` — confirms `scheduleData` as shared planning data element
