## 1. Build the presentational component

- [ ] Create `frontend/src/people/PeopleAvailability.vue`.
- [ ] Define local TypeScript interfaces and typed props for the component title, day labels, people rows, and legend, with deterministic demo defaults matching the four-day design slice.
- [ ] Render a titled card or panel with per-person rows, four day columns, availability badges for `available`, `busy`, `partial`, and `off`, and a visible legend.
- [ ] Use existing `Card`, `Badge`, and `Avatar` primitives with compact inline avatar-plus-name row labels; do not use `PersonChip` or generated API types.
- [ ] Preserve small-screen readability with horizontal overflow for the matrix.

## 2. Integrate the component into the people route

- [ ] Replace the placeholder content in `frontend/src/people/PeopleView.vue` with `PeopleAvailability`.
- [ ] Keep the route presentational only: no dashboard availability query, no planning-window integration, no backend wiring, and no generated-client regeneration.
- [ ] Preserve the existing `/people` route metadata and app-shell/sidebar behavior.

## 3. Add focused frontend verification

- [ ] Add `frontend/tests/people/PeopleAvailability.test.ts` using the existing SSR `renderToString` test pattern.
- [ ] Assert the component renders its title, sample people, four day labels, all four status labels, and the legend.
- [ ] Update `frontend/tests/app-routes-render.test.ts` so `/people` asserts the new component content instead of `People availability foundation`.

## 4. Run repository verification

- [ ] Run `scripts/precommit-run`.
