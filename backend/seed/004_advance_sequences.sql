-- +goose Up
-- +goose StatementBegin
-- The schema migrations (007 rooms_areas_id_seq, 008 backlog_tasks_id_seq,
-- 011 people_id_seq) now run against empty tables, so each sequence sits at its
-- creation baseline. After seeding hardcoded IDs we advance each sequence past
-- the seeded maximum so the next CreatePerson/CreateRoom/CreateTask produces a
-- non-colliding sequential ID (p9, room-9, task-12).
--
-- schedule_task_cards uses a SERIAL whose sequence auto-advances on insert, so
-- it needs no manual setval here.
SELECT setval('people_id_seq', 8);
SELECT setval('rooms_areas_id_seq', 8);
SELECT setval('backlog_tasks_id_seq', 11);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Restore each sequence to the baseline set by its schema migration against an
-- empty table.
SELECT setval('people_id_seq', 0);
SELECT setval('rooms_areas_id_seq', 0);
SELECT setval('backlog_tasks_id_seq', 100, false);
-- +goose StatementEnd
