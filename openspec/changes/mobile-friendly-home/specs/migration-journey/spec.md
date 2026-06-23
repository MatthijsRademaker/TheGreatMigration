## ADDED Requirements

### Requirement: Move-Day Readiness Journey SHALL be hidden on mobile viewports
The `MigrationJourney` component SHALL NOT be rendered on viewports below the `sm` breakpoint (640px). The `HomeView.vue` template SHALL apply `class="hidden sm:block"` (or equivalent) to the `<MigrationJourney>` element so that no DOM is emitted on mobile. On `sm` and above the component SHALL render unchanged.

#### Scenario: Journey is absent from DOM on mobile viewport
- **WHEN** the home route renders on a viewport narrower than 640px
- **THEN** the MigrationJourney card element is not present in the rendered DOM

#### Scenario: Journey renders normally on tablet and desktop
- **WHEN** the home route renders on a viewport of 640px or wider
- **THEN** the MigrationJourney card is present and fully rendered
- **AND** all existing journey behaviour (progress animation, status label, celebration) is unchanged
