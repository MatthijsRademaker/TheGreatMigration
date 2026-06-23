# motion-foundations Specification

## Purpose
Shared motion infrastructure for the frontend: the `useMotionPreference` reduced-motion gate, the standard spring/duration/easing tokens, and availability of the motion libraries (`motion-v`, `@formkit/auto-animate`). Every animated surface consumes this foundation so motion is consistent, centrally tunable, and uniformly gated by `prefers-reduced-motion`.

## Requirements

### Requirement: Centralized motion preference gate

The system SHALL provide a single `useMotionPreference` composable that reports whether non-essential motion is enabled, derived from the user's `prefers-reduced-motion` setting. All animated surfaces SHALL consume this composable rather than checking the media query independently.

#### Scenario: User prefers reduced motion

- **WHEN** the user's OS/browser reports `prefers-reduced-motion: reduce`
- **THEN** `useMotionPreference` reports motion disabled
- **AND** decorative motion (count-ups, springs, confetti, route slides) is skipped while the corresponding state change still applies instantly

#### Scenario: User allows motion

- **WHEN** the user's OS/browser reports `prefers-reduced-motion: no-preference`
- **THEN** `useMotionPreference` reports motion enabled
- **AND** components apply their full playful motion

#### Scenario: Preference changes at runtime

- **WHEN** the user toggles the reduced-motion setting while the app is open
- **THEN** `useMotionPreference` updates reactively without a page reload

### Requirement: Shared motion tokens

The system SHALL define motion tokens (spring configs, durations, easings) in one location, and the motion libraries and CSS SHALL consume those tokens so the playful feel is tunable centrally.

#### Scenario: Single source of truth

- **WHEN** a developer changes a shared spring/duration token
- **THEN** every surface using that token reflects the change without per-component edits

### Requirement: Motion library availability

The system SHALL make `motion-v` and `@formkit/auto-animate` available to the frontend, alongside reuse of the already-installed `@vueuse/core`, `tw-animate-css`, and Vue transition primitives.

#### Scenario: Build succeeds with new dependencies

- **WHEN** the frontend is built and type-checked
- **THEN** the build completes successfully with the motion libraries installed and importable
