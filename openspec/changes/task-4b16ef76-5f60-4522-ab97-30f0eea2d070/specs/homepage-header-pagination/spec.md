# Spec: Homepage Header Pagination

## ADDED Requirements

### Requirement: Homepage pagination composable SHALL own shared page state

A new `useHomePagination` composable (`frontend/src/shared/composables/useHomePagination.ts`) SHALL wrap `usePlanningWindow` and derive a 4-day page window. It SHALL expose:

- `page: Ref<number>` — reactive page ref (1-indexed)
- `daysPerPage: Ref<number>` — reactive days-per-page ref (initialized to 4)
- `totalPages: ComputedRef<number>` — total pages derived from planning window days / daysPerPage
- `start: ComputedRef<string>` — ISO date string for the current page start
- `rangeLabel: ComputedRef<string>` — human-readable range like "2 Jul – 5 Jul, 2024"
- `goPrev(): void` — decrement page if > 1
- `goNext(): void` — increment page if < totalPages
- `goToday(): void` — compute the page containing the real current date within the planning window; if today falls outside the planning window, set page to 1

#### Scenario: Homepage pagination derives from planning window

- **WHEN** the planning window is loaded with 40 days
- **THEN** `totalPages` resolves to `Math.ceil(40 / 4) = 10`
- **AND** `page` starts at 1, `daysPerPage` is 4
- **AND** `start` returns the planning window start date (the first day)

#### Scenario: Page navigation updates start date

- **WHEN** `goNext()` is called from page 1
- **THEN** `page` becomes 2
- **AND** `start` returns the date 4 days after the planning window start

#### Scenario: goToday computes page containing current date

- **WHEN** the planning window spans 5 Jul – 13 Aug 2026 and the real current date is 10 Jul 2026
- **THEN** `goToday()` sets `page` to the page whose visible window contains 10 Jul 2026

#### Scenario: goToday falls back to page 1 when today is outside window

- **WHEN** the planning window spans 1 Jan – 31 Jan 2027 and the real current date is outside that range
- **THEN** `goToday()` sets `page` to 1

### Requirement: AppShell SHALL render a design-aligned homepage header toolbar

`AppShell.vue` SHALL render a timeline navigation toolbar in the header when the current route is `/`. The toolbar SHALL display:

- A formatted date-range label showing the current 4-day page window (e.g., "2 Jul – 5 Jul, 2024")
- A "Today" button that navigates to the page containing the current date
- A Previous chevron button (disabled on page 1)
- A Next chevron button (disabled on the last page)

On non-homepage routes, the header SHALL render the existing simple layout (brand, birds image, `formattedRange`).

The toolbar SHALL be extracted into a dedicated `AppShellTimelineToolbar` sub-component for isolation.

#### Scenario: Home route header shows timeline toolbar

- **WHEN** the current route is `/`
- **THEN** the AppShell header renders the date range label, Today button, Previous chevron, and Next chevron
- **AND** the static `formattedRange` (planning window range) is replaced by the page-specific date range label

#### Scenario: Non-home route header shows simple layout

- **WHEN** the current route is `/calendar`, `/people`, `/tasks`, `/rooms`, or `/settings`
- **THEN** the AppShell header renders the existing simple layout (brand, birds image, `formattedRange`)
- **AND** the timeline toolbar is not rendered

#### Scenario: Previous chevron is disabled on page 1

- **WHEN** `page` is 1
- **THEN** the Previous chevron button is disabled

#### Scenario: Next chevron is disabled on last page

- **WHEN** `page` equals `totalPages`
- **THEN** the Next chevron button is disabled

### Requirement: Text-only design description SHALL be available to non-vision LLMs

A markdown file at `designs/home-header-timeline.md` SHALL describe the homepage header layout in plain language suitable for text-only LLM consumption. The description SHALL cover:

- Header positioning (sticky top bar, full width, left-to-right flex layout)
- Control grouping and order: brand mark → flexible spacer → date range label → Today button → Previous chevron (ChevronLeft) → Next chevron (ChevronRight)
- Spacing and sizing (h-10 height, gap-2 between controls, px-4 horizontal padding)
- Button styling (outline variant, small size for Today; ghost/small for chevrons)
- Date range format (e.g., "2 Jul – 5 Jul, 2024" — day, abbreviated month, year for both start and end)
- Today button behavior (navigates to page containing current date; falls back to page 1 if today outside planning window)
- Chevron behavior (disabled state on boundary pages)

#### Scenario: Design description is readable without image access

- **WHEN** a text-only LLM reads `designs/home-header-timeline.md`
- **THEN** it can understand the header layout, control grouping, spacing, styling, and interactive behavior without inspecting any PNG files

## MODIFIED Requirements

### Requirement: DailySchedule SHALL support explicit pagination suppression

`DailySchedule.vue` SHALL accept a new `hidePagination?: boolean` prop (default `false`). When `true`, the pagination bar (dateRangeLabel, "Page X of Y" text, Previous button, Next button) SHALL be suppressed regardless of `page` and `totalPages` prop values. When `false` or absent, the existing `hasPagination` guard (`page > 0 && totalPages > 0`) SHALL continue to gate the pagination bar as before.

#### Scenario: hidePagination suppresses the pagination bar

- **WHEN** `hidePagination` is `true`
- **THEN** the pagination bar (dateRangeLabel, "Page X of Y", Previous, Next) is not rendered
- **AND** this is true even when `page > 0 && totalPages > 0`

#### Scenario: hidePagination defaults to false

- **WHEN** `hidePagination` is not passed or is `false`
- **THEN** the existing `hasPagination` computed gates the pagination bar as before

### Requirement: useDailySchedule SHALL accept external reactive page ownership

`useDailySchedule` in `frontend/src/calendar/composables/useDailySchedule.ts` SHALL accept optional `pageRef?: Ref<number>` and `daysPerPageRef?: Ref<number>` in its options. When provided:

- The composable SHALL use the provided refs directly instead of creating internal `ref()` calls for `page` and `daysPerPage`.
- The internal `watch` on planning window changes that resets `page` to 1 SHALL be skipped (the external owner is responsible for page resets).
- The `startParam` computed SHALL derive the API start date from `pageRef.value` and `daysPerPageRef.value` as it does today from internal refs.
- The returned `page` and `daysPerPage` SHALL be the provided refs (not new internal refs).

When `pageRef` and `daysPerPageRef` are omitted, the composable SHALL behave identically to the current implementation.

#### Scenario: External pageRef drives pagination

- **WHEN** `useDailySchedule` is called with `pageRef` set to a parent's reactive ref
- **THEN** incrementing the parent's ref changes the composable's `startParam` and triggers a query refetch
- **AND** the composable's returned `page` is the same ref as the parent's

#### Scenario: Omitting pageRef preserves current behavior

- **WHEN** `useDailySchedule` is called without `pageRef` or `daysPerPageRef`
- **THEN** the composable creates internal refs as today
- **AND** the planning-window reset watch fires as today

### Requirement: usePeopleAvailability SHALL accept external reactive page ownership

`usePeopleAvailability` in `frontend/src/shared/composables/usePeopleAvailability.ts` SHALL accept optional `pageRef?: Ref<number>` and `daysPerPageRef?: Ref<number>` in its options. When provided:

- The composable SHALL use the provided refs directly instead of creating internal `ref()` calls for `page` and `daysPerPage`.
- The internal `watch` on planning window changes that resets `page` to 1 SHALL be skipped.
- The `startParam` computed SHALL derive the API start date from `pageRef.value` and `daysPerPageRef.value`.
- The returned `page` and `daysPerPage` SHALL be the provided refs.

When `pageRef` and `daysPerPageRef` are omitted, the composable SHALL behave identically to the current implementation.

#### Scenario: External pageRef drives people availability pagination

- **WHEN** `usePeopleAvailability` is called with `pageRef` set to a parent's reactive ref
- **THEN** changing the parent's ref changes the API query's `start` parameter and triggers a refetch
- **AND** the composable's returned `page` is the same ref as the parent's

#### Scenario: Omitting pageRef preserves current behavior

- **WHEN** `usePeopleAvailability` is called without `pageRef` or `daysPerPageRef`
- **THEN** the composable creates internal refs as today
- **AND** the planning-window reset watch fires as today

### Requirement: HomeView SHALL wire shared pagination to both homepage data sources

`HomeView.vue` SHALL call `useHomePagination()` and pass the resulting `page` ref and `daysPerPage` ref to both `useDailySchedule({ pageRef, daysPerPageRef })` and `usePeopleAvailability({ pageRef, daysPerPageRef })`. It SHALL pass `hidePagination` to `<DailySchedule>` on the homepage. It SHALL NOT pass `page`, `totalPages`, `goToPrevPage`, `goToNextPage`, or `dateRangeLabel` as props or events to `<DailySchedule>`.

#### Scenario: Both homepage sections use the same date window

- **WHEN** the homepage loads with page 1
- **THEN** both the People availability grid and Daily schedule show the same 4-day window
- **AND** the API queries for both sections use the same `start` parameter

#### Scenario: Homepage page navigation updates both sections

- **WHEN** the shared homepage pagination moves to page 2
- **THEN** both the People availability and Daily schedule queries refetch with the new `start` parameter
- **AND** the rendered content updates to show the next 4-day window

#### Scenario: Homepage DailySchedule has no in-card pagination

- **WHEN** the homepage renders DailySchedule
- **THEN** the Pagination bar ("Page X of Y", Previous, Next) is not rendered inside the card
- **AND** the "View by: Day" button and any card header controls remain rendered

## REMOVED Requirements

### Requirement: Homepage DailySchedule SHALL NOT render in-card pagination

**Reason**: Pagination has moved to the global AppShell header on the homepage route. The `hidePagination` prop on `DailySchedule` suppresses the in-card controls.

**Migration**: HomeView passes `:hide-pagination="true"` to `<DailySchedule>` and no longer passes `page`/`totalPages`/`dateRangeLabel` props or `prev-page`/`next-page` event handlers. The DailySchedule card still renders the day columns and task cards — only the pagination bar is removed on the homepage.

#### Scenario: DailySchedule renders day columns but no pagination on homepage

- **WHEN** the homepage renders DailySchedule with `hidePagination`
- **THEN** the day columns and task cards render as before
- **AND** the pagination bar (Previous, Next, "Page X of Y", dateRangeLabel) is absent

### Requirement: Homepage AppShell header SHALL show page range instead of planning-window range

**Reason**: The static `formattedRange` (full planning window span) in the AppShell header is replaced on `/` by the page-specific date range label from `useHomePagination`. Non-homepage routes retain the existing `formattedRange`.

**Migration**: On `/`, AppShell renders the `rangeLabel` from `useHomePagination` instead of `formattedRange` from `usePlanningWindow`. The `formattedRange` is no longer shown on the homepage.

#### Scenario: Homepage header shows page range, not planning-window range

- **WHEN** the current route is `/`
- **THEN** the header date display shows the 4-day page range (e.g., "5 Jul – 8 Jul, 2026")
- **AND** the full planning-window span ("5 Jul – 13 Aug 2026 · 40 days") is not displayed