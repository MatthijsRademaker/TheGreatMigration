-- name: GetDailyScheduleTaskCards :many
SELECT id, title, priority, room_area, people_needed, day_group, sort_order, created_at
FROM schedule_task_cards
ORDER BY day_group, sort_order;

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
