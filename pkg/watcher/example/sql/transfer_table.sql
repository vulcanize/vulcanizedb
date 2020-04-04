CREATE TABLE eth.token_transfers (
  id SERIAL PRIMARY KEY,
  receipt_id INTEGER NOT NULL REFERENCES eth.receipt_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  log_index INTEGER NOT NULL,
  contract_address VARCHAR(66) NOT NULL,
  src VARCHAR(66) NOT NULL,
  dst VARCHAR(66) NOT NULL,
  amount NUMERIC NOT NULL,
  UNIQUE (receipt_id, log_index)
);