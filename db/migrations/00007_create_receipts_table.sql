-- +goose Up
CREATE TABLE public.receipts
(
    id                  SERIAL PRIMARY KEY,
    transaction_id      INTEGER NOT NULL REFERENCES public.transactions (id) ON DELETE CASCADE,
    header_id           INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
    contract_address_id BIGINT  NOT NULL REFERENCES public.addresses (id) ON DELETE CASCADE,
    cumulative_gas_used NUMERIC,
    gas_used            NUMERIC,
    state_root          VARCHAR(66),
    status              INTEGER,
    tx_hash             VARCHAR(66),
    rlp                 BYTEA,
    created        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated        TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (header_id, transaction_id)
);
-- +goose StatementBegin
CREATE FUNCTION set_receipt_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER receipt_updated
    BEFORE UPDATE
    ON public.receipts
    FOR EACH ROW
EXECUTE PROCEDURE set_receipt_updated();

CREATE INDEX receipts_contract_address
    ON public.receipts (contract_address_id);
CREATE INDEX receipts_transaction
    ON public.receipts (transaction_id);

-- +goose Down
DROP TRIGGER receipt_updated ON public.receipts;
DROP FUNCTION set_receipt_updated();
DROP TABLE receipts;
