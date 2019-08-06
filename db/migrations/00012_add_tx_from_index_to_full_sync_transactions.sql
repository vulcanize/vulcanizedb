-- +goose Up
CREATE INDEX tx_from_index ON full_sync_transactions (tx_from);

-- +goose Down
DROP INDEX tx_from_index;
