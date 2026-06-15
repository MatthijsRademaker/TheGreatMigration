# Proposal: Backend Dashboard People Availability API

## Why

The homepage design (`designs/home-page.png`) shows two people-availability surfaces: a summary card (`People available today 6 / 8`) and a multi-day `People Availability` table (2 Jul–5 Jul). Currently `frontend/src/home/HomeView.vue` hardcodes these values as static text, and the backend exposes only `GET /api/hello`. This task adds a backend BFF read endpoint that serves the dashboard people-availability data in a single fetch, enabling the frontend to replace static placeholders with live data.

## What Changes

- Add a new Huma-registered endpoint `GET /api/dashboard/people-availability` with optional query parameters `start` (ISO date, defaults to server-local today) and `days` (positive integer, defaults to 4, inclusive of `start`).
- The endpoint returns a combined JSON payload with range metadata, summary counts (`availableToday`, `totalPeople`), per-person daily availability arrays, and a status legend.
- Availability statuses are constrained to the four design-backed values: `available`, `busy`, `partial`, `off` (matching `docs/design-system-v2.md`).
- `availableToday` counts only people with status `available` on the selected date (documented in the endpoint description).
- The first slice is backed by in-memory seeded data with at least 8 people exercising all four status states.
- New handler and seed data live in `backend/dashboard.go` under package `main`, registered in `backend/main.go`.
- Backend tests in `backend/main_test.go` cover the happy path and at least one range/status assertion.
- The existing `GET /api/hello`, CORS behavior, and `/openapi.json` remain untouched.

## Impact

- **Affected specs**: New `dashboard-people-availability` capability spec.
- **Affected code**: `backend/main.go` (route registration), new `backend/dashboard.go` (handler + seed data), `backend/main_test.go` (tests).
- **Non-breaking**: No existing endpoint, route, CORS configuration, or OpenAPI path is modified or removed.
- **Frontend impact**: None in this task. The endpoint contract is designed so the frontend can later replace the static "Available today" card and add the multi-day grid with a single fetch.
- **OpenAPI**: The new endpoint auto-appears in `/openapi.json` via Huma registration.
- **Compose**: No compose changes required; the existing `app-network` bridge and CORS config cover the new route.