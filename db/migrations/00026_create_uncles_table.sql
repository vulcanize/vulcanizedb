-- +goose Up
CREATE TABLE public.uncles (
  id                    SERIAL PRIMARY KEY,
  hash                  VARCHAR(66) NOT NULL,
  block_id              INTEGER NOT NULL REFERENCES blocks (id) ON DELETE CASCADE,
  block_hash            VARCHAR(66) NOT NULL,
  reward                NUMERIC NOT NULL,
  miner                 VARCHAR(42) NOT NULL,
  raw                   JSONB,
  block_timestamp       NUMERIC,
  eth_node_id           INTEGER,
  eth_node_fingerprint  VARCHAR(128),
  CONSTRAINT eth_nodes_fk FOREIGN KEY (eth_node_id)
  REFERENCES eth_nodes (id)
  ON DELETE CASCADE,
  UNIQUE (block_id, hash)
);

-- +goose Down
DROP TABLE public.uncles;
