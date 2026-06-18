## 1. Specification

- [ ] 1.1 Create `people-availability-design-alignment/spec.md` under `openspec/changes/task-15d3e474-3338-45a3-a7ff-51246d719e92/specs/` with text-only design-alignment requirements following the established OpenSpec pattern.

## 2. Shared Component Refinement

- [ ] 2.1 Remove `CardDescription` from the `PeopleAvailability.vue` template — keep `CardHeader` with `CardTitle` only.
- [ ] 2.2 Remove the summary row (`{{ availableToday }} of {{ totalPeople }} available today`) from the `CardContent` template.
- [ ] 2.3 Remove the "Person" column header `<th>` from the table `<thead>` — keep the first data column (avatar + name) but without a column label.
- [ ] 2.4 Switch day-column headers from abbreviated labels (Mon, Tue…) to date+weekday format (e.g., "Sun 5 Jul") matching the adapter output. Update default `days` prop if needed.
- [ ] 2.5 Move the legend `<div>` from inside the `overflow-x-auto` container to a footer row below the `<table>`, outside the scroll container, preserving the inline-flex layout and Badge variants.
- [ ] 2.6 Verify all editable-mode behaviour is preserved: Popover triggers, status selection, Actions column header, Delete buttons, confirmation dialogs, `deletingPersonId` disabled state, `updating` disabled state.

## 3. Shared Adapter Cleanup

- [ ] 3.1 Remove the hard-coded `description` field from `adaptToComponentProps` in `usePeopleAvailability.ts`.
- [ ] 3.2 Stop passing `availableToday` and `totalPeople` in the adapted props output (keep the raw `DashboardBody` available via the existing `rawData` ref for future use).

## 4. Type Contract Update

- [ ] 4.1 Remove `description` from `PeopleAvailabilityProps` in `frontend/src/people/types.ts`.
- [ ] 4.2 Keep `availableToday` and `totalPeople` in `PeopleAvailabilityProps` (backward-compatible type contract, even though the template no longer renders them).

## 5. Test Suite Updates

- [ ] 5.1 Update `frontend/tests/people/PeopleAvailability.test.ts`: remove assertions on summary row text ("1 of 4 available today", "0 of 0 available today"), remove assertions on description text, update card-title assertions, update date header format expectations, verify "Person" column header is absent, verify legend is present at footer position, preserve all editable-mode tests unchanged.
- [ ] 5.2 Update `frontend/tests/people/PeopleView.test.ts`: verify editable matrix renders without description or summary text leaks; confirm Actions column, Delete buttons, and person names still render; confirm loading/error/empty states are unchanged.
- [ ] 5.3 Update `frontend/tests/app-routes-render.test.ts`: replace `"6 of 8 available today"` content assertion for `/people` route with assertions on legend footer badge variants and matrix person names; update `/people` path description assertion if applicable; verify home route `/` still asserts `"People availability"` heading.

## 6. Verification

- [ ] 6.1 Run `npx vitest run` in `frontend/` and confirm all tests pass.
- [ ] 6.2 Manually verify the home dashboard `/` renders the compact People Availability card without description or summary row.
- [ ] 6.3 Manually verify the `/people` route renders the matrix in editable mode with Actions column, Popover triggers, and Delete functionality intact.