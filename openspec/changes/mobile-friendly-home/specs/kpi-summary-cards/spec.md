## ADDED Requirements

### Requirement: KPI cards SHALL show a compact summary row on mobile with an expand toggle
On viewports below the `sm` breakpoint (640px), `KpiCards.vue` SHALL render a single horizontal summary row in place of the full card grid. The summary row SHALL display all five metrics as compact icon + value pairs inline. A toggle button (chevron) at the end of the row SHALL expand and collapse the full card grid beneath the row. The default state SHALL be collapsed (full grid hidden). On `sm` and above the full card grid SHALL always be visible and the summary row and toggle SHALL NOT be rendered.

#### Scenario: Compact summary row is shown by default on mobile
- **WHEN** the home route renders on a viewport narrower than 640px
- **THEN** the compact summary row is visible
- **AND** all five metric values are present in the summary row (icon + value each)
- **AND** the full card grid is not visible

#### Scenario: Tapping the expand toggle reveals the full card grid on mobile
- **WHEN** the user taps the expand toggle on the compact summary row
- **THEN** the full five-card grid becomes visible below the summary row
- **AND** the toggle icon changes to indicate a collapsed/close action

#### Scenario: Tapping the toggle again hides the full card grid
- **WHEN** the full card grid is expanded and the user taps the toggle again
- **THEN** the full card grid is hidden
- **AND** the compact summary row remains visible

#### Scenario: Full card grid is always visible on tablet and desktop
- **WHEN** the home route renders on a viewport of 640px or wider
- **THEN** the full five-card grid is always visible
- **AND** no compact summary row or expand toggle is rendered in the DOM

#### Scenario: Expanding the grid does not re-fetch queries
- **WHEN** the user expands the full card grid on mobile
- **THEN** no new network requests are made for KPI data
- **AND** the displayed values are the same as those shown in the compact row
