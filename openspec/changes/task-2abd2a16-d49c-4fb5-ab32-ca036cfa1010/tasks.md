## 1. Backend — Add pagination params to dashboard API

- [ ] 1.1 Add `offset` and `limit` fields to `DashboardInput` struct in `backend/api/dashboard.go`. `Offset` SHALL be a non-negative integer with default 0. `Limit` SHALL be a non-negative integer with default 0 (meaning no limit).
- [ ] 1.2 Add `Pagination` struct to `backend/api/dashboard.go` with `TotalPeople` (int), `Page` (int), and `PerPage` (int) fields. This captures the pagination metadata in the `DashboardBody`.
- [ ] 1.3 Add a `Pagination` field of type `Pagination` to the `DashboardBody` struct.

## 2. Backend — Update Store interface and handler

- [ ] 2.1 Update `Store` interface in `backend/api/store.go`: Change `GetPeopleAvailability(ctx, startDate time.Time, days int)` to `GetPeopleAvailability(ctx, startDate time.Time, days int, offset int, limit int)`.
- [ ] 2.2 Update the handler in `backend/api/dashboard.go`: Parse `input.Offset` and `input.Limit` and pass them to `store.GetPeopleAvailability(ctx, startDate, days, offset, limit)`. Compute `DashboardPagination` from offset/limit and the returned `DashboardBody.Summary.TotalPeople`.
- [ ] 2.3 Update the handler to compute `pagination.Page` as 1-indexed from offset and limit: `page = (offset / limit) + 1` when limit > 0, otherwise `page = 1`. Set `pagination.PerPage` to limit when limit > 0, otherwise to `summary.TotalPeople`.

## 3. Backend — Update PgStore implementation

- [ ] 3.1 Update `backend/store.go`: Change `GetPeopleAvailability` signature to accept `offset int, limit int`.
- [ ] 3.2 Add a new sqlc query in `backend/queries/people_availability.sql`: `GetPeoplePaginated` that accepts `offset` and `limit` parameters and returns a subset of people with `LIMIT $1 OFFSET $2`.
- [ ] 3.3 Regenerate sqlc Go code: Run `sqlc generate` in `backend/` (or whatever the project's codegen command is) to produce the new `GetPeoplePaginated` function.
- [ ] 3.4 Update `GetPeopleAvailability` in `backend/store.go`: When `limit > 0`, call the new paginated query. When `limit == 0`, call the existing `GetAllPeople` query. Add a `SELECT COUNT(*) FROM people` query to get `totalPeople` for the response. Wire the pagination metadata into the returned `DashboardBody`.
- [ ] 3.5 Update the `GetPeopleAvailability` response to include pagination metadata in the `DashboardBody.Pagination` field.

## 4. Backend — Update mock store and tests

- [ ] 4.1 Update `backend/store_mock_test.go`: Change the mock `GetPeopleAvailability` signature to match the updated interface. Update the mock implementation to handle offset/limit (slice the full results list).
- [ ] 4.2 Add test scenarios in `backend/main_test.go`:
  - `GET /api/dashboard/people-availability?offset=0&limit=3` returns exactly 3 people with correct pagination metadata
  - `GET /api/dashboard/people-availability?offset=10&limit=3` returns empty people array when offset exceeds total
  - `GET /api/dashboard/people-availability` (no pagination params) returns all people with page=1 and perPage=total
- [ ] 4.3 Run `go test ./...` in `backend/` and verify all tests pass.
- [ ] 4.4 Regenerate OpenAPI snapshot: Run `go run ./cmd/openapi-gen > ../frontend/openapi-snapshot.json` and then regenerate the frontend client with `cd ../frontend && npm run regen:api`.

## 5. Frontend — Update server-state types and composable

- [ ] 5.1 Verify the generated frontend client reflects the new `offset`/`limit` query params and `pagination` field in the response. If the types.gen.ts file is stale, run codegen.
- [ ] 5.2 Update `frontend/src/shared/composables/usePeopleAvailability.ts`:
  - Accept `page` (default 1), `daysPerPage` (default 7), `offset` (default 0), `limit` (default 0) in options.
  - Compute `currentStart = planningWindowStart + (page - 1) * daysPerPage` and pass it as `start` query param.
  - Pass `daysPerPage` as the `days` query param.
  - Pass `offset` and `limit` as query params when they are non-default.
  - Expose `currentPage`, `totalPages`, `daysPerPage`, `totalDays` as computed refs.
  - Expose `totalPeople` from the response's `summary.totalPeople`.
- [ ] 5.3 Verify the composable still works when called without arguments in `HomeView.vue` — it should default to page=1, daysPerPage=7, offset=0, limit=0 (all people, first 7 days).

## 6. Frontend — Update PeopleView with day pagination

- [ ] 6.1 Update `frontend/src/people/PeopleView.vue`:
  - Destructure `currentPage`, `totalPages`, `daysPerPage`, `totalDays` from `usePeopleAvailability()`.
  - Add reactive `page` ref initialized to 1.
  - Add `goToPrevPage` and `goToNextPage` functions that decrement/increment `page` and pass it to the composable.
  - Render navigation bar above the matrix showing:
    - Date range label (e.g., "Sun 5 Jul – Sat 11 Jul")
    - Page indicator ("Page 1 of 6")
    - Previous button (disabled at page 1)
    - Next button (disabled at last page)
- [ ] 6.2 Update `frontend/src/people/types.ts`: Add `PaginationMetadata` type with `totalPeople`, `page`, `perPage` fields. Add optional `pagination` prop to `PeopleAvailabilityProps` if the component needs it (the navigation controls live in PeopleView, not PeopleAvailability, so this may not be needed).

## 7. Frontend — Tests

- [ ] 7.1 Update or add tests in `frontend/src/shared/composables/__tests__/usePeopleAvailability.test.ts`:
  - Page 1 with 7 days produces correct `start` and `days` query params
  - Page 2 produces start offset by 7 days
  - Default page and daysPerPage values are used when not provided
  - totalPages is computed correctly from planning window length
- [ ] 7.2 Update `frontend/src/people/__tests__/PeopleView.spec.ts`:
  - Navigation controls render with correct labels
  - Previous button disabled on page 1
  - Next button disabled on last page
  - Clicking next increments page and re-fetches data
- [ ] 7.3 Run frontend tests and verify all pass.

## 8. Verification

- [ ] 8.1 Run `scripts/precommit-run` and verify all lint, type-check, and test checks pass without errors.
- [ ] 8.2 Verify the change manually by starting backend and frontend, navigating to `/people`, confirming the matrix shows 7 days by default, clicking Next to see the next 7-day slice, and confirming availability upserts work for dates on any page.
