-- +goose Up
CREATE TABLE maker.vat_toll (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           TEXT,
  urn           TEXT,
  take          NUMERIC,
  tx_idx        INTEGER NOT NULL,
  log_idx        INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN vat_toll_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.vat_toll;
ALTER TABLE public.checked_headers
  DROP COLUMN vat_toll_checked;
