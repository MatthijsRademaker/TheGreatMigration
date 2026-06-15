-- name: GetAllPeople :many
SELECT id, name, initials, created_at
FROM people
ORDER BY id;

-- name: GetAvailabilityByDateRange :many
SELECT a.id, a.person_id, a.date, a.status, a.created_at
FROM availability a
WHERE a.date >= sqlc.arg(start_date)::date
  AND a.date <  (sqlc.arg(start_date)::date + sqlc.arg(days)::int * interval '1 day');
