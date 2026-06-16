## Why

The `/calendar` route still renders the planning-window placeholder `Schedule board foundation`, while the product already has a design-backed Daily Schedule board and an existing backend schedule contract. This change turns the supplied design into a text-only implementation brief for a frontend-presentational component so a worker can build the board accurately without needing image interpretation or backend wiring.

## What Changes

- Add `frontend/src/calendar/DailySchedule.vue` as a presentational Vue component under the calendar feature.
- Define local TypeScript interfaces and demo defaults that mirror the backend daily-schedule field names (`date`, `label`, `availablePeopleCount`, `tasks`, `assignedPeople`, `peopleNeeded`, `assignedCount`, `staffingStatus`) while keeping the first slice frontend-only.
- Render the board anatomy from text alone:
  - outer `Card` panel titled `Daily Schedule`
  - static header controls `View by: Day` and `Add task`
  - horizontally scrollable day-column row using `overflow-x-auto` around a flex container, with each column kept readable via `min-w-[280px]` and `shrink-0`
  - four demo day columns exactly matching the design slice: `2 Jul (Tue)`, `3 Jul (Wed)`, `4 Jul (Thu)`, and `5 Jul (Fri)` with `6 available`, `7 available`, `7 available`, and `5 available`
  - task cards rendered as plain `div` items, not nested `Card`s, with `rounded-lg border bg-card shadow-sm p-3` and a `border-l-4` priority accent using semantic tokens for high/medium/low priorities
  - task-card content including title, priority badge, assigned people metadata with compact avatars and names, and staffing count text such as `2 / 2`, `1 / 1`, and an under-staffed example through the same contract
  - dashed `Add task` placeholder control at the bottom of each day column
- Use deterministic demo task data from the design images for visual fidelity: `Painting hall`, `Steam walls`, `Clean up`, `Sanding`, `Bedroom painting`, `Touch up woodwork`, `Living room finishing`, `2nd floor walls`, `Kitchen painting`, and `Final clean`; preserve backend field names and priority/staffing vocabulary so later API wiring stays additive.
- Replace the current `frontend/src/calendar/CalendarView.vue` placeholder, removing its `usePlanningWindow()` dependency, loading/error branches, and `plan-day-column` grid.
- Add an SSR component test for `DailySchedule.vue` and update `frontend/tests/app-routes-render.test.ts` so `/calendar` asserts the new board content instead of the removed placeholder.

## Impact

- `/calendar` becomes a design-backed preview of the Daily Schedule board instead of a planning-window placeholder.
- Future backend integration can pass real `GET /api/dashboard/daily-schedule` data through the local prop contract without rewriting the component shape.
- The change stays bounded to frontend presentation: no backend handler changes, no OpenAPI/client regeneration, no Pinia Colada query wiring, no CRUD or drag/drop behavior, and no home-dashboard replacement in this slice.
- Verification remains frontend-focused through SSR tests plus `scripts/precommit-run`.
