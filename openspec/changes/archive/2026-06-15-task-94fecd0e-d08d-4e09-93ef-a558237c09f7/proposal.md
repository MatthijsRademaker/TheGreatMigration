## Why

Tech-onboarding scaffolding introduced artifacts that do not belong in the product: a developer-only `/showcase` route with a component catalog page, a "Developer" sidebar navigation group, a "Move focus" prose block in the sidebar, placeholder footer copy, and a set of generated `shared/ui` subcomponents that are entirely unused. These artifacts have no counterpart in `designs/design-system.png` (which shows the intended navigation as Dashboard, Tasks, Schedule, People, Rooms / Areas, and Settings) or `designs/components.png` (which documents the product UI component inventory, not a developer catalog). The repo needs a cleanup to remove these onboarding-only surfaces, prune dead shared/ui code, and align the visible navigation shell with the current design boards.

## What Changes

### Removals

- **`/showcase` route and `src/showcase/` directory** — The route registered in `src/app/routes.ts` and the `ShowcaseView.vue` component catalog page. This was created during tech onboarding as a developer discovery surface; the design boards show no developer catalog. The archived change at `openspec/changes/archive/2026-06-15-task-c7c555d5-94d9-4d1a-a31d-dc42fe95651c/` preserves the original implementation.
- **AppSidebar "Developer" navigation group** — The sidebar section containing a `ComponentIcon` link to `/showcase`. No design board includes a developer section.
- **AppSidebar "Move focus" content block** — A static prose panel ("Keep the plan practical…") rendered between the Plan group and Developer group. Not present in any design board.
- **AppSidebar footer copy** — Replaced from `"Local planning app"` to product-aligned copy.
- **Unused sidebar subcomponents** — Seven `.vue` files and their barrel exports: `SidebarGroupAction`, `SidebarInput`, `SidebarMenuAction`, `SidebarMenuSkeleton`, `SidebarMenuSub`, `SidebarMenuSubButton`, `SidebarMenuSubItem`. Grep across `src/**/*.{ts,vue}` confirms zero consumers outside their barrel export in `src/shared/ui/sidebar/index.ts`.
- **`src/shared/ui/skeleton/` directory** — `Skeleton.vue` and its barrel export. After removing `ShowcaseView.vue` and `SidebarMenuSkeleton.vue`, Skeleton has zero product consumers. Removal order: `SidebarMenuSkeleton.vue` first (to break its import of Skeleton), then the skeleton directory.
- **Local `TooltipProvider.vue` wrapper** — `src/shared/ui/tooltip/TooltipProvider.vue` and its barrel export. After removing ShowcaseView (its only consumer), this wrapper is orphaned. `SidebarProvider.vue` already imports `reka-ui`'s `TooltipProvider` directly with `:delay-duration="0"`.

### Updates

- **Sidebar navigation label** — `"Calendar"` → `"Schedule"` in `AppSidebar.vue` to match `designs/design-system.png` navigation. The route path (`/calendar`) and name (`calendar`) remain unchanged — this is a display label change only.
- **Sidebar footer text** — `"Local planning app"` → product-aligned text reflecting "The Great Migration" brand.
- **`docs/design-system-v2.md`** — The out-of-scope clause excluding "compositions from `designs/components.png`" is narrowed to clarify that `designs/design-system.png` governs navigation and layout decisions for this cleanup; the composition exclusion applies only to full dashboard compositions, not to navigation/layout source-of-truth decisions.

### Retained (explicitly NOT removed)

- **Sheet and Tooltip primitives** — `Sidebar.vue` depends on Sheet for mobile sidebar behavior; `SidebarMenuButton.vue` depends on Tooltip/TooltipContent/TooltipTrigger for collapsed-state labels. Neither can be removed without breaking the sidebar.
- **Placeholder product views** — `HomeView.vue`, `TasksView.vue`, `CalendarView.vue`, and `PeopleView.vue` are the current product UI shell. Their content is minimal but they are structurally aligned with the expected route layout from `designs/design-system.png`. Full dashboard reconstruction is deferred to follow-up work per acceptance criteria non-goals.
- **All actively-used shared/ui primitives** — Badge, Button, Card, Input, Separator remain untouched.

### OpenSpec

- **Supersede `openspec/specs/component-showcase/spec.md`** — The canonical spec currently requires the `/showcase` route, Developer sidebar section, and component catalog. This cleanup removes all three. The spec is marked as superseded with a new retirement spec documenting the removal rationale and the archival path.

## Impact

**Affected code**: `src/app/routes.ts` (remove `/showcase` route), `src/showcase/ShowcaseView.vue` (delete), `src/shared/layout/app-sidebar/AppSidebar.vue` (remove Developer group, Move focus block, update footer and nav label), `src/shared/ui/sidebar/index.ts` (remove 7 unused exports), `src/shared/ui/sidebar/SidebarGroupAction.vue` through `SidebarMenuSubItem.vue` (7 file deletions), `src/shared/ui/skeleton/Skeleton.vue` and `src/shared/ui/skeleton/index.ts` (delete directory), `src/shared/ui/tooltip/TooltipProvider.vue` and its export in `src/shared/ui/tooltip/index.ts` (delete).

**Affected docs/specs**: `docs/design-system-v2.md` (narrow out-of-scope clause), `openspec/specs/component-showcase/spec.md` (supersede with retirement spec).

**Breaking changes**: None for production users. The `/showcase` route is developer-only. No external consumers depend on the removed sidebar subcomponents or Skeleton primitive.

**Risk**: Future developers may look for the old component catalog. Mitigation: the retirement spec documents the archival path and rationale.
