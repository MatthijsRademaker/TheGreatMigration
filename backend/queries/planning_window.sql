-- name: GetPlanningWindow :one
SELECT id, start_date, end_date, created_at, updated_at
FROM planning_windows
LIMIT 1;
