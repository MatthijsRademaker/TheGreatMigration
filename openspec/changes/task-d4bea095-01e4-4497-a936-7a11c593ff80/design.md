## Context

The change target is the current Vue 3 + Vite + shadcn-vue + Tailwind v4 application. `components.json` already points shadcn-vue theming at `src/app/styles.css`, and current screens consume semantic tokens through shared primitives rather than page-specific hardcoded colors. The blocker is that the design-system source of truth is still image-first while the executing LLM for later implementation is text-only.

## Goals

- Convert `designs/design-system.png` into a complete text-first design contract.
- Preserve the product direction from the project vision: calm, organized, approachable, practical, and green/nature-inspired for move planning.
- Consolidate the contract into semantic Tailwind v4/CSS-variable tokens in `src/app/styles.css`.
- Define reusable semantics for future generic components without creating those components in this iteration.
- Keep future implementation possible without any image inspection.

## Non-Goals

- Do not build the dashboard/main-screen compositions shown in `designs/components.png`.
- Do not add move-planning product behavior such as scheduling, persistence, authentication, invitations, roles, or collaboration.
- Do not add new reusable/domain Vue component directories.
- Do not bypass shadcn-vue or replace the current Tailwind v4 setup.
- Do not redesign dark mode beyond preserving functional parity and contrast.
- Do not create a new logo/SVG asset in this iteration.

## Decisions

### 1. Authoritative artifact location

The durable text-only design contract lives in the repository as `docs/design-system-v2.md`. OpenSpec artifacts define the change, but the repository Markdown file is the implementation deliverable future text-only executors will read as the canonical visual contract.

### 2. Text-first design contract content

`docs/design-system-v2.md` must capture the full board in text, including the exact palette, typography, control semantics, navigation states, status semantics, spacing/radius/shadow guidance, and explicit scope boundaries.

#### Brand direction

- Product name: **The Great Migration**
- Tone: calm, organized, approachable, practical, green/nature-inspired
- Wordmark: text-forward treatment using the product name
- Motifs: migrating birds, soft botanical sprig, and gentle landscape illustration language
- This iteration documents the brand treatment; it does not create a new asset file

#### Color palette

| Group | Name | Value | Intended semantic use |
| --- | --- | --- | --- |
| Brand greens | Primary Green | `#1E6B3E` | Primary actions and strongest brand emphasis |
| Brand greens | Secondary Green | `#2F8F57` | Secondary brand emphasis and supportive highlights |
| Brand greens | Tertiary Green | `#66B88A` | Soft accent surfaces and lighter positive emphasis |
| Brand greens | Soft Green | `#E6F3EA` | Active fills, subtle emphasis backgrounds, selected navigation states |
| Brand greens | Pale Green | `#F3FAF5` | App/page background and very light supportive surfaces |
| Neutrals | Text Primary | `#111827` | Primary text |
| Neutrals | Text Secondary | `#6B7280` | Secondary/supporting text |
| Neutrals | Border | `#E5E7EB` | Borders and dividers |
| Neutrals | Surface / Card | `#FFFFFF` | Card and floating surface background |
| Neutrals | Background | `#F7FAF8` | Overall page background |
| Status | Danger / High | `#EF4444` | High-priority and destructive states |
| Status | Warning / Medium | `#F59E0B` | Warning and medium-priority states |
| Status | Success / Low | `#22C55E` | Success and low-priority states |
| Status | Info | `#3882F6` | Informational emphasis |

#### Typography scale

Use Inter throughout the system.

| Role | Spec |
| --- | --- |
| Display / Title | `32px`, `600`, letter-spacing `-0.5px` |
| Section Heading | `20px`, `600`, letter-spacing `-0.2px` |
| Card Title | `16px`, `600`, letter-spacing `0px` |
| Body Text | `14px`, `400`, letter-spacing `0px` |
| Caption | `12px`, `400`, letter-spacing `0px` |
| Label | `12px`, `500`, letter-spacing `0.5px` |

#### Control and state guidance

- **Primary button**: filled primary green surface, white text, compact rounded control.
- **Secondary button**: white surface, neutral border, dark text, dropdown-capable treatment.
- **Ghost button**: minimal surface with green text for quiet actions.
- **Filter button**: outline treatment with leading filter icon; prefer existing Button styling before inventing a dedicated component.
- **Date picker / dropdown**: white surface, border, leading icon where applicable, trailing chevron.
- **Checkbox**: green checked state, neutral unchecked state.
- **Toggle**: green “on” track, gray “off” track.
- **Text input / search field**: white surface, subtle border, secondary placeholder text, leading search icon for search.
- **Priority badges**: pill styles for High, Medium, Low.
- **Availability chips**: pill styles for Available, Busy, Partial, Off.
- **Avatar / person chip**: circular initials avatar with adjacent name on neutral surface.
- **Navigation**: active sidebar state uses soft-green background with green icon/text; inactive state is neutral with subtle hover; top bar includes date range, arrow controls, notification badge, and profile chip.

#### Layout, spacing, radius, and elevation guidance

- Radius scale shown on the board: `4px`, `12px`, `16px`
- Spacing rhythm shown on the board: `4px`, `8px`, `12px`, `24px`, `32px`, `40px`
- Grid guidance: `12-column` desktop grid
- Shadow scale shown on the board:
  - `sm`: `0 1px 2px rgba(16, 24, 40, 0.06)`
  - `md`: `0 4px 12px rgba(16, 24, 40, 0.08)`
  - `lg`: `0 8px 24px rgba(16, 24, 40, 0.10)`
- Cards and panels use white surfaces, subtle borders, soft elevation, and rounded corners derived from the documented radius tokens

### 3. Token implementation surface

`src/app/styles.css` remains the only global Tailwind theme surface. The implementation must keep semantic token names, map the documented source hex values into functional OKLCH CSS variables, and expose matching `@theme inline` entries for consumers. Source hex values must remain visible in `docs/design-system-v2.md`, and may also be preserved in comments where helpful.

At minimum, the token contract must cover:

- background / surface / card / popover
- foreground / muted-foreground / label-scale text semantics
- primary / secondary / accent and their foregrounds
- destructive / warning / success / info
- border / input / ring
- sidebar palette parity
- typography scale tokens
- radius tokens derived from the documented scale
- spacing tokens aligned to the documented rhythm
- shadow tokens aligned to the documented elevation levels

### 4. Primitive alignment boundary

No new reusable/domain Vue component directories are created in this change. Component-level work is limited to existing primitives and only to the extent needed to expose the documented foundation cleanly.

- `src/shared/ui/badge/index.ts` may add named priority and availability semantics if token-only usage would otherwise require one-off classes.
- `src/shared/ui/card/Card.vue` may be aligned to the documented radius, spacing, text, border, and shadow tokens.
- `src/shared/ui/button/index.ts` should reuse existing variants for the documented button treatments wherever possible; only add a minimal new variant if the filter treatment cannot be expressed without it.
- Existing primitive APIs must stay compatible unless the change explicitly introduces additive variants.

### 5. Light and dark mode handling

The design board is light-mode only. Light mode is the v2 source of truth. Dark mode must be verified and adjusted only as needed to preserve contrast and current behavior after the light-mode token update. A full dark-mode redesign is intentionally deferred.

## Risks

- **Token cascade risk:** `styles.css` changes affect all current pages and primitives at once.
- **Scope creep risk:** future-dashboard compositions could accidentally slip into the change unless tasks keep the boundary explicit.
- **Primitive-boundary risk:** variant additions must remain additive and must not turn into new component files.
- **Dark-mode regression risk:** changing light tokens without parity checks can break contrast combinations in existing screens.

## Conflict Resolution

- **Durable contract location:** resolved to `docs/design-system-v2.md` to satisfy the text-only executor requirement.
- **Primitive extension ambiguity:** resolved by allowing additive updates to existing primitive variants/classes, while forbidding new Vue component directories.
- **Filter-button ambiguity:** resolved by preferring existing Button variants first; only add a new variant if the documented treatment cannot otherwise be represented.
- **Dark-mode ambiguity:** resolved by preserving parity rather than introducing a new dark-mode art direction.
- **Brand asset ambiguity:** resolved by documenting wordmark/motif guidance only; no SVG/logo implementation is included.

## Traceability

- Task: `d4bea095-01e4-4497-a936-7a11c593ff80`
- Dossier: `2026-06-15T13:52:01.162Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`
- Validated round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial` snapshot for `task-d4bea095-01e4-4497-a936-7a11c593ff80`