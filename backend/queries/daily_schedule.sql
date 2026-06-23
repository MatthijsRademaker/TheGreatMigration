-- name: GetDailyScheduleTaskCards :many
SELECT sc.id, sc.title, sc.priority, sc.area_id, ra.name AS area_name, sc.people_needed, sc.scheduled_date, sc.sort_order, sc.created_at, sc.task_id, sc.completed
FROM schedule_task_cards sc
JOIN rooms_areas ra ON ra.id = sc.area_id
WHERE sc.scheduled_date >= sqlc.arg(start_date)::date
  AND sc.scheduled_date < (sqlc.arg(start_date)::date + sqlc.arg(days)::int * interval '1 day')
ORDER BY sc.scheduled_date, sc.sort_order;

-- name: CreateScheduleCard :one
WITH inserted AS (
    INSERT INTO schedule_task_cards (title, priority, area_id, people_needed, scheduled_date, sort_order, task_id)
    VALUES (sqlc.arg(title), sqlc.arg(priority), sqlc.arg(area_id), sqlc.arg(people_needed), sqlc.arg(scheduled_date), sqlc.arg(sort_order), sqlc.arg(task_id))
    RETURNING id, title, priority, area_id, people_needed, scheduled_date, sort_order, created_at, task_id, completed
)
SELECT i.id, i.title, i.priority, i.area_id, ra.name AS area_name, i.people_needed, i.scheduled_date, i.sort_order, i.created_at, i.task_id, i.completed
FROM inserted i
JOIN rooms_areas ra ON ra.id = i.area_id;

-- name: CreateScheduleAssignment :exec
INSERT INTO schedule_task_assignments (task_card_id, person_id, sort_order)
VALUES (sqlc.arg(task_card_id), sqlc.arg(person_id), sqlc.arg(sort_order));

-- name: DeleteScheduleAssignments :exec
DELETE FROM schedule_task_assignments
WHERE task_card_id = sqlc.arg(task_card_id);

-- name: UpdateScheduleCard :one
WITH updated AS (
    UPDATE schedule_task_cards
    SET title = sqlc.arg(title),
        priority = sqlc.arg(priority),
        area_id = sqlc.arg(area_id),
        people_needed = sqlc.arg(people_needed),
        scheduled_date = sqlc.arg(scheduled_date),
        sort_order = sqlc.arg(sort_order),
        task_id = sqlc.arg(task_id)
    WHERE schedule_task_cards.id = sqlc.arg(id)
    RETURNING id, title, priority, area_id, people_needed, scheduled_date, sort_order, created_at, task_id, completed
)
SELECT u.id, u.title, u.priority, u.area_id, ra.name AS area_name, u.people_needed, u.scheduled_date, u.sort_order, u.created_at, u.task_id, u.completed
FROM updated u
JOIN rooms_areas ra ON ra.id = u.area_id;

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
SELECT bt.id, bt.title, bt.priority, bt.people_needed, bt.area_id, ra.name AS area_name
FROM backlog_tasks bt
JOIN rooms_areas ra ON ra.id = bt.area_id
WHERE bt.id = sqlc.arg(id);

-- name: GetScheduleCardByID :one
SELECT sc.id, sc.title, sc.priority, sc.area_id, ra.name AS area_name, sc.people_needed, sc.scheduled_date, sc.sort_order, sc.created_at, sc.task_id, sc.completed
FROM schedule_task_cards sc
JOIN rooms_areas ra ON ra.id = sc.area_id
WHERE sc.id = sqlc.arg(id);

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
