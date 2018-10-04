CREATE TABLE maker.vat_fold (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk text,
  urn text,
  rate numeric,
  tx_idx        INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN vat_fold_checked BOOLEAN NOT NULL DEFAULT FALSE;
