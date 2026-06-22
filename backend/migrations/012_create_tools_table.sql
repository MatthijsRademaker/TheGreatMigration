-- +goose Up
-- +goose StatementBegin
CREATE TABLE tools (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    brought_by TEXT REFERENCES people(id) ON DELETE SET NULL,
    sort_order INTEGER NOT NULL
);

-- Sequence for atomic tool ID generation, matching the rooms/tasks pattern.
CREATE SEQUENCE IF NOT EXISTS tools_id_seq;

-- Seed a few demo tools so the list, KPI card, and empty-state are exercisable.
-- Two are claimed (crossed off), three are open.
INSERT INTO tools (id, name, brought_by, sort_order) VALUES
    ('tool-1', 'Ladder', NULL, 1),
    ('tool-2', 'Power drill', 'p1', 2),
    ('tool-3', 'Moving dolly', NULL, 3),
    ('tool-4', 'Tarps', NULL, 4),
    ('tool-5', 'Toolbox', 'p3', 5);

-- Advance the sequence past the seeded IDs so the next CreateTool is tool-6.
SELECT setval('tools_id_seq', 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tools;
DROP SEQUENCE IF EXISTS tools_id_seq;
-- +goose StatementEnd
