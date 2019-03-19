-- +goose Up
ALTER TABLE receipts
  ADD COLUMN block_id INT;

UPDATE receipts  
  SET block_id = (
    SELECT block_id FROM full_sync_transactions WHERE full_sync_transactions.id = receipts.transaction_id
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


-- +goose Down
ALTER TABLE receipts
  ADD COLUMN transaction_id INT;

CREATE INDEX transaction_id_index ON receipts (transaction_id);

UPDATE receipts
  SET transaction_id = (
    SELECT id FROM full_sync_transactions WHERE full_sync_transactions.hash = receipts.tx_hash
  );

ALTER TABLE receipts
  ALTER COLUMN transaction_id SET NOT NULL;

ALTER TABLE receipts
  ADD CONSTRAINT transaction_fk
FOREIGN KEY (transaction_id)
REFERENCES full_sync_transactions (id)
ON DELETE CASCADE;

ALTER TABLE receipts
  DROP COLUMN block_id;
