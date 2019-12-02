-- +goose Up
CREATE INDEX node_id_index ON eth_blocks (node_id);

-- +goose Down
DROP INDEX node_id_index;
