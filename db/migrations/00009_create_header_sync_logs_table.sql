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
    tx_hash      VARCHAR(66) REFERENCES header_sync_transactions (hash) ON DELETE CASCADE,
    tx_index     INTEGER,
    log_index    INTEGER,
    raw          JSONB,
    transformed  BOOL    NOT NULL DEFAULT FALSE,
    UNIQUE (header_id, tx_index, log_index)
);

CREATE INDEX header_sync_logs_address
    ON header_sync_logs (address);
CREATE INDEX header_sync_logs_transaction
    ON header_sync_logs (tx_hash);
CREATE INDEX header_sync_logs_untransformed
    ON header_sync_logs (transformed)
    WHERE transformed is false;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX header_sync_logs_transaction;
DROP INDEX header_sync_logs_address;
DROP INDEX header_sync_logs_untransformed;
DROP TABLE header_sync_logs;