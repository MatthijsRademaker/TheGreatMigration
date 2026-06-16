## Context

The sidebar nav items use this rendering chain:

```
SidebarMenuButton (tooltip="People", as-child)
  └─ Tooltip
       └─ TooltipTrigger (as-child)            ← Layer 1
            └─ SidebarMenuButtonChild (as-child)  ← Layer 2 (via Reka UI Primitive)
                 └─ RouterLink (:to="/people")     ← Layer 3 (Vue component)
                      └─ <a href="/people">           ← Layer 4 (actual DOM)
```

Each `as-child` layer must correctly forward `pointer`/`click` event listeners to its child. With Vue components (`RouterLink`) instead of native elements as the child of `Primitive`, attribute forwarding becomes fragile because `RouterLink` is not a plain DOM element — it's a Vue component that applies its own prop/attr handling. A double `as-child` chain (`TooltipTrigger` → `Primitive`) exacerbates this, increasing the surface for lost or duplicated event handlers.

The result: some `<a>` tags in the sidebar fail to receive or dispatch click events reliably, making the nav entry feel "not clickable" or requiring very precise targeting.

Additionally, `SidebarFooter` contains a redundant brand block that duplicates the `SidebarHeader`'s "The Great Migration / House move planner" content. This was introduced as a copy-paste artifact during the initial sidebar implementation.

## Goals / Non-Goals

**Goals:**

- Make all six sidebar nav entries reliably clickable on the first try in both expanded and collapsed modes.
- Remove the duplicate brand card from `SidebarFooter` (reclaiming ~48px of vertical space).
- Preserve the collapsed-mode tooltip (icon labels) without relying on the `Tooltip > TooltipTrigger` event-forwarding chain.
- Preserve native `<a>` semantics (right-click → open in new tab, middle-click, Ctrl+click).
- Add a jsdom-based unit test asserting each nav link is clickable.

**Non-Goals:**

- Changes to shadcn-vue sidebar primitives (`Sidebar.vue`, `SidebarProvider.vue`, `SidebarTrigger.vue`, `SidebarMenuButtonChild.vue`).
- Changes to `SidebarMenuButton.vue` — the fix is in `AppSidebar.vue` only.
- Adding badge counts, real-time data subscriptions, or interactive footer actions.
- Mobile Sheet behavior changes.
- Backend or API changes.

## Decisions

### Decision 1: Use native `title` attribute instead of Reka UI Tooltip for nav item labels

Replace the `SidebarMenuButton` `tooltip` prop with the native HTML `title` attribute on the `RouterLink`:

```
Before:                       After:
───────────────────────────── ─────────────────────────────
SidebarMenuButton             SidebarMenuButton
  tooltip="People"              (no tooltip prop)
  as-child                      (no as-child)
  └─ RouterLink                └─ RouterLink title="People"
```

The `title` attribute provides identical UX in collapsed mode: hovering the icon shows a browser-native tooltip with the item name. In expanded mode, the label text is already visible so no tooltip is needed.

**Rationale:**

| Aspect | Reka UI Tooltip | Native `title` |
|---|---|---|
| Event chain | Adds `TooltipTrigger(as-child)` wrapping | None (zero layers) |
| Click reliability | Fragile due to double `as-child` forwarding | Guaranteed (no forwarding) |
| Tooltip styling | Custom theming | Browser-default |
| Delay/transition | Configurable | Browser-controlled (typically 1-2s) |
| Screen reader | `aria-describedby` | Native accessibility |
| Bundle size | Reka UI Tooltip + TooltipTrigger modules | Zero |

The browser-native `title` tooltip has a slightly different visual appearance and no animation, but this is appropriate for a utility sidebar in a productivity app. The collapsed sidebar is an icon-only mode where fast visual identification via label is the goal, not a polished animation.

### Decision 2: Remove `as-child` from all nav-item `SidebarMenuButton` instances

Without `as-child`, `SidebarMenuButtonChild`'s `Primitive` renders as its own `<button>` element (the default `as: "button"`). The `RouterLink` is then a child element *inside* the button, producing:

```html
<button data-slot="sidebar-menu-button" ...>
  <a href="/people" title="People">
    <svg>...</svg>
    <span>People</span>
  </a>
</button>
```

However, **this creates a nested interactive element** (`<button>` containing `<a>`), which is invalid HTML. Browsers resolve this by making one or both unclickable.

To avoid this, instead of nesting, **remove `RouterLink` entirely and use a plain `<a>` tag with the `href` and navigation handled by the `Primitive` component itself.** Since `SidebarMenuButton` already accepts `as-child`, we can use `as="a"` on the button to render it as an `<a>` tag directly:

```vue
<SidebarMenuButton :is-active="isActive(item.to)" as="a">
```

Then use Vue's `v-bind` or the `@click` handler with `router.push()` to navigate. But this loses `RouterLink`'s built-in navigation handling.

**Better approach: Use `v-bind` on `SidebarMenuButton` to pass `href` as an attribute that gets forwarded to `<a>`:**

```vue
<SidebarMenuButton
  :is-active="isActive(item.to)"
  as="a"
  :href="item.to"
  @click.prevent="router.push(item.to)"
>
  <component :is="item.icon" />
  <span>{{ item.title }}</span>
</SidebarMenuButton>
```

`Primitive` with `as="a"` renders an `<a>` tag. Passing `href` via v-bind attaches it to the element. The `@click.prevent` handler calls `router.push()`. The `RouterLink` is no longer needed.

**But this loses right-click → open in new tab** because `router.push()` only navigates the current tab. To support that, use a native `<a>` with real `href`:

```vue
<a
  data-slot="sidebar-menu-button"
  data-sidebar="menu-button"
  :class="sidebarMenuButtonClasses"
  :href="item.to"
  @click="handleNav($event, item.to)"
>
  <component :is="item.icon" />
  <span>{{ item.title }}</span>
</a>
```

But this bypasses the sidebar primitives entirely.

**Compromise: Keep `SidebarMenuButton` with `as-child`, but pass `as="a"` and the RouterLink's `to` as the `href`:**

```vue
<SidebarMenuButton as-child>
  <a :href="item.to" :title="item.title" @click="handleNav($event, item.to)">
    <component :is="item.icon" />
    <span>{{ item.title }}</span>
  </a>
</SidebarMenuButton>
```

Here, `as-child` is used once (no double forwarding since there's no Tooltip wrapping). The child `<a>` is a native DOM element, not a Vue component — `Primitive` handles this correctly. Navigation works natively through the `href`, plus the `@click` handler calls `router.push()` for SPA navigation.

The `title` attribute provides the collapsed-mode tooltip natively.

### Decision 3: Remove duplicate brand block from SidebarFooter

The `SidebarFooter` template contains:

```vue
<SidebarMenu>
  <SidebarMenuItem>
    <SidebarMenuButton size="lg">
      <div class="flex aspect-square size-9 ...">
        <NotebookTabsIcon />
      </div>
      <div class="grid flex-1 ...">
        <span class="truncate text-sm font-semibold">The Great Migration</span>
        <span class="truncate text-xs text-muted-foreground">House move planner</span>
      </div>
    </SidebarMenuButton>
  </SidebarMenuItem>
</SidebarMenu>
```

This is an exact copy of the `SidebarHeader` branding. It serves no purpose and consumes ~48px of vertical space. It will be removed, with the `PlusIcon` and `CircleHelpIcon` utility buttons remaining below.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|---|---|---|
| Native `title` tooltip has different visual appearance than Reka UI Tooltip | Slight inconsistency in collapsed mode | Acceptable for a productivity app sidebar; the text label appears after a 1-2s delay, same as before |
| Removing `as-child` changes the DOM structure of nav items | Potential CSS breakage in collapsed mode | `data-slot="sidebar-menu-button"` and `data-sidebar="menu-button"` are still on the `<a>` tag; the class string is unchanged. Verify all Tailwind selectors still match with tests. |
| Without `RouterLink`, SPA navigation uses a custom click handler | Right-click → open in new tab must still work via native `href` | The `<a>` tag has a real `href` — right-click works natively. The `@click.prevent` handler only intercepts standard clicks. |
| The `handleNav` function is duplicated across 6 nav items | Minor code repetition | Extract into a shared `navigate` function in the script setup. |
