-- +goose Up
CREATE TABLE public.receipt_cids (
  id                    SERIAL PRIMARY KEY,
  tx_id                 INTEGER NOT NULL REFERENCES transaction_cids (id) ON DELETE CASCADE,
  cid                   TEXT NOT NULL
);

-- +goose Down
DROP TABLE public.receipt_cids;