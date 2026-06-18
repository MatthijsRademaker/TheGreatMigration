# Proposal: Centralize Homepage Pagination in the Header

## Why

The homepage currently distributes timeline pagination across individual dashboard resources:

- **DailySchedule** renders its own `Previous` / `Next` bar inside the card header with "Page X of Y" text.
- **HomeView** calls `usePeopleAvailability()` and `useDailySchedule()` independently with no shared pagination state.
- **AppShell** header shows only a static planning-window range (`formattedRange`), not a page-aware date control.
- The two homepage data sources default to different window sizes: `useDailySchedule` defaults to 4 days, `usePeopleAvailability` defaults to 7 days, so the dashboard cannot present a coherent date window.

This creates confusion — multiple pagination controls on one page, unsynchronized data windows, and no central timeline navigation. The provided design mock shows a unified header-level date-range control with Previous/Next chevrons and a Today button, not duplicate in-card pagination.

## What Changes

- **New shared composable**: `useHomePagination` owns homepage page state, exposes `start`, `rangeLabel`, `page`, `totalPages`, `goPrev`, `goNext`, and `goToday`, and enforces `daysPerPage: 4` for both homepage data sections.
- **Reactive composable upgrade**: `useDailySchedule` and `usePeopleAvailability` gain optional `pageRef` and `daysPerPageRef` parameters accepting Vue `Ref<number>` so a parent composable can drive pagination reactively without re-initializing child composables.
- **AppShell header**: Renders the design-aligned timeline toolbar (date range label, Today button, Previous/Next chevrons) on the `/` route only, using `useRoute()` for route awareness.
- **DailySchedule**: Gains an explicit `hidePagination` boolean prop to suppress in-card pagination when the header owns navigation — used only on the homepage.
- **HomeView**: Wires shared `useHomePagination` state into both `useDailySchedule` and `usePeopleAvailability`, removing passthrough of `page`/`totalPages`/`goToPrevPage`/`goToNextPage` to DailySchedule.
- **Text-only design description**: `designs/home-header-timeline.md` describes the header layout, controls, spacing, and behavior in plain language for non-vision LLM consumption.
- **Test updates**: `app-routes-render.test.ts` asserts centralized header pagination controls on `/` instead of in-card pagination text; `app-shell-layout.test.ts` verifies new header toolbar presence on the home route.

## Impact

| Layer | Impact |
|---|---|
| **Frontend composables** | `useDailySchedule.ts` and `usePeopleAvailability.ts`: add optional `pageRef`/`daysPerPageRef` params. Backward compatible — existing callers without refs behave identically. |
| **Frontend components** | `AppShell.vue`: route-aware header toolbar rendering. `DailySchedule.vue`: new `hidePagination` prop (default false). `HomeView.vue`: shared pagination wiring replaces per-composable page state. |
| **New file** | `frontend/src/shared/composables/useHomePagination.ts` |
| **New file** | `designs/home-header-timeline.md` |
| **Tests** | `app-routes-render.test.ts`: update home route assertions. `app-shell-layout.test.ts`: add header toolbar assertions. |
| **Backend** | No changes. Both dashboard endpoints already accept explicit `start` and `days` query parameters. |
| **Other routes** | No impact. `/calendar`, `/people`, `/tasks`, `/rooms`, `/settings` retain existing pagination and header behavior. |