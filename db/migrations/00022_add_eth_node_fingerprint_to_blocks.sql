-- +goose Up
ALTER TABLE blocks
  ADD COLUMN eth_node_fingerprint VARCHAR(128);

UPDATE blocks
  SET eth_node_fingerprint = (
    SELECT eth_node_id FROM eth_nodes WHERE eth_nodes.id = blocks.eth_node_id
  );

ALTER TABLE blocks
  ALTER COLUMN eth_node_fingerprint SET NOT NULL;


-- +goose Down
ALTER TABLE blocks
  DROP COLUMN eth_node_fingerprint;
