## Why

The People Availability surface already exists as a polished component and the backend already exposes `GET /api/dashboard/people-availability`, but the homepage and `/people` route still render demo defaults instead of backend data. There is also no write surface for adding people or maintaining per-day availability, so the matrix cannot become a fully functional feature without fixing the missing CRUD slice.

## What Changes

- Create a shared `usePeopleAvailability` read path that calls the generated dashboard availability query, adapts `DashboardBody` into the component’s local props contract, converts ISO dates to readable labels, and exposes loading, error, and empty states.
- Wire `PeopleAvailability` on the homepage and `/people` to that shared BFF-backed data, preserving the current dashboard composition and keeping the homepage read-only.
- Add the minimal backend CRUD surface needed to manage the matrix: people create/update/delete plus single-date availability upsert/delete, backed by Huma, sqlc, and the existing Postgres tables.
- Resolve delete conflicts with `409 Conflict` when a person is still referenced by backlog or schedule assignments, and constrain availability writes to canonical statuses and planning-window dates.
- Regenerate the committed OpenAPI snapshot and frontend client artifacts after backend contract changes, then add focused frontend and backend tests covering adapter behavior, route states, CRUD validation, and contract preservation.

## Impact

- The homepage and `/people` stop depending on demo defaults and render the same backend-derived availability data.
- `/people` becomes the management surface for adding people, changing daily statuses, and removing unreferenced people, while the homepage remains a read-only overview.
- The backend gains a documented write surface without changing the existing dashboard availability response contract.
- Verification expands to protect status vocabulary, summary consistency, OpenAPI/client regeneration, and the existing app-shell/sidebar behavior.
