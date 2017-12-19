BEGIN;

ALTER TABLE transactions
  DROP CONSTRAINT fk_test;

ALTER TABLE transactions
  ADD CONSTRAINT blocks_fk
FOREIGN KEY (block_id)
REFERENCES blocks (id)
ON DELETE CASCADE;

COMMIT;
