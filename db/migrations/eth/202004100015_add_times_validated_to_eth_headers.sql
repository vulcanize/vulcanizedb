-- +goose Up
ALTER TABLE eth.header_cids
ADD COLUMN times_validated INTEGER NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE eth.header_cids
DROP COLUMN times_validated;