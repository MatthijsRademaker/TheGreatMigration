-- +goose Up
-- +goose StatementBegin
-- Create a sequence for atomic person ID generation, matching the
-- room and task sequence patterns.
CREATE SEQUENCE IF NOT EXISTS people_id_seq;

-- Set the sequence to the maximum existing person number so the next
-- CreatePerson call produces the next sequential ID. When the table is empty
-- (schema applied before any seed), leave the sequence at its start so the
-- next nextval yields 1 — setval(seq, 0) is rejected (min value is 1).
SELECT setval(
  'people_id_seq',
  GREATEST(COALESCE((SELECT MAX(CAST(SUBSTRING(id FROM 2) AS INTEGER)) FROM people), 0), 1),
  (SELECT EXISTS (SELECT 1 FROM people))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE IF EXISTS people_id_seq;
-- +goose StatementEnd
