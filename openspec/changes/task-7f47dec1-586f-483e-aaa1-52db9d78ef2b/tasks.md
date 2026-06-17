## 1. Refactor card config to drive all four KPI cards from a unified data structure

- [ ] 1.1 Extend the current `BacklogCardConfig` interface to a unified `KpiCardConfig` covering all four cards, adding fields: `id` (for the people/backlog/placeholder source), `accentClass` (e.g., `'bg-info-soft text-info'`), and `status` (for loading/error/ready/empty state)
- [ ] 1.2 Replace the separate `availabilityStatus`/`backlogStatus` computed properties and the inline People-available card with a single `cardConfigs` computed that produces all four KpiCardConfig entries, mapping each to its respective data source and status
- [ ] 1.3 Define the semantic color mapping per card: People → `bg-info-soft text-info`, High priority → `bg-destructive-soft text-destructive`, Unassigned → `bg-warning-soft text-warning`, Rooms → `bg-success-soft text-success`

## 2. Restructure each KPI card's template to the two-column layout

- [ ] 2.1 Replace each card's `CardHeader`/`CardContent` sub-components with a single `div.flex.flex-row` child inside the `Card`, containing a left accent column (`div.flex.w-[72px].shrink-0.items-center.justify-center.self-stretch`) and a right content column (`div.flex.flex-1.flex-col.p-panel`)
- [ ] 2.2 Override the Card's default padding/gap by adding `!py-0 !gap-0` to each Card instance, with a code comment explaining why the override is needed
- [ ] 2.3 Add `rounded-l-lg` to the left accent column so its left edges match the Card's `rounded-lg`, and apply the per-card `accentClass` for background/icon color
- [ ] 2.4 Render the icon at `size-8` (32px) centered in the left accent column, replacing the current `size-10` circle wrapper
- [ ] 2.5 Render the title, KPI value, and subtitle as three stacked rows in the right content column: title at `text-caption text-muted-foreground`, KPI value at `text-3xl font-semibold`, subtitle at `text-sm text-muted-foreground`
- [ ] 2.6 Keep the leaf decoration (`absolute bottom-0 right-0 h-14 w-auto pointer-events-none`) in the same position relative to the Card

## 3. Wire loading, error, and empty states into the unified layout

- [ ] 3.1 Ensure the People-available card's loading state ("Loading…") and error state ("Backend unavailable") render in the KPI value row of the right content column
- [ ] 3.2 Ensure the backlog cards' loading state ("Loading…") and error state ("Backend unavailable") render in the KPI value row, sourced from the unified config's status field
- [ ] 3.3 Ensure the empty-state fallback for People-available (totalPeople === 0 → `—`) renders correctly in the new layout
- [ ] 3.4 Verify that the left accent column always renders (even in loading/error states) so the card's visual framing is consistent

## 4. Clean up and verify

- [ ] 4.1 Remove unused imports: `CardHeader`, `CardContent`, `CardTitle`, `CardDescription` from `@/shared/ui/card`
- [ ] 4.2 Remove the now-unused `backlogCards` computed and `BacklogCardConfig` interface
- [ ] 4.3 Remove the `bgClass` field from any remaining config (replaced by `accentClass`)
- [ ] 4.4 Verify the Rooms-completed card keeps `data-testid="kpi-placeholder-rooms-completed"` and the code comment documenting it as a placeholder

## 5. Update and verify tests

- [ ] 5.1 Run the existing SSR route-render tests (`npx vitest run frontend/tests/app-routes-render.test.ts` or equivalent) to verify label-text assertions still pass without structural changes to assertions
- [ ] 5.2 Run the full test suite to catch any other dependent test breakage
- [ ] 5.3 Run `scripts/precommit-run` and fix any failures
