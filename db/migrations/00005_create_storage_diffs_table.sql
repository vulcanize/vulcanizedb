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
    eth_node_id    INTEGER     NOT NULL REFERENCES public.eth_nodes (id) ON DELETE CASCADE,
    status         diff_status NOT NULL DEFAULT 'new',
    from_backfill  BOOLEAN     NOT NULL DEFAULT FALSE,
    created        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated        TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (block_height, block_hash, address, storage_key, storage_value)
);

CREATE INDEX storage_diff_new_status_index
    ON public.storage_diff (status) WHERE status = 'new';
CREATE INDEX storage_diff_unrecognized_status_index
    ON public.storage_diff (status) WHERE status = 'unrecognized';
CREATE INDEX storage_diff_eth_node
    ON public.storage_diff (eth_node_id);

-- +goose StatementBegin
CREATE FUNCTION set_storage_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER storage_updated
    BEFORE UPDATE
    ON public.storage_diff
    FOR EACH ROW
EXECUTE PROCEDURE set_storage_updated();

-- +goose Down
DROP TYPE public.diff_status CASCADE;
DROP TRIGGER storage_updated ON public.storage_diff;
DROP FUNCTION set_storage_updated();

DROP TABLE public.storage_diff;
