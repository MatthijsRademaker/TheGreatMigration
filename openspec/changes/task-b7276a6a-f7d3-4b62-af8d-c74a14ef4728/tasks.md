## 1. Build the presentational board component

- [x] Create `frontend/src/calendar/DailySchedule.vue`.
- [x] Define local TypeScript interfaces and typed props for schedule days, task cards, and assigned people using backend-compatible field names.
- [x] Provide deterministic demo defaults for the four-day design slice (`2 Jul (Tue)` through `5 Jul (Fri)`), the visible availability counts, and representative high/medium/low tasks plus staffing metadata.
- [x] Render a `Daily Schedule` panel with static `View by: Day` and `Add task` header controls, four horizontally scrollable day columns, priority-accented task cards, compact assignee metadata, staffing counts, and dashed per-column `Add task` placeholders.
- [x] Use existing `Card`, `Badge`, `Button`, and `Avatar` primitives plus semantic tokens only; render task cards as plain `div` items rather than nested `Card` components.
- [x] Preserve small-screen readability with `overflow-x-auto`, a horizontal flex row, and day columns that keep a minimum readable width.

## 2. Integrate the component into the calendar route

- [x] Replace the placeholder content in `frontend/src/calendar/CalendarView.vue` with `DailySchedule`.
- [x] Remove the `usePlanningWindow()` import, loading/error branches, and `plan-day-column` grid from `CalendarView.vue`.
- [x] Keep the existing `/calendar` route metadata and app-shell/sidebar behavior unchanged.
- [x] Keep the route presentational only: no daily-schedule API query, no generated-client imports, and no backend wiring.

## 3. Add focused frontend verification

- [x] Add `frontend/tests/calendar/DailySchedule.test.ts` using the existing SSR `renderToString` pattern.
- [x] Assert the component renders the panel title, header controls, four day labels, availability counts, representative task titles, priority labels, assignee metadata, staffing counts, and `Add task` placeholders.
- [x] Cover custom-prop rendering so the typed contract can replace demo data without template changes.
- [x] Update `frontend/tests/app-routes-render.test.ts` so `/calendar` asserts the new Daily Schedule content instead of `Schedule board foundation` and no longer counts `data-testid="plan-day-column"`.

## 4. Run repository verification

- [x] Run `scripts/precommit-run`.
