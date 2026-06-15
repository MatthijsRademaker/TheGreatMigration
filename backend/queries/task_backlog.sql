-- name: GetTaskBacklog :many
SELECT id, title, priority, people_needed, room, status
FROM backlog_tasks
ORDER BY sort_order;

-- name: GetTaskBacklogAssignments :many
SELECT task_id, person_id
FROM backlog_task_assignments
ORDER BY task_id, sort_order;
