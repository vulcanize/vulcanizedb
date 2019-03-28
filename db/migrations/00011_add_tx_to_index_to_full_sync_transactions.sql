-- +goose Up
CREATE INDEX tx_to_index ON full_sync_transactions(tx_to);

-- +goose Down
DROP INDEX tx_to_index;
