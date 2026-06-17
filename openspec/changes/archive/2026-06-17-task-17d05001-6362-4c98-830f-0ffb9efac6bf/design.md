## Context

The current KPI cards component (`frontend/src/home/components/KpiCards.vue`) renders four summary cards in a grid. The component was built against the original `openspec/specs/kpi-summary-cards/spec.md` spec that was itself an approximation of the design. Now that the component exists with working data wiring and tests, the refinement room has identified a concrete set of deltas between the current implementation and the design target in `designs/components.png` / `designs/home-page.png`. This design document crystallizes the architectural decisions needed to close those deltas.

The key structural difference: the current implementation uses a **full-height accent column** (`w-[72px]` with centered icon) that occupies the entire left edge of each card, while the design shows a **thin left accent border** (`border-l-4`) with a **compact icon chip** positioned in the upper-left of the content area. This is a template-level restructure, not a style tweak.

## Goals

- Restructure all four KPI cards from full-height accent columns to thin-left-accent-border + compact-icon-chip anatomy.
- Adopt the CardHeader/CardContent/CardTitle/CardDescription primitives mandated by the canonical OpenSpec rather than bypassing them with `!py-0 !gap-0` overrides.
- Reorder cards to High priority → People → Unassigned → Rooms as shown in the design.
- Switch the people card value from `X of Y available` to `X / Y` fraction format with a date-context subtitle from `range.selectedDate`.
- Replace card descriptions with design-like outcome copy optimized for scanability.
- Commit to specific Lucide icon names for all four cards.
- Preserve Rooms completed as an isolated placeholder with no business logic.
- Update unit tests and SSR tests to assert the new layout, order, copy, and state behavior.
- Update the canonical OpenSpec to reflect all changes in lockstep with the implementation.

## Non-Goals

- Redesigning the rest of `HomeView.vue` or adjacent dashboard panels.
- Inventing a backend room-progress API or business rule.
- Adding new shared UI primitive directories or new global design tokens.
- Changing task-backlog or people-availability backend contracts beyond consuming fields that already exist.
- Adding skeleton/shimmer loading states (the current text fallback pattern is preserved).
- Modifying the `/people` route SSR assertions.

## Decisions

### 1. Card order: High priority first, per design

**What**: Cards render in this order: High priority tasks → People available today → Unassigned jobs → Rooms completed.
**Why**: The design artifacts (`designs/components.png`, `designs/home-page.png`) show High priority as the leftmost card. All three refinement reviewers (architect, lead-dev, reviewer) agree on this order. The canonical spec currently lists People first and must be updated.
**Evidence**: Dossier affectedAreas area 2, architect requirement 1, lead-dev requirement 1, reviewer recommendation.

### 2. Layout anatomy: Thin left accent border + compact icon chip

**What**: Replace the `w-[72px]` full-height accent column with a `border-l-4` (4px) semantic left border on the Card itself (using `border-destructive`, `border-info`, `border-warning`, `border-success`). Place a compact `size-8 rounded-lg` icon chip with semantic soft background (`bg-destructive-soft text-destructive`, etc.) in the upper-left content area (CardHeader). The leaf decoration (`/images/leaf.png`, absolute positioned bottom-right) is preserved.
**Why**: The design shows a thin accent line on the left edge of each card with a separate compact icon tile, not a full-height colored column. The current implementation's `w-[72px]` column with centered icon is a fundamentally different visual structure. All three reviewers identified this as the primary structural gap.
**Evidence**: Dossier affectedAreas area 2, architect requirement 2, lead-dev requirement 2, reviewer requirement 2.

### 3. Card primitives: Adopt CardHeader/CardContent structure

**What**: Use CardHeader to contain the icon chip and caption label (`[font-size:var(--text-caption)] text-muted-foreground`). Use CardContent to contain the value (`text-3xl font-semibold`) and subtitle (`text-sm text-muted-foreground`). Remove the `!py-0 !gap-0` overrides. The thin left accent border is applied to the Card itself via `border-l-4` with the semantic color class.
**Why**: The canonical OpenSpec explicitly requires CardHeader/CardContent/CardTitle/CardDescription usage. The current implementation bypasses them. Adopting the primitives satisfies the spec contract while the thin-border + icon-chip approach makes the primitives work with the design anatomy.
**Evidence**: Architect requirement 3, architect blocker 1, lead-dev blocker 1, reviewer risk.

### 4. People value format: Fraction with date subtitle

**What**: The people card displays `X / Y` (e.g., "6 / 8") as the primary value line. The subtitle displays `available on <MMM D>` sourced from `range.selectedDate` (e.g., "available on Jul 5"). If `selectedDate` is undefined, fall back to `available today`. The `displayAvailableToday` computed property continues to clamp `rawAvailableToday` to `totalPeople` defensively.
**Why**: The design PNGs show a fraction/slash treatment with date context. The current "X of Y available" format is verbose and generic. All three reviewers agreed on the fraction format. The `range.selectedDate` field already exists in the availability API response and is available in the SSR mock, so no backend change is needed.
**Evidence**: Architect requirement 4, lead-dev requirement 3, lead-dev blocker 2, reviewer requirement 1, reviewer requirement 5.

### 5. Icon selection: FlagIcon, UsersRoundIcon, BriefcaseIcon, CheckCircleIcon

**What**: High priority uses `FlagIcon` (replaces `TriangleAlertIcon`). People uses `UsersRoundIcon` (unchanged). Unassigned uses `BriefcaseIcon` (replaces `HammerIcon`). Rooms uses `CheckCircleIcon` (replaces `Building2Icon`). All imported from `@lucide/vue`.
**Why**: The current icons were chosen as placeholders during the original implementation. The design intent (inferred from typical dashboard iconography patterns and the dossier's open questions) suggests flag for priority, briefcase for unassigned jobs, and check-circle for completion. All three reviewers flagged icon ambiguity as a gap. This is a binding decision — the implementation must use these exact icons.
**Evidence**: Architect requirement 6, reviewer requirement 4, reviewer blocker 4.

### 6. Subtitle copy: Outcome-focused, scanable text

**What**:
- High priority: `high priority tasks need attention`
- People: `available on <MMM D>` (dynamic, with `available today` fallback)
- Unassigned: `jobs that need assignment`
- Rooms: `rooms fully packed and cleared` (static placeholder)
**Why**: Current descriptions are generically explanatory (e.g., "Tasks marked as high priority", "Helpers with confirmed availability"). The design target and the dashboard product intent (`http://project-docs/vision-and-goals.md`) optimize for scanability — users should absorb the outcome at a glance, not read definitions. All three reviewers called for concrete copy targets.
**Evidence**: Architect requirement 5, lead-dev requirement 4, reviewer requirement 3, reviewer blocker 3.

### 7. Rooms completed: Isolated placeholder, no business logic

**What**: The Rooms completed card remains a hardcoded placeholder: value `—`, subtitle `rooms fully packed and cleared`, `data-testid="kpi-placeholder-rooms-completed"`, no backend query, no derived counts from other data sources. Code comments document it as a placeholder for a future room-progress contract. The card uses the same layout anatomy as the other three (thin accent, icon chip, leaf decoration).
**Why**: The canonical OpenSpec explicitly requires this isolation, and all three reviewers agreed. The dossier assumptions state there is no current backend contract for room progress. Inventing business logic would be scope creep. A follow-up change can introduce a real room-progress data source.
**Evidence**: Architect requirement 7, lead-dev requirement 8, reviewer requirement 6, reviewer blocker 5.

### 8. Styling: Existing semantic tokens only

**What**: All styling uses semantic token classes from `frontend/src/app/styles.css`. Accent borders: `border-destructive` (high priority), `border-info` (people), `border-warning` (unassigned), `border-success` (rooms). Icon chips: `bg-destructive-soft text-destructive`, `bg-info-soft text-info`, `bg-warning-soft text-warning`, `bg-success-soft text-success`. Typography: caption role for labels, `text-3xl font-semibold` for values, `text-sm text-muted-foreground` for subtitles. No raw hex, no new tokens, no new shared directories.
**Why**: The Design System v2 contract (`docs/design-system-v2.md`) and the canonical spec require exclusive use of existing semantic tokens. The dossier acceptance criteria explicitly forbid raw hex hacks and new shared UI directories. The current implementation already uses the correct semantic classes for icon backgrounds.
**Evidence**: Architect requirement 10, design system v2 docs, dossier acceptance criteria item 5.

## Risks

- **SSR test breakage (high)**: Changing the people KPI format from "6 of 8 available" to "6 / 8" will break the home-route assertions `expect(html).toContain("of 8")` and `expect(html).toContain("available")` in `app-routes-render.test.ts`. Mitigation: update assertions alongside the component change. The `/people` route assertion must NOT be modified.
- **Unit test breakage (high)**: All assertions locking in `w-[72px]`, `X of Y available`, and the old card order will fail. Mitigation: rewrite tests alongside implementation, asserting new layout markers and copy.
- **Card primitive friction (medium)**: The Card.vue primitive uses `gap-panel`, `py-panel`, `flex-col` defaults that may need `!` overrides to accommodate the thin-accent-border + compact-icon-chip anatomy. Mitigation: accept limited overrides as a special composition pattern, document with code comments, and do not add new Card props.
- **`range.selectedDate` availability (low)**: If the availability API response omits `range.selectedDate` in some environments, the people card subtitle will fall back to `available today`. The defensive fallback is built into the design.
- **Icon choice without vision model (low)**: The chosen icons (FlagIcon, BriefcaseIcon, CheckCircleIcon) are best approximations inferred from the design context. If design review later requires different icons, it's a trivial swap of Lucide imports with no structural impact.

## Traceability

- **Task**: `17d05001-6362-4c98-830f-0ffb9efac6bf`
- **Dossier**: `2026-06-17T19:38:02.819Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Round outputs**: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- **Artifact base**: `initial` snapshot for `task-17d05001-6362-4c98-830f-0ffb9efac6bf`

## Conflict Resolution

- **Card primitive adoption vs. flat flex layout**: The architect flagged that the canonical spec mandates CardHeader/CardContent primitives while the current code bypasses them. **Resolved**: Adopt CardHeader/CardContent structurally. The thin accent border goes on the Card root (`border-l-4`), the icon chip goes in CardHeader, and the value/subtitle stack goes in CardContent. This satisfies the spec requirement and achieves the design anatomy without requiring a spec amendment.
- **Card order: People-first (canonical spec) vs. High-priority-first (design)**: The canonical spec at `openspec/specs/kpi-summary-cards/spec.md` lists People available today as the first card. The design shows High priority first. All three reviewers agree the design is the authoritative source of truth. **Resolved**: Follow the design. The canonical spec is updated in this change.
- **People value format: "X of Y available" vs. "X / Y" fraction**: The current code and canonical spec use the "X of Y" format. The design shows fraction/slash. All three reviewers recommend the fraction format. **Resolved**: Use "X / Y" fraction with a date-context subtitle. This provides more information (value, total, date) in less space.
- **Rooms completed: Placeholder-only vs. follow-up change**: The canonical spec mandates placeholder-only, and the dossier open questions ask whether a separate room-progress contract should be proposed. **Resolved**: Keep as placeholder in this task. A follow-up change can introduce a real room-progress data source when the backend contract is defined. This task does not propose that follow-up.
