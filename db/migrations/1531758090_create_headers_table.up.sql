CREATE TABLE public.headers (
  id                    SERIAL PRIMARY KEY,
  hash                  VARCHAR(66),
  block_number          BIGINT,
  raw                   JSONB,
  block_timestamp       NUMERIC,
  eth_node_id           INTEGER,
  eth_node_fingerprint  VARCHAR(128),
  CONSTRAINT eth_nodes_fk FOREIGN KEY (eth_node_id)
  REFERENCES eth_nodes (id)
  ON DELETE CASCADE
);