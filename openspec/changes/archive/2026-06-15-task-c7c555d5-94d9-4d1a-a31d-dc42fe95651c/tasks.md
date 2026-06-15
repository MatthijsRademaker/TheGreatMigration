## 1. Route and navigation

- [x] 1.1 Register `/showcase` route in `src/app/router.ts` with lazy-loaded `@/showcase/ShowcaseView.vue`, meta title `Component showcase`, and meta description `Developer preview of reusable UI components and patterns.`
- [x] 1.2 Add a secondary "Developer" `SidebarGroup` with label "Developer" to `AppSidebar.vue`, below the existing `SidebarSeparator`, containing a single `SidebarMenuItem` linking to `/showcase` with `ComponentIcon` from `lucide-vue`

## 2. Showcase view and catalog

- [x] 2.1 Create `src/showcase/ShowcaseView.vue` following the feature-folder convention with `<script setup lang="ts">` and `<template>` structure matching existing route views
- [x] 2.2 Define a curated catalog data structure grouping components by category (Actions, Inputs, Feedback, Layout, Overlays) with per-component metadata (name, description, variants to render)
- [x] 2.3 Render the Actions group: Button with all 6 variants (default, outline, secondary, ghost, destructive, link) × 3 key sizes (default, sm, lg), each labeled with its variant+size combination
- [x] 2.4 Render the Inputs group: Input with a single default state example and usage note
- [x] 2.5 Render the Feedback group: Badge with all 6 variants (default, secondary, destructive, outline, ghost, link), Skeleton with a single pulse row and loading-placeholder explanation
- [x] 2.6 Render the Layout group: Card subcomponents (Card, CardHeader, CardTitle, CardDescription, CardContent) as a composed example card, Separator with both horizontal and vertical orientations
- [x] 2.7 Render the Overlays group: Sheet with a trigger button + SheetContent demo (title + description + close), Tooltip with a hover-trigger demo using TooltipProvider wrapper
- [x] 2.8 Add a prose-only entry for the Sidebar component family documenting it as layout-coupled and referencing the live AppSidebar as the canonical usage example

## 3. Verification

- [x] 3.1 Navigate to `/showcase` and confirm the AppShell header shows "Component showcase" title and the showcase description
- [x] 3.2 Confirm the "Developer" sidebar section appears below the separator with the showcase entry and active styling when on the route
- [x] 3.3 Verify all rendered components match their actual shared/ui appearance (no style regressions from component imports)
- [x] 3.4 Verify existing routes (`/`, `/tasks`, `/calendar`, `/people`) continue to render unchanged
- [x] 3.5 Run `scripts/precommit-run` to verify the change passes project checks
