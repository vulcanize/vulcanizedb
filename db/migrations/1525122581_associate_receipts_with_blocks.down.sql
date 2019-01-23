BEGIN;

ALTER TABLE receipts
  ADD COLUMN transaction_id INT;

CREATE INDEX transaction_id_index ON receipts (transaction_id);

UPDATE receipts
  SET transaction_id = (
    SELECT id FROM transactions WHERE transactions.hash = receipts.tx_hash
  );

ALTER TABLE receipts
  ALTER COLUMN transaction_id SET NOT NULL;

ALTER TABLE receipts
  ADD CONSTRAINT transaction_fk
FOREIGN KEY (transaction_id)
REFERENCES transactions (id)
ON DELETE CASCADE;

ALTER TABLE receipts
  DROP COLUMN block_id;

COMMIT;
