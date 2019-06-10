-- +goose Up
ALTER TABLE eth_blocks
  ADD COLUMN eth_node_fingerprint VARCHAR(128);

UPDATE eth_blocks
  SET eth_node_fingerprint = (
    SELECT eth_node_id FROM eth_nodes WHERE eth_nodes.id = eth_blocks.eth_node_id
  );

ALTER TABLE eth_blocks
  ALTER COLUMN eth_node_fingerprint SET NOT NULL;


-- +goose Down
ALTER TABLE eth_blocks
  DROP COLUMN eth_node_fingerprint;
