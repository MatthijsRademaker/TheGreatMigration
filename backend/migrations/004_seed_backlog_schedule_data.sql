-- +goose Up
-- +goose StatementBegin
-- Backlog tasks (11 rows, matching seedTasks in tasks.go exactly).
INSERT INTO backlog_tasks (id, title, priority, people_needed, room, status, sort_order) VALUES
    ('task-1', 'Disconnect kitchen appliances', 'high', 3, 'Kitchen', 'backlog', 1),
    ('task-2', 'Wrap living room furniture', 'high', 2, 'Living Room', 'ready', 2),
    ('task-3', 'Pack kitchen fragile items', 'high', 2, 'Kitchen', 'assigned', 3),
    ('task-4', 'Disassemble bedroom furniture', 'high', 2, 'Bedroom 1', 'assigned', 4),
    ('task-5', 'Sort and label moving boxes', 'medium', 3, 'Living Room', 'backlog', 5),
    ('task-6', 'Clear garage shelving', 'medium', 1, 'Garage', 'ready', 6),
    ('task-7', 'Move bedroom wardrobe', 'medium', 3, 'Bedroom 2', 'assigned', 7),
    ('task-8', 'Sweep garage floor', 'low', 1, 'Garage', 'backlog', 8),
    ('task-9', 'Dust living room shelves', 'low', 2, 'Living Room', 'ready', 9),
    ('task-10', 'Wipe down kitchen counters', 'low', 2, 'Kitchen', 'assigned', 10),
    ('task-11', 'Inventory bedroom closet', 'medium', 3, 'Bedroom 1', 'ready', 11);

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
INSERT INTO schedule_task_cards (title, priority, room_area, people_needed, day_group, sort_order) VALUES
    -- day_group 0
    ('Kitchen deep clean', 'high', 'Kitchen', 2, 0, 0),
    ('Window washing', 'medium', 'Living Room', 2, 0, 1),
    -- day_group 1
    ('Furniture assembly', 'high', 'Bedroom', 1, 1, 0),
    ('Packing supplies inventory', 'low', 'Storage', 3, 1, 1),
    -- day_group 2
    ('Electronics setup', 'medium', 'Office', 2, 2, 0),
    ('Declutter garage', 'low', 'Garage', 1, 2, 1),
    -- day_group 3
    ('Box labeling', 'medium', 'Storage', 2, 3, 0),
    ('Curtain removal', 'high', 'Living Room', 2, 3, 1);

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
