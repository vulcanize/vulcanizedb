-- +goose Up
CREATE TYPE public.diff_status AS ENUM (
    'new',
    'transformed',
    'unrecognized',
    'noncanonical',
    'unwatched'
    );

CREATE TABLE public.storage_diff
(
    id             BIGSERIAL PRIMARY KEY,
    block_height   BIGINT,
    block_hash     BYTEA,
    hashed_address BYTEA,
    storage_key    BYTEA,
    storage_value  BYTEA,
    eth_node_id    INTEGER     NOT NULL REFERENCES public.eth_nodes (id) ON DELETE CASCADE,
    status         diff_status NOT NULL DEFAULT 'new',
    from_backfill  BOOLEAN     NOT NULL DEFAULT FALSE,
    UNIQUE (block_height, block_hash, hashed_address, storage_key, storage_value)
);

CREATE INDEX storage_diff_new_status_index
    ON public.storage_diff (status) WHERE status = 'new';
CREATE INDEX storage_diff_eth_node
    ON public.storage_diff (eth_node_id);

-- +goose Down
DROP TYPE public.diff_status CASCADE;
DROP TABLE public.storage_diff;
