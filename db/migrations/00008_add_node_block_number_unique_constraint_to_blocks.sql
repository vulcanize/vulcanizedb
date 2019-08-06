-- +goose Up
ALTER TABLE blocks
    ADD CONSTRAINT node_id_block_number_uc UNIQUE (number, node_id);

-- +goose Down
ALTER TABLE blocks
    DROP CONSTRAINT node_id_block_number_uc;
