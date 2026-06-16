-- +goose Up
-- +goose StatementBegin
-- Create a sequence for atomic person ID generation, matching the
-- room and task sequence patterns.
CREATE SEQUENCE IF NOT EXISTS people_id_seq;

-- Set the sequence to the maximum existing person number so the next
-- CreatePerson call produces the next sequential ID.
SELECT setval('people_id_seq',
  COALESCE((SELECT MAX(CAST(SUBSTRING(id FROM 2) AS INTEGER)) FROM people), 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE IF EXISTS people_id_seq;
-- +goose StatementEnd
