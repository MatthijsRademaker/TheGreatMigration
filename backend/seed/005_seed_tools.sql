-- +goose Up
-- +goose StatementBegin
-- Seed a few demo tools so the list, KPI card, and empty-state are exercisable.
-- Two are claimed (brought_by p1/p3, seeded in 001), three are open.
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
DELETE FROM tools;
SELECT setval('tools_id_seq', 1, false);
-- +goose StatementEnd
