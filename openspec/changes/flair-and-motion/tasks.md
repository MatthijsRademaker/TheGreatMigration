## 1. Motion foundations

- [ ] 1.1 Add `motion-v` and `@formkit/auto-animate` to `frontend/package.json` and install
- [ ] 1.2 Define shared motion tokens (spring configs, durations, easings) co-located with the existing design tokens
- [ ] 1.3 Create `src/shared/composables/useMotionPreference.ts` wrapping `@vueuse/core` reduced-motion detection, exposing a reactive boolean and pre-resolved motion configs
- [ ] 1.4 Verify build/type-check passes with the new dependencies (`scripts/build-dashboard`)

## 2. Animated dashboard sweep

- [ ] 2.1 Add spring count-up to `KpiCards.vue` using `@vueuse/core` `useTransition`, gated by `useMotionPreference`
- [ ] 2.2 Apply AutoAnimate FLIP to the tasks list (add/remove/reorder), gated by motion preference
- [ ] 2.3 Apply AutoAnimate FLIP to the people list, gated by motion preference
- [ ] 2.4 Apply AutoAnimate FLIP to the tools list, gated by motion preference
- [ ] 2.5 Create reusable shimmer skeleton component(s) using `tw-animate-css`; replace plain "Fetching…"/loading text in dashboard panels (schedule, people, tools, tasks)
- [ ] 2.6 Add route crossfade/slide transition in the app layout/router view, gated by motion preference
- [ ] 2.7 Add spring hover/press micro-interactions to interactive cards and buttons, gated by motion preference

## 3. Reward moments

- [ ] 3.1 Implement fully-staffed reward (pop + checkmark morph) triggered on reconciled `staffingStatus === 'fullyStaffed'`, gated by motion preference
- [ ] 3.2 Implement day-complete celebration triggered on 100% readiness, gated by motion preference

## 4. Drag-and-drop schedule board

- [ ] 4.1 Confirm `updateScheduleCard` assignment/reschedule semantics by exercising the endpoint (constraint enforcement, duplicate handling)
- [ ] 4.2 Make people-rail entries draggable and task cards drop targets using `motion-v` drag affordances (lift/tilt/highlight), gated by motion preference
- [ ] 4.3 Implement optimistic assign-on-drop: append person to `assignedTo`, call `updateScheduleCard`, reconcile with response, invalidate the schedule query
- [ ] 4.4 Implement optimistic reschedule-on-drop between day columns: set `scheduledDate`, call `updateScheduleCard`, FLIP-animate the move, reconcile
- [ ] 4.5 Implement rollback + non-blocking error surfacing for failed assign/reschedule mutations
- [ ] 4.6 Guard against duplicate assignment (no duplicate, no redundant mutation)
- [ ] 4.7 Ensure home-page read-only `DailySchedule` exposes no drag handles or drop targets

## 5. Migration journey (signature)

- [ ] 5.1 Implement a derived readiness selector aggregating per-day staffing completeness and rooms/tools coverage from existing queries (no new state)
- [ ] 5.2 Build the Move-Day Readiness Journey visualization (old place → in transit → new place path) consuming the derived readiness
- [ ] 5.3 Animate progress position changes (including after DnD-driven query invalidation) and the arrived/complete state, gated by motion preference
- [ ] 5.4 Place the journey on the home page hero (confirm placement during build)

## 6. Verification

- [ ] 6.1 Add/extend component tests for count-up, list transitions, DnD assign/reschedule (including rollback), and reduced-motion gating
- [ ] 6.2 Manually verify all motion honors `prefers-reduced-motion` (reduced shows instant final state, no decorative animation)
- [ ] 6.3 Run `scripts/precommit-run` and confirm build, type-check, and tests pass; confirm no significant bundle-size regression
