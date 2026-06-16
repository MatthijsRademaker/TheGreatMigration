## 1. Fix the composable's pending signal

- [x] 1.1 In `frontend/src/shared/composables/useUpdatePlanningWindow.ts`, destructure `isLoading` from `useMutation()` and rename it to `isPending`, replacing the `status`-derived computed

## 2. Update the mock in settings-view tests

- [x] 2.1 In `frontend/tests/settings-view.test.ts`, change the mock's `controlledIsPending` ref to `controlledIsLoading` and wire it to `isLoading` instead of the old `isPending`

## 3. Add mutation lifecycle tests

- [x] 3.1 In `frontend/tests/settings-view.test.ts`, add a test that verifies `isPending` is `false` initially (before any mutation call)
- [x] 3.2 Add a test that verifies `isPending` is `true` during mutation execution and returns to `false` after completion

## 4. Verify the fix

- [x] 4.1 Run `scripts/precommit-run` and confirm all checks pass
