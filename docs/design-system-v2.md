# Design System v2

This document is the authoritative repository design contract for Design System v2. Future UI work should be implementable from this file and `src/app/styles.css` alone, without opening `designs/design-system.png`.

## Scope and source of truth

- Product: **The Great Migration**
- Tone: calm, organized, approachable, practical, green and nature-inspired
- Brand treatment: text-forward wordmark with gentle moving-birds, botanical, and landscape illustration language
- Global theme surface: `src/app/styles.css`
- Light mode is the v2 source of truth; dark mode only preserves parity and contrast

## Explicitly out of scope

This change does **not** include:

- compositions from `designs/components.png`
- new reusable Vue component directories
- new domain component directories
- dashboard implementation work beyond the existing primitives and semantic-token foundation
- a new SVG or logo asset

## Palette

All source values below are contractual and must remain visible in repository text.

| Group | Token intent | Hex | Usage |
| --- | --- | --- | --- |
| Brand green | Primary green | `#1E6B3E` | Primary actions, strongest emphasis |
| Brand green | Secondary green | `#2F8F57` | Focus rings, supportive brand emphasis |
| Brand green | Tertiary green | `#66B88A` | Softer positive emphasis, dark-mode primary lift |
| Brand green | Soft green | `#E6F3EA` | Active nav fills, subtle emphasis backgrounds |
| Brand green | Pale green | `#F3FAF5` | Calm surfaces, secondary fills, sidebar background |
| Neutral | Text primary | `#111827` | Primary text |
| Neutral | Text secondary | `#6B7280` | Supporting text |
| Neutral | Border | `#E5E7EB` | Borders, dividers, input chrome |
| Neutral | Surface | `#FFFFFF` | Cards, popovers, controls on light mode |
| Neutral | Page background | `#F7FAF8` | App background |
| Status | Danger / High | `#EF4444` | Destructive actions, high priority |
| Status | Warning / Medium | `#F59E0B` | Warnings, medium priority, partial availability |
| Status | Success / Low | `#22C55E` | Success, low priority, available state |
| Status | Info | `#3882F6` | Informational emphasis |

## Typography

Use **Inter** everywhere.

| Role | Size | Weight | Letter spacing | Guidance |
| --- | --- | --- | --- | --- |
| Display / title | `32px` | `600` | `-0.5px` | Large page hero moments |
| Section heading | `20px` | `600` | `-0.2px` | Page sections and major panels |
| Card title | `16px` | `600` | `0px` | Card headings |
| Body | `14px` | `400` | `0px` | Default UI copy |
| Caption | `12px` | `400` | `0px` | Supporting metadata |
| Label | `12px` | `500` | `0.5px` | Buttons, badges, compact controls |

## Control guidance

### Buttons and controls

- **Primary button**: filled primary green with white text, compact height, rounded medium corners, subtle shadow.
- **Secondary button**: white or pale-green surface, neutral border, dark text, suitable for dropdown triggers.
- **Ghost button**: quiet action with minimal surface and green text.
- **Filter-style control**: stays within the existing Button primitive; use outline styling with a leading filter icon.
- **Dropdown / date control**: white surface, neutral border, leading icon when useful, trailing chevron.
- **Text input / search field**: white surface, subtle border, secondary placeholder text, optional leading search icon.
- **Checkbox**: green checked state, neutral unchecked state.
- **Toggle**: green active track, subdued inactive track.

### Status, badge, and chip semantics

- **Priority badges**
  - High: destructive color family
  - Medium: warning color family
  - Low: success color family
- **Availability chips**
  - Available: success
  - Busy: destructive
  - Partial: warning
  - Off: muted/neutral
- Supporting chips should prefer soft semantic surfaces rather than ad hoc opacity mixes.

### Avatar and person chip

- Circular initials avatar
- Adjacent person name
- Neutral or pale-green supporting surface
- Compact, calm presentation over decorative treatment

## Navigation states

- Sidebar background uses the pale-green family.
- Active sidebar item uses a soft-green fill with green icon and text.
- Inactive sidebar items stay neutral and gain subtle hover emphasis.
- Top bar keeps date-range controls, directional controls, notification affordances, and profile/person chips visually quiet and practical.

## Layout and rhythm

### Grid

- Desktop layout uses a **12-column grid**.
- Cards and panels should align to that grid instead of inventing local widths.

### Spacing rhythm

- `4px`
- `8px`
- `12px`
- `24px`
- `32px`
- `40px`

Use the smaller values inside controls and chips, and the larger values for cards, sections, and page spacing.

### Radius hierarchy

- `4px`: tight corners for very compact controls
- `12px`: default control radius
- `16px`: cards, panels, and larger surfaces

### Shadow levels

- `sm`: `0 1px 2px rgba(16, 24, 40, 0.06)`
- `md`: `0 4px 12px rgba(16, 24, 40, 0.08)`
- `lg`: `0 8px 24px rgba(16, 24, 40, 0.10)`

Cards and floating surfaces should stay soft and restrained.

## Implementation contract

- Keep semantic naming aligned with shadcn-vue and Tailwind semantics.
- `src/app/styles.css` is the only global theme surface for Design System v2 tokens.
- Theme tokens must expose semantic color, typography, radius, spacing, and shadow utilities so pages and primitives do not rely on raw hex classes.
- Primitive updates stay inside the existing Button, Badge, and Card boundaries.
- This document should remain sufficient for later generic-component iterations without image inspection.
