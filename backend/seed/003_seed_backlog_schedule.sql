-- +goose Up
-- +goose StatementBegin
-- Backlog tasks (11 rows, matching seedTasks in tasks.go exactly).
-- area_id references rooms_areas seeded in 002_seed_rooms_areas.sql:
-- room-1 Kitchen, room-2 Living Room, room-3 Bedroom 1, room-4 Bedroom 2,
-- room-5 Garage, room-6 Bedroom, room-7 Storage, room-8 Office.
INSERT INTO backlog_tasks (id, title, priority, people_needed, area_id, status, sort_order) VALUES
    ('task-1', 'Disconnect kitchen appliances', 'high', 3, 'room-1', 'backlog', 1),
    ('task-2', 'Wrap living room furniture', 'high', 2, 'room-2', 'ready', 2),
    ('task-3', 'Pack kitchen fragile items', 'high', 2, 'room-1', 'assigned', 3),
    ('task-4', 'Disassemble bedroom furniture', 'high', 2, 'room-3', 'assigned', 4),
    ('task-5', 'Sort and label moving boxes', 'medium', 3, 'room-2', 'backlog', 5),
    ('task-6', 'Clear garage shelving', 'medium', 1, 'room-5', 'ready', 6),
    ('task-7', 'Move bedroom wardrobe', 'medium', 3, 'room-4', 'assigned', 7),
    ('task-8', 'Sweep garage floor', 'low', 1, 'room-5', 'backlog', 8),
    ('task-9', 'Dust living room shelves', 'low', 2, 'room-2', 'ready', 9),
    ('task-10', 'Wipe down kitchen counters', 'low', 2, 'room-1', 'assigned', 10),
    ('task-11', 'Inventory bedroom closet', 'medium', 3, 'room-3', 'ready', 11);

-- Backlog task assignments (matching the assignedTo arrays in seedTasks).
INSERT INTO backlog_task_assignments (task_id, person_id, sort_order) VALUES
    ('task-3', 'p1', 0),
    ('task-4', 'p2', 0),
    ('task-4', 'p3', 1),
    ('task-5', 'p4', 0),
    ('task-7', 'p5', 0),
    ('task-7', 'p6', 1),
    ('task-7', 'p7', 2),
    ('task-8', 'p8', 0),
    ('task-10', 'p1', 0),
    ('task-10', 'p2', 1),
    ('task-11', 'p3', 0);

-- Schedule task cards (8 rows, matching seedTasksForDay groups 0-3).
-- SERIAL ids will be 1-8 in insertion order.
-- This seed runs after all schema migrations, so the post-014 schema applies:
-- the dropped room_area column is replaced by an area_id foreign key.
INSERT INTO schedule_task_cards (title, priority, area_id, people_needed, scheduled_date, sort_order) VALUES
    -- day_group 0 -> 2026-07-05
    ('Kitchen deep clean', 'high', 'room-1', 2, '2026-07-05', 0),
    ('Window washing', 'medium', 'room-2', 2, '2026-07-05', 1),
    -- day_group 1 -> 2026-07-06
    ('Furniture assembly', 'high', 'room-6', 1, '2026-07-06', 0),
    ('Packing supplies inventory', 'low', 'room-7', 3, '2026-07-06', 1),
    -- day_group 2 -> 2026-07-07
    ('Electronics setup', 'medium', 'room-8', 2, '2026-07-07', 0),
    ('Declutter garage', 'low', 'room-5', 1, '2026-07-07', 1),
    -- day_group 3 -> 2026-07-08
    ('Box labeling', 'medium', 'room-7', 2, '2026-07-08', 0),
    ('Curtain removal', 'high', 'room-2', 2, '2026-07-08', 1);

-- Schedule task assignments (matching seedTasksForDay assigneeIds).
-- Card ids: 1=Kitchen deep clean, 2=Window washing, 3=Furniture assembly,
-- 4=Packing supplies inventory, 5=Electronics setup, 6=Declutter garage,
-- 7=Box labeling, 8=Curtain removal
INSERT INTO schedule_task_assignments (task_card_id, person_id, sort_order) VALUES
    (1, 'p1', 0),
    (1, 'p2', 1),
    (2, 'p3', 0),
    (3, 'p4', 0),
    (4, 'p5', 0),
    (5, 'p6', 0),
    (5, 'p7', 1),
    (7, 'p8', 0),
    (8, 'p1', 0),
    (8, 'p3', 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM schedule_task_assignments;
DELETE FROM schedule_task_cards;
DELETE FROM backlog_task_assignments;
DELETE FROM backlog_tasks;
-- +goose StatementEnd
