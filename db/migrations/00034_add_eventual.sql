-- +goose Up
ALTER TABLE eth.state_cids
ADD COLUMN eventual bool NOT NULL DEFAULT false;

ALTER TABLE eth.storage_cids
ADD COLUMN eventual bool NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE eth.storage_cids
DROP COLUMN eventual;

ALTER TABLE eth.state_cids
DROP COLUMN eventual;