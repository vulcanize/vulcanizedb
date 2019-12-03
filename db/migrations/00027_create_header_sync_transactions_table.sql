-- +goose Up
CREATE TABLE public.header_sync_transactions
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    hash       VARCHAR(66),
    gas_limit  NUMERIC,
    gas_price  NUMERIC,
    input_data BYTEA,
    nonce      NUMERIC,
    raw        BYTEA,
    tx_from    VARCHAR(44),
    tx_index   INTEGER,
    tx_to      VARCHAR(44),
    "value"    NUMERIC,
    UNIQUE (header_id, hash)
);

CREATE INDEX header_sync_transactions_header
    ON public.header_sync_transactions (header_id);

CREATE INDEX header_sync_transactions_tx_index
    ON public.header_sync_transactions (tx_index);

-- +goose Down
DROP INDEX public.header_sync_transactions_header;
DROP INDEX public.header_sync_transactions_tx_index;

DROP TABLE header_sync_transactions;
