## Context

The repository has a tokenized design foundation (`docs/design-system-v2.md`, `frontend/src/app/styles.css`), an inventory of shadcn-vue sidebar primitives in `frontend/src/shared/ui/sidebar/`, and a working but incomplete sidebar in `AppSidebar.vue`. The sidebar currently defines four navigation items under a single `Plan` group (Dashboard, Tasks, Schedule, People) and a generic footer chip (`GM` / `The Great Migration`). The design artifacts `designs/components.png` and `designs/home-page.png` show a two-group navigation structure plus lower sidebar content (project card, utility actions). An archived cleanup change (`task-94fecd0e`) removed onboarding sidebar clutter but deliberately deferred the design-derived nav expansion to this change.

The semantic sidebar tokens are already in place: `--sidebar` (`var(--secondary)` = pale green `#F3FAF5`), `--sidebar-accent` (`var(--muted)` = soft green `#E6F3EA`), and `--sidebar-accent-foreground` (`var(--primary)` = brand green `#1E6B3E`). These implement exactly the design-system-v2 navigation state contract (pale-green background, soft-green active fill, green icon/text, neutral inactive with hover) without any new CSS.

## Goals / Non-Goals

### Goals
- Reorganize `AppSidebar.vue` to render six nav items across two groups matching the design-derived information architecture.
- Add lightweight placeholder routes for Rooms / Areas and Settings so every nav item is clickable.
- Replace the `GM` footer chip with the design-derived project card and utility actions.
- Remove hardcoded badge counts on Tasks (12) and People (6) since real data wiring is out of scope.
- Update route-render tests to assert the expanded sidebar contract.
- Reuse existing sidebar primitives, semantic tokens, and lucide icons exclusively — no new component directories.

### Non-Goals
- Building full Rooms / Areas or Settings product features beyond minimal placeholder shells.
- Rebuilding the dashboard, task board, people matrix, or top bar from the full-page design comps.
- Creating a new logo asset or illustration system.
- Wiring sidebar badges or utility actions to real backend data or mutations.
- Changing existing route paths (`/calendar` stays `/calendar`).
- Adding a component showcase, Storybook, or developer catalog surface.
- Modifying the shared sidebar primitives in `frontend/src/shared/ui/sidebar/`.

## Decisions

### D1: Two-group navigation with SidebarSeparator
**Decision**: Organize navigation into two `SidebarGroup` sections separated by a `SidebarSeparator`: primary destinations labeled "Plan" (Dashboard, Tasks, Schedule, People) and secondary destinations labeled "Organization" (Rooms / Areas, Settings).
**Rationale**: The design hierarchy shows logical separation between core workflow items and secondary management items. The existing `SidebarGroup`/`SidebarSeparator` primitives support this without new components. The lead-dev's proposed group labels align with the design intent.
**Evidence**: `designs/components.png` section 2 nav structure; `designs/home-page.png` left rail composition; architect R1 suggested requirement #1; lead-dev R1 suggested requirement #1.

### D2: Lightweight placeholder routes for Rooms / Areas and Settings
**Decision**: Add `/rooms` → `RoomsView.vue` and `/settings` → `SettingsView.vue` with lazy-loaded minimal placeholder views that render through the shared `AppShell` shell with route metadata.
**Rationale**: Making every nav item clickable satisfies acceptance criteria #3 and avoids dead-end UX. The lazy-loaded pattern matches existing routes. Placeholder content keeps scope minimal — full feature work is explicitly excluded.
**Evidence**: Dossier acceptance criteria #3, #4; architect R1 suggested requirement #2; lead-dev R1 suggested requirement #3.

### D3: Hardcoded badges removed
**Decision**: Remove the hardcoded `badge: '12'` on Tasks and `badge: '6'` on People. Add a comment documenting that badges should be re-added when real data subscriptions exist.
**Rationale**: The architect's position is the stronger one: fake counts that never update erode user trust more than no counts at all. The non-goal forbids wiring to real data. The lead-dev's mitigation (code comment flagging removal) is incorporated.
**Evidence**: Architect R1 suggested requirement #4; dossier non-goals; reviewer R1 risk about badge integrity.

### D4: Project card replaces footer chip; utility actions as display-only
**Decision**: Replace the `GM` footer chip with a project card showing project name and subtitle from existing branding. Add utility actions (Add note with `PlusIcon`, Help & Support with `CircleHelpIcon`) as non-interactive display-only items in the footer. Add a comment noting interactivity is deferred to follow-up.
**Rationale**: Both architect and lead-dev include utility actions in scope. The reviewer's concern (scope creep into wiring backend) is addressed by making them display-only. The footer chip is a direct replacement target per acceptance criteria #2.
**Evidence**: `designs/home-page.png` lower sidebar; architect R1 suggested requirement #3; lead-dev R1 suggested requirement #2; reviewer R1 suggested requirement #4.

### D5: Keep existing text-forward brand header
**Decision**: Preserve the existing `SidebarHeader` with `NotebookTabsIcon`, "The Great Migration" title, and "House move planner" subtitle. Do not simplify to text-only.
**Rationale**: The design doc explicitly forbids new logo assets. Two of three participants (architect, lead-dev) support keeping the existing brand block. It already implements the text-forward wordmark approach and provides visual brand presence without new assets.
**Evidence**: `docs/design-system-v2.md` brand treatment; architect recommendation; lead-dev recommendation.

### D6: Icon choices for new nav items
**Decision**: Use `Building2Icon` (from lucide-vue) for Rooms / Areas and `SettingsIcon` for Settings.
**Rationale**: `Building2Icon` evokes physical rooms/spaces and is already available in lucide-vue. `SettingsIcon` (cog) is the standard settings metaphor. These icons are semantically appropriate and exist in the installed `@lucide/vue` package.
**Evidence**: `lucide-vue` ^1.18.0 confirmed installed; reviewer R1 question about icon choices.

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| Placeholder routes create expectation for real feature content | Medium | Low | Views clearly labeled as placeholders ("Feature coming soon"). Non-goals explicitly exclude full feature work. |
| Project card and utility actions challenged as scope creep | Low | Low | Both map directly to design artifacts. Utility actions are display-only — no backend wiring. |
| Adding sidebar items causes overflow on small viewports | Low | Medium | Sidebar primitives handle scroll overflow natively. `SidebarContent` supports scroll. Confirm with visual check. |
| Visual spacing around SidebarSeparator differs from design images | Medium | Low | Use existing sidebar token spacing. Design images are not pixel-contract — `design-system-v2.md` and token definitions are authoritative. |
| Route-render test false-pass for new placeholder views | Low | Medium | Placeholder views must render at least one unique content string for each route. Add content assertions alongside nav label assertions. |

## Traceability

- **Task**: `dcb90e61-3e80-48c9-9a07-3813a882d623`
- **Design contract**: `docs/design-system-v2.md` — navigation states, palette, typography, spacing
- **Theme surface**: `frontend/src/app/styles.css` — semantic CSS variables for sidebar tokens
- **Design artifacts**: `designs/components.png` section 2, `designs/home-page.png`
- **Existing sidebar**: `frontend/src/shared/layout/app-sidebar/AppSidebar.vue`
- **Sidebar primitives**: `frontend/src/shared/ui/sidebar/index.ts`
- **Dossier**: `2026-06-16T04:32:25.491Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Rounds**: 1 (architect, lead-dev, reviewer)
- **Preceding cleanup**: `openspec/changes/archive/2026-06-15-task-94fecd0e-.../design.md`
