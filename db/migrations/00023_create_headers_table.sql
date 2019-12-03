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
    eth_node_fingerprint VARCHAR(128),
    UNIQUE (block_number, hash, eth_node_fingerprint)
);

CREATE INDEX headers_block_number
    ON public.headers (block_number);

CREATE INDEX headers_block_timestamp
    ON public.headers (block_timestamp);

-- +goose Down
DROP INDEX public.headers_block_number;
DROP INDEX public.headers_block_timestamp;

DROP TABLE public.headers;
