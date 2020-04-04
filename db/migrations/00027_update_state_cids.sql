-- +goose Up
ALTER TABLE eth.state_cids
ADD COLUMN state_path BYTEA;

ALTER TABLE eth.state_cids
DROP COLUMN leaf;

ALTER TABLE eth.state_cids
ADD COLUMN node_type INTEGER;

ALTER TABLE eth.state_cids
ALTER COLUMN state_key DROP NOT NULL;

ALTER TABLE eth.state_cids
DROP CONSTRAINT state_cids_header_id_state_key_key;

ALTER TABLE eth.state_cids
ADD CONSTRAINT state_cids_header_id_state_path_key UNIQUE (header_id, state_path);

-- +goose Down
ALTER TABLE eth.state_cids
ADD CONSTRAINT state_cids_header_id_state_key_key UNIQUE (header_id, state_key);

ALTER TABLE eth.state_cids
DROP CONSTRAINT state_cids_header_id_state_path_key;

ALTER TABLE eth.state_cids
ALTER COLUMN state_key SET NOT NULL;

ALTER TABLE eth.state_cids
DROP COLUMN node_type;

ALTER TABLE eth.state_cids
ADD COLUMN leaf BOOLEAN NOT NULL;

ALTER TABLE eth.state_cids
DROP COLUMN state_path;