## MODIFIED Requirements

### Requirement: The home page SHALL render five KPI cards in design order

The `/` route SHALL render the `KpiCards` component at `frontend/src/home/components/KpiCards.vue`. The component SHALL display five KPI cards in this exact order: **High priority tasks** (first), **People available today** (second), **Unassigned jobs** (third), **Rooms completed** (fourth), **Tools covered** (fifth). The **Tools covered** card SHALL display a fraction value of claimed tools over total tools, sourced from `GET /api/tools` (`summary.claimed` / `summary.total`).

#### Scenario: Home route renders five KPI cards in the correct order

- **WHEN** the `/` route renders
- **THEN** the top summary row contains exactly five cards
- **AND** the first card is labeled `High priority tasks`
- **AND** the second card is labeled `People available today`
- **AND** the third card is labeled `Unassigned jobs`
- **AND** the fourth card is labeled `Rooms completed`
- **AND** the fifth card is labeled `Tools covered`

#### Scenario: Tools covered card shows a claimed-over-total fraction

- **WHEN** the KpiCards component renders and the tools summary reports 4 claimed of 7 total
- **THEN** the Tools covered card displays the fraction `4 / 7`

#### Scenario: KPI cards use Card primitives

- **WHEN** the KpiCards component renders
- **THEN** each card uses Card, CardHeader, and CardContent from `@/shared/ui/card`

### Requirement: Each KPI card SHALL use a thin left accent border and compact icon chip layout

Each card SHALL render a 4px left accent border (`border-l-4`) using a semantic border color class matching the card's intent. Each card SHALL render a compact rounded icon chip (`size-8 rounded-lg`) with a semantic soft background class in the upper-left of the content area (CardHeader). The current `w-[72px]` full-height left accent column SHALL be removed entirely. The leaf decoration image (`/images/leaf.png`) SHALL remain as an absolute-positioned element in the lower-right corner.

#### Scenario: Each card has a left accent border with the correct semantic color

- **WHEN** the KpiCards component renders
- **THEN** the High priority tasks card has a left border with semantic class `border-destructive`
- **AND** the People available today card has a left border with semantic class `border-info`
- **AND** the Unassigned jobs card has a left border with semantic class `border-warning`
- **AND** the Rooms completed card has a left border with semantic class `border-success`
- **AND** the Tools covered card has a left border with semantic class `border-info`

#### Scenario: Each card has a compact icon chip with the correct semantic background

- **WHEN** the KpiCards component renders
- **THEN** the High priority tasks card renders a `size-8 rounded-lg` icon chip with classes `bg-destructive-soft text-destructive` containing a `FlagIcon`
- **AND** the People available today card renders a `size-8 rounded-lg` icon chip with classes `bg-info-soft text-info` containing a `UsersRoundIcon`
- **AND** the Unassigned jobs card renders a `size-8 rounded-lg` icon chip with classes `bg-warning-soft text-warning` containing a `BriefcaseIcon`
- **AND** the Rooms completed card renders a `size-8 rounded-lg` icon chip with classes `bg-success-soft text-success` containing a `CheckCircleIcon`
- **AND** the Tools covered card renders a `size-8 rounded-lg` icon chip with classes `bg-info-soft text-info` containing a `WrenchIcon`

#### Scenario: Card content uses correct typography roles

- **WHEN** the KpiCards component renders
- **THEN** card labels use `[font-size:var(--text-caption)] text-muted-foreground`
- **AND** card values use `text-3xl font-semibold`
- **AND** card subtitles use `text-sm text-muted-foreground`

#### Scenario: No full-height accent column remains

- **WHEN** the KpiCards component renders
- **THEN** no card contains a `w-[72px]` element
- **AND** no card uses `flex-row` to create a two-column layout with a left accent column
