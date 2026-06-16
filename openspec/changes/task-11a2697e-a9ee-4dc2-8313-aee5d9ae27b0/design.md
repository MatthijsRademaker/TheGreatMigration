## Context

The repository already has a presentational `TaskManagementPanel.vue`, a generated frontend query for `GET /api/tasks/backlog`, and a Postgres-backed backlog read model. The remaining gap is end-to-end wiring: the panel still renders local fixtures, the `/tasks` route has no mutation path, the home dashboard uses the same component without a read-only distinction, and the backend Store/API surface does not yet support task writes.

## Goals

- Replace fixture-backed task rendering with backend-derived backlog data on `/tasks` and the home dashboard.
- Add the minimal backlog task CRUD and assignment-write surface needed to make `/tasks` functional.
- Keep `GET /api/tasks/backlog` as the canonical task read model after every write.
- Preserve the home dashboard as a read-only overview while `/tasks` owns management interactions.
- Refresh committed OpenAPI/client artifacts and add focused verification across backend, frontend, generated client, and route rendering.

## Non-Goals

- Authentication, roles, collaboration, notifications, or audit history.
- Schedule-board CRUD, drag-and-drop scheduling, or backlog-to-schedule moves.
- Normalizing the task `room` field into room entities or changing room CRUD behavior.
- Move-notes CRUD, room-progress metrics, or a broader redesign of shared primitives.
- Server-side filtering or unrelated project-management scope beyond the existing backlog task surface.

## Decisions

### 1. Preserve the backlog endpoint as the canonical read model

`GET /api/tasks/backlog` remains the single source for task rows, summary counts, priority legend, and status legend. The write endpoints mutate persistence, and the frontend refreshes by invalidating or refetching that same backlog query.

### 2. Add a minimal task CRUD API following existing backend patterns

The backend adds `POST /api/tasks`, `PUT /api/tasks/{id}`, and `DELETE /api/tasks/{id}` using the same Huma registration, Store delegation, and sqlc-backed persistence pattern already used for people CRUD. Create and update bodies carry `title`, `priority`, `peopleNeeded`, `room`, `status`, and `assignedTo`.

### 3. Replace assignments transactionally in the task body

Assignment changes are handled through the task create/update body instead of a dedicated assignment endpoint. Store writes replace the full `assignedTo` set transactionally, and task deletion removes the task row and its assignment rows in the same transaction.

### 4. Generate task IDs and append sort order server-side

New backlog tasks receive stable server-assigned `task-*` identifiers and the next append-only `sort_order` value from the database side. The accepted trade-off is the simple max-plus-one approach used for this single-organizer app, with concurrency risk documented as a future concern rather than expanded scope for this change.

### 5. Keep the panel presentational and separate home from `/tasks` by variant

Introduce a shared `useTaskBacklog` composable modeled on `usePeopleAvailability`. `TaskManagementPanel` stays the rendering surface and receives a documented `readOnly` variant so `HomeView.vue` can hide filter/add/edit/delete controls while `TasksView.vue` owns the full management flows.

### 6. Use focused Sheet-based task forms and typed client mutations on `/tasks`

`/tasks` owns create, edit, delete, and assignment updates through the generated frontend client and existing shared `Sheet`, `Input`, `Select`, and `Button` primitives. Successful mutations invalidate the backlog query so rows, legends, and summary-derived information stay aligned with backend-confirmed state.

### 7. Treat fixture leakage and loose task typing as part of the integration work

The task feature must stop depending on `frontend/src/tasks/fixtures.ts` for successful rendering. Frontend task rows use the canonical priority/status vocabulary from the generated client or matching union types, and the displayed `Unassigned` pill remains a derived UI label from empty `assignedTo` arrays rather than a new backend status.

## Risks

- Server-side `task-*` ID and `sort_order` generation via max-plus-one has a concurrency race, although accepted as sufficient for the current lightweight app.
- Assignment replacement must be transactional to avoid partial writes or orphaned rows.
- Existing SSR route tests currently mock an empty backlog but assert fixture content; they will fail unless updated to realistic BFF-shaped payloads and explicit state assertions.
- OpenAPI snapshot and generated client artifacts can drift if regeneration is missed after backend task-write changes.

## Conflict Resolution

- **Assignment API shape:** resolved to `assignedTo` in the task create/update body with transactional replacement, not a separate assignment endpoint.
- **Home versus `/tasks` behavior:** resolved to a shared panel with a `readOnly` home variant and full management controls on `/tasks` only.
- **Delete semantics:** resolved to hard-delete backlog tasks and their assignment rows in one transaction because current schema references only backlog assignments.
- **Create/edit UX:** resolved to a focused Sheet-based form using existing shared primitives rather than an unrelated redesign.

## Traceability

- Task: `11a2697e-a9ee-4dc2-8313-aee5d9ae27b0`
- Dossier: `2026-06-16T11:06:50.746Z`
- Decisions: `1-swarm-architect-recommendation`, `1-swarm-lead-dev-recommendation`, `1-swarm-reviewer-recommendation`
- Round outputs: `round:1:agent:swarm-architect`, `round:1:agent:swarm-lead-dev`, `round:1:agent:swarm-reviewer`
- Artifact base: `initial`