## Why

The `/` route still stops at a partial placeholder: `KpiCards` renders, but the rest of `HomeView.vue` is still a `Today’s plan` card plus a static `Move notes` card. The repository now already contains reusable dashboard panels for task backlog, people availability, and daily schedule, and the archived ASCII homepage blueprint defines the intended row structure for those panels. This change composes the home page from those existing pieces without pulling backend wiring, shell changes, or new domain components into scope.

## What Changes

- Update `frontend/src/home/HomeView.vue` to act as a composition shell for the dashboard.
- Keep `<KpiCards />` as the first row and preserve its existing live-query behavior.
- Replace the `Today’s plan` placeholder with a responsive row that reuses:
  - `TaskManagementPanel` from `@/tasks/components/TaskManagementPanel.vue`
  - `PeopleAvailability` from `@/people/PeopleAvailability.vue`
- Add a second responsive row that reuses:
  - `DailySchedule` from `@/calendar/DailySchedule.vue`
  - the existing static `Move notes` card content from `HomeView.vue`
- Remove the placeholder-only `upcomingWork` array and the now-unused `Badge` import.
- Use two independent grid containers with the existing `xl:grid-cols-[1.4fr_0.9fr]` layout pattern so the wider left-column panels mirror the ASCII blueprint without changing shell or route structure.
- Update `frontend/tests/app-routes-render.test.ts` so the `/` route asserts the composed dashboard sections and no longer depends on `Today’s plan`.

## Impact

- The home page becomes a composed dashboard surface aligned with the archived ASCII layout and the current frontend component inventory.
- The change remains frontend-only and compositional: no backend wiring, no generated-client regeneration, no route changes, no `AppShell.vue` or `AppSidebar.vue` changes, no Move Notes CRUD/component work, and no canonical OpenSpec edits.
- Verification stays bounded to frontend route rendering plus the project-standard `scripts/precommit-run`.
