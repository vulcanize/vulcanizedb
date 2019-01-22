-- +goose Up
ALTER TABLE blocks
    ADD COLUMN is_final BOOLEAN;

-- +goose Down
ALTER TABLE blocks
    DROP COLUMN is_final;