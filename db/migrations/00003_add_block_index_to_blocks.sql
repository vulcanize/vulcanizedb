-- +goose Up
CREATE INDEX number_index ON eth_blocks (number);


-- +goose Down
DROP INDEX number_index;
