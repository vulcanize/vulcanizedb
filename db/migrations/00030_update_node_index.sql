-- +goose Up
BEGIN;

ALTER TABLE nodes
  DROP CONSTRAINT node_uc;

ALTER TABLE nodes
  ADD CONSTRAINT node_uc UNIQUE (genesis_block, network_id, node_id);

COMMIT;

-- +goose Down
-- +goose Up
BEGIN;

ALTER TABLE nodes
  DROP CONSTRAINT node_uc;

ALTER TABLE nodes
  ADD CONSTRAINT node_uc UNIQUE (genesis_block, network_id);

COMMIT;
