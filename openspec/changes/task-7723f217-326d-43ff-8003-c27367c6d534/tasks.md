## 1. Database Migration

- [ ] 1.1 Create `backend/migrations/009_add_people_id_sequence.sql` with `CREATE SEQUENCE IF NOT EXISTS people_id_seq` and seed from `MAX(CAST(SUBSTRING(id FROM 2) AS INTEGER))` from the `people` table (matching the room-sequence pattern in `007_add_rooms_areas_sequence.sql`)
- [ ] 1.2 Update `backend/queries/people_availability.sql`: change `CreatePerson` from `:exec` with `$1` for id to `:one` with `RETURNING` using `'p' || nextval('people_id_seq')` for the id — matching the room pattern

## 2. Regenerate sqlc

- [ ] 2.1 Run sqlc from `src/` to regenerate `backend/db/people_availability.sql.go`

## 3. Backend Store Interface

- [ ] 3.1 Update `backend/api/store.go`: change `CreatePerson(ctx, id, name, initials string) error` to `CreatePerson(ctx, name, initials string) (string, error)` — returns the generated ID

## 4. Backend Store Implementation

- [ ] 4.1 Update `backend/store.go`: update `PgStore.CreatePerson` to call the new sqlc query and return the generated ID

## 5. Backend Test Stores

- [ ] 5.1 Update `backend/store_mock_test.go`: change `mockStore.CreatePerson` signature to `(ctx, name, initials string) (string, error)`, return a fake ID
- [ ] 5.2 Update `backend/main_test.go`: change `peopleTestStore.CreatePerson` signature to `(ctx, name, initials string) (string, error)`, auto-generate sequential `p{N}` IDs
- [ ] 5.3 Update `nilPlanningWindowStore.CreatePerson`, `failingStore.CreatePerson`, and `partialFailingStore.CreatePerson` in `backend/main_test.go` to match the new interface signature

## 6. Backend API Handler

- [ ] 6.1 Update `backend/api/people.go`: remove `ID` field from `CreatePersonInput.Body`, update the POST handler to call the new `store.CreatePerson(ctx, name, initials)` signature, remove the duplicate-ID check (sequence guarantees uniqueness), construct the response Person using the returned ID
- [ ] 6.2 Run `scripts/precommit-run` from `src/` to verify backend compiles

## 7. Backend Tests

- [ ] 7.1 Update `TestCreatePerson` in `backend/main_test.go`: send body without `id`, verify response contains an auto-generated `id` matching `p{N}` pattern
- [ ] 7.2 Update `TestCreatePersonDuplicate` in `backend/main_test.go`: seed via store, send body without `id`, verify the store-generated IDs don't collide (the duplicate case becomes impossible — delete this test or repurpose)
- [ ] 7.3 Update `TestCreatePersonMissingFields` in `backend/main_test.go`: send body with only `name` and `initials` (no `id`), verify 422 on empty `name`/`initials`
- [ ] 7.4 Update any other backend tests that call `store.CreatePerson` with the new `(ctx, name, initials)` signature and handle the returned ID
- [ ] 7.5 Run `go test ./...` from `src/` to verify all tests pass

## 8. OpenAPI Snapshot & Frontend Client

- [ ] 8.1 Refresh `frontend/openapi-snapshot.json` from the running backend
- [ ] 8.2 Regenerate the frontend API client types from `frontend/` to update `CreatePersonData` schema (remove `id` from request body)

## 9. Frontend Form

- [ ] 9.1 Update `frontend/src/people/PeopleView.vue`: remove the ID input field (`newId` ref, label, Input component), update `handleCreate` to send only `name` and `initials` in the mutation body, remove the empty-ID validation
- [ ] 9.2 Run `scripts/precommit-run` from `src/` to verify frontend compiles and lints

## 10. Frontend Tests

- [ ] 10.1 Update any frontend tests that reference the person creation form or ID field to match the new form shape
- [ ] 10.2 Run `scripts/precommit-run` from `src/` to verify all checks pass

## 11. Final Verification

- [ ] 11.1 Run `scripts/precommit-run` from `src/` for a clean full verification pass
