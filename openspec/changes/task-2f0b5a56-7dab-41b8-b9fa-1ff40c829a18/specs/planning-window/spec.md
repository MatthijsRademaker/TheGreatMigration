## MODIFIED Requirements

### Requirement: PeopleView SHALL host the People Availability component while still deferring planning-window integration

`frontend/src/people/PeopleView.vue` SHALL replace the previous placeholder helper cards with the presentational `PeopleAvailability` component. This change SHALL remain independent of `planWindowDays`, generated dashboard availability queries, and backend-derived availability data. The `/people` route SHALL continue to render deterministically from local demo data or explicit props until a later data-wiring change is scheduled.

#### Scenario: PeopleView shows the new presentational panel

- **WHEN** `PeopleView` renders after this change
- **THEN** it displays the People Availability panel instead of the old helper availability placeholder cards
- **AND** the route continues to render inside the existing app shell

#### Scenario: PeopleView still does not depend on planning-window or backend availability data

- **WHEN** the PeopleView implementation is inspected
- **THEN** it does not import `planWindowDays`
- **AND** it does not call the dashboard people-availability query
- **AND** it renders without backend access
