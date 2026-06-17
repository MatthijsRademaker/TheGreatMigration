## Context

The app uses a sidebar layout from shadcn-vue with a `SidebarProvider` (`flex min-h-svh`) wrapping a fixed `AppSidebar` on desktop (Sheet overlay on mobile) and a `SidebarInset` (`flex-1 flex-col`) for the main content area. Inside `SidebarInset`, a `<header>` contains the page title, route description, planning-window date range, and the mobile sidebar trigger button.

Current scroll behavior: `SidebarInset` has `overflow-hidden`, and no view component provides its own scroll container. On mobile viewports, when content exceeds the remaining viewport height, it is clipped rather than scrollable. This means the app effectively has broken scrolling on mobile.

## Goals / Non-Goals

**Goals:**

- Enable vertical scrolling for content that exceeds the viewport height on mobile
- Make the header sticky at the top so navigation context (title, date range, sidebar toggle) remains visible during scroll
- Keep the change minimal — a single file modified (`AppShell.vue`), no new components, no backend

**Non-Goals:**

- No new composables, state management, or API changes
- No view-specific layout changes — every view inherits scrolling from the shell
- Desktop behavior is not a primary concern, but the solution applies universally for consistency

## Decisions

### Decision: Sticky via `position: sticky` rather than `position: fixed`

**Rationale:** `position: fixed` positions relative to the viewport, which conflicts with the `SidebarInset` layout — on desktop, the inset has margins (`md:peer-data-[variant=inset]:m-2`) that a fixed element would ignore. `position: sticky` is relative to the nearest scroll container (the `SidebarInset` itself) and stays within the inset's boundaries, respecting margins and padding on all breakpoints.

### Decision: Enable scroll on `SidebarInset` rather than on each view

**Rationale:** Adding `overflow-y: auto` to `SidebarInset` makes it a single scroll container for all content, including the header. The header sticks within this container. The alternative — making each view its own scroll container — is fragile, repetitive, and creates nested-scroll problems. One container, one source of truth.

### Decision: Apply sticky on all breakpoints, not just mobile

**Rationale:** The header already has `backdrop-blur` and `bg-background/90` — these are visual cues designed for sticky behavior. Applying the same behavior on desktop is a zero-cost consistency improvement. There is no semantic difference in behavior across breakpoints.

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| **Horizontal scroll inside views (Calendar `overflow-x-auto`) may conflict with parent `overflow-y`** | `overflow-x` and `overflow-y` are independent. Inner `overflow-x: auto` on DailySchedule's flex row works correctly inside a parent with `overflow-y: auto`. Verified by the CSS spec — `overflow-x` and `overflow-y` shorthand maps to separate longhands. |
| **Z-index conflict with mobile sidebar Sheet** | The Sheet (from Reka/Vaul) creates its own stacking context via a portal. The header only needs `z-10` to stay above scrolling content within SidebarInset — it won't compete with the overlay layer. |
| **Desktop inset margins may cause visual gap above sticky header** | The sticky header is inside SidebarInset. On desktop, SidebarInset has `md:peer-data-[variant=inset]:m-2` margin. The sticky header sits at `top: 0` within the inset, which is below the margin — this is correct and maintains the inset visual style. |
