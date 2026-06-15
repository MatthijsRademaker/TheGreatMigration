## 1. Backend endpoint implementation

- [x] Create `backend/dashboard.go` under package `main` with:
  - Request/response data types matching the contract (`DashboardInput`, `DashboardOutput`, `Range`, `Summary`, `Person`, `AvailabilityEntry`, `StatusLegend`).
  - Huma operation registration tags for `GET /api/dashboard/people-availability`.
  - In-memory seed data: at least 8 people with stable IDs (`p1`–`p8`), names, initials, and daily availability exercising all four statuses (`available`, `busy`, `partial`, `off`) across the default 4-day range.
  - Handler logic: parse `start` and `days` query params with defaults, filter seed data to the requested range, compute `summary.availableToday` (count of `available` status on `selectedDate`), and return the combined payload including a full status legend.

- [x] Register the new endpoint in `backend/main.go` by calling the handler from `dashboard.go` via `huma.Register` following the existing hello-world pattern.

- [x] Verify `go vet ./...` and `go build ./...` pass in `backend/` with the new file.

## 2. Backend tests

- [x] Add `TestDashboardPeopleAvailability` in `backend/main_test.go` covering:
  - Happy path: 200 OK, valid JSON structure, correct content-type.
  - Response shape: `range.startDate`, `range.endDate`, `range.days`, `range.selectedDate`, `summary.availableToday`, `summary.totalPeople`, `people[]` with `id`, `name`, `initials`, `availability[]` with `{date, status}`, and `statuses[]`.
  - Status validation: all availability entries use only canonical statuses (`available`, `busy`, `partial`, `off`).
  - Range/status assertion: the count of people with `available` on `selectedDate` equals `summary.availableToday`.
  - Existing `/api/hello` still returns 200 with the correct message.

- [x] Run `go test ./...` in `backend/` and confirm all tests pass.

## 3. Verification

- [x] Run `scripts/precommit-run` from the repo root.
- [x] Confirm existing behavior unchanged: `GET /api/hello` returns `{"message": "Hello from the backend!"}`, CORS headers are present for allowed origins, and `/openapi.json` includes both `/api/hello` and the new `/api/dashboard/people-availability` endpoint.
