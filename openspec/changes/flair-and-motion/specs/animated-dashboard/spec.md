## ADDED Requirements

### Requirement: Animated KPI counters

KPI summary cards SHALL animate their numeric values with a spring count-up (with slight overshoot) when the value first loads or changes, gated by motion preference.

#### Scenario: Value animates on load

- **WHEN** a KPI card receives its value and motion is enabled
- **THEN** the displayed number counts up from a baseline to the final value with a spring overshoot that settles on the exact value

#### Scenario: Reduced motion shows final value

- **WHEN** a KPI card receives its value and motion is disabled
- **THEN** the final value is shown immediately with no count-up

### Requirement: Animated list transitions

The tasks, people, and tools lists SHALL animate item addition, removal, and reordering with FLIP transitions so items never teleport, gated by motion preference.

#### Scenario: Item added

- **WHEN** an item is added to an animated list and motion is enabled
- **THEN** the new item animates in and existing items smoothly shift to accommodate it

#### Scenario: Item removed or reordered

- **WHEN** an item is removed or the list reorders and motion is enabled
- **THEN** affected items animate to their new positions rather than jumping

### Requirement: Skeleton loading states

Loading states for dashboard panels SHALL render shimmer skeletons that approximate the eventual content layout, replacing plain "Fetching…" text.

#### Scenario: Panel is loading

- **WHEN** a dashboard panel's query is in flight
- **THEN** a shimmer skeleton matching the panel's layout is shown instead of placeholder text

#### Scenario: Reduced motion skeleton

- **WHEN** a panel is loading and motion is disabled
- **THEN** a static (non-shimmering) skeleton placeholder is shown

### Requirement: Route transitions

Navigating between top-level routes SHALL apply a crossfade/slide transition, gated by motion preference.

#### Scenario: Navigating between pages

- **WHEN** the user navigates from one route to another and motion is enabled
- **THEN** the outgoing view transitions out and the incoming view transitions in smoothly

#### Scenario: Reduced motion navigation

- **WHEN** the user navigates and motion is disabled
- **THEN** the new view appears immediately with no transition

### Requirement: Playful micro-interactions and reward moments

Interactive cards and buttons SHALL provide spring-based hover/press feedback, and the system SHALL play a reward animation when a task becomes fully staffed and when a day reaches 100% readiness, all gated by motion preference.

#### Scenario: Hover and press feedback

- **WHEN** the user hovers or presses an interactive card or button and motion is enabled
- **THEN** the element responds with a spring-based scale/elevation change

#### Scenario: Task becomes fully staffed

- **WHEN** a task's staffing status transitions to `fullyStaffed` and motion is enabled
- **THEN** a celebratory pop with a checkmark morph plays on that card

#### Scenario: Reduced motion suppresses rewards

- **WHEN** a reward condition is met and motion is disabled
- **THEN** the status change is reflected with no celebratory animation
