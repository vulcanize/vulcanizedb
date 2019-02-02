-- +goose Up
ALTER TABLE blocks
  DROP CONSTRAINT node_fk;

ALTER TABLE blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id)
REFERENCES nodes (id)
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE blocks
  DROP CONSTRAINT node_fk;

ALTER TABLE blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id)
REFERENCES nodes (id);
