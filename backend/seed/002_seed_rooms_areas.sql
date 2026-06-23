-- +goose Up
-- +goose StatementBegin
INSERT INTO rooms_areas (id, name, type) VALUES
    ('room-1', 'Kitchen', 'room'),
    ('room-2', 'Living Room', 'room'),
    ('room-3', 'Bedroom 1', 'room'),
    ('room-4', 'Bedroom 2', 'room'),
    ('room-5', 'Garage', 'area'),
    ('room-6', 'Bedroom', 'room'),
    ('room-7', 'Storage', 'area'),
    ('room-8', 'Office', 'room');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM rooms_areas;
-- +goose StatementEnd
