## Why

The home dashboard KPI cards row (`frontend/src/home/components/KpiCards.vue`) does not match the design target defined in `designs/components.png` section 3 and `designs/home-page.png`. The current implementation uses a different card order (People first), a full-height `w-[72px]` left accent column instead of the design's thin left accent border + compact icon chip, generic descriptive copy instead of design-like outcome text, and different Lucide icons than the design intends. The backlog task description "make the KPI cards component closer to the design" is image-dependent and cannot be implemented by a text-only LLM. This change translates the visual delta into precise, text-only implementation requirements and updates the canonical OpenSpec, implementation, and tests in lockstep.

## What Changes

- Reorder the four KPI cards to match the design: **High priority tasks** (first), **People available today** (second), **Unassigned jobs** (third), **Rooms completed** (fourth).
- Restructure each KPI card layout from the current `w-[72px]` full-height left accent column + centered icon to a thin left accent border (`border-l-4` with semantic color) + compact rounded icon chip (`size-8 rounded-lg` with semantic soft background) positioned at the top-left of the content area.
- Adopt Card primitives (CardHeader, CardContent, CardTitle, CardDescription) structurally instead of bypassing them with `!py-0 !gap-0` overrides.
- Update the **People available today** card to display `X / Y` fraction format (e.g., "6 / 8") with a date-context subtitle (`available on Jul 5`) sourced from `range.selectedDate` in the availability API response. Remove the current "X of Y available" wording.
- Update card subtitle copy to design-like outcome text: `high priority tasks need attention` (high priority), `available on <MMM D>` dynamic from `range.selectedDate` (people), `jobs that need assignment` (unassigned), `rooms fully packed and cleared` (rooms, static placeholder).
- Replace current Lucide icons: `TriangleAlertIcon` → `FlagIcon` (high priority), `HammerIcon` → `BriefcaseIcon` (unassigned), `Building2Icon` → `CheckCircleIcon` (rooms). Keep `UsersRoundIcon` for people.
- Preserve the **Rooms completed** card as an isolated placeholder with `—` value, `data-testid="kpi-placeholder-rooms-completed"`, and no backend query or derived business logic.
- Update `frontend/src/home/__tests__/KpiCards.spec.ts` to assert new card order, new fraction format, new subtitle copy, thin accent border + icon chip anatomy, and loading/error state rendering inside the new layout. Remove assertions for the obsolete `w-[72px]` column and `X of Y available` wording.
- Update `frontend/tests/app-routes-render.test.ts` home-route assertions: replace `"of 8"` and `"available"` expectations with "6 / 8" and new people-card copy. The `/people` route assertion `"6 of 8 available today"` is unrelated and must NOT be modified.
- Update the canonical spec `openspec/specs/kpi-summary-cards/spec.md` to reflect the new card order, layout anatomy, value format, copy, and icon choices.
- Use only existing semantic tokens from `frontend/src/app/styles.css`. No raw hex values, no new global theme tokens, no new shared UI directories.
- Preserve existing loading and error state handling for people-availability and task-backlog queries within the new layout.

## Impact

- **KpiCards.vue**: Template, card configs, data destructuring, computed properties, and icon imports all change. Adding `range.selectedDate` access is a new data dependency within the same query response — no backend contract change.
- **Canonical spec** (`openspec/specs/kpi-summary-cards/spec.md`): Updated to reflect new card order, anatomy, value format, copy, and icon choices. Spec update is part of this change, not a separate dependency.
- **Unit tests** (`KpiCards.spec.ts`): Most assertions break and must be rewritten to match the new layout and copy contract.
- **SSR tests** (`app-routes-render.test.ts`): Home-route assertions for people KPI values break and must be updated. The `/people` route assertions remain unchanged.
- **Non-impacted**: Backend services, persistence, sidebar, top bar, lower HomeView sections (`Today's plan`, `Move notes`), people matrix, schedule board, rooms/areas page, settings page, design system token definitions, shared UI primitives.
