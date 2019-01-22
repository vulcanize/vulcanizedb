-- +goose Up
BEGIN;

ALTER TABLE blocks
  DROP CONSTRAINT node_fk;

ALTER TABLE blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id)
REFERENCES nodes (id)
ON DELETE CASCADE;

COMMIT;

-- +goose Down
BEGIN;

ALTER TABLE blocks
  DROP CONSTRAINT node_fk;

ALTER TABLE blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id)
REFERENCES nodes (id);

COMMIT;