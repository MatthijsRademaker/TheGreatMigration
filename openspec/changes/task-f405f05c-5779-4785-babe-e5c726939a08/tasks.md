## 1. Mobile Trigger in AppShell Header

- [ ] 1.1 Add `SidebarTrigger` import and `useRoute` import to `AppShell.vue`
- [ ] 1.2 Insert a `SidebarTrigger` element in the `<header>` before the date-range/title `<div>`, with class `md:hidden`
- [ ] 1.3 Verify the header layout compresses naturally on narrow viewports without overflow

## 2. Hide In-Sidebar Trigger on Mobile

- [ ] 2.1 Import `useSidebar` from `@/shared/ui/sidebar` in `AppSidebar.vue` and destructure `isMobile`
- [ ] 2.2 Wrap the existing `SidebarTrigger` `SidebarMenuItem` in the sidebar header with `v-if="!isMobile"`
- [ ] 2.3 Verify the in-sidebar trigger still renders on desktop viewports

## 3. Auto-Close Mobile Sheet on Navigation

- [ ] 3.1 In `AppSidebar.vue`, add a `watch` on `() => route.path` that calls `setOpenMobile(false)` when `isMobile` is true and `openMobile` is true
- [ ] 3.2 Guard against initial mount: the watcher must not fire `setOpenMobile(false)` on initial page load when `openMobile` is already false

## 4. Mobile Sidebar Component Tests

- [ ] 4.1 Create `frontend/tests/sidebar-mobile.test.ts` with `// @vitest-environment jsdom` directive
- [ ] 4.2 Add a top-level mock that stubs `window.matchMedia` to return `matches: true` (mobile viewport)
- [ ] 4.3 Write a test that mounts a test wrapper containing `SidebarProvider` + `AppSidebar` and asserts the mobile `SidebarTrigger` is visible
- [ ] 4.4 Write a test that clicks the mobile-sidebar trigger and asserts the Sheet opens (`openMobile` becomes true)
- [ ] 4.5 Write a test that, with the Sheet open, triggers a route navigation and asserts `openMobile` becomes false
- [ ] 4.6 Verify existing SSR tests in `app-routes-render.test.ts` still pass without modification

## 5. OpenSpec Spec Extension

- [ ] 5.1 Add mobile-availability requirements to `openspec/specs/sidebar-navigation/spec.md` covering trigger visibility on sub-768px viewports and Sheet open/close behavior
- [ ] 5.2 Add auto-close-on-navigation requirements covering the route watcher behavior
- [ ] 5.3 Add a requirement preserving desktop trigger behavior (in-sidebar trigger stays visible on non-mobile viewports)

## 6. Verification

- [ ] 6.1 Run `scripts/precommit-run` and ensure all checks pass
- [ ] 6.2 Confirm `vue-tsc --noEmit` passes with zero type errors
- [ ] 6.3 Confirm `vitest run` passes all existing and new tests
- [ ] 6.4 Confirm existing SSR route-render tests remain unchanged and continue to pass
- [ ] 6.5 Visually verify the mobile trigger appears at ≤768px and the desktop sidebar is unchanged at >768px