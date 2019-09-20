-- +goose Up
ALTER TABLE full_sync_logs
    ADD COLUMN receipt_id INT;

ALTER TABLE full_sync_logs
    ADD CONSTRAINT receipts_fk
        FOREIGN KEY (receipt_id)
            REFERENCES full_sync_receipts (id)
            ON DELETE CASCADE;


-- +goose Down
ALTER TABLE full_sync_logs
    DROP CONSTRAINT receipts_fk;

ALTER TABLE full_sync_logs
    DROP COLUMN receipt_id;
