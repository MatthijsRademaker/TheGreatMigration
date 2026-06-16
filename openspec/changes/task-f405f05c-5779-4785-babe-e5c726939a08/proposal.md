## Why

The shared app shell (`AppShell.vue`) provides navigation through a desktop sidebar but has no mobile navigation entry point. On phone-width viewports, the sidebar is hidden and cannot be opened — making primary navigation undiscoverable. The sidebar also stays open after selecting a nav destination, covering the new content with the open drawer.

The shadcn-vue sidebar primitives already ship with mobile `Sheet` support (`Sidebar.vue` renders a `Sheet` when `isMobile` is true), `SidebarProvider.vue` owns `openMobile` state, and `SidebarTrigger.vue` correctly dispatches `toggleSidebar()` → `setOpenMobile()` on mobile. The gap is orchestration, not infrastructure: no always-visible mobile trigger exists in the shared shell, and nothing closes the mobile drawer after navigation.

## What Changes

Four concrete, additive changes — zero modifications to the shared UI sidebar primitives:

1. **Mobile trigger in AppShell.vue header** — Add a `SidebarTrigger` in the `AppShell.vue` header row with class `md:hidden`, placed left-aligned before the date range container. This makes the sidebar-open control universally visible on viewports ≤ 768px.

2. **Hide in-sidebar SidebarTrigger on mobile** — Conditionally render the existing `SidebarTrigger` inside `AppSidebar.vue`'s header with `v-if="!isMobile"` to avoid presenting a redundant close-toggle inside the already-open mobile Sheet.

3. **Auto-close on navigation** — Add a `watch` on `route.path` in `AppSidebar.vue` that calls `setOpenMobile(false)` when the route changes and `isMobile` is true. The watcher guards against the initial mount so it only fires on actual navigation. `AppSidebar.vue` is the correct location because it already imports `useRoute`, sits inside `SidebarProvider` context, and avoids introducing a vue-router dependency into the generic UI primitive.

4. **jsdom component tests** — Add `frontend/tests/sidebar-mobile.test.ts` using `@vue/test-utils` mount, mocking `window.matchMedia` to simulate the mobile breakpoint, and asserting: trigger visibility on mobile, click-opens-sheet, and navigation-closes-sheet. Existing SSR route-render tests in `app-routes-render.test.ts` remain unchanged as the desktop regression guard.

**Spec extension**: `openspec/specs/sidebar-navigation/spec.md` gains new requirements covering mobile-availability (trigger presence on sub-768px viewports), auto-close behavior (drawer dismisses on route change when mobile), and preservation of existing desktop controls.

## Impact

- **Affected code**: `frontend/src/shared/layout/app-shell/AppShell.vue`, `frontend/src/shared/layout/app-sidebar/AppSidebar.vue`, `frontend/tests/sidebar-mobile.test.ts` (new), `openspec/specs/sidebar-navigation/spec.md`
- **Not affected**: All sidebar primitives in `frontend/src/shared/ui/sidebar/`, all routed views, all existing routes, all existing tests
- **Dependencies**: No new packages. Reuses existing `SidebarTrigger`, `useSidebar()`, `useRoute`, `@vue/test-utils`, and `vitest`
- **Non-breaking**: Desktop sidebar (inset variant, icon collapsible, rail, grouped nav, footer content) is preserved. The `md:hidden` class on the mobile trigger keeps it visually hidden on desktop; the `v-if="!isMobile"` guard on the in-sidebar trigger is a no-op on desktop
- **Out of scope**: Redesigning navigation structure, adding new routes, changing sidebar styling/tokens, implementing interactive footer actions, full mobile responsive redesign of individual views