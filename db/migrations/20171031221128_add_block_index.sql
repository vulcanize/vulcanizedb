-- +goose Up
CREATE INDEX block_number_index ON blocks (block_number);


-- +goose Down
DROP INDEX block_number_index;