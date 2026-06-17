## Context

The Daily Schedule component (`frontend/src/calendar/DailySchedule.vue`) is the shared presentational board for viewing and managing scheduled task cards across a 4-day window. It is consumed by two routes:

- **HomeView** (`/`): Read-only — passes `read-only` prop plus pagination props/emits.
- **CalendarView** (`/calendar`): Editable — passes pagination props/emits plus `@add-task`, `@edit-task`, `@delete-task` handlers.

Both routes use the `useDailySchedule` composable for backend-driven data, pagination state, and date-range labeling.

The component currently renders a three-layer visual hierarchy (pagination bar above CardHeader, CardHeader with CardTitle, CardContent with header-controls row) that diverges from the approved design artifacts. Task cards are denser than designed, showing avatars, a "From backlog" badge, and explicit Edit/Delete buttons. Day headers use a vertical stack instead of the design's compact horizontal layout.

The design system (`docs/design-system-v2.md`) defines the calm green visual language, spacing rhythm, badge semantics, and card styling that govern this panel. The historical design brief (`task-b7276a6a/design.md`) already translated the Daily Schedule design images into text.

## Goals / Non-Goals

### Goals

1. Restructure DailySchedule.vue to render a single compact panel header matching the design artifacts.
2. Simplify task cards to match design density: title, priority pill, compact assignee text, staffing ratio.
3. Tighten the four-column board layout for standard desktop viewports.
4. Formally document component-layer pagination as the settled architecture.
5. Update frontend tests to assert the refined structure.
6. Produce a self-contained text-only implementation brief.

### Non-Goals

- Redesigning the backend daily-schedule API, planning-window semantics, or schedule-card persistence.
- Adding new scheduling capabilities (drag-and-drop, filters, search, alternate board modes, URL-synced page state).
- Reworking unrelated dashboard panels, global layout, or the broader design system.
- Silently changing product behavior around read-only vs. editable actions.
- Re-architecting pagination to the route layer.
- Modifying the useDailySchedule composable, its data adaptation logic, or the API query contract.

## Decisions

### 1. Pagination stays in DailySchedule.vue (component-layer)

The current component accepts `page`, `totalPages`, `dateRangeLabel` props and emits `prev-page`/`next-page`. This architecture is the settled contract, implemented in a prior task. The archived `daily-schedule-pagination` spec requiring route-layer pagination is formally superseded by this decision. Moving pagination to routes would require refactoring 4+ files and is outside visual-alignment scope.

### 2. Single compact panel header row

The pagination bar (currently above CardHeader), panel title (currently in CardHeader), and header controls (currently inside CardContent) are merged into one compact row. Layout: title on the left; date-range label, page indicator, prev/next buttons, "View by: Day" label, and conditional "Add task" button on the right.

### 3. Simplified task-card anatomy

Remove the "From backlog" badge from default card rendering. Simplify assignee display from individual Avatar+name rows to comma-separated initials or names. Retain: priority accent border-left, priority badge, staffing ratio text (`X / Y`), under-staffed indicator ("— needs help"). Edit/Delete buttons stay as conditional ghost buttons on the editable variant (`v-if="!readOnly"`).

### 4. Home read-only hides Add task entirely

The home dashboard continues to suppress all Add task controls (header button and per-column dashed placeholders) via the `readOnly` prop. No visible-but-disabled affordances are introduced. This matches the dossier's initial acceptance criteria and both architect and lead-dev recommendations.

### 5. Compact day headers

Each day column header renders the date label (e.g., "1 Aug (Sat)") and available-people count (e.g., "3 available") in a horizontal layout rather than the current vertical stack, matching the design artifacts' compact format.

### 6. Four-column board with tightened grid

Preserve the scrollable `flex` layout with `min-w-[280px]` per column, but tighten the spacing so the default 4-day slice fits without horizontal scroll on standard desktop viewports. Overflow remains available for non-default day counts.

### 7. "View by: Day" stays as static label

Keep the "View by: Day" control as a decorative/static label (current behavior). It is present in the design artifacts and has no wired behavior. No dropdown or modal is introduced.

### 8. Design-system primitives preserved

All markup continues to use existing shadcn-vue primitives (Card, Button, Badge, Avatar) and Design System v2 semantic tokens. No new low-level components are created.

## Conflict Resolution

**Pagination ownership (Architect vs. Lead-dev):** The architect recommended keeping component-layer pagination as the settled architecture, while the lead-dev recommended extracting pagination to the route layer per the archived spec. Resolution: component-layer pagination is accepted. The archived spec was superseded by the current implementation in a prior task. Re-architecting pagination now would consume scope that should go toward visual alignment. The lead-dev's blocker was classified as non-blocking in the room. The pagination bar is visually integrated into the compact header row rather than positioned above the Card shell.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Test breakage from markup changes | Medium — current tests assert specific data-slot strings, button text, and header control presence | Update all affected assertions in DailySchedule.test.ts and app-routes-render.test.ts as part of the implementation |
| Removing "From backlog" badge loses card-source distinction | Low — users may not be able to tell which cards are linked to backlog tasks | If needed, surface the backlog relationship as a subtle tooltip or icon in a follow-up; removing the prominent badge matches the design contract |
| Simplified assignee display reduces information density | Low — full names are replaced by compact initials/names | The design artifacts show compact assignee presentation; full detail remains available in the edit modal |
| Pagination-in-component divergence from archived spec | Low — future maintainers may discover the archived spec and question the architecture | This decision is explicitly documented in the design.md; the archived spec is formally superseded |

## Traceability

- **Task**: `29534381-9c4d-4aa3-9fb0-4ecff9210a15`
- **Dossier**: `2026-06-17T19:20:29.406Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Non-blocking findings**: `round-1-swarm-lead-dev-blocker-1`, `round-1-swarm-reviewer-blocker-1`, `round-1-swarm-reviewer-blocker-2`, `round-1-swarm-reviewer-blocker-3`, `round-1-swarm-reviewer-blocker-4`
- **Reference**: Historical design brief `task-b7276a6a/design.md`, design system `docs/design-system-v2.md`, canonical component spec `openspec/specs/daily-schedule-component/spec.md`
