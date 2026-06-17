## 1. Text-Only Implementation Brief

- [x] 1.1 Confirm the refined `proposal.md` and `design.md` spell out the Daily Schedule panel header, day-header, task-card, add-placeholder, and read-only-vs-editable anatomy in text — no image inspection required.

## 2. Daily Schedule Component Visual Refinement

- [x] 2.1 Collapse the three-layer header hierarchy into a single compact panel header row: merge pagination bar + CardHeader/CardTitle + CardContent controls row into one header area with title left, controls right.
- [x] 2.2 Integrate pagination controls (date-range label, page indicator "Page X of Y", prev/next buttons) into the compact header row alongside "View by: Day" and conditional "Add task".
- [x] 2.3 Simplify day headers from vertical stack (`h3` label + `p` count) to a compact horizontal layout.
- [x] 2.4 Remove the "From backlog" badge from task-card rendering.
- [x] 2.5 Simplify assignee display from individual Avatar+name rows to comma-separated initials or names.
- [x] 2.6 Retain priority accent border-left, priority badge, staffing ratio (`X / Y`), under-staffed indicator, and conditional Edit/Delete ghost buttons.
- [x] 2.7 Ensure Add task placeholder (`+ Add task` in dashed-border card) is visible only when `!readOnly`.
- [x] 2.8 Tighten column grid spacing so 4 columns fit without horizontal scroll on standard desktop, preserving overflow for non-default day counts.
- [x] 2.9 Keep "View by: Day" as a static/decorative label — no new behavior wired.

## 3. Preserve Existing Contracts

- [x] 3.1 Keep all existing component props (`days`, `readOnly`, `page`, `totalPages`, `dateRangeLabel`) and emits (`add-task`, `edit-task`, `delete-task`, `prev-page`, `next-page`) unchanged.
- [x] 3.2 Do not modify `useDailySchedule` composable, its data adaptation logic, or the API query contract.
- [x] 3.3 Do not change `HomeView.vue` or `CalendarView.vue` props/emits wiring — pagination stays in the component.
- [x] 3.4 Home read-only continues to hide Add task controls entirely (header button and per-column placeholders).

## 4. Test Updates

- [x] 4.1 Update `DailySchedule.test.ts`: replace assertions for `data-slot="card-title"` and stacked-header structure with assertions for the compact header row, compact day headers, simplified card anatomy, and pagination-in-header integration.
- [x] 4.2 Update `app-routes-render.test.ts`: adjust read-only assertions for `/` (no Add task, no Edit, no Delete) and editable assertions for `/calendar` (Add task present, Edit/Delete present on cards).
- [x] 4.3 Verify all existing test scenarios (day labels, availability counts, task titles, priority variants, staffing counts, under-staffed indicator, pagination bar, read-only hiding) pass with the refined structure.
- [x] 4.4 Add new test assertions for the compact header anatomy (single row containing title + controls) and compact assignee display (comma-separated format).
