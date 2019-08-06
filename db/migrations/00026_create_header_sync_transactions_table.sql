-- +goose Up
CREATE TABLE header_sync_transactions
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

-- +goose Down
DROP TABLE header_sync_transactions;
