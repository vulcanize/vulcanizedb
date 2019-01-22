-- +goose Up
ALTER TABLE blocks
  ADD COLUMN node_id INTEGER NOT NULL,
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id)
REFERENCES nodes (id);

-- +goose Down
ALTER TABLE blocks
  DROP COLUMN node_id;