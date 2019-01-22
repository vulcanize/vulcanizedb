-- +goose Up
BEGIN;
ALTER TABLE transactions
  RENAME COLUMN tx_hash TO hash;

ALTER TABLE transactions
  RENAME COLUMN tx_nonce TO nonce;

ALTER TABLE transactions
  RENAME COLUMN tx_gaslimit TO gaslimit;

ALTER TABLE transactions
  RENAME COLUMN tx_gasprice TO gasprice;

ALTER TABLE transactions
  RENAME COLUMN tx_value TO value;

ALTER TABLE transactions
  RENAME COLUMN tx_input_data TO input_data;
COMMIT;


-- +goose Down
-- +goose Up
BEGIN;
ALTER TABLE transactions
  RENAME COLUMN hash TO tx_hash;

ALTER TABLE transactions
  RENAME COLUMN nonce TO tx_nonce;

ALTER TABLE transactions
  RENAME COLUMN gaslimit TO tx_gaslimit;

ALTER TABLE transactions
  RENAME COLUMN gasprice TO tx_gasprice;

ALTER TABLE transactions
  RENAME COLUMN value TO tx_value;

ALTER TABLE transactions
  RENAME COLUMN input_data TO tx_input_data;
COMMIT;
