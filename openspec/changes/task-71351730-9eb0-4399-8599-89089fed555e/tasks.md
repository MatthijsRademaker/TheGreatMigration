## 1. Fix sidebar nav-item clickability

- [ ] 1.1 Open `frontend/src/shared/layout/app-sidebar/AppSidebar.vue` and replace all `SidebarMenuButton` nav-item templates: remove `tooltip` prop, keep `as-child` prop, replace `<RouterLink>` with a native `<a>` element that has `:href` and custom Vue directives `v-click-nav` and `v-title` for navigation and tooltip respectively (directives survive Primitive VNode cloning where `@click` and `:title` bindings are dropped from merged props)
- [ ] 1.2 Add a `vClickNav` custom Vue directive (with `Directive<HTMLElement, string>`) that registers a `click` event listener via `addEventListener` in its `mounted` hook, calling `router.push(to)` for standard left-clicks (button === 0 and no modifier keys) with a `.catch()` fallback to `window.location.assign(to)`. Store the handler in a `WeakMap<HTMLElement, Handler>` and remove it in the `unmounted` hook for proper cleanup.
- [ ] 1.3 Remove `RouterLink` import from the script setup of `AppSidebar.vue` since it is no longer used
- [ ] 1.4 Remove `import { Tooltip, TooltipContent, TooltipTrigger } from 'reka-ui'` or any Reka UI Tooltip import from `SidebarMenuButton.vue`'s nav-item usage — verify `SidebarMenuButton.vue` itself does not import tooltip modules; if unused Tooltip wrapper code remains in `SidebarMenuButton.vue`, it is out of scope for this change (only `AppSidebar.vue` template changes)
- [ ] 1.5 Verify the `isActive` helper still works: the active nav item correctly receives accent styling via `data-active` in both expanded and collapsed states

## 2. Remove duplicate branding from SidebarFooter

- [ ] 2.1 In `AppSidebar.vue`, remove the first `<SidebarMenu>` block inside `<SidebarFooter>` — the one containing the `size="lg"` `SidebarMenuButton` with `NotebookTabsIcon`, "The Great Migration", and "House move planner"
- [ ] 2.2 Verify the `SidebarHeader`'s branding remains unchanged (it still shows "The Great Migration / House move planner" with `NotebookTabsIcon`)
- [ ] 2.3 Remove the `NotebookTabsIcon` import if it is no longer used after removing the footer block (check if it's still needed in the header)

## 3. Update sidebar specs — explicitly stop using `tooltip` prop

- [ ] 3.1 In `AppSidebar.vue`, remove the `tooltip` prop from all `SidebarMenuButton` nav-item instances (the `title` attribute on the `<a>` element replaces it)
- [ ] 3.2 Verify no `<SidebarMenuButton>` for a navigation item passes a `tooltip` prop (the `as-child` prop is kept for correct DOM structure)

## 4. Add clickability unit tests

- [ ] 4.1 Open `frontend/tests/sidebar-mobile.test.ts` and add a new test case inside the existing `describe` block: assert that clicking a sidebar nav `<a>` element triggers `router.push` and navigates to the correct route
- [ ] 4.2 Add a test case asserting that each of the six nav `<a>` elements has a non-empty `href` attribute matching its expected route
- [ ] 4.3 Add a test case asserting that each nav `<a>` element has a `title` attribute set to the item's label text (for collapsed-mode tooltip)
- [ ] 4.4 Add a test case asserting that the `SidebarFooter` does NOT contain the text "The Great Migration" or "House move planner"
- [ ] 4.5 Run `npm run test:unit` (or `npx vitest run`) in the `frontend/` directory and confirm all tests pass

## 5. Verify pre-commit checks

- [ ] 5.1 Run `scripts/precommit-run` from the project root and confirm lint, type-check, and test checks all pass
