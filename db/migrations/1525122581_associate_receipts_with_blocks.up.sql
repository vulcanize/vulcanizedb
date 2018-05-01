BEGIN;

ALTER TABLE receipts
  ADD COLUMN block_id INT;

UPDATE receipts  
  SET block_id = (
    SELECT block_id FROM transactions WHERE transactions.id = receipts.transaction_id
  );

ALTER TABLE receipts
  ALTER COLUMN block_id SET NOT NULL;

ALTER TABLE receipts
  ADD CONSTRAINT blocks_fk
FOREIGN KEY (block_id)
REFERENCES blocks (id)
ON DELETE CASCADE;

ALTER TABLE receipts
  DROP COLUMN transaction_id;

COMMIT;
