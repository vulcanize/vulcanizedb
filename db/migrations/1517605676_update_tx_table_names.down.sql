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
