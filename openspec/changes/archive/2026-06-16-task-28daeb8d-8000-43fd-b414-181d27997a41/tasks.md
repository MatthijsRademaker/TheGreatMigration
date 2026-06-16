## 1. CalendarView Room Select

- [x] 1.1 Import `listRoomsQuery` from the generated Pinia Colada client in `CalendarView.vue`
- [x] 1.2 Add `useQuery(listRoomsQuery())` call to fetch room data alongside the existing schedule query
- [x] 1.3 Replace the free-form `Input` for "Room / Area" in the modal template with a `Select` populated from the rooms query data
- [x] 1.4 Add loading placeholder state for the room Select while the rooms query is pending
- [x] 1.5 Add error message with retry control for the room Select when the rooms query fails

## 2. CalendarView DatePicker

- [x] 2.1 Import the shared `DatePicker` component from `@/shared/ui/date-picker`
- [x] 2.2 Replace the free-form text `Input` for "Scheduled date" in the modal template with the `DatePicker` component
- [x] 2.3 Wire the DatePicker v-model to `formScheduledDate` ensuring `YYYY-MM-DD` format output

## 3. TasksView Room Select

- [x] 3.1 Import `listRoomsQuery` from the generated Pinia Colada client in `TasksView.vue`
- [x] 3.2 Add `useQuery(listRoomsQuery())` call to fetch room data alongside the existing backlog query
- [x] 3.3 Replace the free-form `Input` for "Room / Area" in the modal template with a `Select` populated from the rooms query data
- [x] 3.4 Add loading placeholder state for the room Select while the rooms query is pending
- [x] 3.5 Add error message with retry control for the room Select when the rooms query fails

## 4. Frontend Tests

- [x] 4.1 Update CalendarView modal tests to assert room Select is rendered with room names from mocked `listRoomsQuery` response
- [x] 4.2 Add CalendarView test asserting room Select loading state
- [x] 4.3 Add CalendarView test asserting room Select error state with retry
- [x] 4.4 Update CalendarView modal tests to assert DatePicker is rendered instead of date text Input
- [x] 4.5 Update TasksView modal tests to assert room Select is rendered with room names from mocked `listRoomsQuery` response
- [x] 4.6 Add TasksView test asserting room Select loading state
- [x] 4.7 Add TasksView test asserting room Select error state with retry
- [x] 4.8 Run `scripts/precommit-run` and verify all checks pass
