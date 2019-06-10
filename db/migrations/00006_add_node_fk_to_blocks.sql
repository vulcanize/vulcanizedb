-- +goose Up
ALTER TABLE eth_blocks
  ADD COLUMN node_id INTEGER NOT NULL,
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id)
REFERENCES nodes (id)
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE eth_blocks
  DROP COLUMN node_id;
