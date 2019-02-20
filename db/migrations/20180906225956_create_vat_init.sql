-- +goose Up
CREATE TABLE maker.vat_init (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id),
  log_idx       INTEGER NOT NUll,
  tx_idx        INTEGER NOT NUll,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN vat_init_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.vat_init;

ALTER TABLE public.checked_headers
  DROP COLUMN vat_init_checked;