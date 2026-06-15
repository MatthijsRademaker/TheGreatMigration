## Why

The project has a tokenized design foundation (`docs/design-system-v2.md`, `frontend/src/app/styles.css`) and a small shared-ui inventory (`button`, `badge`, `card`, `input`, `separator`, `sheet`, `sidebar`, `tooltip`), but it still lacks the reusable utility controls shown in section 8 of `designs/components.png`. Upcoming dashboard, header, and task work is blocked or at risk of one-off markup for dropdown/select, date picker, search field, checkbox, person/avatar chips, and help affordances. This change adds the missing utilities as additive, token-driven shared-ui components without rebuilding full screen compositions or introducing a showcase route.

## What Changes

Six new shared-ui utility components are added to `frontend/src/shared/ui/`, plus four supporting Reka-backed primitives, all consuming the existing semantic design tokens from `frontend/src/app/styles.css`:

- **Checkbox** — Reka `CheckboxRoot`/`CheckboxIndicator` wrapper with green-checked / neutral-unchecked token styling.
- **Select** — Reka `SelectRoot`/`SelectTrigger`/`SelectContent`/`SelectItem` dropdown with white surface, neutral border, and trailing chevron per the design contract.
- **Popover** — Reka `PopoverRoot`/`PopoverTrigger`/`PopoverContent` wrapper, a required building block for date-picker and future compositions.
- **Calendar** — Reka `CalendarRoot` wrapper providing the date grid surface, styled with existing popover/calendar color and typography tokens.
- **Avatar** — Reka `AvatarRoot`/`AvatarImage`/`AvatarFallback` wrapper for circular initials rendering, used by the person chip.
- **SearchField** — Compose the existing `Input` primitive with a leading Lucide `Search` icon, supporting `v-model` and placeholder text.
- **PersonChip** — Compose `Avatar` + person name label, display-only by default with optional `as` prop for interactive use via `PopoverTrigger` or `button`.
- **HelpButton** — Compose `Button` (ghost variant, `icon-sm` size) wrapping a Lucide `CircleHelp` icon, wrapped with the existing `Tooltip` for hover context.
- **DatePicker** — Compose `Popover` + `Calendar` + trigger `Button` for single-date selection via Reka's `DateValue`/`CalendarDate` type system.

Existing `Button` variants cover Primary and Secondary treatments; no duplicate button components are introduced. `@internationalized/date` is added as an explicit dependency (required peer of Reka Calendar/DatePicker). Barrel exports are updated so all utilities are importable through `@/shared/ui`. SSR/Vitest contract tests extend the existing `design-system-primitives.test.ts` pattern. No showcase/storybook surface, no AppShell wiring, and no dashboard composition work is included.

## Impact

- **Dependencies**: `@internationalized/date` added to `frontend/package.json`.
- **Affected code**: New directories under `frontend/src/shared/ui/` (`checkbox/`, `select/`, `popover/`, `calendar/`, `avatar/`, `search-field/`, `person-chip/`, `help-button/`, `date-picker/`). Barrel export updates in existing `index.ts` files.
- **Tests**: `frontend/tests/design-system-primitives.test.ts` extended with SSR contract coverage for new utilities.
- **Non-breaking**: All new components are additive. Existing primitives, routes, and screens are unchanged.
- **Out of scope**: AppShell header integration, dashboard composition rebuilds, showcase/storybook surfaces, date-range support, backend behavior, domain workflows.