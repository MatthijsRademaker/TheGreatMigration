-- name: GetDailyScheduleTaskCards :many
SELECT id, title, priority, room_area, people_needed, scheduled_date, sort_order, created_at, task_id, completed
FROM schedule_task_cards
WHERE scheduled_date >= sqlc.arg(start_date)::date
  AND scheduled_date < (sqlc.arg(start_date)::date + sqlc.arg(days)::int * interval '1 day')
ORDER BY scheduled_date, sort_order;

-- name: CreateScheduleCard :one
INSERT INTO schedule_task_cards (title, priority, room_area, people_needed, scheduled_date, sort_order, task_id)
VALUES (sqlc.arg(title), sqlc.arg(priority), sqlc.arg(room_area), sqlc.arg(people_needed), sqlc.arg(scheduled_date), sqlc.arg(sort_order), sqlc.arg(task_id))
RETURNING id, title, priority, room_area, people_needed, scheduled_date, sort_order, created_at, task_id, completed;

-- name: CreateScheduleAssignment :exec
INSERT INTO schedule_task_assignments (task_card_id, person_id, sort_order)
VALUES (sqlc.arg(task_card_id), sqlc.arg(person_id), sqlc.arg(sort_order));

-- name: DeleteScheduleAssignments :exec
DELETE FROM schedule_task_assignments
WHERE task_card_id = sqlc.arg(task_card_id);

-- name: UpdateScheduleCard :one
UPDATE schedule_task_cards
SET title = sqlc.arg(title),
    priority = sqlc.arg(priority),
    room_area = sqlc.arg(room_area),
    people_needed = sqlc.arg(people_needed),
    scheduled_date = sqlc.arg(scheduled_date),
    sort_order = sqlc.arg(sort_order),
    task_id = sqlc.arg(task_id)
WHERE id = sqlc.arg(id)
RETURNING id, title, priority, room_area, people_needed, scheduled_date, sort_order, created_at, task_id, completed;

-- name: DeleteScheduleCard :exec
DELETE FROM schedule_task_cards
WHERE id = sqlc.arg(id);

-- name: TaskExists :one
SELECT EXISTS (
  SELECT 1 FROM backlog_tasks WHERE id = sqlc.arg(id)
);

-- name: TaskHasScheduleCards :one
SELECT EXISTS (
  SELECT 1 FROM schedule_task_cards WHERE task_id = sqlc.arg(task_id)
);

-- name: GetTaskByIDForRef :one
SELECT id, title, priority, people_needed, room
FROM backlog_tasks
WHERE id = sqlc.arg(id);

-- name: GetScheduleCardByID :one
SELECT id, title, priority, room_area, people_needed, scheduled_date, sort_order, created_at, task_id, completed
FROM schedule_task_cards
WHERE id = sqlc.arg(id);

-- name: GetMaxScheduleSortOrder :one
SELECT COALESCE(MAX(sort_order), 0)::int AS max_sort_order
FROM schedule_task_cards;

-- name: SetScheduleCardCompleted :exec
UPDATE schedule_task_cards
SET completed = sqlc.arg(completed)
WHERE id = sqlc.arg(id);

-- name: GetDailyScheduleAssignments :many
SELECT task_card_id, person_id
FROM schedule_task_assignments
ORDER BY task_card_id, sort_order;

-- name: GetAvailableCountByDate :many
SELECT a.date, COUNT(*)::int AS available_count
FROM availability a
WHERE a.status = 'available'
  AND a.date >= sqlc.arg(start_date)::date
  AND a.date < (sqlc.arg(start_date)::date + sqlc.arg(days)::int * interval '1 day')
GROUP BY a.date
ORDER BY a.date;
