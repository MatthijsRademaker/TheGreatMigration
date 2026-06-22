## Context

The frontend is Vue 3.5 + Vite 6 + Tailwind 4 + shadcn-vue (reka-ui), feature-sliced under `src/` (home, tasks, people, tools, calendar, rooms). It already ships `@vueuse/core`, `tw-animate-css`, and Vue's built-in transition primitives, but uses almost none of their motion capability. Data flows through Pinia Colada queries/mutations generated from an OpenAPI spec (`src/client/`). The schedule board reads `getDashboardDailySchedule` and can write via `updateScheduleCard` (PUT `/api/schedule/cards/{id}`), whose body accepts `assignedTo: string[]` and `scheduledDate`, returning the recomputed `TaskCard` (`assignedCount`, `staffingStatus` derived server-side).

This change layers a coherent, playful motion system on top without altering the backend or the underlying data model.

## Goals / Non-Goals

**Goals:**
- One reusable motion foundation (tokens + reduced-motion gate) that every animated surface consumes, so motion is consistent and centrally tunable.
- Broad, cheap polish across the dashboard (counters, list FLIP, skeletons, route transitions, micro-interactions).
- A drag-and-drop schedule board that performs real, persisted mutations with optimistic UI and clear failure recovery.
- One signature, on-theme visualization (the Move-Day Readiness Journey).
- Full `prefers-reduced-motion` compliance from the first commit.

**Non-Goals:**
- No backend/API changes. No new persistence.
- No touch/mobile drag support in this change (desktop pointer-first; the board degrades to existing controls on touch).
- No global redesign of the visual language (colors, spacing, typography stay as-is).
- The home page's read-only `DailySchedule` stays non-interactive.

## Decisions

### Decision: Library set — `motion-v` + `@formkit/auto-animate`, reuse the rest
- `@formkit/auto-animate` for list add/remove/reorder (tasks, people, tools) — one-line, zero-config FLIP, highest value-to-effort.
- `motion-v` (Vue's Motion/Framer port) for spring physics, gesture/drag affordances on the board, layout transitions, and reward moments — the declarative model fits Vue and shadcn-vue cleanly.
- Reuse `@vueuse/core`'s `useTransition` for numeric KPI count-ups, `tw-animate-css` for skeleton shimmer keyframes, and Vue `<Transition>`/`<TransitionGroup>` for route transitions.
- **Alternatives considered**: GSAP (rejected — imperative, heavier, overkill without complex timelines); `@vueuse/motion` (older, directive-based — superseded by `motion-v`); building everything on bare `<Transition>` (rejected — no spring/gesture/drag ergonomics for the board).

### Decision: Centralized reduced-motion gate via `useMotionPreference`
A single composable wraps `@vueuse/core`'s `usePreferredReducedMotion` and exposes both a boolean and pre-resolved motion configs (durations → 0, springs → instant) so components never branch ad hoc. Reduced motion disables decorative/playful motion (count-up, confetti, hover springs, route slide) while preserving instant state changes and essential affordances.
- **Alternative considered**: per-component `@media (prefers-reduced-motion)` CSS — rejected because JS-driven motion (`motion-v`, `useTransition`) can't be gated by CSS alone, and scattering the check invites drift.

### Decision: DnD interactions are optimistic, then reconciled
On drop, update local card state immediately (append person / change date), fire `updateScheduleCard`, and reconcile with the returned `TaskCard`. On error, roll back to pre-drop state and surface a non-blocking error. Reward moments (fully-staffed pop) fire on the reconciled `staffingStatus === 'fullyStaffed'`, not optimistically, to avoid celebrating a failed write.
- **Alternative considered**: pessimistic (await server before moving) — rejected as it makes drag feel laggy, defeating the tactile goal.

### Decision: Migration Journey is derived/read-only
The readiness visualization computes purely from existing query data (per-day staffing completeness, rooms/tools coverage) — no new state. It animates as the underlying queries update (including after a DnD-driven mutation invalidates them).

### Decision: Motion tokens live alongside the design system
Spring/duration/easing constants are defined once (co-located with the existing design tokens) and imported by `motion-v` configs and CSS, keeping the playful feel tunable in one place.

## Risks / Trade-offs

- **Optimistic DnD diverges from server truth on failure** → roll back to captured pre-drop snapshot and invalidate the schedule query; never leave the board in a fabricated state (aligns with house rule "fail loudly, no fake success").
- **Two new dependencies increase bundle size** → both are small and tree-shakeable; route transitions/skeletons reuse already-installed libs; measure build size in verification.
- **Playful motion annoying on repeat / under time pressure** → conservative spring settings, short durations, no autoplay loops; everything gates off under reduced motion.
- **`updateScheduleCard` may reject invalid assignments (e.g. over-staffing, person already assigned)** → treat server error as authoritative, roll back, show the reason; do not mask with a silent default.
- **AutoAnimate + `motion-v` both touching the same list could conflict** → assign each list exactly one motion mechanism; board cards use `motion-v` layout, simple lists use AutoAnimate.
- **Desktop-only drag excludes touch users** → board remains operable via the existing (non-drag) controls; documented as a non-goal, not a regression.

## Open Questions

- Does `updateScheduleCard` enforce assignment constraints (max = `peopleNeeded`, duplicate person) server-side, or must the client guard before the call? Resolve during implementation by exercising the endpoint.
- Should the Journey live on the home page, the calendar page, or both? Default: home page hero; confirm during build.
