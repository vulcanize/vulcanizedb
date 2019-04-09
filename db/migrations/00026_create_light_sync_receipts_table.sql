-- +goose Up
CREATE TABLE light_sync_receipts(
  id                  SERIAL PRIMARY KEY,
  transaction_id      INTEGER NOT NULL REFERENCES light_sync_transactions(id) ON DELETE CASCADE,
  header_id           INTEGER NOT NULL REFERENCES headers(id) ON DELETE CASCADE,
  contract_address    VARCHAR(42),
  cumulative_gas_used NUMERIC,
  gas_used            NUMERIC,
  state_root          VARCHAR(66),
  status              INTEGER,
  tx_hash             VARCHAR(66),
  rlp                 BYTEA,
  UNIQUE(header_id, transaction_id)
);


-- +goose Down
DROP TABLE light_sync_receipts;
