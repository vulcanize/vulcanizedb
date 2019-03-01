-- +goose Up
ALTER TABLE nodes
  DROP CONSTRAINT node_uc;

ALTER TABLE nodes
  ADD CONSTRAINT node_uc UNIQUE (genesis_block, network_id, node_id);


-- +goose Down
ALTER TABLE nodes
  DROP CONSTRAINT node_uc;

ALTER TABLE nodes
  ADD CONSTRAINT node_uc UNIQUE (genesis_block, network_id);
