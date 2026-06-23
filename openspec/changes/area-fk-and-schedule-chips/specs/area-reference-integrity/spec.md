## ADDED Requirements

### Requirement: Backlog tasks and schedule cards SHALL reference rooms_areas by foreign key

`backlog_tasks` and `schedule_task_cards` SHALL store their room/area as a non-null `area_id TEXT` column constrained by a foreign key to `rooms_areas(id)`. The free-text `backlog_tasks.room` and `schedule_task_cards.room_area` columns SHALL be removed. Inserting or updating a task or schedule card with an `area_id` that does not exist in `rooms_areas` SHALL fail.

#### Scenario: A task cannot be created with an unknown area id

- **WHEN** a backlog task or schedule card is written with an `area_id` not present in `rooms_areas`
- **THEN** the write fails (foreign-key violation surfaced as a 400/404 at the API)
- **AND** no row is persisted

#### Scenario: The legacy free-text room columns no longer exist

- **WHEN** the schema is inspected after migration
- **THEN** `backlog_tasks.room` and `schedule_task_cards.room_area` are absent
- **AND** both tables expose a non-null `area_id` foreign key to `rooms_areas(id)`

### Requirement: The migration SHALL backfill existing room strings without data loss

The migration SHALL resolve every existing `backlog_tasks.room` and `schedule_task_cards.room_area` value to a `rooms_areas` id. When the string matches one or more catalog names, the row with the lowest `id` SHALL be chosen deterministically. When the string matches no catalog name, the migration SHALL insert a new `rooms_areas` row with that name and `type = 'area'`, and link to it. No task or schedule card SHALL be left without a valid `area_id`.

#### Scenario: Existing string matching a catalog name is linked to that area

- **WHEN** a card's `room_area` equals an existing `rooms_areas.name`
- **THEN** its `area_id` is set to that catalog row's id

#### Scenario: Duplicate catalog names resolve deterministically

- **WHEN** a string matches more than one `rooms_areas` row with the same name
- **THEN** the `area_id` is set to the matching row with the lowest `id`
- **AND** repeating the backfill yields the same mapping

#### Scenario: Orphan string auto-creates a catalog area

- **WHEN** a card's `room_area` matches no `rooms_areas.name`
- **THEN** a new `rooms_areas` row is created with that name and `type = 'area'`
- **AND** the card's `area_id` references the newly created row
- **AND** no card or task is dropped or left with a null `area_id`
