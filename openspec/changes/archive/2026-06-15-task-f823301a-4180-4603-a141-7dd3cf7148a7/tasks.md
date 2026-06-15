## 1. Dependency Setup

- [x] 1.1 Add `@internationalized/date` as explicit dependency in `frontend/package.json`
- [x] 1.2 Run `bun install` to lock the new dependency
- [x] 1.3 Verify `vue-tsc --noEmit` passes after dependency addition

## 2. Reka-Backed Primitives

- [ ] 2.1 Add `Checkbox` primitive wrapping Reka `CheckboxRoot`/`CheckboxIndicator` with green-checked / neutral-unchecked token styling, `cva` variants, and barrel export
- [ ] 2.2 Add `Select` primitive wrapping Reka `SelectRoot`/`SelectTrigger`/`SelectContent`/`SelectItem`/`SelectValue` with white-surface trigger, neutral border, trailing chevron, and barrel export
- [ ] 2.3 Add `Popover` primitive wrapping Reka `PopoverRoot`/`PopoverTrigger`/`PopoverContent` with token styling and barrel export
- [ ] 2.4 Add `Calendar` primitive wrapping Reka `CalendarRoot` with calendar subcomponents (header, grid, cell) and barrel export
- [ ] 2.5 Add `Avatar` primitive wrapping Reka `AvatarRoot`/`AvatarImage`/`AvatarFallback` for circular initials with barrel export

## 3. Composed Utilities

- [ ] 3.1 Add `SearchField` utility composing existing `Input` with leading Lucide `Search` icon, forwarding `v-model` and placeholder
- [ ] 3.2 Add `PersonChip` utility composing `Avatar` + person name label with optional `as` prop and barrel export
- [ ] 3.3 Add `HelpButton` utility composing Button (ghost, `icon-sm`) with Lucide `CircleHelp` icon, wrapped with existing `Tooltip`
- [ ] 3.4 Add `DatePicker` utility composing `Popover` + `Calendar` + trigger `Button` for single-date selection with barrel export

## 4. Barrel Export Updates

- [ ] 4.1 Update `frontend/src/shared/ui/index.ts` (or equivalent barrel) to export all new primitives and composed utilities

## 5. Test Coverage

- [ ] 5.1 Extend `frontend/tests/design-system-primitives.test.ts` with SSR `renderToString` coverage for `Checkbox` (checked/unchecked state classes)
- [ ] 5.2 Add SSR coverage for `Select` trigger rendering and variant classes
- [ ] 5.3 Add SSR coverage for `SearchField` (icon presence, `v-model` contract, placeholder)
- [ ] 5.4 Add SSR coverage for `PersonChip` (initials fallback, name rendering, tokenized surface classes)
- [ ] 5.5 Add SSR coverage for `HelpButton` (tooltip trigger, icon rendering, ghost variant classes)
- [ ] 5.6 Add SSR coverage for `DatePicker` trigger rendering and variant classes
- [ ] 5.7 Add component instantiation tests (mount-without-error) for portal-based components (`SelectContent`, `DatePicker` calendar panel, `PopoverContent`)

## 6. Verification

- [ ] 6.1 Run `scripts/precommit-run` and ensure all checks pass
- [ ] 6.2 Confirm `vue-tsc --noEmit` passes with zero type errors
- [ ] 6.3 Confirm `vitest run` passes all existing and new tests
