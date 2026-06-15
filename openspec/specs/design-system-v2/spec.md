# design-system-v2 Specification

## Purpose
TBD - created by archiving change task-d4bea095-01e4-4497-a936-7a11c593ff80. Update Purpose after archive.
## Requirements
### Requirement: Design System v2 SHALL publish a text-only repository contract
The change SHALL add `docs/design-system-v2.md` as the authoritative repository design contract for Design System v2. The document MUST let a text-only implementer reproduce the visual foundation without opening any image file. It MUST include the documented brand direction, the exact palette values (`#1E6B3E`, `#2F8F57`, `#66B88A`, `#E6F3EA`, `#F3FAF5`, `#111827`, `#6B7280`, `#E5E7EB`, `#FFFFFF`, `#F7FAF8`, `#EF4444`, `#F59E0B`, `#22C55E`, `#3882F6`), the full Inter typography scale (`32px/600/-0.5px`, `20px/600/-0.2px`, `16px/600`, `14px/400`, `12px/400`, `12px/500/0.5px`), control treatments, navigation states, priority and availability semantics, avatar/person-chip guidance, spacing rhythm, radius hierarchy, shadow levels, and 12-column layout guidance.

#### Scenario: A future executor implements from text only
- **WHEN** a worker opens `docs/design-system-v2.md` to implement the next UI iteration
- **THEN** the worker can recover the complete palette, typography, state, and layout contract without inspecting `designs/design-system.png`

#### Scenario: Scope boundaries are explicit in the contract
- **WHEN** a worker reads the design-system contract for this change
- **THEN** the document states that `designs/components.png` compositions and new reusable/domain Vue component directories are out of scope for this iteration

### Requirement: The global theme SHALL expose Design System v2 semantics through `src/app/styles.css`
`src/app/styles.css` SHALL remain the single global Tailwind v4 theme surface for this foundation update. The implementation MUST map the documented source hex values to functional semantic CSS variables using the existing OKLCH-based approach, and it MUST expose matching `@theme inline` entries for the full light-mode design contract. The theme MUST cover background, card/surface, foreground, muted text, primary/secondary/accent, border/input/ring, destructive/warning/success/info, sidebar parity, typography-scale tokens, radius tokens, spacing tokens, and shadow tokens.

#### Scenario: Existing pages consume the updated tokens semantically
- **WHEN** current screens render after the token update
- **THEN** they continue to use semantic `bg-*`, `text-*`, `border-*`, and related theme utilities instead of page-specific raw color rewrites

#### Scenario: Dark mode preserves parity instead of receiving a redesign
- **WHEN** light-mode Design System v2 tokens are introduced
- **THEN** the `.dark` token block is only adjusted as needed to preserve contrast and current behavior rather than introducing a separate dark-mode redesign

### Requirement: Existing primitives SHALL align to the documented semantics without adding new component directories
The implementation SHALL keep component-level changes inside existing primitives only. It MUST NOT add new reusable or domain-specific Vue component directories for this change. Existing primitive files MAY receive additive alignment updates so future generic components can consume the foundation cleanly: Badge semantics MUST support priority and availability use cases if token-only usage would otherwise require ad hoc styling, Card styling MUST align to the documented radius/border/spacing/shadow/text tokens where needed, and Button styling MUST reuse existing variants first and only add a minimal additive variant if the documented filter treatment cannot otherwise be represented while preserving API compatibility.

#### Scenario: Future priority or availability UI reuses the foundation
- **WHEN** a later generic component needs priority badges or availability chips
- **THEN** it can consume the documented semantic tokens and any additive Badge semantics without introducing a separate chip component

#### Scenario: Filter-style controls stay within the primitive boundary
- **WHEN** the documented filter-button treatment is implemented
- **THEN** the change reuses the existing Button primitive surface and does not introduce a new button-like component file

### Requirement: Verification SHALL prove the foundation update does not regress current screens
The implementation MUST verify that `AppShell`, `AppSidebar`, `HomeView`, `TasksView`, `CalendarView`, and `PeopleView` continue to render through semantic tokens after the Design System v2 token and primitive alignment work. The change MUST run `scripts/precommit-run` before completion.

#### Scenario: Repository verification passes
- **WHEN** the implementation is complete
- **THEN** `scripts/precommit-run` passes successfully

#### Scenario: Current screens do not require domain-specific rewrites
- **WHEN** the token and primitive updates are applied
- **THEN** existing screens remain compatible without adding dashboard-specific styling work just to survive the design-system change
