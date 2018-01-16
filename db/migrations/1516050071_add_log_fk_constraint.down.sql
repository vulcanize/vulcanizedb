BEGIN;

ALTER TABLE logs
  DROP CONSTRAINT receipts_fk;

ALTER TABLE logs
  DROP COLUMN receipt_id;

ALTER TABLE logs
  ADD CONSTRAINT log_uc UNIQUE (block_number, index);

COMMIT;