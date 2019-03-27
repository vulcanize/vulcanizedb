-- +goose Up
CREATE TABLE light_sync_transactions (
  id          SERIAL PRIMARY KEY,
  header_id   INTEGER NOT NULL REFERENCES headers(id) ON DELETE CASCADE,
  hash        TEXT,
  gaslimit    NUMERIC,
  gasprice    NUMERIC,
  input_data  BYTEA,
  nonce       NUMERIC,
  raw         BYTEA,
  tx_from     TEXT,
  tx_index    INTEGER,
  tx_to       TEXT,
  "value"     NUMERIC,
  UNIQUE (header_id, hash)
);

-- +goose Down
DROP TABLE light_sync_transactions;
