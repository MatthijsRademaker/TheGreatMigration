-- +goose Up
-- +goose StatementBegin
INSERT INTO planning_windows (start_date, end_date) VALUES ('2026-07-05', '2026-08-13');

INSERT INTO people (id, name, initials) VALUES
    ('p1', 'Sophia Chen', 'SC'),
    ('p2', 'Marcus Rivera', 'MR'),
    ('p3', 'Elena Kowalski', 'EK'),
    ('p4', 'James Okafor', 'JO'),
    ('p5', 'Priya Nair', 'PN'),
    ('p6', 'Thomas Berg', 'TB'),
    ('p7', 'Amara Diallo', 'AD'),
    ('p8', 'Noah Larsson', 'NL');

-- Generate 320 availability rows (8 people x 40 dates).
-- p1-p6: always 'available'
-- p7: always 'busy'
-- p8: cycle through off/partial/busy/available
INSERT INTO availability (person_id, date, status)
SELECT
    p.id,
    d.date,
    CASE
        WHEN p.id IN ('p1','p2','p3','p4','p5','p6') THEN 'available'
        WHEN p.id = 'p7' THEN 'busy'
        WHEN p.id = 'p8' THEN
            CASE (d.date::date - '2026-07-05'::date) % 4
                WHEN 0 THEN 'off'
                WHEN 1 THEN 'partial'
                WHEN 2 THEN 'busy'
                WHEN 3 THEN 'available'
            END
    END AS status
FROM people p
CROSS JOIN generate_series('2026-07-05'::date, '2026-08-13'::date, '1 day'::interval) AS d(date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM availability;
DELETE FROM people;
DELETE FROM planning_windows;
-- +goose StatementEnd
