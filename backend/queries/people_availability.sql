-- name: GetAllPeople :many
SELECT id, name, initials, created_at
FROM people
ORDER BY id;

-- name: GetAvailabilityByDateRange :many
SELECT a.id, a.person_id, a.date, a.status, a.created_at
FROM availability a
WHERE a.date >= sqlc.arg(start_date)::date
  AND a.date <  (sqlc.arg(start_date)::date + sqlc.arg(days)::int * interval '1 day');

-- name: GetPerson :one
SELECT id, name, initials, created_at
FROM people
WHERE id = $1;

-- name: PersonExists :one
SELECT EXISTS (
    SELECT 1 FROM people WHERE id = $1
) AS exists;

-- name: CreatePerson :one
INSERT INTO people (id, name, initials)
VALUES ('p' || nextval('people_id_seq'), $1, $2)
RETURNING id;

-- name: UpdatePerson :exec
UPDATE people
SET name = $2, initials = $3
WHERE id = $1;

-- name: DeletePerson :exec
DELETE FROM people
WHERE id = $1;

-- name: UpsertAvailability :exec
INSERT INTO availability (person_id, date, status)
VALUES ($1, $2, $3)
ON CONFLICT (person_id, date)
DO UPDATE SET status = EXCLUDED.status;

-- name: DeleteAvailability :exec
DELETE FROM availability
WHERE person_id = $1 AND date = $2;

-- name: CheckPersonBacklogReferences :one
SELECT EXISTS (
    SELECT 1 FROM backlog_task_assignments WHERE person_id = $1
) AS has_refs;

-- name: CheckPersonScheduleReferences :one
SELECT EXISTS (
    SELECT 1 FROM schedule_task_assignments WHERE person_id = $1
) AS has_refs;
