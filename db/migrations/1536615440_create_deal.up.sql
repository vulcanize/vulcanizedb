CREATE TABLE maker.deal (
  id                SERIAL PRIMARY KEY,
  header_id         INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id            NUMERIC NOT NULL,
  contract_address  VARCHAR,
  log_idx           INTEGER NOT NUll,
  tx_idx            INTEGER NOT NUll,
  raw_log           JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN deal_checked BOOLEAN NOT NULL DEFAULT FALSE;