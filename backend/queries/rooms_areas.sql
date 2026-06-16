-- name: ListRooms :many
SELECT id, name, type, created_at, updated_at
FROM rooms_areas
ORDER BY name;

-- name: GetRoomByID :one
SELECT id, name, type, created_at, updated_at
FROM rooms_areas
WHERE id = $1;

-- name: CreateRoom :one
INSERT INTO rooms_areas (id, name, type)
VALUES ('room-' || nextval('rooms_areas_id_seq'), $1, $2)
RETURNING id, name, type, created_at, updated_at;

-- name: UpdateRoom :one
UPDATE rooms_areas
SET name = $2, type = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, name, type, created_at, updated_at;

-- name: DeleteRoom :one
DELETE FROM rooms_areas
WHERE id = $1
RETURNING id;
