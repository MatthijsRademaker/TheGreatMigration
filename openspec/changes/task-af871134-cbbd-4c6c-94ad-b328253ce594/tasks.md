## 1. Types and event contract

- [ ] 1.1 Add `editable?: boolean` field to the `PeopleAvailabilityProps` interface in `frontend/src/people/types.ts`
- [ ] 1.2 Add `CellChangePayload` type (`{ personId: string; dayIndex: number; status: AvailabilityStatus | null }`) and export it from `frontend/src/people/types.ts`

## 2. Editable mode in PeopleAvailability component

- [ ] 2.1 Import `Popover`, `PopoverContent`, `PopoverTrigger` from `@/shared/ui/popover` and `Separator` from `@/shared/ui/separator` in `PeopleAvailability.vue`
- [ ] 2.2 Add `const emit = defineEmits<{ 'update-cell': [payload: CellChangePayload]; 'delete-person': [personId: string] }>()` to the component script
- [ ] 2.3 Add a `const statusOptions: AvailabilityStatus[] = ['available', 'busy', 'partial', 'off']` constant in the script section
- [ ] 2.4 In the template, conditionally wrap each status cell badge with `<Popover>` when `editable` is true: use `<PopoverTrigger as-child>` on the Badge, and render the inline status picker inside `<PopoverContent>`
- [ ] 2.5 Build the inline status picker template inside PopoverContent: four clickable options (each showing a Badge with the status variant and label), a `<Separator>`, and a "Clear (reset to off)" button
- [ ] 2.6 Wire click handlers: selecting a status emits `update-cell` with the status value; clicking "Clear" emits `update-cell` with `status: null`
- [ ] 2.7 When `editable` is true, add an extra `<th>` column header for actions in `<thead>` and an extra `<td>` in each person row containing a delete button (variant="destructive" size="sm" or icon button) that emits `delete-person` on click
- [ ] 2.8 When `editable` is false or not provided, render the existing read-only template exactly as-is (no Popovers, no action column)

## 3. PeopleView redistribution

- [ ] 3.1 Pass `:editable="true"` to the `<PeopleAvailability v-bind="availabilityData" />` call in PeopleView.vue
- [ ] 3.2 Wire `@update-cell` on the component: map `dayIndex` to an ISO date via `getISODate(payload.dayIndex)`, then call the upsert mutation if `payload.status` is non-null or the delete-availability mutation if `payload.status` is null
- [ ] 3.3 Wire `@delete-person` on the component to the existing `handleDelete` function
- [ ] 3.4 Remove the "Manage people" `<Card>` block (the entire per-person management controls section with the status picker, person list, and duplicate availability badges)
- [ ] 3.5 Keep the create-person form card above the matrix (no change to its template)
- [ ] 3.6 Keep the error display (`createError`, `statusError`, `deleteError`, `clearAvailabilityError`) and render them above the matrix so mutation errors are visible
- [ ] 3.7 Clean up unused variables: `editingCell`, `handleStatusUpdate`, `handleClearAvailability`, `statusOptions`, `statusError`, `clearAvailabilityError`, `getISODate` (still needed by the event handler), and any imports that are no longer referenced

## 4. Test updates

- [ ] 4.1 Add client-render tests (using `mount` not `renderToString`) in `PeopleAvailability.test.ts` for editable mode: render with `editable={true}`, assert that status cells are wrapped in Popover triggers, simulate status selection and assert emitted events
- [ ] 4.2 Add a test in `PeopleAvailability.test.ts` for read-only mode: render without `editable` prop and assert no Popover attributes appear in the output
- [ ] 4.3 Update `PeopleView.test.ts` to remove assertions about the old "Manage people" card content (status editing section, person list with duplicate badges, status picker UI)
- [ ] 4.4 Add assertions in `PeopleView.test.ts` that the matrix receives `editable={true}` and that the create-person form still renders above the matrix
- [ ] 4.5 Run `scripts/precommit-run` and verify all tests pass
