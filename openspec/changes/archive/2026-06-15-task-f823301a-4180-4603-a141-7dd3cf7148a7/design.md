## Context

The repository now has a design-system-v2 token foundation (`docs/design-system-v2.md`, `frontend/src/app/styles.css`) and a core shared-ui inventory (`button`, `badge`, `card`, `input`, `separator`, `sheet`, `sidebar`, `tooltip`). The next step is to add the utility-level controls described in section 8 of `designs/components.png` so that upcoming dashboard, header, and task features can compose from existing pieces rather than building one-off markup.

This design covers the architecture and decisions for implementing those missing utilities as additive, token-driven shared-ui components inside `frontend/src/shared/ui/`.

## Goals / Non-Goals

### Goals
- Add reusable shared-ui utility components for: checkbox, dropdown/select, date picker, search field, person/avatar chip, and help affordance.
- Reuse existing `Button`, `Input`, and `Tooltip` primitives where applicable.
- Consume existing semantic design tokens from `frontend/src/app/styles.css` exclusively — no raw hex values or page-specific CSS.
- Follow the established shared-ui pattern: Reka `Primitive` wrappers, `cva` variant definitions, `cn()` utility composition, and barrel exports.
- Verify through SSR/Vitest contract tests extending the existing `design-system-primitives.test.ts` pattern.

### Non-Goals
- AppShell header integration or wiring utilities into existing app surfaces.
- Rebuilding the full dashboard, task table, sidebar, or header compositions from `designs/components.png`.
- Adding a `/showcase`, Storybook, or any developer-only catalog surface.
- Backend behavior, persistence, filtering logic, notifications, or person/task domain workflows.
- Date-range selection support (single-date only for this change).

## Decisions

### D1: Two-layer implementation — primitives then compositions
All new components follow a two-layer model: first, thin Reka-backed primitives (`Checkbox`, `Select`, `Popover`, `Calendar`, `Avatar`) that wrap the corresponding `reka-ui` roots; then, composed utilities (`SearchField`, `PersonChip`, `HelpButton`, `DatePicker`) that combine primitives and/or existing shared-ui components. This keeps each file focused, minimizes duplication, and mirrors the existing `Tooltip` family pattern.

### D2: Dropdown maps to full Reka Select (not styled button trigger)
The design contract describes "Dropdown / date control: white surface, neutral border, leading icon when useful, trailing chevron." While the text also conflates dropdown with date controls, the consensus across refinement participants is that the image context warrants a full `Select` control with listbox, option items, and keyboard navigation — not just a chevron-style button. The `Select` component wraps Reka's `SelectRoot`/`SelectTrigger`/`SelectContent`/`SelectItem` primitives, with the trigger styled using token classes matching the design contract (white surface, neutral border, trailing chevron).

### D3: Date picker is single-date Popover+Calendar composition
Single-date selection only. Date range is out of scope — it roughly doubles complexity (two calendars, start/end linking, overlap handling). Consumers needing range can compose two `DatePicker` instances or a dedicated `DateRangePicker` can be added later if the top-bar design demands it. The composition uses Reka's `PopoverRoot` + `CalendarRoot` + a trigger `Button`.

### D4: PersonChip is display-only with optional `as` prop for future interactivity
The design contract describes "Circular initials avatar, Adjacent person name. Compact, calm presentation." No click behavior is specified. The component defaults to a `span` root but accepts an optional `as` prop (following the existing `Button`/`Badge` Reka `Primitive` pattern) so consumers can later render it as a `button` or `PopoverTrigger` without refactoring.

### D5: HelpButton is icon-only Button + Tooltip
No standalone help popover, menu, or routing component. The simplest interpretation — a ghost `Button` (size `icon-sm`) wrapping Lucide `CircleHelp`, wrapped with the existing `Tooltip` for hover text — keeps scope minimal and matches the calm/approachable tone. The Tooltip already exists in `shared/ui/tooltip/`.

### D6: SearchField composes existing Input with leading icon
Rather than modifying `Input.vue` or creating a parallel styling path, `SearchField` wraps the existing `Input` component and prepends a Lucide `Search` icon. It forwards `v-model` through the underlying `Input` and passes through `placeholder` and other standard input attributes.

### D7: Calendar requires a standalone primitive before DatePicker composition
The Reka `CalendarRoot` is a distinct primitive that needs its own wrapper (`Calendar.vue`, calendar cell/header subcomponents) before the `DatePicker` composition can consume it. This follows the same pattern as `Popover` being a prerequisite for `DatePicker`.

### D8: @internationalized/date is added as an explicit dependency
Reka UI's Calendar/DatePicker primitives depend on `@internationalized/date` types (`CalendarDate`, `DateValue`). While it may be a transitive dependency through `reka-ui`, build tooling (particularly `vue-tsc`) requires it as an explicit dependency for type resolution. Adding it to `frontend/package.json` ensures reliable builds.

### D9: SSR tests cover triggers; portal content uses mount-without-error coverage
Reka `Popover`/`Portal`-based content (`SelectContent`, `DatePicker` calendar panel) renders only its trigger in SSR. The test strategy uses `renderToString` for trigger rendering and variant class verification, and a separate `createSSRApp` + mount instantiation test that confirms portal-containing components don't throw during setup. This matches the risk mitigation in R3 and preserves the existing test pattern.

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| `@internationalized/date` type resolution fails at build | Medium | High — blocks date-picker and calendar work | Added as explicit dependency (D8); verify with `vue-tsc --noEmit` |
| PersonChip needs image avatars in future | Medium | Medium — current initials-only approach needs refactoring | Architecture allows additive `Avatar` updates; current contract only specifies initials |
| Reka Select naming collision with native `<select>` | Low | Low | PascalCase exports (`Select`, `SelectTrigger`) follow existing shadcn-vue convention |
| Portal-based component tests false-pass on SSR | Medium | Medium — interactive states (open/close, selection) untested | Documented as accepted limitation (D9); interactive tests deferred to Playwright/component test coverage in future work |
| `vue-tsc` strictness rejects Reka Calendar types | Low | Medium | Reka-ui 2.9.10 has stable Calendar types; existing Tooltip pattern proves Reka types work with the project's TS config |

## Traceability

- Task: `f823301a-4180-4603-a141-7dd3cf7148a7`
- Design contract: `docs/design-system-v2.md` — control guidance, palette, typography, spacing, radius, shadows
- Theme surface: `frontend/src/app/styles.css` — semantic CSS variable tokens
- Existing primitives: `frontend/src/shared/ui/button/` (Button variants), `frontend/src/shared/ui/input/` (Input), `frontend/src/shared/ui/tooltip/` (Tooltip family)
- Test pattern: `frontend/tests/design-system-primitives.test.ts`
- Dossier: `2026-06-15T19:26:56.044Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Rounds: 1 (architect, lead-dev, reviewer)
