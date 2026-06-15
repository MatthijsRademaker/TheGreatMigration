-- +goose Up
-- +goose StatementBegin
CREATE TABLE planning_windows (
    id SERIAL PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE people (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    initials TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE availability (
    id SERIAL PRIMARY KEY,
    person_id TEXT NOT NULL REFERENCES people(id),
    date DATE NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('available','busy','partial','off')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(person_id, date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS availability;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS planning_windows;
-- +goose StatementEnd
