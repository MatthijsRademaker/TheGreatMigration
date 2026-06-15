## Purpose

This specification defines the shared-ui utility components that implement the control treatments described in `docs/design-system-v2.md` and correspond to section 8 of `designs/components.png`. These utilities are additive, token-driven, Reka-backed or token-composed Vue components in `frontend/src/shared/ui/` that consume the existing semantic design tokens from `frontend/src/app/styles.css`. They fill the gap between the current primitives (`Button`, `Badge`, `Card`, `Input`, `Tooltip`) and future dashboard/header/task features.

## ADDED Requirements

### Requirement: Checkbox SHALL provide tokenized checked and unchecked states

The `Checkbox` component SHALL wrap Reka UI's `CheckboxRoot` and `CheckboxIndicator` primitives. The checked state SHALL use the `bg-primary` token (design-system primary green `#1E6B3E`) for the fill background. The unchecked state SHALL use the `border-border` token (`#E5E7EB`) for the border. The component SHALL render a Lucide `Check` icon as the checked indicator and SHALL accept standard `modelValue` / `update:modelValue` v-model binding.

#### Scenario: Checkbox renders unchecked by default

- **WHEN** a `Checkbox` component is rendered without a `modelValue` or with `modelValue` set to `false`
- **THEN** the rendered HTML contains `border-border` class on the checkbox root element
- **AND** the Lucide `Check` indicator icon is not present in the output

#### Scenario: Checkbox renders checked state with green fill

- **WHEN** a `Checkbox` component is rendered with `modelValue` set to `true`
- **THEN** the rendered HTML contains `bg-primary` class on the checkbox root element
- **AND** the Lucide `Check` indicator icon is rendered

### Requirement: Select SHALL provide a dropdown control with listbox options

The `Select` component SHALL wrap Reka UI's `SelectRoot`, `SelectTrigger`, `SelectValue`, `SelectContent`, `SelectItem`, and related primitives. The `SelectTrigger` SHALL use `bg-background`, `border-border`, and `text-foreground` classes matching the design contract's "white surface, neutral border" treatment. It SHALL render a trailing Lucide `ChevronDown` icon. The `SelectContent` portal SHALL use `bg-popover`, `border-border`, and `text-popover-foreground` tokens.

#### Scenario: Select trigger renders with white surface and neutral border

- **WHEN** a `Select` component is rendered with any placeholder and no selection
- **THEN** the trigger element contains `bg-background` and `border-border` classes
- **AND** a trailing `ChevronDown` icon is present in the trigger

#### Scenario: Select content renders in a portal

- **WHEN** a `Select` component is mounted with option items
- **THEN** the component instantiates without throwing an error
- **AND** the `SelectContent` element uses `bg-popover` and `border-border` token classes

### Requirement: Popover SHALL provide a floating content container

The `Popover` component SHALL wrap Reka UI's `PopoverRoot`, `PopoverTrigger`, and `PopoverContent` primitives. The `PopoverContent` SHALL use `bg-popover`, `border-border`, `text-popover-foreground`, and `shadow-md` token classes. This component is a building block for the `DatePicker` composition and future popover-based controls.

#### Scenario: Popover trigger renders as a slot wrapper

- **WHEN** a `Popover` component is rendered with a trigger child element
- **THEN** the trigger element is rendered in SSR output
- **AND** the popover content (portal) does not appear in SSR

#### Scenario: Popover component mounts without error

- **WHEN** a `Popover` component is mounted with trigger and content slots
- **THEN** the Vue component instantiates without throwing an error

### Requirement: Calendar SHALL provide a date grid surface

The `Calendar` component SHALL wrap Reka UI's `CalendarRoot` primitive and include subcomponents for calendar header, grid, and day cells. It SHALL accept a `modelValue` of Reka's `DateValue` type for single-date selection and SHALL emit `update:modelValue` on date change. The calendar grid SHALL use token classes for typography (`text-body`, `text-label`), surfaces (`bg-popover`, `border-border`), and selected-state highlighting (`bg-primary`, `text-primary-foreground`).

#### Scenario: Calendar renders month grid with tokenized typography

- **WHEN** a `Calendar` component is rendered with a `DateValue` model value
- **THEN** the calendar header (month/year) uses `text-body` token class
- **AND** day labels use `text-label` token class
- **AND** the component instantiates without error

#### Scenario: Selected date uses primary green highlight

- **WHEN** a `Calendar` renders with a selected date matching the model value
- **THEN** the selected day cell has `bg-primary` and `text-primary-foreground` classes

### Requirement: Avatar SHALL render circular initial-based person representation

The `Avatar` component SHALL wrap Reka UI's `AvatarRoot`, `AvatarImage`, and `AvatarFallback` primitives. When no image `src` is provided, it SHALL render the first character of a `name` prop inside a `rounded-full` container using `bg-muted` background and `text-label` typography. When a valid `src` is provided, `AvatarImage` SHALL render instead of the fallback.

#### Scenario: Avatar renders initials fallback when no image src is set

- **WHEN** an `Avatar` component is rendered with `name` set to "Alex" and no `src` prop
- **THEN** the rendered HTML contains a `rounded-full` element with the character "A"
- **AND** the element uses `bg-muted` and `text-label` token classes

#### Scenario: Avatar renders image when src is provided

- **WHEN** an `Avatar` component is rendered with `src` set to a valid URL and `name` set to "Alex"
- **THEN** the `AvatarImage` subcomponent is present in the rendered output
- **AND** the `AvatarFallback` initials are not visible when the image loads

### Requirement: SearchField SHALL compose Input with a leading search icon

The `SearchField` component SHALL compose the existing `Input` primitive with a leading Lucide `Search` icon. It SHALL support `v-model` binding through the underlying `Input` and SHALL pass through `placeholder` and other standard input attributes. The leading icon SHALL use `text-muted-foreground` token class.

#### Scenario: SearchField renders with leading search icon

- **WHEN** a `SearchField` component is rendered with a placeholder
- **THEN** the rendered HTML contains a Lucide `Search` icon preceding the input element
- **AND** the icon uses `text-muted-foreground` class

#### Scenario: SearchField supports v-model

- **WHEN** a `SearchField` component is rendered with `modelValue` set to "query text"
- **THEN** the underlying input element has value "query text"

### Requirement: PersonChip SHALL render an avatar with adjacent person name

The `PersonChip` component SHALL compose the `Avatar` component with an adjacent person name text span. It SHALL use `bg-secondary` or `bg-muted` surface, `rounded-full` avatar, and `text-body` typography for the name. It SHALL accept an optional `as` prop (defaulting to `"span"`, following the Reka `Primitive` pattern) so consumers can render it as a `button` or `PopoverTrigger` for interactive use.

#### Scenario: PersonChip renders display-only by default

- **WHEN** a `PersonChip` component is rendered with `name` set to "Alex Smith" and no `as` prop
- **THEN** the root element is a `<span>`
- **AND** the rendered HTML contains a circular avatar element with the initial "A"
- **AND** the rendered HTML contains the text "Alex Smith" adjacent to the avatar
- **AND** the root element uses `bg-secondary` or `bg-muted` token class

#### Scenario: PersonChip renders as a button when as prop is set

- **WHEN** a `PersonChip` component is rendered with `as` set to `"button"`
- **THEN** the root element is a `<button>`
- **AND** the avatar and name are rendered inside the button

### Requirement: HelpButton SHALL provide an icon button with tooltip help text

The `HelpButton` component SHALL compose the existing `Button` (ghost variant, `icon-sm` size) wrapping a Lucide `CircleHelp` icon, wrapped with the existing `Tooltip` and `TooltipTrigger`. The `TooltipContent` SHALL render the value of a `tooltip` prop as help text. The component SHALL use token classes from the existing `buttonVariants` and `Tooltip` components without introducing new styling.

#### Scenario: HelpButton renders icon button with tooltip wrapper

- **WHEN** a `HelpButton` component is rendered with `tooltip` set to "How to use this feature"
- **THEN** the rendered HTML contains a `button` element with `data-variant="ghost"` and `data-size="icon-sm"` attributes
- **AND** the Lucide `CircleHelp` icon is rendered inside the button
- **AND** the tooltip trigger wraps the button

#### Scenario: HelpButton tooltip content renders the help text

- **WHEN** a `HelpButton` component is mounted with `tooltip` set to "Click for help"
- **THEN** the `TooltipContent` element contains the text "Click for help"

### Requirement: DatePicker SHALL provide single-date selection via popover calendar

The `DatePicker` component SHALL compose `Popover` + `Calendar` + a trigger `Button`. The trigger `Button` SHALL display the selected date formatted as a readable string or a `placeholder` prop when no date is selected. It SHALL use `bg-background`, `border-border`, and `text-foreground` token classes matching the "white surface, neutral border" treatment. The calendar SHALL open on trigger click via `Popover`. Single-date selection SHALL be supported; date range is out of scope.

#### Scenario: DatePicker trigger shows placeholder when no date selected

- **WHEN** a `DatePicker` component is rendered with `placeholder` set to "Select date" and no `modelValue`
- **THEN** the trigger button displays "Select date"
- **AND** the trigger uses `bg-background` and `border-border` token classes

#### Scenario: DatePicker trigger shows selected date

- **WHEN** a `DatePicker` component is rendered with a `CalendarDate` modelValue
- **THEN** the trigger button displays the formatted date string
- **AND** the component mounts without error

#### Scenario: DatePicker calendar panel renders on mount

- **WHEN** a `DatePicker` component is mounted
- **THEN** the `Calendar` subcomponent is present in the component tree
- **AND** the component instantiates without throwing an error

### Requirement: New utilities SHALL NOT introduce raw hex styling or page-specific CSS

All new utility components SHALL consume only the semantic design tokens defined in `frontend/src/app/styles.css` (e.g., `bg-primary`, `text-foreground`, `border-border`, `shadow-sm`, `text-label` token classes). No raw hex color values, inline `style` attributes with colors, or page-specific `<style>` blocks SHALL be introduced.

#### Scenario: Utility component classes use semantic tokens only

- **WHEN** any new utility component's source file is inspected
- **THEN** the Tailwind classes reference only theme token names (e.g., `bg-primary`, `text-muted-foreground`, `border-border`)
- **AND** no raw hex values (e.g., `#1E6B3E`) appear in component `<template>` or `<style>` blocks

### Requirement: New utilities SHALL be importable through the @/shared/ui barrel

All new utility components SHALL be re-exported from their parent barrel export(s) so consumers can import them through the established `@/shared/ui` alias path. No consumer SHALL need to use deep import paths like `@/shared/ui/checkbox/Checkbox.vue`.

#### Scenario: Utilities are importable from barrel exports

- **WHEN** a consumer imports `Checkbox`, `Select`, `Popover`, `Calendar`, `Avatar`, `SearchField`, `PersonChip`, `HelpButton`, or `DatePicker` from `@/shared/ui`
- **THEN** the import resolves to the correct component file without errors

### Requirement: Tests SHALL verify utility contracts using SSR and instantiation patterns

SSR `renderToString` tests SHALL verify trigger rendering and tokenized class output for each utility following the existing `design-system-primitives.test.ts` pattern. Portal-based components (`SelectContent`, `DatePicker` calendar panel, `PopoverContent`) SHALL be verified with component instantiation tests (mount without error) since portal content is client-side only. Existing tests for `Button`, `Badge`, and `Card` SHALL continue to pass.

#### Scenario: SSR tests verify tokenized classes for non-portal utilities

- **WHEN** the test suite is executed
- **THEN** `renderToString` tests exist for `Checkbox`, `Select` trigger, `SearchField`, `PersonChip`, `HelpButton`, and `DatePicker` trigger
- **AND** each test asserts at least one semantic token class in the rendered output

#### Scenario: Portal-based components instantiate without error

- **WHEN** the test suite is executed
- **THEN** at least one test creates a Vue app containing a portal-based component (`Select` with content, `DatePicker` with calendar, or `Popover` with content)
- **AND** the component instantiates without throwing an error

### Requirement: No showcase or demo route SHALL be added

The implementation SHALL NOT add a `/showcase`, `/components`, Storybook, or any other developer-only catalog route or page. Component verification SHALL rely exclusively on automated tests.

#### Scenario: No catalog surface exists after implementation

- **WHEN** the project's route configuration and page files are inspected after implementation
- **THEN** no new route, page, or layout file exists that serves as a component catalog or showcase
