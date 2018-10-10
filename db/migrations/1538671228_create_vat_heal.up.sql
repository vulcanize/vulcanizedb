CREATE TABLE maker.vat_heal (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  urn           varchar,
  v             varchar,
  rad           int,
  tx_idx        INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx)
);

ALTER TABLE public.checked_headers
    ADD COLUMN vat_heal_checked BOOLEAN NOT NULL DEFAULT FALSE;
