## Why

The dashboard is functionally complete but visually static — cards, grids, and KPIs render with no motion, loading states show plain "Fetching…" text, and the schedule board requires form-based editing to assign helpers or move work between days. For a tool named "The Great Migration" that coordinates a stressful, fast-moving event, the interface should feel alive and tactile while staying fast and legible. This change adds a coherent, playful motion language — functional animation that conveys state and change — plus one signature on-theme visualization and a drag-and-drop schedule board that turns helper assignment and rescheduling into direct manipulation.

## What Changes

- Add a shared motion foundation: install `motion-v` and `@formkit/auto-animate`, and introduce a `useMotionPreference` composable that gates all non-essential motion behind `prefers-reduced-motion`.
- Animate the dashboard broadly (the "sweep"): KPI count-ups with spring overshoot, AutoAnimate FLIP transitions on the tasks/people/tools lists, skeleton-shimmer loading states replacing plain "Fetching…" text, and route-change crossfade/slide transitions.
- Add a playful micro-interaction layer: spring-based hover/press feedback on cards and buttons, and reward moments (a pop + checkmark morph when a task becomes fully staffed; a celebration when a day reaches 100% readiness).
- Add a drag-and-drop schedule board: drag a person from the people rail onto a task card to assign them (appends to `assignedTo` via `updateScheduleCard`), and drag a card between day columns to reschedule it (changes `scheduledDate` via `updateScheduleCard`). The home page's read-only `DailySchedule` remains non-draggable.
- Add the signature "Move-Day Readiness Journey" visualization: an on-theme progress path showing the move travelling from old place → in transit → new place as cards get staffed and days complete.

## Capabilities

### New Capabilities
- `motion-foundations`: Shared motion infrastructure — library setup, the `useMotionPreference` composable, reduced-motion gating contract, and the standard spring/easing tokens all other motion uses.
- `animated-dashboard`: The app-wide motion sweep — KPI count-ups, animated list transitions, skeleton loading states, route transitions, and playful micro-interactions/reward moments.
- `dnd-schedule-board`: Drag-and-drop assignment and rescheduling on the schedule board, backed by `updateScheduleCard`, including drag affordances, optimistic state, and error handling.
- `migration-journey`: The signature Move-Day Readiness Journey visualization that aggregates staffing and completion state into an animated, on-theme progress path.

### Modified Capabilities
<!-- None. Motion is an additive enhancement layer; the DnD board introduces new interaction behavior captured as its own new capability rather than altering existing daily-schedule requirements. -->

## Impact

- **Dependencies**: adds `motion-v` and `@formkit/auto-animate`. Reuses already-installed `@vueuse/core` (`useTransition`), `tw-animate-css`, and Vue's built-in `<Transition>`/`<TransitionGroup>`.
- **Frontend code**: `src/home/components/KpiCards.vue`, `src/tasks/`, `src/people/`, `src/tools/` lists; `src/calendar/DailySchedule.vue` and its composable; `src/app/` router/layout for route transitions; new `src/shared/composables/useMotionPreference.ts` and shared motion components/tokens.
- **API**: no backend changes — DnD uses the existing `updateScheduleCard` (PUT `/api/schedule/cards/{id}`, `assignedTo` + `scheduledDate`).
- **Accessibility/perf**: all playful motion must honor `prefers-reduced-motion`; animations must not regress interaction latency on the dashboard.
