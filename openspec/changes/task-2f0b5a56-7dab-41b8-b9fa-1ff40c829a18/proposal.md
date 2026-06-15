## Why

The `/people` route still renders a placeholder "People availability foundation" card, while the design contract already defines a denser People Availability panel as a key dashboard question. This change turns that placeholder into a presentational component that matches `designs/components.png` and `designs/home-page.png` without pulling backend or generated-client work into scope.

## What Changes

- Add `frontend/src/people/PeopleAvailability.vue` as a presentational Vue component for the People Availability matrix.
- Give the component typed props with deterministic demo defaults for its title, days, people, and legend so it renders immediately now and can accept real data later.
- Compose the UI from existing design-system primitives and semantics: `Card`, `Badge` availability variants, and `Avatar` with compact inline name labels instead of `PersonChip`.
- Integrate the component into `frontend/src/people/PeopleView.vue` in place of the current placeholder content.
- Add focused SSR/frontend tests for the component and update the `/people` route render assertion to the new visible content.
- Keep the change strictly frontend-only: no backend endpoint wiring, no generated client imports or regeneration, and no backend or canonical OpenSpec file edits.

## Impact

- The `/people` route becomes a visually representative preview of the final people-availability surface.
- Future backend wiring can pass real data through the established props contract instead of rewriting the component structure.
- Existing shared UI primitives gain a real composition use case without adding duplicate low-level components or bespoke styling surfaces.
- Backend handlers, persistence, generated API files, and unrelated dashboard panels remain unchanged.
