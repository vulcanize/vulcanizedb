-- +goose Up
CREATE TABLE public.header_cids (
  id                    SERIAL PRIMARY KEY,
  block_number          BIGINT NOT NULL,
  block_hash            VARCHAR(66) NOT NULL,
  cid                   TEXT NOT NULL,
  uncle                 BOOLEAN NOT NULL,
  UNIQUE (block_number, block_hash)
);

-- +goose Down
DROP TABLE public.header_cids;