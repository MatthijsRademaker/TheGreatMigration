-- +goose Up
-- +goose StatementBegin
-- Replace the free-text room/room_area columns with a real foreign key to
-- rooms_areas. Existing values are backfilled by name (lowest matching id on
-- duplicates); strings with no catalog match auto-create a rooms_areas row
-- (type 'area') so no task or card is left without a valid area_id.

-- 1. Add nullable area_id to both tables.
ALTER TABLE backlog_tasks ADD COLUMN area_id TEXT;
ALTER TABLE schedule_task_cards ADD COLUMN area_id TEXT;

-- 2. Auto-create catalog rows for orphan strings (distinct across both tables).
INSERT INTO rooms_areas (id, name, type)
SELECT 'room-' || nextval('rooms_areas_id_seq'), s.name, 'area'
FROM (
    SELECT DISTINCT room AS name FROM backlog_tasks
    UNION
    SELECT DISTINCT room_area AS name FROM schedule_task_cards
) s
WHERE NOT EXISTS (SELECT 1 FROM rooms_areas ra WHERE ra.name = s.name);

-- 3. Backfill area_id by name, choosing the numerically lowest matching id.
UPDATE backlog_tasks bt
SET area_id = (
    SELECT ra.id FROM rooms_areas ra
    WHERE ra.name = bt.room
    ORDER BY CAST(SUBSTRING(ra.id FROM 6) AS INTEGER) ASC
    LIMIT 1
);

UPDATE schedule_task_cards sc
SET area_id = (
    SELECT ra.id FROM rooms_areas ra
    WHERE ra.name = sc.room_area
    ORDER BY CAST(SUBSTRING(ra.id FROM 6) AS INTEGER) ASC
    LIMIT 1
);

-- 4. Enforce NOT NULL and the foreign key.
ALTER TABLE backlog_tasks ALTER COLUMN area_id SET NOT NULL;
ALTER TABLE backlog_tasks ADD CONSTRAINT fk_backlog_area
    FOREIGN KEY (area_id) REFERENCES rooms_areas(id);

ALTER TABLE schedule_task_cards ALTER COLUMN area_id SET NOT NULL;
ALTER TABLE schedule_task_cards ADD CONSTRAINT fk_card_area
    FOREIGN KEY (area_id) REFERENCES rooms_areas(id);

-- 5. Drop the legacy free-text columns.
ALTER TABLE backlog_tasks DROP COLUMN room;
ALTER TABLE schedule_task_cards DROP COLUMN room_area;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE backlog_tasks ADD COLUMN room TEXT;
ALTER TABLE schedule_task_cards ADD COLUMN room_area TEXT;

UPDATE backlog_tasks bt
SET room = (SELECT ra.name FROM rooms_areas ra WHERE ra.id = bt.area_id);
UPDATE schedule_task_cards sc
SET room_area = (SELECT ra.name FROM rooms_areas ra WHERE ra.id = sc.area_id);

ALTER TABLE backlog_tasks ALTER COLUMN room SET NOT NULL;
ALTER TABLE schedule_task_cards ALTER COLUMN room_area SET NOT NULL;

ALTER TABLE backlog_tasks DROP CONSTRAINT fk_backlog_area;
ALTER TABLE backlog_tasks DROP COLUMN area_id;
ALTER TABLE schedule_task_cards DROP CONSTRAINT fk_card_area;
ALTER TABLE schedule_task_cards DROP COLUMN area_id;
-- +goose StatementEnd
