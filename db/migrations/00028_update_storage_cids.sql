-- +goose Up
ALTER TABLE eth.storage_cids
ADD COLUMN storage_path BYTEA;

ALTER TABLE eth.storage_cids
DROP COLUMN leaf;

ALTER TABLE eth.storage_cids
ADD COLUMN node_type INTEGER NOT NULL;

ALTER TABLE eth.storage_cids
ALTER COLUMN storage_key DROP NOT NULL;

ALTER TABLE eth.storage_cids
DROP CONSTRAINT storage_cids_state_id_storage_key_key;

ALTER TABLE eth.storage_cids
ADD CONSTRAINT storage_cids_state_id_storage_path_key UNIQUE (state_id, storage_path);

-- +goose Down
ALTER TABLE eth.storage_cids
DROP CONSTRAINT storage_cids_state_id_storage_path_key;

ALTER TABLE eth.storage_cids
ADD CONSTRAINT storage_cids_state_id_storage_key_key UNIQUE (state_id, storage_key);

ALTER TABLE eth.storage_cids
ALTER COLUMN storage_key SET NOT NULL;

ALTER TABLE eth.storage_cids
DROP COLUMN node_type;

ALTER TABLE eth.storage_cids
ADD COLUMN leaf BOOLEAN NOT NULL;

ALTER TABLE eth.storage_cids
DROP COLUMN storage_path;