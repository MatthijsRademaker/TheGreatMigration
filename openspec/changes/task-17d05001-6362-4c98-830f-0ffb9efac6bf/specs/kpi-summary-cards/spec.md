# kpi-summary-cards Specification

## Purpose

Defines the home-dashboard KPI summary card component layout, copy, iconography, and data wiring for the four summary cards (High priority tasks, People available today, Unassigned jobs, Rooms completed). This spec updates the original KPI cards spec to align the card order, visual anatomy, value format, and copy with the design target in `designs/components.png` section 3 and `designs/home-page.png`.

## MODIFIED Requirements

### Requirement: The home page SHALL render four KPI cards in design order

The `/` route SHALL render the `KpiCards` component at `frontend/src/home/components/KpiCards.vue`. The component SHALL display four KPI cards in this exact order: **High priority tasks** (first), **People available today** (second), **Unassigned jobs** (third), **Rooms completed** (fourth).

#### Scenario: Home route renders four KPI cards in the correct order

- **WHEN** the `/` route renders
- **THEN** the top summary row contains exactly four cards
- **AND** the first card is labeled `High priority tasks`
- **AND** the second card is labeled `People available today`
- **AND** the third card is labeled `Unassigned jobs`
- **AND** the fourth card is labeled `Rooms completed`

#### Scenario: KPI cards use Card primitives

- **WHEN** the KpiCards component renders
- **THEN** each card uses Card, CardHeader, CardContent, CardTitle, and CardDescription from `@/shared/ui/card`

### Requirement: Each KPI card SHALL use a thin left accent border and compact icon chip layout

Each card SHALL render a 4px left accent border (`border-l-4`) using a semantic border color class matching the card's intent. Each card SHALL render a compact rounded icon chip (`size-8 rounded-lg`) with a semantic soft background class in the upper-left of the content area (CardHeader). The current `w-[72px]` full-height left accent column SHALL be removed entirely. The leaf decoration image (`/images/leaf.png`) SHALL remain as an absolute-positioned element in the lower-right corner.

#### Scenario: Each card has a left accent border with the correct semantic color

- **WHEN** the KpiCards component renders
- **THEN** the High priority tasks card has a left border with semantic class `border-destructive`
- **AND** the People available today card has a left border with semantic class `border-info`
- **AND** the Unassigned jobs card has a left border with semantic class `border-warning`
- **AND** the Rooms completed card has a left border with semantic class `border-success`

#### Scenario: Each card has a compact icon chip with the correct semantic background

- **WHEN** the KpiCards component renders
- **THEN** the High priority tasks card renders a `size-8 rounded-lg` icon chip with classes `bg-destructive-soft text-destructive` containing a `FlagIcon`
- **AND** the People available today card renders a `size-8 rounded-lg` icon chip with classes `bg-info-soft text-info` containing a `UsersRoundIcon`
- **AND** the Unassigned jobs card renders a `size-8 rounded-lg` icon chip with classes `bg-warning-soft text-warning` containing a `BriefcaseIcon`
- **AND** the Rooms completed card renders a `size-8 rounded-lg` icon chip with classes `bg-success-soft text-success` containing a `CheckCircleIcon`

#### Scenario: Card content uses correct typography roles

- **WHEN** the KpiCards component renders
- **THEN** card labels use `[font-size:var(--text-caption)] text-muted-foreground`
- **AND** card values use `text-3xl font-semibold`
- **AND** card subtitles use `text-sm text-muted-foreground`

#### Scenario: No full-height accent column remains

- **WHEN** the KpiCards component renders
- **THEN** no card contains a `w-[72px]` element
- **AND** no card uses `flex-row` to create a two-column layout with a left accent column

### Requirement: The People available today card SHALL display a fraction value with date-context subtitle

The card SHALL consume `getDashboardPeopleAvailabilityQuery`. The primary value SHALL display `availableToday / totalPeople` as a fraction (e.g., `6 / 8`). The subtitle SHALL display `available on <MMM D>` where the date is sourced from `range.selectedDate` in the API response. If `range.selectedDate` is undefined, the subtitle SHALL fall back to `available today`. The `availableToday` value SHALL be clamped to never exceed `totalPeople` (defensive guard). The card SHALL handle loading and error states.

#### Scenario: Card displays fraction value and date subtitle from API

- **WHEN** `getDashboardPeopleAvailabilityQuery` resolves with `summary.availableToday = 6`, `summary.totalPeople = 8`, and `range.selectedDate = "2026-07-05"`
- **THEN** the People available today card displays `6 / 8` as the primary value
- **AND** the subtitle displays `available on Jul 5`

#### Scenario: Card falls back when selectedDate is undefined

- **WHEN** `getDashboardPeopleAvailabilityQuery` resolves with `summary.availableToday = 6` and `summary.totalPeople = 8` but `range.selectedDate` is undefined
- **THEN** the People available today card displays `6 / 8` as the primary value
- **AND** the subtitle displays `available today`

#### Scenario: Card shows loading state while query is pending

- **WHEN** `getDashboardPeopleAvailabilityQuery` is in-flight
- **THEN** the People available today card displays a loading indicator (`Loading…`) in the value slot rather than a stale or zero value

#### Scenario: Card shows error state on query failure

- **WHEN** `getDashboardPeopleAvailabilityQuery` fails
- **THEN** the People available today card displays a graceful error state (`Backend unavailable`) in the value slot

### Requirement: The High priority tasks and Unassigned jobs cards SHALL display live task-backlog data with design-like copy

The cards SHALL consume `getTasksBacklogQuery`. The High priority tasks card SHALL display `summary.highPriorityTasks` as its value with subtitle `high priority tasks need attention`. The Unassigned jobs card SHALL display `summary.unassignedTasks` as its value with subtitle `jobs that need assignment`. Both cards SHALL handle loading and error states.

#### Scenario: High priority tasks card displays count and outcome copy from API

- **WHEN** `getTasksBacklogQuery` resolves with `summary.highPriorityTasks = 4`
- **THEN** the High priority tasks card displays the value `4`
- **AND** the subtitle displays `high priority tasks need attention`

#### Scenario: Unassigned jobs card displays count and outcome copy from API

- **WHEN** `getTasksBacklogQuery` resolves with `summary.unassignedTasks = 3`
- **THEN** the Unassigned jobs card displays the value `3`
- **AND** the subtitle displays `jobs that need assignment`

#### Scenario: Task-backlog cards show loading state while query is pending

- **WHEN** `getTasksBacklogQuery` is in-flight
- **THEN** both the High priority tasks and Unassigned jobs cards display a loading indicator (`Loading…`) in their value slots

#### Scenario: Task-backlog cards show error state on query failure

- **WHEN** `getTasksBacklogQuery` fails
- **THEN** both task-backlog KPI cards display a graceful error state (`Backend unavailable`) in their value slots

### Requirement: The Rooms completed card SHALL be an isolated placeholder with no derived business logic

The fourth card SHALL use the same Card primitive layout, thin accent border, compact icon chip, and leaf decoration as the other three cards. It SHALL display a non-numeric placeholder value (`—`) labeled `Rooms completed` with subtitle `rooms fully packed and cleared`. It SHALL include `data-testid="kpi-placeholder-rooms-completed"` and a code comment documenting it as a placeholder for a future room-progress contract. It SHALL NOT derive values from task backlog, schedule, people availability, or any other backend data.

#### Scenario: Rooms completed card renders as styled placeholder

- **WHEN** the KpiCards component renders
- **THEN** the fourth card displays the label `Rooms completed` and the value `—`
- **AND** the subtitle displays `rooms fully packed and cleared`
- **AND** the card renders with a `border-success` left accent border
- **AND** the card renders an icon chip with `CheckCircleIcon` in `bg-success-soft text-success`
- **AND** the card renders `data-testid="kpi-placeholder-rooms-completed"`

#### Scenario: Rooms completed card does not query any backend endpoint

- **WHEN** the KpiCards component renders
- **THEN** no network request is made for room-progress data
- **AND** the card content is purely static

### Requirement: The KPI card styling SHALL use only existing design system tokens

All card styling SHALL use semantic token classes from the Design System v2 token surface exposed in `frontend/src/app/styles.css`. No raw hex values, no new global theme tokens, no new shared primitive directories (`frontend/src/shared/ui/`), and no page-specific raw color hacks SHALL be introduced. The component SHALL remain at `frontend/src/home/components/KpiCards.vue` and SHALL NOT create new directories.

#### Scenario: All accent and icon colors use semantic token classes

- **WHEN** the KpiCards component renders
- **THEN** no card uses raw hex color values in its styling
- **AND** all accent borders use `border-destructive`, `border-info`, `border-warning`, or `border-success`
- **AND** all icon chips use `bg-*-soft text-*` classes from the semantic token surface

#### Scenario: No new shared component directories are created

- **WHEN** the KPI card update is complete
- **THEN** no new directories are added under `frontend/src/shared/ui/` for this change

### Requirement: Unit tests SHALL assert the new card layout, order, copy, and state behavior

`frontend/src/home/__tests__/KpiCards.spec.ts` SHALL be updated to remove all assertions for the obsolete `w-[72px]` accent column, `X of Y available` wording, and old card descriptions. It SHALL add assertions for the thin left accent border, compact icon chip, `X / Y` fraction format, date-context subtitle, new outcome-focused copy, and new card order. Existing coverage for loading, error, and empty states SHALL be preserved.

#### Scenario: Tests assert new card order

- **WHEN** the unit tests run
- **THEN** the test asserting card count verifies the first card is `High priority tasks` and the second is `People available today`

#### Scenario: Tests assert new people card format

- **WHEN** the unit tests run with availability data `availableToday = 5, totalPeople = 8`
- **THEN** the test asserts the people card displays `5 / 8` and the date subtitle `available on <formatted date>`

#### Scenario: Tests assert new layout anatomy

- **WHEN** the unit tests run
- **THEN** tests verify the left accent border class is present (e.g., `border-info` on the people card)
- **AND** tests verify the compact icon chip renders with `size-8 rounded-lg` dimensions
- **AND** no test asserts the presence of a `w-[72px]` element

#### Scenario: Tests assert new subtitle copy

- **WHEN** the unit tests run
- **THEN** tests assert `high priority tasks need attention` appears for the high priority card
- **AND** tests assert `jobs that need assignment` appears for the unassigned card
- **AND** tests assert `rooms fully packed and cleared` appears for the rooms card

#### Scenario: Loading and error state coverage is preserved

- **WHEN** the unit tests run with pending or failed queries
- **THEN** tests still assert `Loading…` and `Backend unavailable` text appears within the new CardContent layout

### Requirement: SSR route-render test SHALL assert updated home-route KPI card content

`frontend/tests/app-routes-render.test.ts` SHALL be updated to replace the obsolete home-route assertions `"of 8"` and `"available"` with assertions matching the new people-card fraction format and subtitle. Existing KPI card label assertions and all non-home route assertions SHALL continue to pass. The `/people` route assertion `"6 of 8 available today"` SHALL NOT be modified.

#### Scenario: Home route SSR test asserts updated people card format

- **WHEN** the SSR route-render test runs for `/`
- **THEN** the rendered HTML contains `6 / 8` (matching the mock data)
- **AND** the rendered HTML does NOT assert `of 8` as a standalone string

#### Scenario: Home route SSR test asserts KPI card labels still appear

- **WHEN** the SSR route-render test runs for `/`
- **THEN** the rendered HTML contains `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`

#### Scenario: People route assertion is unchanged

- **WHEN** the SSR route-render test runs for `/people`
- **THEN** the rendered HTML still contains `6 of 8 available today` (unchanged from current)