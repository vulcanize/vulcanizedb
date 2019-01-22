-- +goose Up
ALTER TABLE blocks ADD COLUMN id SERIAL PRIMARY KEY;

-- +goose Down
ALTER TABLE blocks DROP id;
