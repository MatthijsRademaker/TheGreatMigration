-- +goose Up
-- +goose StatementBegin
-- Create a sequence for atomic room ID generation, replacing the
-- MAX+1 CTE approach which races under concurrent inserts.
CREATE SEQUENCE IF NOT EXISTS rooms_areas_id_seq;

-- Set the sequence to the maximum existing room number so the next
-- CreateRoom call produces the next sequential ID.
SELECT setval('rooms_areas_id_seq',
  COALESCE((SELECT MAX(CAST(SUBSTRING(id FROM 6) AS INTEGER)) FROM rooms_areas), 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE IF EXISTS rooms_areas_id_seq;
-- +goose StatementEnd
