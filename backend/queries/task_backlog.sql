-- name: GetTaskBacklog :many
SELECT id, title, priority, people_needed, room, status
FROM backlog_tasks
ORDER BY sort_order;

-- name: GetTaskBacklogAssignments :many
SELECT task_id, person_id
FROM backlog_task_assignments
ORDER BY task_id, sort_order;

-- name: GetTaskByID :one
SELECT id, title, priority, people_needed, room, status
FROM backlog_tasks
WHERE id = $1;

-- name: GetTaskAssignments :many
SELECT task_id, person_id
FROM backlog_task_assignments
WHERE task_id = $1
ORDER BY sort_order;

-- name: GetMaxSortOrder :one
SELECT COALESCE(MAX(sort_order), 0)::int AS max_sort_order
FROM backlog_tasks;

-- name: CreateTask :one
INSERT INTO backlog_tasks (id, title, priority, people_needed, room, status, sort_order)
VALUES ('task-' || nextval('backlog_tasks_id_seq'), $1, $2, $3, $4, $5, $6)
RETURNING id, title, priority, people_needed, room, status;

-- name: CreateTaskAssignment :exec
INSERT INTO backlog_task_assignments (task_id, person_id, sort_order)
VALUES ($1, $2, $3);

-- name: UpdateTask :one
UPDATE backlog_tasks
SET title = $2, priority = $3, people_needed = $4, room = $5, status = $6
WHERE id = $1
RETURNING id, title, priority, people_needed, room, status;

-- name: DeleteTaskAssignments :exec
DELETE FROM backlog_task_assignments
WHERE task_id = $1;

-- name: DeleteTask :exec
DELETE FROM backlog_tasks
WHERE id = $1;
