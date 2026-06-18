# Design: Centralize Homepage Pagination in the Header

## Context

The homepage (`/`) is the primary dashboard surface. It composes People availability, Daily schedule, Task backlog (read-only), KPI cards, and Move notes. Today, pagination/navigation for the timeline is fragmented:

- `useDailySchedule` manages its own `page` ref and `daysPerPage: 4`.
- `usePeopleAvailability` manages its own `page` ref and `daysPerPage: 7`.
- `DailySchedule.vue` renders a pagination bar with Previous/Next buttons and "Page X of Y" text inside its card header.
- `AppShell.vue`'s sticky header displays a static `formattedRange` (the full planning window), not the current page range.

The provided design mock (`designs/home-page.png`, `designs/components.png`, `designs/design-system.png`) shows a header-level timeline control with a date range display (e.g., "2 Jul – 5 Jul, 2024"), a Today button, and Previous/Next chevrons, with no duplicate pagination inside the dashboard cards. The design also shows notification and profile affordances in the header that are out of scope for this task.

Both backend dashboard endpoints (`GET /api/dashboard/daily-schedule` and `GET /api/dashboard/people-availability`) already accept explicit `start` and `days` query parameters, enabling frontend-only synchronization without backend changes.

## Goals / Non-Goals

**Goals:**

- Provide one homepage-level date-range navigation control in the header instead of resource-level pagination inside homepage cards.
- Synchronize the homepage People availability and Daily schedule sections to the same 4-day visible window and page movement.
- Align the homepage header UI with the supplied design assets using existing shared UI primitives (Button, lucide-vue icons).
- Extend `useDailySchedule` and `usePeopleAvailability` to accept reactive external page state via `pageRef`/`daysPerPageRef` so a parent composable can drive both data sources without violating Vue composable rules.
- Create a text-only design description artifact at `designs/home-header-timeline.md` for LLM consumption.
- Update automated tests to assert the new centralized header controls and removal of in-card pagination on the homepage.

**Non-Goals:**

- Redesign the task backlog, move notes, or KPI cards.
- Change backend dashboard contracts — current endpoints already support explicit `start`/`days`.
- Remove or change pagination on the dedicated `/calendar` or `/people` routes.
- Make the header toolbar appear on non-homepage routes.
- Add notification chip, profile avatar, or other header affordances shown in the design mock beyond the timeline toolbar.
- Introduce speculative timeline features not evidenced by the task or designs.

## Decisions

### Decision 1: Reactive external page ownership via `pageRef` / `daysPerPageRef`

**Chosen: Add optional `pageRef: Ref<number>` and `daysPerPageRef: Ref<number>` parameters to `useDailySchedule` and `usePeopleAvailability`.**

When provided, these refs replace the composable's internal `page`/`daysPerPage` refs, enabling the parent composable (`useHomePagination`) to drive pagination reactively. When omitted, internal refs are used (backward compatible).

**Rationale:**

- Both composables currently accept `options.start` as a static string capture. Passing a computed `start` from a parent composable would not trigger query refetches on page change because Vue composables are called once during setup.
- Adding `pageRef`/`daysPerPageRef` is the minimal backward-compatible change that enables true shared ownership without restructuring component hierarchies.
- The `startParam` computed in both composables already derives the API `start` from `page` and `daysPerPage`; making those values external refs means the parent controls the page and the child automatically fetches the right window.

**Alternatives considered:**

- **Pass `start` and `daysPerPage` as static values from HomeView**: Rejected because `options.start` is captured at call time and does not react to parent page changes without re-initializing the composable (which violates Vue composable rules requiring unconditional invocation during setup).
- **Use `provide`/`inject` for shared page state**: More implicit, couples child composables to a specific injection key, and is unnecessary when explicit ref parameters achieve the same goal with less magic.

### Decision 2: `hidePagination` prop on DailySchedule

**Chosen: Add an explicit `hidePagination: boolean` prop (default `false`) to `DailySchedule.vue`.**

When `true`, the pagination bar (dateRangeLabel, "Page X of Y", Previous, Next) is suppressed regardless of `page`/`totalPages` values. HomeView sets `hidePagination` on the homepage.

**Rationale:**

- The current `hasPagination` computed hides pagination when `page <= 0 || totalPages <= 0`. This is fragile — it relies on the caller knowing to pass zero values, and a non-zero page could accidentally re-enable pagination.
- An explicit `hidePagination` prop is self-documenting and cannot be accidentally overridden by prop values.
- Default `false` means zero impact on other routes (`/calendar`) that use DailySchedule with its own pagination.

### Decision 3: Route-aware header in AppShell

**Chosen: Use `useRoute()` in `AppShell.vue` to conditionally render the timeline toolbar when `route.path === '/'`, via a well-named `isHomeRoute` computed.**

The existing simple header (brand, birds image, `formattedRange`) remains on non-homepage routes.

The timeline toolbar is extracted into a sub-component (`AppShellTimelineToolbar`) for isolation.

**Rationale:**

- `useRoute` is already imported and used in `AppSidebar.vue` (same layouts package), establishing the pattern.
- Extracting the toolbar into its own component isolates route-gating logic and prevents AppShell from becoming bloated.
- A `route.path === '/'` check is simple and sufficient given the task scope is explicitly homepage-only.

**Future extension:** If other routes later need the same toolbar, the condition can be extended or the toolbar can be converted to a slot. The isolated sub-component makes this straightforward.

### Decision 4: `goToday` behavior when current date is outside the planning window

**Chosen: Compute the page containing the real current date within the planning window. If today falls outside the planning window entirely, fall back to page 1 (the planning window start).**

The Today button is always enabled — it always navigates to a valid page within the planning window.

**Rationale:**

- Always navigating to a valid page is simpler and never leaves the user stuck.
- Falling back to page 1 when today is outside the window is consistent with the "window has moved past today" interpretation (show the earliest available page).
- A disabled button that requires explanation tooltips is more complex UX for marginal benefit in the narrow case where today truly falls outside the window.

**Alternatives considered:**

- **Disable the button when today is outside the window with a tooltip explaining why**: More user-friendly for the edge case but adds complexity (tooltip rendering, accessibility considerations) without a strong user need.
- **No-op when today is outside the window**: Frustrating — the button appears functional but does nothing.

### Decision 5: Text-only design description location

**Chosen: `designs/home-header-timeline.md`**

Colocated with the existing design images (`designs/home-page.png`, `designs/components.png`, `designs/design-system.png`) so all design artifacts for the homepage header are in one place.

### Decision 6: No changes to PeopleAvailability component pagination

**Confirmed: `PeopleAvailability.vue` has no pagination UI on the homepage today.**

The dossier's claim of "remove duplicated pagination bars/buttons inside the content area" for PeopleAvailability is incorrect — the component renders a grid of people × days with no pagination controls. Pagination for people data exists only on the dedicated `/people` route (`PeopleView.vue`), which is out of scope.

The actual synchronization work is: ensure both `useDailySchedule` and `usePeopleAvailability` query the same 4-day window driven by the shared `useHomePagination` composable.

## Risks / Trade-offs

| Risk | Severity | Mitigation |
|---|---|---|
| **Route-gating coupling in AppShell**: Hardcoded `/` check couples a shared layout component to route knowledge. | Low | Extract toolbar into isolated sub-component; use well-named `isHomeRoute` computed. Future routes can extend the condition or convert to a slot without restructuring AppShell. |
| **Reactivity gap in existing composables**: If `pageRef`/`daysPerPageRef` integration is missed, the homepage will render stale data independent of header navigation. | High (mitigated by design) | Both composables are explicitly modified to accept reactive refs. The `hidePagination` prop eliminates accidental double-navigation rendering. |
| **Test assertion breakage**: `app-routes-render.test.ts` currently asserts "Page 1 of 10", "Previous", "Next" in home route content. These strings move to the header. | Medium | Update test assertions to target the header container and verify the same strings appear in the new location. Preserve non-homepage route assertions unchanged. |
| **AppShell test fragility**: `app-shell-layout.test.ts` mocks `usePlanningWindow` and checks simple header structure. New toolbar rendering must be testable with the same mock pattern. | Low | Run the test suite after changes to catch regressions; update mock to include `useRoute` if needed. |
| **CalendarView shares DailySchedule**: Any change to DailySchedule's pagination props must not break `/calendar`. | Low | `hidePagination` defaults to `false`; CalendarView does not pass it. The existing `hasPagination` guard is preserved internally for non-hidden cases. |
| **Planning window reset watch conflict**: Both composables have a `watch` on planning window changes that resets `page` to 1. When `pageRef` is externally provided, that watch should not fire or should fire against the external ref. | Medium | When `pageRef` is provided, disable or skip the internal planning-window reset watch. When omitted (existing callers), the watch continues as before. |

## Traceability

All requirements in the spec trace to:

- **Task input**: "Make pagination global in the header on the homepage" and "make the description of the desired design ingestable for a text only llm."
- **Exploration dossier**: problem framing, goals, non-goals, acceptance criteria, affected areas
- **Architect recommendation (round 1)**: `pageRef`/`daysPerPageRef` reactivity pattern, `goToday` fallback behavior, `useHomeTimelinePagination` composable
- **Lead-dev recommendation (round 1)**: `hidePagination` prop, `useHomePagination` composable API, `AppShell` route-gating with `isHomeRoute` computed, design doc at `designs/home-header-timeline.md`
- **Reviewer recommendation (round 1)**: daysPerPage sync enforcement, fragility analysis of `hasPagination` guard, PeopleAvailability pagination scope correction
- **Design assets**: `designs/home-page.png`, `designs/components.png`, `designs/design-system.png`
- **Codebase evidence**: `AppShell.vue`, `HomeView.vue`, `DailySchedule.vue`, `useDailySchedule.ts`, `usePeopleAvailability.ts`, `usePlanningWindow.ts`, `PeopleAvailability.vue`, test files, backend API files