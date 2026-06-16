## Why

The People Availability matrix renders beautiful read-only status pills, but editing a person's daily availability requires scrolling past the matrix to a separate "Manage people" card that duplicates the same data in a list view. This two-presentation pattern adds cognitive overhead and maintenance cost — users should be able to click any pill in the matrix to change its status.

## What Changes

- `PeopleAvailability.vue` gains an `editable` prop that makes status cells clickable
- Clicking a cell in editable mode opens an inline popover with the four canonical status options (available, busy, partial, off) plus a clear action
- The popover dispatches the same upsert/delete mutations already wired in `PeopleView.vue`
- The separate "Manage people" card is removed — create-person moves above the matrix, per-person delete moves into the person row
- `HomeView.vue` stays read-only (passes no `editable` prop, default is `false`)
- No backend changes — the existing CRUD endpoints (`PUT /api/people/{id}/availability/{date}`, `DELETE /api/people/{id}/availability/{date}`, `DELETE /api/people/{id}`, `POST /api/people`) already cover all operations

## Capabilities

### New Capabilities

*(none — this change modifies existing capabilities, it doesn't introduce new ones)*

### Modified Capabilities

- `people-availability-component`: The component contract gains an `editable` boolean prop and emits a cell-change event. Cells render a clickable badge that opens an inline status picker when `editable` is true. The creation form and per-person delete actions are hosted by the parent view, not the component itself.
- `people-availability-integration`: `PeopleView.vue` wires the component's cell events to the existing upsert/delete availability mutations. The "Manage people" card is replaced by a create-person form above the matrix and a delete button per row in the matrix. `HomeView.vue` remains read-only.

## Impact

| Area | Impact |
|------|--------|
| `frontend/src/people/PeopleAvailability.vue` | Add `editable` prop, click handler on cells, inline popover component, emit cell-change event |
| `frontend/src/people/PeopleView.vue` | Remove "Manage people" card, move create form above matrix, wire cell-change events to mutations, move delete into matrix row |
| `frontend/src/people/types.ts` | Add types for the cell-change event payload and the editable interaction |
| `frontend/src/home/HomeView.vue` | No change (already passes `v-bind` data, `editable` defaults to false) |
| Backend | No change |
| Tests | Update PeopleView tests for new layout; add PeopleAvailability tests for editable mode; remove assertions about the old management card |
