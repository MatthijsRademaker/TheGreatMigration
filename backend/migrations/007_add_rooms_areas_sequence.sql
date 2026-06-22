-- +goose Up
-- +goose StatementBegin
-- Create a sequence for atomic room ID generation, replacing the
-- MAX+1 CTE approach which races under concurrent inserts.
CREATE SEQUENCE IF NOT EXISTS rooms_areas_id_seq;

-- Set the sequence to the maximum existing room number so the next
-- CreateRoom call produces the next sequential ID. When the table is empty
-- (schema applied before any seed), leave the sequence at its start so the
-- next nextval yields 1 — setval(seq, 0) is rejected (min value is 1).
SELECT setval(
  'rooms_areas_id_seq',
  GREATEST(COALESCE((SELECT MAX(CAST(SUBSTRING(id FROM 6) AS INTEGER)) FROM rooms_areas), 0), 1),
  (SELECT EXISTS (SELECT 1 FROM rooms_areas))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE IF EXISTS rooms_areas_id_seq;
-- +goose StatementEnd
