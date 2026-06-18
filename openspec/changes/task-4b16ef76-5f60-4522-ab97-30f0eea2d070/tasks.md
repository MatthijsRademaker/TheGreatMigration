# Tasks: Centralize Homepage Pagination in the Header

## 1. Shared Homepage Pagination Composable

- [ ] 1.1 Create `frontend/src/shared/composables/useHomePagination.ts` that wraps `usePlanningWindow` and derives a 4-day page window, exposing `page` (Ref<number>), `daysPerPage` (Ref<number>), `totalPages`, `start` (computed ISO date string), `rangeLabel`, `goPrev`, `goNext`, and `goToday` (computes page containing today's date, falls back to page 1 if outside window).

## 2. Reactive External State in Data Composables

- [ ] 2.1 Add optional `pageRef?: Ref<number>` and `daysPerPageRef?: Ref<number>` parameters to `UseDailyScheduleOptions` in `useDailySchedule.ts`. When provided, use them directly instead of creating internal `ref()` calls. Skip the planning-window reset watch when `pageRef` is provided (external owner handles resets).
- [ ] 2.2 Add optional `pageRef?: Ref<number>` and `daysPerPageRef?: Ref<number>` parameters to `UsePeopleAvailabilityOptions` in `usePeopleAvailability.ts`. When provided, use them directly instead of creating internal `ref()` calls. Skip the planning-window reset watch when `pageRef` is provided.

## 3. DailySchedule `hidePagination` Prop

- [ ] 3.1 Add `hidePagination?: boolean` prop to `DailySchedule.vue` (default `false`). When `true`, suppress the pagination bar (dateRangeLabel, "Page X of Y", Previous, Next) in the card header regardless of `page`/`totalPages` values.

## 4. AppShell Header Timeline Toolbar

- [ ] 4.1 Create `AppShellTimelineToolbar.vue` sub-component (in `frontend/src/shared/layout/app-shell/`) that renders: a formatted date range label, a Today button, a Previous chevron button (ChevronLeft), and a Next chevron button (ChevronRight). Accepts props for `rangeLabel`, `canGoPrev`, `canGoNext` and emits `today`, `prev`, `next`.
- [ ] 4.2 Update `AppShell.vue` to import `useRoute`, compute `isHomeRoute`, import `useHomePagination` (conditionally or via a wrapper), and render `AppShellTimelineToolbar` in the header when `isHomeRoute` is `true`. On non-homepage routes, render the existing simple header (brand, birds image, `formattedRange`).

## 5. HomeView Shared Pagination Wiring

- [ ] 5.1 Update `HomeView.vue` to call `useHomePagination()` and pass `pageRef`/`daysPerPageRef` to both `useDailySchedule` and `usePeopleAvailability`.
- [ ] 5.2 Remove `page`, `totalPages`, `goToPrevPage`, `goToNextPage`, and `dateRangeLabel` passthrough from HomeView to `<DailySchedule>`. Add `:hide-pagination="true"` to the DailySchedule component on the homepage.

## 6. Text-Only Design Description

- [ ] 6.1 Create `designs/home-header-timeline.md` with a plain-text description of the homepage header layout: positioning (sticky top bar, left-to-right flex flow), control grouping (brand → spacer → date range label → Today button → Previous chevron → Next chevron), spacing details, button styling (outline variant, small size), date range format (e.g., "2 Jul – 5 Jul, 2024"), Today button behavior, and chevron icon choice (ChevronLeft/ChevronRight from lucide-vue).

## 7. Test Updates

- [ ] 7.1 Update `frontend/tests/app-routes-render.test.ts` home route assertions: remove expectations for in-card "Page 1 of 10"/"Previous"/"Next" inside DailySchedule content. Add assertions verifying the new header toolbar renders on `/` with date range label, Today button, and navigation chevrons. Preserve `/calendar` route pagination assertions unchanged.
- [ ] 7.2 Update `frontend/tests/app-shell-layout.test.ts` to mock `useRoute` and verify the timeline toolbar renders when route is `/` and the simple header renders on non-homepage routes. Update or extend the `usePlanningWindow` mock if needed.
- [ ] 7.3 Run the full test suite and fix any regressions in existing tests (e.g., CalendarView tests, PeopleView tests) that may be affected by DailySchedule prop changes.