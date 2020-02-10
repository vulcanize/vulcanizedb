-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE public.event_logs
(
    id           SERIAL PRIMARY KEY,
    header_id    INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address      INTEGER NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    topics       BYTEA[],
    data         BYTEA,
    block_number BIGINT,
    block_hash   VARCHAR(66),
    tx_hash      VARCHAR(66) REFERENCES public.transactions (hash) ON DELETE CASCADE,
    tx_index     INTEGER,
    log_index    INTEGER,
    raw          JSONB,
    transformed  BOOL    NOT NULL DEFAULT FALSE,
    UNIQUE (header_id, tx_index, log_index)
);

CREATE INDEX event_logs_address
    ON event_logs (address);
CREATE INDEX event_logs_transaction
    ON event_logs (tx_hash);
CREATE INDEX event_logs_untransformed
    ON event_logs (transformed)
    WHERE transformed is false;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX event_logs_transaction;
DROP INDEX event_logs_address;
DROP INDEX event_logs_untransformed;
DROP TABLE event_logs;