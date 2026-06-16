-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE backlog_tasks_id_seq START 100;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE IF EXISTS backlog_tasks_id_seq;
-- +goose StatementEnd
