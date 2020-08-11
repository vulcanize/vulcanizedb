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
    address        BYTEA,
    block_height   BIGINT,
    block_hash     BYTEA,
    storage_key    BYTEA,
    storage_value  BYTEA,
<<<<<<< HEAD
    eth_node_id    INTEGER     NOT NULL REFERENCES public.eth_nodes (id) ON DELETE CASCADE,
    status         diff_status NOT NULL DEFAULT 'new',
    from_backfill  BOOLEAN     NOT NULL DEFAULT FALSE,
    UNIQUE (block_height, block_hash, hashed_address, storage_key, storage_value)
=======
    eth_node_id    INTEGER NOT NULL REFERENCES public.eth_nodes (id) ON DELETE CASCADE,
    checked        BOOLEAN NOT NULL DEFAULT FALSE,
    from_backfill  BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (address, block_height, block_hash, storage_key, storage_value)
>>>>>>> fa657d32... Fixes diff repository and migrations
);

CREATE INDEX storage_diff_new_status_index
    ON public.storage_diff (status) WHERE status = 'new';
CREATE INDEX storage_diff_unrecognized_status_index
    ON public.storage_diff (status) WHERE status = 'unrecognized';
CREATE INDEX storage_diff_eth_node
    ON public.storage_diff (eth_node_id);

-- +goose Down
DROP TYPE public.diff_status CASCADE;
DROP TABLE public.storage_diff;
