-- +goose Up
CREATE TABLE public.uncle_cids (
  id                    SERIAL PRIMARY KEY,
  header_id             INTEGER NOT NULL REFERENCES header_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  block_hash            VARCHAR(66) NOT NULL,
  parent_hash           VARCHAR(66) NOT NULL,
  cid                   TEXT NOT NULL,
  UNIQUE (header_id, block_hash)
);

-- +goose Down
DROP TABLE public.uncle_cids;