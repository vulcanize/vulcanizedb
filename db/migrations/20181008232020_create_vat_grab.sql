-- +goose Up
CREATE TABLE maker.vat_grab (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id),
  urn           TEXT,
  v             TEXT,
  w             TEXT,
  dink          NUMERIC,
  dart          NUMERIC,
  log_idx       INTEGER NOT NULL,
  tx_idx        INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN vat_grab_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.vat_grab;
ALTER TABLE public.checked_headers
  DROP COLUMN vat_grab_checked;
