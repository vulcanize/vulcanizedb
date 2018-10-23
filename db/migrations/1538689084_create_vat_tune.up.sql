CREATE TABLE maker.vat_tune (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           TEXT,
  urn           TEXT,
  v             TEXT,
  w             TEXT,
  dink          NUMERIC,
  dart          NUMERIC,
  tx_idx        INTEGER NOT NULL,
  log_idx       INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN vat_tune_checked BOOLEAN NOT NULL DEFAULT FALSE;