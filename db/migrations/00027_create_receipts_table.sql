-- +goose Up
CREATE TABLE receipts
(
  id                  SERIAL PRIMARY KEY,
  transaction_id      INTEGER NOT NULL,
  contract_address    VARCHAR(42),
  cumulative_gas_used NUMERIC,
  gas_used            NUMERIC,
  state_root          VARCHAR(66),
  status              INTEGER,
  tx_hash             VARCHAR(66),
  CONSTRAINT transaction_fk FOREIGN KEY (transaction_id)
  REFERENCES transactions (id)
  ON DELETE CASCADE
);




-- +goose Down
DROP TABLE receipts;

