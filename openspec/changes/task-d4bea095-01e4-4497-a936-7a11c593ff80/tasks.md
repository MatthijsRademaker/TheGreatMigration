## 1. Publish the text-only design contract

- [ ] Create `docs/design-system-v2.md` as the authoritative repository design-system document.
- [ ] Transcribe the design board into text so implementation no longer depends on `designs/design-system.png`.
- [ ] Include the exact palette values: `#1E6B3E`, `#2F8F57`, `#66B88A`, `#E6F3EA`, `#F3FAF5`, `#111827`, `#6B7280`, `#E5E7EB`, `#FFFFFF`, `#F7FAF8`, `#EF4444`, `#F59E0B`, `#22C55E`, `#3882F6`.
- [ ] Include the exact typography scale: Inter; display/title `32px 600 -0.5px`; section heading `20px 600 -0.2px`; card title `16px 600`; body `14px 400`; caption `12px 400`; label `12px 500 0.5px`.
- [ ] Document button/control guidance, navigation states, priority badges, availability chips, avatar/person-chip treatment, spacing rhythm, radius hierarchy, shadow levels, and 12-column layout guidance.
- [ ] State explicitly that `designs/components.png` and new reusable/domain component creation are out of scope for this change.

## 2. Align the global Tailwind theme surface

- [ ] Update `src/app/styles.css` as the single Tailwind v4 theme surface for the design-system foundation.
- [ ] Map the documented source hex values to functional semantic OKLCH variables while preserving the source hex values in the contract documentation and comments where useful.
- [ ] Add missing semantic tokens for status (`destructive`, `warning`, `success`, `info`), availability/supporting chip use, sidebar parity, typography scale, radius, spacing, and shadow levels.
- [ ] Extend the `@theme inline` block so semantic utilities can consume the new token set without raw color classes.
- [ ] Keep semantic naming aligned with shadcn-vue conventions rather than introducing brand-specific utility names.

## 3. Align existing primitives without creating new components

- [ ] Audit `src/shared/ui/button/index.ts`, `src/shared/ui/badge/index.ts`, and `src/shared/ui/card/Card.vue` against the documented contract.
- [ ] Reuse existing Button variants for primary, secondary, ghost, and filter-style controls wherever possible; only add an additive variant if the documented filter treatment cannot be expressed otherwise.
- [ ] Add named Badge semantics for priority and availability only if token-only usage would otherwise require one-off styling.
- [ ] Update Card styling to consume the documented radius, border, spacing, text, and shadow tokens where current hardcoded values would drift from the contract.
- [ ] Do not add any new reusable/domain Vue component directories or implement the generic/dashboard components from `designs/components.png`.

## 4. Preserve compatibility across current screens

- [ ] Verify `src/shared/layout/app-shell/AppShell.vue`, `src/shared/layout/app-sidebar/AppSidebar.vue`, `src/home/HomeView.vue`, `src/tasks/TasksView.vue`, `src/calendar/CalendarView.vue`, and `src/people/PeopleView.vue` continue to render through semantic tokens after the foundation update.
- [ ] Keep existing primitive APIs compatible except for clearly additive variant extensions.
- [ ] Treat light mode as the design-system source of truth and only adjust `.dark` tokens as needed to preserve contrast and current behavior.

## 5. Verify the change

- [ ] Run `scripts/precommit-run` after the token and primitive alignment work is complete.
- [ ] Confirm the final change leaves a reusable, text-only design foundation for later generic-component iterations.