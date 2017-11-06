ALTER TABLE transactions
  ADD COLUMN block_id INTEGER NOT NULL,
  ADD CONSTRAINT fk_test
  FOREIGN KEY (block_id)
  REFERENCES blocks (id)
