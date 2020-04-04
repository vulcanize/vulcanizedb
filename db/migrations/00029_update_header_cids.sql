-- +goose Up
ALTER TABLE eth.header_cids
ADD COLUMN state_root VARCHAR(66);

ALTER TABLE eth.header_cids
ADD COLUMN tx_root VARCHAR(66);

ALTER TABLE eth.header_cids
ADD COLUMN receipt_root VARCHAR(66);

ALTER TABLE eth.header_cids
ADD COLUMN uncle_root VARCHAR(66);

ALTER TABLE eth.header_cids
ADD COLUMN bloom BYTEA;

ALTER TABLE eth.header_cids
ADD COLUMN timestamp NUMERIC;

-- +goose Down
ALTER TABLE eth.header_cids
DROP COLUMN timestamp;

ALTER TABLE eth.header_cids
DROP COLUMN bloom;

ALTER TABLE eth.header_cids
DROP COLUMN uncle_root;

ALTER TABLE eth.header_cids
DROP COLUMN receipt_root;

ALTER TABLE eth.header_cids
DROP COLUMN tx_root;

ALTER TABLE eth.header_cids
DROP COLUMN state_root;