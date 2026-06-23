-- name: GetTaskBacklog :many
SELECT bt.id, bt.title, bt.priority, bt.people_needed, bt.area_id, ra.name AS area_name, bt.status
FROM backlog_tasks bt
JOIN rooms_areas ra ON ra.id = bt.area_id
ORDER BY bt.sort_order;

-- name: GetTaskBacklogAssignments :many
SELECT task_id, person_id
FROM backlog_task_assignments
ORDER BY task_id, sort_order;

-- name: GetTaskByID :one
SELECT bt.id, bt.title, bt.priority, bt.people_needed, bt.area_id, ra.name AS area_name, bt.status
FROM backlog_tasks bt
JOIN rooms_areas ra ON ra.id = bt.area_id
WHERE bt.id = $1;

-- name: GetTaskAssignments :many
SELECT task_id, person_id
FROM backlog_task_assignments
WHERE task_id = $1
ORDER BY sort_order;

-- name: GetMaxSortOrder :one
SELECT COALESCE(MAX(sort_order), 0)::int AS max_sort_order
FROM backlog_tasks;

-- name: CreateTask :one
WITH inserted AS (
    INSERT INTO backlog_tasks (id, title, priority, people_needed, area_id, status, sort_order)
    VALUES ('task-' || nextval('backlog_tasks_id_seq'), $1, $2, $3, $4, $5, $6)
    RETURNING id, title, priority, people_needed, area_id, status
)
SELECT i.id, i.title, i.priority, i.people_needed, i.area_id, ra.name AS area_name, i.status
FROM inserted i
JOIN rooms_areas ra ON ra.id = i.area_id;

-- name: CreateTaskAssignment :exec
INSERT INTO backlog_task_assignments (task_id, person_id, sort_order)
VALUES ($1, $2, $3);

-- name: UpdateTask :one
WITH updated AS (
    UPDATE backlog_tasks
    SET title = $2, priority = $3, people_needed = $4, area_id = $5, status = $6
    WHERE backlog_tasks.id = $1
    RETURNING id, title, priority, people_needed, area_id, status
)
SELECT u.id, u.title, u.priority, u.people_needed, u.area_id, ra.name AS area_name, u.status
FROM updated u
JOIN rooms_areas ra ON ra.id = u.area_id;

-- name: DeleteTaskAssignments :exec
DELETE FROM backlog_task_assignments
WHERE task_id = $1;

-- name: DeleteTask :exec
DELETE FROM backlog_tasks
WHERE id = $1;
