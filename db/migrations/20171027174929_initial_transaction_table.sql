-- +goose Up
CREATE TABLE transactions
(
  id SERIAL PRIMARY KEY,
  tx_hash VARCHAR(66),
  tx_nonce NUMERIC,
  tx_to varchar(66),
  tx_gaslimit NUMERIC,
  tx_gasprice NUMERIC,
  tx_value NUMERIC
)

-- +goose Down
DROP TABLE transactions