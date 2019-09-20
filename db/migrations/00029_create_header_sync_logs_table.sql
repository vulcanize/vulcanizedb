-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE header_sync_logs
(
    id           SERIAL PRIMARY KEY,
    header_id    INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
    address      INTEGER NOT NULL REFERENCES addresses (id) ON DELETE CASCADE,
    topics       BYTEA[],
    data         BYTEA,
    block_number BIGINT,
    block_hash   VARCHAR(66),
    tx_hash      VARCHAR(66),
    tx_index     INTEGER,
    log_index    INTEGER,
    raw          JSONB,
    transformed  BOOL    NOT NULL DEFAULT FALSE,
    UNIQUE (header_id, tx_index, log_index)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE header_sync_logs;