-- +goose Up
-- +goose StatementBegin
-- Add scheduled_date column (nullable initially to backfill).
ALTER TABLE schedule_task_cards ADD COLUMN scheduled_date DATE;

-- Backfill scheduled_date from day_group relative to the planning-window
-- start date (2026-07-05).  day_group 0 → 2026-07-05, 1 → 2026-07-06, etc.
UPDATE schedule_task_cards
SET scheduled_date = '2026-07-05'::date + (day_group * interval '1 day');

-- Enforce NOT NULL after backfill.
ALTER TABLE schedule_task_cards ALTER COLUMN scheduled_date SET NOT NULL;

-- Drop the now-redundant day_group column.
ALTER TABLE schedule_task_cards DROP COLUMN day_group;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE schedule_task_cards ADD COLUMN day_group INTEGER;
UPDATE schedule_task_cards
SET day_group = EXTRACT(DAY FROM scheduled_date - '2026-07-05'::date)::int;
ALTER TABLE schedule_task_cards ALTER COLUMN day_group SET NOT NULL;
ALTER TABLE schedule_task_cards DROP COLUMN scheduled_date;
-- +goose StatementEnd
