-- +goose Up
-- +goose StatementBegin
CREATE TABLE tools (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    brought_by TEXT REFERENCES people(id) ON DELETE SET NULL,
    sort_order INTEGER NOT NULL
);

-- Sequence for atomic tool ID generation, matching the rooms/tasks pattern.
-- Demo tools are seeded by the separate seed dataset (backend/seed), not here.
CREATE SEQUENCE IF NOT EXISTS tools_id_seq;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tools;
DROP SEQUENCE IF EXISTS tools_id_seq;
-- +goose StatementEnd
