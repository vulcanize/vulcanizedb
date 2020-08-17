-- +goose Up
CREATE TABLE public.transactions
(
    id         SERIAL PRIMARY KEY,
    header_id  INTEGER            NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    hash       VARCHAR(66) UNIQUE NOT NULL,
    gas_limit  NUMERIC,
    gas_price  NUMERIC,
    input_data BYTEA,
    nonce      NUMERIC,
    raw        BYTEA,
    tx_from    VARCHAR(44),
    tx_index   INTEGER,
    tx_to      VARCHAR(44),
    "value"    NUMERIC,
    created    TIMESTAMP          NOT NULL DEFAULT NOW(),
    updated    TIMESTAMP          NOT NULL DEFAULT NOW()
);

-- +goose StatementBegin
CREATE FUNCTION set_transaction_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER transaction_updated
    BEFORE UPDATE
    ON public.transactions
    FOR EACH ROW
EXECUTE PROCEDURE set_transaction_updated();

CREATE INDEX transactions_header
    ON transactions (header_id);

-- +goose Down
DROP TRIGGER transaction_updated ON public.transactions;
DROP FUNCTION set_transaction_updated();

DROP TABLE transactions;
