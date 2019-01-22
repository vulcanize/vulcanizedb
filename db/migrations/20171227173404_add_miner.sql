-- +goose Up
ALTER TABLE blocks
    ADD COLUMN block_miner VARCHAR(42);

-- +goose Down
ALTER TABLE blocks
    DROP COLUMN block_miner;