-- +goose Up
CREATE TABLE header_sync_transactions (
  id          SERIAL PRIMARY KEY,
  header_id   INTEGER NOT NULL REFERENCES headers(id) ON DELETE CASCADE,
  hash        VARCHAR(66) UNIQUE NOT NULL,
  gas_limit   NUMERIC,
  gas_price   NUMERIC,
  input_data  BYTEA,
  nonce       NUMERIC,
  raw         BYTEA,
  tx_from     VARCHAR(44),
  tx_index    INTEGER,
  tx_to       VARCHAR(44),
  "value"     NUMERIC
);

CREATE INDEX header_sync_transactions_header
    ON header_sync_transactions (header_id);

-- +goose Down
DROP INDEX header_sync_transactions_header;
DROP TABLE header_sync_transactions;
