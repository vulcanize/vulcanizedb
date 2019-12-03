-- +goose Up
CREATE TABLE public.headers
(
    id                   SERIAL PRIMARY KEY,
    hash                 VARCHAR(66),
    block_number         BIGINT,
    raw                  JSONB,
    block_timestamp      NUMERIC,
    check_count          INTEGER NOT NULL DEFAULT 0,
    eth_node_id          INTEGER NOT NULL REFERENCES eth_nodes (id) ON DELETE CASCADE,
    UNIQUE (block_number, hash, eth_node_id)
);

-- Index is removed when table is
CREATE INDEX headers_block_number ON public.headers (block_number);


-- +goose Down
DROP TABLE public.headers;
