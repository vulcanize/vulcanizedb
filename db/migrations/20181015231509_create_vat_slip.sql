-- +goose Up
CREATE TABLE maker.vat_slip (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id),
  guy           TEXT,
  rad           NUMERIC,
  tx_idx        INTEGER NOT NULL,
  log_idx       INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN vat_slip_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.vat_slip;
ALTER TABLE public.checked_headers
  DROP COLUMN vat_slip_checked;
