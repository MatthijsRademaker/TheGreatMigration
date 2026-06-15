## 1. Remove Developer Showcase Surface

- [ ] 1.1 Delete `/showcase` route entry from `src/app/routes.ts` (the route object including lazy import and meta)
- [ ] 1.2 Delete `src/showcase/ShowcaseView.vue`
- [ ] 1.3 Remove `src/showcase/` directory (confirm it is empty or contains only deleted file)

## 2. Clean Up AppSidebar

- [ ] 2.1 Remove "Move focus" `SidebarGroup` block (lines with `SidebarGroupLabel`, prose `div`, and its surrounding `SidebarSeparator`)
- [ ] 2.2 Remove "Developer" `SidebarGroup` block (lines with `SidebarGroupLabel`, `SidebarMenu`, `SidebarMenuItem`, `SidebarMenuButton`, and its preceding `SidebarSeparator`)
- [ ] 2.3 Update navigation label from `"Calendar"` to `"Schedule"` in `primaryNavigation` array
- [ ] 2.4 Update `SidebarFooter` text from `"Local planning app"` to `"The Great Migration"`
- [ ] 2.5 Remove unused `ComponentIcon` import from lucide-vue imports

## 3. Remove Unused Sidebar Subcomponents

- [ ] 3.1 Delete `src/shared/ui/sidebar/SidebarGroupAction.vue` and remove its export from `index.ts`
- [ ] 3.2 Delete `src/shared/ui/sidebar/SidebarInput.vue` and remove its export from `index.ts`
- [ ] 3.3 Delete `src/shared/ui/sidebar/SidebarMenuAction.vue` and remove its export from `index.ts`
- [ ] 3.4 Delete `src/shared/ui/sidebar/SidebarMenuSkeleton.vue` and remove its export from `index.ts` (MUST be deleted before Skeleton — it imports Skeleton)
- [ ] 3.5 Delete `src/shared/ui/sidebar/SidebarMenuSub.vue` and remove its export from `index.ts`
- [ ] 3.6 Delete `src/shared/ui/sidebar/SidebarMenuSubButton.vue` and remove its export from `index.ts`
- [ ] 3.7 Delete `src/shared/ui/sidebar/SidebarMenuSubItem.vue` and remove its export from `index.ts`

## 4. Remove Skeleton Primitive

- [ ] 4.1 Delete `src/shared/ui/skeleton/Skeleton.vue`
- [ ] 4.2 Delete `src/shared/ui/skeleton/index.ts`
- [ ] 4.3 Remove `src/shared/ui/skeleton/` directory

## 5. Remove Local TooltipProvider Wrapper

- [ ] 5.1 Delete `src/shared/ui/tooltip/TooltipProvider.vue`
- [ ] 5.2 Remove `TooltipProvider` export line from `src/shared/ui/tooltip/index.ts`

## 6. Update Design Contract

- [ ] 6.1 Amend `docs/design-system-v2.md` out-of-scope clause: clarify that `designs/design-system.png` governs navigation and layout decisions; the composition exclusion from `designs/components.png` applies only to full dashboard compositions

## 7. Supersede Component Showcase OpenSpec Spec

- [ ] 7.1 Mark `openspec/specs/component-showcase/spec.md` as superseded by this change
- [ ] 7.2 Commit the retirement spec at `openspec/changes/task-94fecd0e-d08d-4e09-93ef-a558237c09f7/specs/component-showcase-retirement/spec.md` documenting removal rationale

## 8. Verification

- [ ] 8.1 Run `scripts/precommit-run` and confirm all checks pass
- [ ] 8.2 Manually verify sidebar renders correctly with: Dashboard, Tasks, Schedule, People navigation items only
- [ ] 8.3 Manually verify `/showcase` returns 404 (route removed)
- [ ] 8.4 Manually verify `/tasks`, `/calendar`, `/people` still render their placeholder views
- [ ] 8.5 Confirm `src/shared/ui/skeleton/` directory no longer exists
- [ ] 8.6 Confirm `src/shared/ui/tooltip/TooltipProvider.vue` no longer exists
- [ ] 8.7 Confirm all 7 unused sidebar `.vue` files are deleted and their exports removed from `index.ts`