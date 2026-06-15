## Context

The project is a Vue 3 + Vite SPA with a shadcn-vue design system. The `src/shared/ui/` directory contains nine reusable component families (badge, button, card, input, separator, sheet, skeleton, sidebar, tooltip) with CVA-based variant systems and standard `index.ts` barrel exports. The app uses vue-router with lazy-loaded feature-folder views (`src/<feature>/<Feature>View.vue`), route meta for AppShell header content, and an AppSidebar with hard-coded navigation arrays.

There is no existing component discovery page. Developers must inspect source files or browse existing views to understand the design-system inventory. This change adds a manually curated showcase route inspired by Storybook's catalog concepts without adopting its tooling.

## Goals / Non-Goals

**Goals:**
- Add a `/showcase` route that renders live examples of all currently available `shared/ui` components in grouped sections
- Use route meta to provide distinct showcase header titling in the existing AppShell
- Add a secondary "Developer" sidebar section with a low-visual-weight showcase entry, visually separated from the primary "Plan" navigation
- Render concrete variants for each showcased component (all button variants × sizes, all badge variants, Card subcomponent compositions, Separator orientations, Skeleton example, Sheet trigger+content demo, Tooltip hover demo)
- Document the Sidebar component family as layout-coupled (not rendered standalone)
- Follow existing feature-folder, route-registration, and shared/ui import conventions

**Non-Goals:**
- No Storybook installation or integration
- No automated component discovery from the filesystem
- No code snippets, prop playgrounds, or copy-import-path affordances in v1
- No environment-based gating or auth guards
- No sheet- or tooltip-interaction matrices beyond a single trigger demo
- No modification of existing planning routes or views

## Decisions

### 1. Manual curated catalog (not auto-discovery)

**Choice**: A static catalog object in `ShowcaseView.vue` that explicitly lists each component group and its examples.

**Alternatives considered**:
- Filesystem glob + export introspection: Requires build-tooling conventions and processing that the project does not currently have. Rejected — violates non-goal #2 of the exploration dossier.
- Barrel-export analysis: Would tightly couple the showcase to the exact export shape of index.ts files and would not handle variant enumeration.

**Rationale**: Manual curation is simple, predictable, and keeps the maintenance cost low. Adding a new component requires one catalog entry — no build changes.

### 2. Sidebar: secondary "Developer" section (not direct-URL-only, not primary nav)

**Choice**: Add a `SidebarSeparator` + `SidebarGroup` with label "Developer" containing a single showcase entry using `ComponentIcon`.

**Alternatives considered**:
- Direct-URL-only: Would hide the showcase from developers, defeating its purpose as a discovery surface.
- Primary "Plan" nav entry: Would clutter the core move-planning workflow. Rejected per AC #7.

**Rationale**: The existing AppSidebar already establishes the secondary-section pattern with "Move focus" (lines 64-66). Adding a "Developer" section below it is a natural extension that provides discoverability without intruding on the primary flow.

### 3. Sheet: inline trigger demo (not omitted)

**Choice**: Include Sheet in the Overlays group with a live `<SheetTrigger>` + `<SheetContent>` demo inside a container.

**Alternatives considered**:
- Omit entirely: Would leave a gap in the catalog — Sheet is a valid reusable component.
- Static description only: Would not satisfy AC #2 which requires "rendered examples, not just text links."

**Rationale**: The showcase is a full interactive Vue SPA, not a static document. Sheet's trigger+content pattern works correctly inline without conflicting with app-level providers.

### 4. Sidebar family: prose documentation, no standalone render

**Choice**: Add a catalog note explaining the Sidebar family is layout-coupled (requires `SidebarProvider` context) and referencing the live `AppSidebar.vue` usage as the canonical example. No standalone render attempt.

**Rationale**: Rendering a second `<Sidebar>` inside the showcase would conflict with the app-level `SidebarProvider`. All three refinement agents agree a standalone sidebar demo is impractical.

### 5. No code snippets in v1

**Choice**: Each showcased component shows live rendered examples with state labels and a one-sentence usage note. No import-path copies or code blocks.

**Rationale**: Manual code snippets drift as implementations change. Rendering the actual components with their live styling is more valuable for visual review and always accurate.

## Risks / Trade-offs

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Catalog drift: new shared/ui components added without showcase entry | Medium | Low — catalog becomes incomplete | Document the convention in the change: adding a component to shared/ui requires adding one catalog entry in ShowcaseView. Consider a PR checklist note. |
| Skeleton has limited visual value in isolation | Low | Low — shows as a single pulsing bar | Render one skeleton example with a label explaining its loading-placeholder purpose. |
| Tooltip best-effort on touch devices | Low | Low — the showcase is developer-facing | The underlying reka-ui TooltipRoot handles this; document the limitation in the usage note. |
| Sheet demo may look odd in a confined card | Low | Low — functional but visually constrained | Render the Sheet demo in its own section with adequate space. The trigger button opens the sheet from the edge as expected. |

## Open Questions

- Should future iterations add code-snippet copying or a prop playground? (Out of scope for this change — deferred to a future task.)
- Should the showcase be gated by environment in production builds? (Out of scope — the codebase has no env-based route pattern today. This is a separate concern requiring its own change.)

## Traceability

- Task: `c7c555d5-94d9-4d1a-a31d-dc42fe95651c`
- Dossier: `2026-06-15T14:15:21.206Z`
- Decisions: `1-swarm-architect-recommendation` (sidebar + groupings), `1-swarm-lead-dev-recommendation` (no code snippets, inline demos), `1-swarm-reviewer-recommendation` (sidebar resolution, per-component strategy)
- Round 1 outputs: swarm-architect, swarm-lead-dev, swarm-reviewer