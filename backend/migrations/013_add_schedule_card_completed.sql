-- +goose Up
-- +goose StatementBegin
ALTER TABLE schedule_task_cards ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE schedule_task_cards DROP COLUMN completed;
-- +goose StatementEnd
