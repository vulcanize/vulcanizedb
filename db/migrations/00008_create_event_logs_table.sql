-- +goose Up
CREATE TABLE public.event_logs
(
    id           BIGSERIAL PRIMARY KEY,
    header_id    INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    address      BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    topics       BYTEA[],
    data         BYTEA,
    block_number BIGINT,
    block_hash   VARCHAR(66),
    tx_hash      VARCHAR(66) REFERENCES public.transactions (hash) ON DELETE CASCADE,
    tx_index     INTEGER,
    log_index    INTEGER,
    raw          JSONB,
    transformed  BOOL      NOT NULL DEFAULT FALSE,
    created      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated      TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (header_id, tx_index, log_index)
);
-- +goose StatementBegin
CREATE FUNCTION set_event_log_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER event_log_updated
    BEFORE UPDATE
    ON public.event_logs
    FOR EACH ROW
EXECUTE PROCEDURE set_event_log_updated();

CREATE INDEX event_logs_address
    ON event_logs (address);
CREATE INDEX event_logs_transaction
    ON event_logs (tx_hash);
CREATE INDEX event_logs_untransformed
    ON event_logs (transformed)
    WHERE transformed = false;

-- +goose Down
DROP TRIGGER event_log_updated ON public.event_logs;
DROP FUNCTION set_event_log_updated();

DROP TABLE event_logs;