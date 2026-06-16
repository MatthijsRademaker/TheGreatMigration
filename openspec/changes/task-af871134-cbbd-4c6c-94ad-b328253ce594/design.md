## Context

The People Availability matrix is rendered by `PeopleAvailability.vue` — a pure presentational component that accepts props and renders a table of status pills. It is used in two contexts:

- **HomeView.vue** (read-only dashboard) — stays read-only
- **PeopleView.vue** (management route) — already has create-person form, delete mutation, upsert/delete availability mutations wired via `@pinia/colada`

The management mutations exist in PeopleView but are accessed through a separate "Manage people" card below the matrix, which duplicates the same data in a list view. The matrix pills themselves are not clickable. The backend already exposes all needed CRUD endpoints (`PUT /api/people/{id}/availability/{date}`, `DELETE /api/people/{id}/availability/{date}`, `POST /api/people`, `DELETE /api/people/{id}`).

A Popover component already exists in the shared UI library (`frontend/src/shared/ui/popover/`) based on Reka UI. All four canonical status badge variants (`available`, `busy`, `partial`, `off`) also exist.

## Goals / Non-Goals

**Goals:**

- Make matrix pills clickable in editable mode, opening an inline status picker
- Wire cell interactions to the existing upsert/delete mutations in PeopleView
- Remove the duplicate "Manage people" card, redistributing its concerns
- Keep HomeView strictly read-only (no change)
- No backend changes — the existing CRUD contract is sufficient
- Existing tests pass; new tests cover editable mode and the updated PeopleView layout

**Non-Goals:**

- Creating new UI primitives (Popover already exists)
- Changing the backend API, Store interface, or database schema
- Changing the `usePeopleAvailability` composable or adapter
- Adding drag-and-drop, bulk editing, or date navigation to the matrix
- Changing the HomeView layout or the daily schedule component
- Opening the status picker via keyboard shortcuts (follow-up concern)

## Decisions

### 1. Cell interaction: inline Popover anchored to the badge

Clicking a status pill in editable mode opens a Reka UI Popover positioned below the clicked badge. The popover displays the four canonical status options as clickable badges plus a "Clear" action. Selecting a status closes the popover and emits an `update-cell` event to the parent.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Dropdown select in cell | Limited to single-selection, no room for "Clear" action, harder to style consistently with badges |
| Full-width management card stays | Duplicates data, adds cognitive overhead, doesn't fulfill the "editable matrix" ask |
| Modal/dialog on cell click | Over-engineered for a single status change; heavier UX |
| Side panel | Same as modal — too heavy for a one-click status change |

**Why Popover wins:** It is already in the shared UI library, positions relative to the clicked element automatically (Reka UI portal), supports any content (badges + clear button), and dismisses on click-outside by default.

### 2. Component event contract

`PeopleAvailability.vue` gains a single emitted event:

```ts
defineEmits<{
  'update-cell': [payload: { personId: string; dayIndex: number; status: AvailabilityStatus | null }]
  'delete-person': [personId: string]
}>()
```

- `dayIndex` is the index within the `days` array, not an ISO date string. This avoids coupling the component to ISO date logic. The parent maps `dayIndex` to an ISO date using `rawData.range.startDate` via the existing `getISODate()` helper.
- `status` is `AvailabilityStatus` (one of the four canonical values) when the user picks a new status, or `null` when the user picks "Clear" (maps to the `deleteAvailabilityMutation` on the parent).
- `delete-person` is emitted when the user clicks a delete trigger in the person's row (only shown in editable mode).

### 3. Component prop additions

```ts
export interface PeopleAvailabilityProps {
  // ... existing props unchanged
  editable?: boolean
}
```

When `editable` is `true`:

- Badge `<td>` cells wrap the Badge in `<Popover>` with `<PopoverTrigger as-child>`
- The popover content shows the four status badges as clickable options plus a "Clear (reset to off)" action
- Each person row `<tr>` gains an additional `<td>` with a delete button (visible on row hover or always visible)
- The legend is unchanged

When `editable` is `false` or undefined (default):

- Current read-only behavior is preserved exactly
- No Popovers, no delete buttons, no click handlers

### 4. Popover content: inline status picker

```html
<PopoverContent class="w-48 p-2">
  <div class="flex flex-col gap-1">
    <p class="text-xs font-medium text-muted-foreground mb-1">Set status</p>
    <button
      v-for="s in statusOptions"
      :key="s"
      class="flex items-center gap-2 rounded-md px-2 py-1 text-sm
             hover:bg-muted cursor-pointer w-full text-left"
      @click="handleSelect(s)"
    >
      <Badge :variant="s" class="pointer-events-none">
        {{ s.charAt(0).toUpperCase() + s.slice(1) }}
      </Badge>
    </button>
    <Separator class="my-1" />
    <button
      class="text-xs text-muted-foreground hover:text-foreground
             px-2 py-1 w-full text-left cursor-pointer rounded-md hover:bg-muted"
      @click="handleClear"
    >
      Clear (reset to off)
    </button>
  </div>
</PopoverContent>
```

### 5. PeopleView redistribution

The "Manage people" Card is removed entirely. Its pieces are redistributed:

| Concern | New home |
|---|---|
| Create person form | Stays as its own Card above the matrix (same as now, just the management Card goes away) |
| Status editing per cell | Moves into the matrix via Popover (controlled by `editable` prop on PeopleAvailability) |
| Delete person | Moves into the matrix person row (delete button in a new `<td>` column, shown only when `editable` is true) |
| Error display | Error messages stay in PeopleView, rendered above the matrix or inline near the affected section |

The event wiring in PeopleView maps component events to mutations:

```ts
function handleUpdateCell(payload: { personId: string; dayIndex: number; status: AvailabilityStatus | null }) {
  const date = getISODate(payload.dayIndex)
  if (!date) return
  if (payload.status) {
    // upsert
    handleStatusUpdateFor(payload.personId, date, payload.status)
  } else {
    // clear
    handleClearAvailabilityFor(payload.personId, date)
  }
}
```

### 6. HomeView stays unchanged

HomeView passes `v-bind="availabilityData"` to `<PeopleAvailability>`. Since `editable` defaults to `false`, the component renders in read-only mode. No change needed.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|---|---|---|
| Popover positioning in a scrollable `<div class="overflow-x-auto">` | Popover might clip or position incorrectly in a scrollable container | Reka UI Popover uses a portal, rendering outside the scroll container. Test with horizontal scroll on narrow viewports. |
| Many Popover instances in the DOM | Performance concern for large teams (20+ people × 30+ days) | Popover is lazy — only renders content when open. Only one can be open at a time. Acceptable for the planning window scale (4-14 days). |
| Delete button in matrix row adds visual density | Rows become busier, especially on narrow screens | Use a compact icon button (trash icon) with `size="icon" variant="ghost"`. Only show on hover in the row. |
| SSR mismatch: Popover is client-only | SSR render will not include Popover content | Badge still renders as a `<span>` in SSR. Popover only hydrates on the client. Tests for editable mode should use a client-render test, not SSR. |
| Error state management: parent owns errors | Popover closes after click, user may not see the error | Display inline error text above the matrix in PeopleView when a mutation fails. The popover reopening by re-clicking the cell will show the error. Could also add a brief toast notification in a follow-up. |
