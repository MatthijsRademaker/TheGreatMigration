## Why

The shared Daily Schedule board is functionally wired to backend data, pagination, and calendar CRUD, but its current markup and content density drift from the repository's approved design artifacts (`designs/home-page.png`, `designs/components.png`, `docs/design-system-v2.md`). The backlog request is also too vague for a text-only implementation worker: it names "closer to the design" without translating the image-based Daily Schedule header, column, card, and read-only/editable differences into explicit text requirements.

## What Changes

1. **Single compact panel header** — Collapse the current three-layer hierarchy (pagination bar above CardHeader, CardHeader with CardTitle, CardContent controls row) into one compact header row containing the panel title, date-range label, page indicator, prev/next pagination, decorative "View by: Day" label, and conditional "Add task" button.

2. **Simplified task-card anatomy** — Remove the "From backlog" badge from default card display. Simplify assignee rows from individual Avatar+name pairs to a compact text representation (comma-separated initials or names). Retain the priority accent border, priority badge, staffing ratio, and under-staffed indicator. Keep Edit/Delete as conditional ghost buttons on the editable calendar variant only.

3. **Compact day headers** — Render the day label and available-people count in a horizontal layout matching the design artifacts' compact format.

4. **Four-column board tightening** — Tighten the column grid spacing and min-width so the default 4-day slice fits without horizontal scroll on standard desktop viewports, while preserving overflow for non-default day counts.

5. **Pagination owned by component** — Formally accept the current component-layer pagination architecture as the settled contract. The archived pagination spec (task-213bc642) requiring route-layer pagination is superseded. Pagination controls remain inside DailySchedule.vue but visually integrate into the compact header row rather than sitting above the Card shell.

6. **Home read-only contract preserved** — The home dashboard continues to hide Add task controls entirely (header button and per-column placeholders). No visible-but-disabled affordances are introduced.

7. **Updated frontend tests** — Component and route-render tests are updated to assert the new compact header structure, simplified card anatomy, and explicit read-only/editable behavior differences.

8. **Text-only implementation brief** — The task description (this document and the design.md) spells out every Daily Schedule anatomy element in text so a text-only LLM implementer does not need to inspect design PNGs.

## Impact

- **Affected files**: `frontend/src/calendar/DailySchedule.vue` (template and style restructuring), `frontend/tests/calendar/DailySchedule.test.ts` (assertion updates), `frontend/tests/app-routes-render.test.ts` (read-only/editable assertion updates).
- **Not affected**: `HomeView.vue` and `CalendarView.vue` props/emits wiring (unchanged — pagination stays in component), `useDailySchedule` composable, backend API, data contracts.
- **Risk**: Test assertions for `data-slot="card-title"`, specific string matches for Edit/Delete button text, and "Add task" header visibility will change and require coordinated updates.
- **Supersedes**: The archived `daily-schedule-pagination` spec (openspec/changes/archive/2026-06-17-task-213bc642-2291-4127-9e9c-bb0d83456413/specs/daily-schedule-pagination/spec.md) requirement that pagination controls render at the route layer — component-layer pagination is the settled architecture.