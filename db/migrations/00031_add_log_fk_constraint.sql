-- +goose Up
BEGIN;
ALTER TABLE logs
    DROP CONSTRAINT log_uc;

ALTER TABLE logs
  ADD COLUMN receipt_id INT;

ALTER TABLE logs
  ADD CONSTRAINT receipts_fk
FOREIGN KEY (receipt_id)
REFERENCES receipts (id)
ON DELETE CASCADE;

COMMIT;

-- +goose Down
BEGIN;

ALTER TABLE logs
  DROP CONSTRAINT receipts_fk;

ALTER TABLE logs
  DROP COLUMN receipt_id;

ALTER TABLE logs
  ADD CONSTRAINT log_uc UNIQUE (block_number, index);

COMMIT;