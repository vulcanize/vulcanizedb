-- +goose Up
CREATE TABLE public.storage_diff
(
    id             BIGSERIAL PRIMARY KEY,
    block_height   BIGINT,
    block_hash     BYTEA,
    hashed_address BYTEA,
    storage_key    BYTEA,
    storage_value  BYTEA,
    UNIQUE (block_height, block_hash, hashed_address, storage_key, storage_value)
);

-- +goose Down
DROP TABLE public.storage_diff;