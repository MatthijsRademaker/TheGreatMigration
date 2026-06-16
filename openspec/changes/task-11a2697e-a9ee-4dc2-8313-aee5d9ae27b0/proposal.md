## Why

The product already has a Task Management panel and a canonical backlog read endpoint, but `/tasks` and the home dashboard still depend on static fixtures while the backend exposes only `GET /api/tasks/backlog`. The organizer cannot create, edit, delete, or reassign backlog work, so the task surface is not yet functional.

## What Changes

- Add the minimal backend task write surface required for a functional backlog: `POST /api/tasks`, `PUT /api/tasks/{id}`, and `DELETE /api/tasks/{id}` following the established Huma + Store + sqlc pattern.
- Preserve `GET /api/tasks/backlog` as the canonical read model for task rows, summary counts, priority legend, and status legend; successful writes are reflected by invalidating and refetching that same query.
- Replace fixture-backed task rendering with a shared `useTaskBacklog` frontend read path that consumes the generated client, adapts backlog data into the existing panel contract, and shows explicit loading, backend-unavailable, and empty states.
- Keep the home dashboard read-oriented through a `readOnly` Task Management panel variant, while `/tasks` owns create, edit, delete, and assignment management through a focused Sheet-based form.
- Refresh `frontend/openapi-snapshot.json` and committed `frontend/src/client/` artifacts after the backend contract change, then expand backend and frontend tests around CRUD behavior, BFF-backed rendering, and route mocks that use real backlog-shaped payloads.

## Impact

- `/tasks` becomes the functional task-management surface for a single household organizer.
- The home dashboard shows the same backend-derived task data without exposing destructive controls.
- The backend gains a documented task write surface without replacing or weakening the existing backlog read contract.
- Verification protects assignment persistence, summary consistency, generated-client compatibility, and SSR route rendering against BFF data instead of fixtures.