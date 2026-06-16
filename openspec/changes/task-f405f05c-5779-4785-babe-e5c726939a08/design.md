## Context

The repository provides a shared application shell (`AppShell.vue`) that wraps every routed view with a sidebar and header. The sidebar is built on shadcn-vue sidebar primitives (`frontend/src/shared/ui/sidebar/`) with semantic CSS tokens (`--sidebar-*`) defined in `frontend/src/app/styles.css`. The `AppSidebar.vue` component composes these primitives into the actual navigation: six items across two groups, a brand header, footer content, and a `SidebarRail`.

The primitives already support mobile rendering: `Sidebar.vue` switches to a `Sheet` component when `isMobile` is true (breakpoint: `max-width: 768px`), `SidebarProvider.vue` owns `openMobile` state and provides `setOpenMobile()`/`toggleSidebar()`, and `SidebarTrigger.vue` already dispatches `toggleSidebar()` which correctly routes to `setOpenMobile()` on mobile. However, no mobile trigger exists in the app shell, and no route-change hook closes the mobile drawer.

The existing canonical spec (`openspec/specs/sidebar-navigation/spec.md`) covers desktop navigation structure, active state, badge policy, footer content, and variant preservation — but has zero mobile requirements.

## Goals / Non-Goals

### Goals
- Provide an always-visible mobile sidebar trigger in the shared `AppShell.vue` header on viewports ≤ 768px.
- Reuse the existing `SidebarTrigger` component and `toggleSidebar()` dispatch — no new trigger components.
- Wire the mobile sidebar to automatically close when the user navigates to any route.
- Hide the in-sidebar `SidebarTrigger` on mobile to avoid a redundant close control inside the open Sheet.
- Preserve every existing desktop sidebar behavior: inset variant, icon collapsible, rail, grouped nav items, footer content, cookie-persisted open/closed state.
- Add jsdom component tests for the mobile open/close flow using the existing `@vue/test-utils` mount pattern.
- Extend `openspec/specs/sidebar-navigation/spec.md` with mobile-availability and auto-close requirements.

### Non-Goals
- Redesigning the six-item navigation structure or changing route inventory.
- Implementing new Rooms, Settings, or footer utility features.
- Reworking page-level mobile layouts outside the shared navigation entry point.
- Changing sidebar color tokens or broader design-system styling.
- Adding a floating action button or alternative navigation paradigm.
- Implementing a full mobile responsive redesign of each routed view.
- Wiring footer utility actions to backend functionality.

## Conflict Resolution

The architect (Round 1) proposed a route watcher in `SidebarProvider.vue`; the lead-dev (Round 1) proposed the same mechanism in `AppSidebar.vue`. Both are accepted decisions; the reviewer flagged this as needing resolution.

**Resolution**: Place the route watcher in `AppSidebar.vue`.

**Rationale**: `SidebarProvider.vue` is a generic UI primitive that manages sidebar state (open/closed, mobile/desktop). It currently has zero router dependencies. Introducing `useRoute` and a `watch` on `route.path` into the provider would couple a generic UI primitive to vue-router's API surface — architecturally impure. `AppSidebar.vue` is the composition layer that already imports `useRoute`, already sits inside `SidebarProvider` context (so `useSidebar()` is available), and is the natural orchestration point between navigation events and sidebar state. The watcher reads `route.path` and calls `setOpenMobile(false)` — a clean delegation from the composition layer to the state provider without breaking encapsulation.

## Decisions

### D1: Mobile trigger in AppShell.vue header, left-aligned, md:hidden
**Decision**: Add a `SidebarTrigger` component in `AppShell.vue`'s `<header>` element, placed before the date-range/title flex container, with class `md:hidden` so it is only visible below the 768px breakpoint.
**Rationale**: `AppShell.vue` is the single shared layout wrapper rendered by every route — placing the trigger there guarantees universal availability without touching any routed view. The `SidebarTrigger` already calls `toggleSidebar()` which dispatches to `setOpenMobile()` when `isMobile` is true. Left-alignment before the date range matches standard mobile navigation patterns (hamburger on the left). The `md:hidden` class uses the same 768px breakpoint as `useMediaQuery` in `SidebarProvider.vue`.
**Evidence**: Dossier initial acceptance criteria #1; architect Round 1 suggested requirement #1; lead-dev Round 1 point (1); reviewer Round 1 blocker-2 (resolved nonBlocking).

### D2: Route watcher in AppSidebar.vue for auto-close
**Decision**: Add a `watch(() => route.path, ...)` in `AppSidebar.vue` that calls `setOpenMobile(false)` when the route changes and `isMobile` is true. Guard against the initial mount by only closing when `openMobile` is already true.
**Rationale**: Resolved via conflict resolution (see above). `AppSidebar.vue` is the correct place — it already imports `useRoute` and `useSidebar()`, and sits at the composition boundary between routing and sidebar state. The guard (`if (isMobile.value && openMobile.value)`) prevents the watcher from firing on initial page load when the sheet is already closed.
**Evidence**: Dossier initial acceptance criteria #3; architect Round 1 suggested requirement #3; lead-dev Round 1 point (2); reviewer Round 1 blocker-1 (resolved nonBlocking); Conflict Resolution section above.

### D3: Hide in-sidebar SidebarTrigger on mobile
**Decision**: Wrap the existing `SidebarTrigger` inside `AppSidebar.vue`'s `SidebarHeader` in `v-if="!isMobile"` to hide it on mobile viewports.
**Rationale**: On desktop the in-sidebar trigger provides the collapse/expand control. On mobile, the same trigger rendered inside the open Sheet creates a confusing UX — a "close the thing you're inside" button using a `PanelLeftIcon` that suggests collapse, not close. Hiding it prevents the redundant control while keeping the Sheet's built-in close mechanism (swipe, overlay click, Escape key) as the dismissal path.
**Evidence**: Dossier open questions; architect Round 1 suggested requirement #2; lead-dev Round 1 point (3); reviewer Round 1 blocker-3 (resolved nonBlocking).

### D4: jsdom component tests for mobile sidebar behavior
**Decision**: Create `frontend/tests/sidebar-mobile.test.ts` with `// @vitest-environment jsdom` directive, using `@vue/test-utils` mount. Mock `window.matchMedia` via `Object.defineProperty` to simulate the mobile breakpoint. Mount sidebar components and assert: (a) trigger element exists when `isMobile` is true, (b) clicking trigger opens the Sheet, (c) navigating to a new route closes the Sheet. Existing SSR tests in `app-routes-render.test.ts` are preserved unchanged as desktop regression guards.
**Rationale**: The existing SSR tests use `renderToString` and cannot exercise interactive behavior, responsive breakpoints, or Sheet open/close state. The project already has a jsdom component test pattern (`settings-view.test.ts`) demonstrating `@vue/test-utils` mount, composable mocking, and DOM assertions. The `matchMedia` mock addresses the jsdom limitation (jsdom does not implement `window.matchMedia`).
**Evidence**: Dossier initial acceptance criteria #5; architect Round 1 suggested requirement #5; lead-dev Round 1 point (4); reviewer Round 1 risk-1 (test infrastructure gap).

### D5: Desktop sidebar behavior preserved
**Decision**: No changes to `SidebarProvider.vue`, `Sidebar.vue`, `SidebarTrigger.vue`, or any other shared UI primitive. Desktop sidebar configuration (`collapsible="icon"`, `variant="inset"`, `SidebarRail`, grouped nav, footer content) remains untouched. The mobile trigger's `md:hidden` class keeps it invisible on desktop.
**Rationale**: The change is additive — all desktop behavior is preserved by not modifying existing contracts. Existing SSR tests already assert desktop-sidebar chrome and serve as regression guards.
**Evidence**: Dossier non-goals; architect Round 1 suggested requirement #4; existing spec requirement "Sidebar SHALL preserve existing variant and rail configuration".

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Test infrastructure: jsdom lacks `window.matchMedia` | High | Medium | Mock via `Object.defineProperty(window, 'matchMedia', ...)` before mount. Following the established mock pattern from `settings-view.test.ts`. |
| Layout crowding: header may overflow at 320px viewports | Low | Low | `SidebarTrigger` is a compact icon button (size-icon-sm, ~32px). The header uses `gap-3` and the trigger sits before the `flex-1` flex container, compressing naturally. |
| Route watcher fires on initial page load | Low | Medium | Guard `setOpenMobile(false)` with `if (openMobile.value)` — only close when already open. |
| Desktop regression from `v-if="!isMobile"` guard | Low | Low | On desktop, `isMobile` is always false, so `!isMobile` is always true — the in-sidebar trigger renders identically to current behavior. |
| Spec drift: change touches AppShell, AppSidebar, and new test file | Low | Low | Zero changes to shared UI primitives. State change lives in AppSidebar, structural change in AppShell, test coverage in new file. |
| SSR hydration concern: mobile trigger renders in SSR output | Low | Low | The `md:hidden` class only affects visual display, not DOM presence. Existing SSR tests do not assert against trigger absence, so no existing assertions break. |

## Traceability

- **Task**: `f405f05c-5779-4785-babe-e5c726939a08`
- **Dossier**: `2026-06-16T15:16:48.594Z`
- **Decisions**: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- **Rounds**: 1 (architect, lead-dev, reviewer)
- **Existing spec**: `openspec/specs/sidebar-navigation/spec.md` — preserved and extended
- **Existing sidebar primitives**: `frontend/src/shared/ui/sidebar/` — zero modifications
- **Existing test infrastructure**: `settings-view.test.ts` — jsdom component test pattern
- **Design system**: `docs/design-system-v2.md` — sidebar/top-bar visual-state contract