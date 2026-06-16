## Context

The repository already has the design-system primitives and semantic tokens needed for a people-availability matrix, but `frontend/src/people/PeopleView.vue` still renders a placeholder helper card layout. The task scope is explicitly component-only: match the People Availability panel from `designs/components.png`, keep data local or prop-driven, and defer all backend wiring.

## Goals

- Add a presentational People Availability component under `frontend/src/people/`.
- Match the design intent with a titled panel, person rows, four day columns, status cells, and a visible legend.
- Reuse existing `Card`, `Badge`, and `Avatar` primitives with design-system tokens.
- Replace the `/people` placeholder content with the new component.
- Add focused SSR/frontend tests for the component and the `/people` route.

## Non-Goals

- Wiring `GET /api/dashboard/people-availability` into the component or route.
- Importing generated query artifacts or generated client types.
- Regenerating OpenAPI/client files or changing backend handlers, persistence, or endpoint contracts.
- Building other dashboard panels, filters, editing flows, CRUD, drag/drop, or date navigation.

## Decisions

### 1. Component location and route integration

Place the component at `frontend/src/people/PeopleAvailability.vue` and render it from `frontend/src/people/PeopleView.vue`, replacing the current placeholder card. This uses the existing `/people` route as the intentional preview surface instead of creating a new route or leaving the component unmounted.

### 2. Data contract

Expose typed props for the component's displayed content, including people rows and day labels, and provide deterministic demo defaults. This keeps the component visually testable now while making future backend wiring additive instead of structural.

### 3. Visual composition

Compose the surface from existing primitives: `Card` for the panel, `Badge` variants `available`, `busy`, `partial`, and `off` for status chips, and `Avatar` for initials. Person rows use a compact inline avatar-plus-name treatment rather than `PersonChip`, whose pill styling is too bulky for a dense matrix.

### 4. Responsive behavior

Use a matrix layout that preserves readability on small screens through horizontal overflow instead of stacking each person into separate cards or compressing labels.

### 5. Testing

Add a dedicated SSR component test at `frontend/tests/people/PeopleAvailability.test.ts` following the existing `renderToString` pattern, and update `frontend/tests/app-routes-render.test.ts` so `/people` asserts new component content instead of the removed placeholder text.

## Risks

- Pixel-level spacing and legend alignment still need implementation review against the design images because the accepted decisions are based on the dossier's text interpretation of those assets.
- Dense matrix layouts can become visually heavy if default spacing from reused primitives is not kept compact inside the card.
- Route tests will fail until the old `/people` placeholder assertion is updated with the new rendered content.

## Conflict Resolution

- **Row label treatment:** Resolved to compact `Avatar` plus inline name text, not `PersonChip`.
- **Props vs. internal data:** Resolved to typed props with demo defaults.
- **File naming:** Resolved to `PeopleAvailability.vue` because the domain concept matters more than the `Card` wrapper implementation detail.
- **Route integration:** Resolved to replace the existing `PeopleView.vue` placeholder in this task rather than leaving integration optional.

## Traceability

- Task: `2f0b5a56-7dab-41b8-b9fa-1ff40c829a18`
- Dossier: `2026-06-15T21:05:38.000Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial`
