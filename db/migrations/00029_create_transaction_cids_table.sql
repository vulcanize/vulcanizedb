-- +goose Up
CREATE TABLE public.transaction_cids (
  id                    SERIAL PRIMARY KEY,
  header_id             INTEGER NOT NULL REFERENCES header_cids (id) ON DELETE CASCADE,
  tx_hash               VARCHAR(66) NOT NULL,
  cid                   TEXT NOT NULL,
  dst                   VARCHAR(66) NOT NULL,
  src                   VARCHAR(66) NOT NULL,
  UNIQUE (header_id, tx_hash)
);

-- +goose Down
DROP TABLE public.transaction_cids;
