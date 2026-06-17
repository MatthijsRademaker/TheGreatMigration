# kpi-summary-cards Specification

## ADDED Requirements

### Requirement: The KPI card icon-accent column SHALL use per-card semantic color mapping

Each KPI card SHALL render a full-height left accent column colored with a per-card semantic token. The icon SHALL be centered inside this column at `size-8` (32px). The color mapping SHALL be: People available today → `info`, High priority tasks → `destructive`, Unassigned jobs → `warning`, Rooms completed → `success`. The column background SHALL use the `-soft` variant of each token (e.g., `bg-destructive-soft`), and the icon SHALL use the full-strength token (e.g., `text-destructive`). No new global theme tokens SHALL be introduced.

#### Scenario: People available today card shows info-accented left column

- **WHEN** the KpiCards component renders
- **THEN** the People available today card has a left accent column with `bg-info-soft` and a `UsersRoundIcon` in `text-info`

#### Scenario: High priority tasks card shows destructive-accented left column

- **WHEN** the KpiCards component renders
- **THEN** the High priority tasks card has a left accent column with `bg-destructive-soft` and a `TriangleAlertIcon` in `text-destructive`

#### Scenario: Unassigned jobs card shows warning-accented left column

- **WHEN** the KpiCards component renders
- **THEN** the Unassigned jobs card has a left accent column with `bg-warning-soft` and a `HammerIcon` in `text-warning`

#### Scenario: Rooms completed card shows success-accented left column

- **WHEN** the KpiCards component renders
- **THEN** the Rooms completed card has a left accent column with `bg-success-soft` and a `Building2Icon` in `text-success`

## MODIFIED Requirements

### Requirement: The home page SHALL render a KPI summary card grid driven by contract-backed backend data

The `/` route SHALL render a feature-local `KpiCards` component in place of the current static summary row. The component SHALL display four KPI cards in a two-column layout per card: a full-height left accent column containing the icon, and a right content column stacking title (row 1), KPI value (row 2), and subtitle (row 3). The component SHALL use the base `Card` primitive from the shared UI library as the outer shell, with a custom `<div>`-based inner layout replacing `CardHeader`, `CardContent`, `CardTitle`, and `CardDescription`. The component SHALL live at `frontend/src/home/components/KpiCards.vue`.

#### Scenario: Home route renders four KPI cards

- **WHEN** the `/` route renders
- **THEN** the top summary row contains exactly four cards with labels `People available today`, `High priority tasks`, `Unassigned jobs`, and `Rooms completed`

#### Scenario: KPI cards use a two-column flex layout with Card as the outer shell

- **WHEN** the KpiCards component renders
- **THEN** each card uses the `Card` primitive from `@/shared/ui/card` as the outer shell
- **AND** each card's inner layout is a flex row with a left accent column (icon) and a right content column (title, KPI value, subtitle)
- **AND** `CardHeader`, `CardContent`, `CardTitle`, and `CardDescription` are NOT used in the per-card template

#### Scenario: Lower HomeView sections are preserved

- **WHEN** the `/` route renders after the KPI card change
- **THEN** the `Task Management`, `People availability`, `Daily Schedule`, and `Move notes` sections remain present below the KPI row

### Requirement: The KPI card styling SHALL use existing design system tokens with a left-accent-column layout

The cards SHALL use semantic accent classes from the Design System v2 token surface exposed in `frontend/src/app/styles.css`. Card icon columns SHALL use the `*-soft` background token class paired with the full `text-*` foreground class (e.g., `bg-destructive-soft text-destructive`). No new global theme tokens, shared primitive directories, or page-specific raw color hacks SHALL be introduced.

#### Scenario: Card accent columns use semantic token classes

- **WHEN** the KpiCards component renders
- **THEN** each card's left accent column uses a semantic `bg-*-soft` plus `text-*` class pair (e.g., `bg-info-soft text-info`) rather than raw color values

#### Scenario: No new shared component directories are created

- **WHEN** the KPI card implementation is complete
- **THEN** no new directories are added under `frontend/src/shared/ui/` for this change
