## Why

The "Add a person" form currently requires the user to manually type an ID (e.g. `p9`). This forces users to invent a unique slug, check for duplicates manually, and learn the `p{N}` convention. Every other entity in the system — rooms (`room-1`, `room-2`) and tasks (`task-100`, `task-101`) — already auto-generates its ID server-side via a PostgreSQL sequence. The person ID should follow the same pattern: the user provides name and initials, the backend assigns the ID.

## What Changes

- **BREAKING**: `POST /api/people` no longer accepts `id` in the request body. The ID is generated server-side.
- New PostgreSQL sequence `people_id_seq` for atomic ID generation, seeded from the max existing `p{N}` value.
- The `CreatePerson` SQL query switches to `'p' || nextval('people_id_seq')`.
- The frontend form drops the ID input field — only Name and Initials remain.
- The generated frontend API client types reflect the new request body (no `id`).

## Capabilities

### New Capabilities

*(none — this is a modification of existing capability)*

### Modified Capabilities

- `people-management-api`: The current requirement "SHALL use client-supplied stable person IDs suitable for name-derived slugs" changes to "SHALL assign server-generated sequential IDs prefixed with `p-`". The API contract for `POST /api/people` changes: `id` is removed from the request body, and the response still returns the full person record including the server-assigned `id`.

## Impact

- **API**: `POST /api/people` body schema changes — `id` field removed, only `name` and `initials` required
- **Backend**: New migration (`people_id_seq`), updated sqlc query, updated API input type and handler
- **Frontend**: Remove ID field from `PeopleView.vue`, update generated client types
- **OpenAPI snapshot**: Must be refreshed after API schema change
- **Seed data**: Unchanged — existing `p1`..`p8` coexist with new sequential IDs starting at 9
- **Tests**: Update backend test that sends `id` in body; update any frontend test that references the ID field
