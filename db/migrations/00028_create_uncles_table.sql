-- +goose Up
CREATE TABLE public.uncles (
  id                    SERIAL PRIMARY KEY,
  hash                  VARCHAR(66) NOT NULL,
  block_id              INTEGER NOT NULL REFERENCES eth_blocks (id) ON DELETE CASCADE,
  reward                NUMERIC NOT NULL,
  miner                 VARCHAR(42) NOT NULL,
  raw                   JSONB,
  block_timestamp       NUMERIC,
  eth_node_id           INTEGER NOT NULL REFERENCES eth_nodes (id) ON DELETE CASCADE,
  eth_node_fingerprint  VARCHAR(128),
  UNIQUE (block_id, hash)
);

-- +goose Down
DROP TABLE public.uncles;
