-- +goose Up
CREATE INDEX node_id_index ON blocks (node_id);

-- +goose Down
DROP INDEX node_id_index;
