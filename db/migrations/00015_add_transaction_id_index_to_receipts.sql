-- +goose Up
CREATE INDEX transaction_id_index ON full_sync_receipts (transaction_id);

-- +goose Down
DROP INDEX transaction_id_index;
