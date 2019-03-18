-- +goose Up
CREATE TABLE public.queued_storage (
  id            SERIAL PRIMARY KEY,
  block_height  BIGINT,
  block_hash    BYTEA,
  contract      BYTEA,
  storage_key   BYTEA,
  storage_value BYTEA
);

-- +goose Down
DROP TABLE public.queued_storage;
