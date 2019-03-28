-- +goose Up
CREATE INDEX block_id_index ON transactions (block_id);

-- +goose Down
DROP INDEX block_id_index;
