-- +goose Up
CREATE TABLE transactions (
  id          SERIAL PRIMARY KEY,
  block_id    INTEGER NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
  input_data  VARCHAR,
  tx_from     VARCHAR(66),
  gaslimit    NUMERIC,
  gasprice    NUMERIC,
  hash        VARCHAR(66),
  nonce       NUMERIC,
  tx_to       VARCHAR(66),
  "value"     NUMERIC
);

-- +goose Down
DROP TABLE transactions;