## Context

The `PeopleAvailability` component is a shared Vue component rendered on both the home dashboard (`/`) in read-only mode and the `/people` management route in editable mode. The design references `designs/home-page.png` and `designs/components.png` define a target visual anatomy that is denser and more compact than the current implementation. The existing `docs/design-system-v2.md` contract governs token usage, badge semantics, spacing, and typography. Two prior design-alignment specs (`task-management-design-alignment`, `daily-schedule-design-alignment`) provide the established repository pattern for converting visual PNG-driven deltas into text-only implementation requirements.

The current component diverges from the design targets in five specific areas: (1) extraneous `CardDescription` and summary row, (2) verbose matrix anatomy with "Person" column header and abbreviated day labels, (3) legend inside the scroll container, (4) hard-coded adapter copy for description and summary, and (5) test assertions that lock in the current layout.

## Goals / Non-Goals

### Goals
- Produce a text-only design-alignment spec that describes the target People Availability panel well enough for a non-vision LLM to implement without inspecting PNG files.
- Strip `CardDescription` and summary row from the shared component template.
- Compact the matrix: remove "Person" column header, use date+weekday column headers, and use compact status pills per `design-system-v2.md`.
- Relocate the legend to a single footer row below the matrix, outside the scroll container.
- Update the shared adapter to stop producing now-unused `description`, `availableToday`, and `totalPeople` props.
- Preserve all editable-mode behaviour: Popover triggers, status selection, Actions column, Delete buttons, and confirmation dialogs.
- Update the three affected test suites to assert the new design structure.

### Non-Goals
- Changing backend endpoints, OpenAPI output, persistence, or availability status vocabulary.
- Redesigning the `/people` route chrome: create-person form, pagination bar, loading/error/empty state cards, or mutation wiring.
- Redesigning unrelated dashboard regions: Tasks Backlog, Daily Schedule, Move Notes, sidebar, or toolbar.
- Adding new product behaviour such as filters, search, bulk editing, drag-and-drop, or person pagination controls.
- Changing the Avatar primitive to support pastel background colours (aspirational polish, out of scope).
- Adding new props or emits to the `PeopleAvailability` component.
- Modifying the `KpiCards.vue` component on the home dashboard.

## Decisions

All decisions are resolved from the refinement-room consensus across architect, lead-dev, and reviewer recommendations.

### Decision 1: Summary row removal
**Resolution**: Remove the `X of Y available today` summary row from the shared component template unconditionally. Keep `availableToday` and `totalPeople` in the `PeopleAvailabilityProps` type interface to avoid breaking the adapter's type contract, but stop passing them as visible props in the adapted output. The `KpiCards.vue` component on the home dashboard already surfaces the available-today count in a richer format (`6 / 8 available on Jul 5`), so the summary row is redundant on both routes.

**Rationale**: All three refinement agents converge on removal. The architect and lead-dev explicitly recommend stripping it; the reviewer flagged its absence as a blocker requiring resolution. The home dashboard KPI cards already render this information. On `/people`, the matrix itself is the primary data surface and the summary row adds visual noise that the design references do not include.

### Decision 2: CardDescription removal
**Resolution**: Remove the `CardDescription` from the shared component template. Remove the hard-coded `description` field from `adaptToComponentProps` in `usePeopleAvailability.ts`. Remove `description` from `PeopleAvailabilityProps` in `types.ts`.

**Rationale**: The design references show a title-only card header. The `/people` route's own loading/error/empty state cards use separate `CardDescription` instances that are not affected. The architect and lead-dev both recommend removal; the reviewer flagged it as a blocker needing explicit resolution.

### Decision 3: Legend placement and format
**Resolution**: Render the legend as a single inline-flex footer row below the matrix table, outside the `overflow-x-auto` container. Order: Available, Busy, Partial, Off (matching current order and design-system-v2.md semantics). Each legend item is a compact `Badge` with its corresponding variant.

**Rationale**: The `designs/home-page.png` reference shows a footer legend row. The current inline legend inside the scroll container is visually inconsistent with the design. The lead-dev recommendation explicitly calls for relocating to a footer row.

### Decision 4: Date header format
**Resolution**: Use the date+weekday format produced by the adapter's `formatPlanDayLabel` (e.g., "Sun 5 Jul"). Update the component's default `days` prop to empty (or keep short labels as demo-only defaults — the adapter always overrides them in production).

**Rationale**: The dossier describes "date+weekday column labels" from the design reference. The adapter already produces this format. The lead-dev recommendation calls for date+weekday column headers. The architect confirms the adapter-derived format should be canonical.

### Decision 5: Person column header
**Resolution**: Remove the "Person" `<th>` label. The first column retains avatar+name cells but has no visible column header text.

**Rationale**: The design references show no column label for the person column — only compact avatar/name rows. The lead-dev explicitly recommends removing the Person column header.

### Decision 6: Status cell rendering
**Resolution**: Keep the current capitalized-text `Badge` elements ("Available", "Busy", "Partial", "Off") with compact sizing. Do NOT switch to textless colour-only pills. Maintain screen-reader-visible text labels as currently rendered.

**Rationale**: `design-system-v2.md` defines availability chip semantics as "soft semantic surfaces" — not textless dots. The existing accessibility contract (screen-reader-visible text) is correct. The architect questioned textless pills but recommended keeping labels under the design-system contract.

### Decision 7: Title casing
**Resolution**: Keep the lowercase form "People availability" to match existing route metadata, canonical spec headings, and the adapter default. Do NOT change to "People Availability".

**Rationale**: Changing the casing would ripple through route metadata (`app-routes-render.test.ts` path metadata), multiple test assertions, and the adapter default. The canonical specs use lowercase. The risk/cost of casing change outweighs any visual fidelity benefit.

### Decision 8: Avatar background colours
**Resolution**: Pastel-variant avatar backgrounds shown in the design PNGs are aspirational visual polish and NOT in scope for this task. The current `Avatar` primitive uses `bg-muted` only and remains unchanged.

**Rationale**: All agents agree this is polish, not core design alignment. The Avatar primitive would need a separate enhancement cycle.

### Decision 9: Scope boundary — `/people` route chrome
**Resolution**: The `/people` route's create-person form, pagination bar, loading/error/empty state `Card` wrappers, and mutation wiring are explicitly OUT of scope. Only the shared `PeopleAvailability` component and its adapter are changed.

**Rationale**: The dossier non-goals explicitly exclude "redesigning unrelated dashboard regions." The reviewer flagged scope creep into management chrome as a risk and recommended fencing it out — which this decision formalises.

### Decision 10: Adapter prop cleanup
**Resolution**: Remove `description` from the adapted output and from `PeopleAvailabilityProps`. Remove `availableToday` and `totalPeople` from the adapted output but keep them in the type interface (for backward compatibility of the type contract, even if unused in the template).

**Rationale**: The lead-dev recommendation calls for adapter cleanup. The architect recommends keeping the props in the type contract while removing template usage. This balanced approach avoids introducing type errors in consumer code while eliminating dead template code.

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| Test breakage in three suites if assertions not updated atomically | Medium | List all three test files in the spec with exact assertion changes; update them in the same change set. |
| Adapter output change breaks consumer `v-bind` spread if props are removed from types.ts | Medium | Keep `availableToday` and `totalPeople` in the interface even while removing them from the adapted output; this avoids prop-mismatch warnings. |
| Editable-mode Popover/Dialog tests break if template changes touch their DOM structure | Low | The editable section of the template (v-if="editable" blocks) is not being restructured — only the surrounding card shell and static rows change. |
| `app-routes-render.test.ts` `/people` route assertion `6 of 8 available today` must be replaced with a stable alternative | Medium | Replace with assertions on the legend footer badges and matrix structure (person names, status variants). |

## Traceability

| Source | Decision/Requirement |
|--------|---------------------|
| Dossier acceptance criteria | Non-placeholder text-only brief; explicit delta mapping; frontend-only scope; preserved vocabulary; documented open questions |
| Architect round 1 | Design-alignment spec pattern; summary row removal; CardDescription removal; type contract preservation; test file enumeration |
| Lead-dev round 1 | Strip CardDescription + summary; compact matrix anatomy; date+weekday headers; footer legend; adapter cleanup; preserve editable mode |
| Reviewer round 1 | Five blockers resolved: summary row, CardDescription, test fallout, legend placement, date header format |
| `docs/design-system-v2.md` | Availability chip semantics (Available=success, Busy=destructive, Partial=warning, Off=muted); spacing/typography/radius tokens |
| `openspec/specs/task-management-design-alignment/spec.md` | Established spec pattern: per-element requirements with scenario-based acceptance criteria |
| `openspec/specs/daily-schedule-design-alignment/spec.md` | Established spec pattern: explicit prop/emit contracts, test assertion updates |
| `openspec/specs/dashboard-people-availability/spec.md` | Canonical status vocabulary and backend data contract |
| `openspec/specs/people-availability-integration/spec.md` | Shared frontend data flow and read-only vs editable route split |
