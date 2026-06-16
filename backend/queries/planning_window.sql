-- name: GetPlanningWindow :one
SELECT id, start_date, end_date, created_at, updated_at
FROM planning_windows
LIMIT 1;

-- name: UpsertPlanningWindow :one
INSERT INTO planning_windows (id, start_date, end_date)
VALUES (1, $1, $2)
ON CONFLICT (id)
DO UPDATE SET start_date = $1, end_date = $2, updated_at = NOW()
RETURNING id, start_date, end_date, created_at, updated_at;
