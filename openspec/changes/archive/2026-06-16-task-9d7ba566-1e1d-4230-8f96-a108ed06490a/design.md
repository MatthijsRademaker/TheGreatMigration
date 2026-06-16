## Context

The repository already has three key pieces in place: a presentational `PeopleAvailability` Vue component, a read-only `GET /api/dashboard/people-availability` endpoint, and generated frontend query artifacts for that endpoint. The gap is between those layers. `HomeView.vue` and `PeopleView.vue` still rely on the component’s deterministic defaults, while the backend has no write endpoints for managing people or per-date availability.

## Goals

- Render `PeopleAvailability` on the homepage from the existing dashboard availability endpoint.
- Reuse the same backend-derived availability data on `/people` with explicit loading, error, and empty states.
- Add the smallest write surface needed to manage people and daily availability without expanding into broader scheduling or collaboration scope.
- Keep the homepage read-only and place availability management controls on `/people`.
- Regenerate committed OpenAPI/client artifacts and add focused tests that protect the existing contracts.

## Non-Goals

- Authentication, invites, roles, permissions, or collaboration features.
- Task backlog CRUD, schedule editing, room progress CRUD, or move-notes editing.
- A redesign beyond the existing People Availability dashboard/component intent.
- Bulk availability updates or date-range editing in this slice.
- Replacing Pinia Colada, the generated Hey API client, Huma, sqlc, or the current Postgres-backed architecture.

## Decisions

### 1. Use a shared adapter/composable for the read path

Introduce `usePeopleAvailability` as the single bridge between the generated dashboard response and the component’s narrow local props contract. It owns the `DashboardBody` to component mapping, null-safe handling for generated arrays, ISO-date-to-display-label conversion, canonical status validation, and route-friendly loading/error/empty state exposure.

### 2. Keep the component presentational and the homepage read-only

`PeopleAvailability.vue` remains the rendering surface, not the fetch owner. `HomeView.vue` and `PeopleView.vue` consume the shared composable, but only `/people` exposes management actions. The homepage continues to act as a read-oriented dashboard panel in the existing grid layout.

### 3. Add a minimal people-and-availability write surface

The backend adds `POST`, `PUT`, and `DELETE` operations for people plus single-date `PUT`/`DELETE` operations for availability. Person IDs are client-supplied stable slugs. Availability writes are limited to planning-window dates and the canonical statuses `available`, `busy`, `partial`, and `off`.

### 4. Resolve delete conflicts with explicit `409` behavior

Deleting a person must not silently break assignment data. If the person is referenced by backlog or schedule assignments, the delete path returns `409 Conflict`. Missing people return `404`, and malformed statuses or dates return `400`.

### 5. Regenerate typed client artifacts and use SDK-backed mutations on `/people`

The read path continues to use the generated Pinia Colada query for `GET /api/dashboard/people-availability`. After the backend contract changes, the committed OpenAPI snapshot and generated client artifacts are refreshed so `/people` can call the new write endpoints type-safely. Query generation remains the established pattern; mutations can be wired from the generated SDK where the Colada output remains query-oriented.

### 6. Sequence the work in two implementation stages inside one change

To reduce coordination risk, implementation should wire the shared read path first and then add the CRUD slice, while still keeping CRUD in scope for this change because the task explicitly asks to fix forward any missing endpoints needed for a fully functional component.

## Risks

- The change spans backend handlers, sqlc generation, Store methods, OpenAPI regeneration, frontend adapter code, and mutation state; missing any one layer can break the end-to-end flow.
- Existing frontend tests assert demo names, day labels, and summary copy, so they must be updated carefully when backend data replaces demo defaults on the routes.
- Delete handling must translate assignment-reference failures into `409 Conflict` rather than leaking a generic backend error.
- The generated Pinia Colada artifacts currently prove the read-query pattern, but write flows still require deliberate mutation wiring on `/people`.

## Conflict Resolution

- **CRUD scope:** Kept in scope for this change because the task asks for a fully functional component; risk is managed by sequencing read wiring before CRUD work.
- **Route ownership:** The homepage stays read-only while `/people` owns create, update, and delete interactions.
- **Delete semantics:** Resolved to `409 Conflict` when backlog or schedule assignment references block deletion.
- **Component contract mismatch:** Resolved with a shared adapter/composable that translates generated API data into the component’s local display contract.

## Traceability

- Task: `9d7ba566-1e1d-4230-8f96-a108ed06490a`
- Dossier: `2026-06-16T07:14:16.547Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial`
