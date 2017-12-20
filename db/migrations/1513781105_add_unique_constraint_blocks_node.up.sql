ALTER TABLE blocks
  ADD CONSTRAINT node_id_block_number_uc UNIQUE (block_number, node_id);