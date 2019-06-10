-- +goose Up
ALTER TABLE full_sync_receipts
  ADD COLUMN block_id INT;

UPDATE full_sync_receipts
  SET block_id = (
    SELECT block_id FROM full_sync_transactions WHERE full_sync_transactions.id = full_sync_receipts.transaction_id
  );

ALTER TABLE full_sync_receipts
  ALTER COLUMN block_id SET NOT NULL;

ALTER TABLE full_sync_receipts
  ADD CONSTRAINT eth_blocks_fk
FOREIGN KEY (block_id)
REFERENCES eth_blocks (id)
ON DELETE CASCADE;

ALTER TABLE full_sync_receipts
  DROP COLUMN transaction_id;


-- +goose Down
ALTER TABLE full_sync_receipts
  ADD COLUMN transaction_id INT;

CREATE INDEX transaction_id_index ON full_sync_receipts (transaction_id);

UPDATE full_sync_receipts
  SET transaction_id = (
    SELECT id FROM full_sync_transactions WHERE full_sync_transactions.hash = full_sync_receipts.tx_hash
  );

ALTER TABLE full_sync_receipts
  ALTER COLUMN transaction_id SET NOT NULL;

ALTER TABLE full_sync_receipts
  ADD CONSTRAINT transaction_fk
FOREIGN KEY (transaction_id)
REFERENCES full_sync_transactions (id)
ON DELETE CASCADE;

ALTER TABLE full_sync_receipts
  DROP COLUMN block_id;
