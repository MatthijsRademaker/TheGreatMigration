## 1. Restructure KpiCards.vue layout and adopt Card primitives

- [ ] Replace the `w-[72px]` full-height left accent column with a `border-l-4` semantic left border on each Card root element (using `border-destructive`, `border-info`, `border-warning`, `border-success` classes).
- [ ] Place a compact `size-8 rounded-lg` icon chip with semantic soft background (`bg-destructive-soft text-destructive`, etc.) inside CardHeader at the top-left of the content area.
- [ ] Adopt CardHeader (icon chip + label) and CardContent (value + subtitle) primitives from `@/shared/ui/card`. Remove the `!py-0 !gap-0` overrides that bypass them.
- [ ] Use `[font-size:var(--text-caption)] text-muted-foreground` for card labels in CardHeader, `text-3xl font-semibold` for values in CardContent, and `text-sm text-muted-foreground` for subtitles in CardContent.
- [ ] Preserve the leaf decoration (`<img src="/images/leaf.png" ... absolute bottom-0 right-0 h-14 w-auto pointer-events-none />`).
- [ ] Verify the new card anatomy renders correctly in light and dark modes.

## 2. Reorder KPI cards to match the design

- [ ] Change `cardConfigs` computed array order to: `high-priority` (first), `people` (second), `unassigned` (third), `rooms` (fourth).
- [ ] Update all `id` references if any downstream logic depends on card index/order.

## 3. Update People card: fraction format + date subtitle from range.selectedDate

- [ ] Add `range` destructuring from `availabilityQuery.data.value` to access `range.selectedDate`.
- [ ] Create a computed property to format `selectedDate` as "MMM D" (e.g., "Jul 5"). Fall back to `undefined` if the field is missing.
- [ ] Update the people card value template from `{{ card.value }} of {{ totalPeople }} available` to `{{ card.value }} / {{ totalPeople }}`.
- [ ] Update the people card subtitle from `Helpers with confirmed availability` to `available on <formattedDate>` (dynamic). Fall back to `available today` when `selectedDate` is undefined.
- [ ] Verify that `displayAvailableToday` computed (clamping to `totalPeople`) still applies correctly.

## 4. Update card subtitle copy to outcome-focused text

- [ ] Change high priority subtitle from `Tasks marked as high priority` to `high priority tasks need attention`.
- [ ] Change unassigned subtitle from `Tasks with no one assigned yet` to `jobs that need assignment`.
- [ ] Keep rooms subtitle as `rooms fully packed and cleared` (static placeholder text).

## 5. Replace Lucide icon imports

- [ ] Import `FlagIcon`, `BriefcaseIcon`, `CheckCircleIcon` from `@lucide/vue` (in addition to `UsersRoundIcon`, which stays).
- [ ] Remove imports for `TriangleAlertIcon`, `HammerIcon`, `Building2Icon`.
- [ ] Update `cardConfigs` icon mapping: high priority → `FlagIcon`, unassigned → `BriefcaseIcon`, rooms → `CheckCircleIcon`, people → `UsersRoundIcon` (unchanged).

## 6. Preserve Rooms completed as isolated placeholder

- [ ] Keep `status: 'empty'` and `value: 0` in card config.
- [ ] Keep `data-testid="kpi-placeholder-rooms-completed"`.
- [ ] Keep code comment documenting it as placeholder for a future room-progress contract.
- [ ] Ensure no backend query or derived business logic is added for this card.

## 7. Update unit tests in KpiCards.spec.ts

- [ ] Remove the test `"renders two-column layout with left accent column for each card"` that asserts `w-[72px]`, `flex-row`, and `bg-info-soft text-info` accent column classes. Replace with a test that asserts the thin accent border (`border-l-4` with semantic class) and compact icon chip (`size-8 rounded-lg` with semantic soft background).
- [ ] Remove the test `"shows 'X of Y available' format for the people card"`. Replace with a test asserting the `X / Y` fraction format (e.g., `"5 / 8"`) and subtitle showing date context.
- [ ] Update the test `"renders four cards"` to assert new card order (check first card is High priority, second is People).
- [ ] Update the test `"shows backlog cards with correct values"` to assert new subtitle copy for high-priority and unassigned cards.
- [ ] Update the test `"shows loading state when backlog is pending"` to assert Loading text appears inside the new CardContent layout.
- [ ] Update mock helpers if `range.selectedDate` requires new mock fields.
- [ ] Preserve existing coverage for: loading state on availability pending, error state on query failure, empty state when totalPeople is zero, rooms placeholder with data-testid.

## 8. Update SSR route-render test in app-routes-render.test.ts

- [ ] Update the home-route (`/`) assertion block: remove `expect(html).toContain("of 8")` and `expect(html).toContain("available")`.
- [ ] Add assertion for the new people-card format (e.g., `expect(html).toContain("6 / 8")`).
- [ ] Add assertion for the date-context subtitle if the SSR mock's `range.selectedDate` renders.
- [ ] Verify existing KPI card label assertions (`People available today`, `High priority tasks`, `Unassigned jobs`, `Rooms completed`) still pass.
- [ ] Verify all existing route-render assertions for other routes continue to pass.
- [ ] Verify the `/people` route assertion `"6 of 8 available today"` is NOT modified.

## 9. Update the canonical OpenSpec spec

- [ ] Update `openspec/specs/kpi-summary-cards/spec.md` to reflect: new card order (High priority first), new layout anatomy (thin accent border + icon chip), new people value format (fraction with date context), new subtitle copy, and new icon choices.
- [ ] Ensure all scenarios remain consistent with the updated requirements.

## 10. Verify the change end-to-end

- [ ] Run `npm run check` (vue-tsc type-checking) from `frontend/`.
- [ ] Run `npm run test` (Vitest unit + SSR tests) from `frontend/`.
- [ ] Run `scripts/precommit-run` from the repo root.