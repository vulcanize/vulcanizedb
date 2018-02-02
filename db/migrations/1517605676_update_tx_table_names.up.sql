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
