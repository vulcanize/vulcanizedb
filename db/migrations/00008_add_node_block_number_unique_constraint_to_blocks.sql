-- +goose Up
ALTER TABLE eth_blocks
  ADD CONSTRAINT node_id_block_number_uc UNIQUE (number, node_id);

-- +goose Down
ALTER TABLE eth_blocks
  DROP CONSTRAINT node_id_block_number_uc;
