-- +goose Up
CREATE TABLE public.queued_storage (
  id            SERIAL PRIMARY KEY,
  block_height  BIGINT,
  block_hash    BYTEA,
  contract      BYTEA,
  storage_key   BYTEA,
  storage_value BYTEA,
  UNIQUE (block_height, block_hash, contract, storage_key, storage_value)
);

-- +goose Down
DROP TABLE public.queued_storage;
