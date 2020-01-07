-- +goose Up
CREATE TABLE public.headers
(
    id                   SERIAL PRIMARY KEY,
    hash                 VARCHAR(66) NOT NULL,
    block_number         BIGINT NOT NULL,
    raw                  JSONB,
    block_timestamp      NUMERIC,
    check_count          INTEGER NOT NULL DEFAULT 0,
    eth_node_id          INTEGER NOT NULL REFERENCES eth_nodes (id) ON DELETE CASCADE,
    UNIQUE (block_number, hash, eth_node_id)
);

-- Index is removed when table is
CREATE INDEX headers_block_number ON public.headers (block_number);
CREATE INDEX headers_check_count ON public.headers (check_count);
CREATE INDEX headers_eth_node ON public.headers (eth_node_id);


-- +goose Down
DROP INDEX headers_block_number;
DROP INDEX headers_check_count;
DROP INDEX headers_eth_node;

DROP TABLE public.headers;
