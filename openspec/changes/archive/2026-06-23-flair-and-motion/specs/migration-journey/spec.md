## ADDED Requirements

### Requirement: Move-Day Readiness Journey visualization

The system SHALL render a signature "Move-Day Readiness Journey" — an on-theme progress path representing the move travelling from old place → in transit → new place. Its progress SHALL be derived purely from existing query data (per-day staffing completeness and rooms/tools coverage) with no new persisted state.

#### Scenario: Journey reflects current readiness

- **WHEN** the journey renders with current dashboard data
- **THEN** the progress indicator position corresponds to the aggregate readiness (staffing/coverage) of the move

#### Scenario: Journey updates after a change

- **WHEN** an underlying query updates (for example after a drag-and-drop assignment invalidates the schedule)
- **THEN** the journey animates to its new progress position (subject to motion preference)

#### Scenario: Reduced motion journey

- **WHEN** the journey renders or updates and motion is disabled
- **THEN** the correct progress position is shown immediately with no travelling animation

#### Scenario: Move fully ready

- **WHEN** aggregate readiness reaches 100%
- **THEN** the journey shows the arrived/complete state and plays the day-complete celebration (subject to motion preference)
