-- +goose Up
ALTER TABLE eth.state_cids
RENAME COLUMN state_key TO state_leaf_key;

ALTER TABLE eth.storage_cids
RENAME COLUMN storage_key TO storage_leaf_key;

-- +goose Down
ALTER TABLE eth.storage_cids
RENAME COLUMN storage_leaf_key TO storage_key;

ALTER TABLE eth.state_cids
RENAME COLUMN state_leaf_key TO state_key;