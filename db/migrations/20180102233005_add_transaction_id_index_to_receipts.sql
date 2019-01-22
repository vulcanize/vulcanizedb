-- +goose Up
CREATE INDEX transaction_id_index ON receipts (transaction_id);

-- +goose Down
DROP INDEX transaction_id_index;
