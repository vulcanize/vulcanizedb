BEGIN;

ALTER TABLE transactions
  DROP CONSTRAINT blocks_fk;

ALTER TABLE transactions
  ADD CONSTRAINT fk_test
FOREIGN KEY (block_id)
REFERENCES blocks (id);

COMMIT;