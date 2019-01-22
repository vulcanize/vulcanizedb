-- +goose Up
ALTER TABLE nodes
  ADD COLUMN node_id VARCHAR(128),
  ADD COLUMN client_name VARCHAR;

-- +goose Down
ALTER TABLE nodes
  DROP COLUMN node_id,
  DROP COLUMN client_name;