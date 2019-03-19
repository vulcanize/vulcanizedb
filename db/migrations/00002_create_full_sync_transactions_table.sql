-- +goose Up
CREATE TABLE full_sync_transactions (
  id          SERIAL PRIMARY KEY,
  block_id    INTEGER NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
  gaslimit    NUMERIC,
  gasprice    NUMERIC,
  hash        VARCHAR(66),
  input_data  VARCHAR,
  nonce       NUMERIC,
  raw         BYTEA,
  tx_from     VARCHAR(66),
  tx_index    INTEGER,
  tx_to       VARCHAR(66),
  "value"     NUMERIC
);

-- +goose Down
DROP TABLE full_sync_transactions;