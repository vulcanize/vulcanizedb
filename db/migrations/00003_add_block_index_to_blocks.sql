-- +goose Up
CREATE INDEX number_index ON blocks (number);


-- +goose Down
DROP INDEX number_index;
