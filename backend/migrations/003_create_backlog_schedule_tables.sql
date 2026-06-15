-- +goose Up
-- +goose StatementBegin
CREATE TABLE backlog_tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    priority TEXT NOT NULL CHECK (priority IN ('high', 'medium', 'low')),
    people_needed INTEGER NOT NULL CHECK (people_needed >= 1),
    room TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('backlog', 'ready', 'assigned')),
    sort_order INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE backlog_task_assignments (
    task_id TEXT NOT NULL REFERENCES backlog_tasks(id),
    person_id TEXT NOT NULL REFERENCES people(id),
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (task_id, person_id)
);

CREATE TABLE schedule_task_cards (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    priority TEXT NOT NULL CHECK (priority IN ('high', 'medium', 'low')),
    room_area TEXT NOT NULL,
    people_needed INTEGER NOT NULL CHECK (people_needed >= 1),
    day_group INTEGER NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schedule_task_assignments (
    task_card_id INTEGER NOT NULL REFERENCES schedule_task_cards(id),
    person_id TEXT NOT NULL REFERENCES people(id),
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (task_card_id, person_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schedule_task_assignments;
DROP TABLE IF EXISTS schedule_task_cards;
DROP TABLE IF EXISTS backlog_task_assignments;
DROP TABLE IF EXISTS backlog_tasks;
-- +goose StatementEnd
