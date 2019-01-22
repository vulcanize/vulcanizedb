-- +goose Up
ALTER TABLE blocks
    ADD COLUMN block_extra_data VARCHAR;

-- +goose Down
ALTER TABLE blocks
    DROP COLUMN block_extra_data;