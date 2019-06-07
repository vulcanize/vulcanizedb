-- +goose Up
CREATE TABLE public.receipt_cids (
  id                    SERIAL PRIMARY KEY,
  tx_id                 INTEGER NOT NULL REFERENCES transaction_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  cid                   TEXT NOT NULL,
  topic0s               VARCHAR(66)[]
);

-- +goose Down
DROP TABLE public.receipt_cids;