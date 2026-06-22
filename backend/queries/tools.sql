-- name: GetTools :many
SELECT id, name, brought_by, sort_order
FROM tools
ORDER BY sort_order;

-- name: GetToolByID :one
SELECT id, name, brought_by, sort_order
FROM tools
WHERE id = $1;

-- name: GetMaxToolSortOrder :one
SELECT COALESCE(MAX(sort_order), 0)::int AS max_sort_order
FROM tools;

-- name: CreateTool :one
INSERT INTO tools (id, name, sort_order)
VALUES ('tool-' || nextval('tools_id_seq'), $1, $2)
RETURNING id, name, brought_by, sort_order;

-- name: UpdateTool :one
UPDATE tools
SET name = $2, sort_order = $3
WHERE id = $1
RETURNING id, name, brought_by, sort_order;

-- name: DeleteTool :one
DELETE FROM tools
WHERE id = $1
RETURNING id;

-- name: SetToolBringer :one
UPDATE tools
SET brought_by = $2
WHERE id = $1
RETURNING id, name, brought_by, sort_order;

-- name: ClearToolBringer :one
UPDATE tools
SET brought_by = NULL
WHERE id = $1
RETURNING id, name, brought_by, sort_order;
