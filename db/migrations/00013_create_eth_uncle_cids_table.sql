-- +goose Up
CREATE TABLE eth.uncle_cids (
  id                    SERIAL PRIMARY KEY,
  header_id             INTEGER NOT NULL REFERENCES eth.header_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  block_hash            VARCHAR(66) NOT NULL,
  parent_hash           VARCHAR(66) NOT NULL,
  cid                   TEXT NOT NULL,
  UNIQUE (header_id, block_hash)
);

-- +goose Down
DROP TABLE eth.uncle_cids;