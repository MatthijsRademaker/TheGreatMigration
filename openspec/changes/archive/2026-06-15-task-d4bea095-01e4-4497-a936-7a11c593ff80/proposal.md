## Why

The requested change is not to build the dashboard or generic components yet. It is to turn `designs/design-system.png` into an implementation-ready, text-only design foundation that a future executor can apply without reading images. The current Vue 3 + Vite + shadcn-vue + Tailwind v4 app already has semantic tokens and reusable primitives, but the durable visual contract is still split between image references and partial CSS variables.

## What Changes

- Publish a standalone repository design-system contract at `docs/design-system-v2.md` as the authoritative text-only source of truth.
- Transcribe the design-system board into text: brand direction, exact palette values, typography scale, button/control treatments, navigation states, status/badge semantics, avatar/person-chip guidance, spacing rhythm, radius hierarchy, shadow levels, and grid/layout guidance.
- Keep `src/app/styles.css` as the single global Tailwind v4 theme surface and align it to the documented contract by adding or updating semantic CSS variables and `@theme inline` entries for palette, typography, status, availability, sidebar, spacing, radius, and shadow tokens.
- Align existing primitives only where the documented semantics require it: update existing `Badge`, `Card`, and only-if-needed `Button` variants/classes so future generic components can consume the foundation without ad hoc styling.
- Explicitly keep `designs/components.png` and any new reusable/domain Vue component directories out of scope for this change.

## Impact

- Future component iterations can be implemented from repository text and source alone, without image inspection.
- Token updates will cascade through existing semantic consumers such as `AppShell`, `AppSidebar`, `HomeView`, `TasksView`, `CalendarView`, and `PeopleView`, so the change must preserve API compatibility and verify current rendering paths.
- Dark mode is preserved for parity and regression safety, but a full dark-mode v2 visual redesign is deferred.
- Verification must include `scripts/precommit-run` after the token and primitive alignment work is complete.
