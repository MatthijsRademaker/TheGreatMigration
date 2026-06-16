-- +goose Up
-- +goose StatementBegin
-- Add nullable task_id FK to backlog_tasks so schedule cards can reference backlog tasks.
ALTER TABLE schedule_task_cards ADD COLUMN task_id TEXT REFERENCES backlog_tasks(id);

-- Index for efficient lookups when checking task references before delete.
CREATE INDEX idx_schedule_task_cards_task_id ON schedule_task_cards(task_id);

-- Existing cards have no backlog reference; task_id remains NULL for those.
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_schedule_task_cards_task_id;
ALTER TABLE schedule_task_cards DROP COLUMN IF EXISTS task_id;
-- +goose StatementEnd
