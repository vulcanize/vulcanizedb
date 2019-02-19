-- +goose Up
CREATE TABLE maker.pit_file_ilk (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           TEXT,
  what          TEXT,
  data          NUMERIC,
  log_idx       INTEGER NOT NUll,
  tx_idx        INTEGER NOT NUll,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

CREATE TABLE maker.pit_file_stability_fee (
  id SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  what      TEXT,
  data      TEXT,
  log_idx   INTEGER NOT NULL,
  tx_idx    INTEGER NOT NULL,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

CREATE TABLE maker.pit_file_debt_ceiling (
  id SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  what      TEXT,
  data      NUMERIC,
  log_idx   INTEGER NOT NULL,
  tx_idx    INTEGER NOT NULL,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN pit_file_debt_ceiling_checked BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.checked_headers
  ADD COLUMN pit_file_ilk_checked BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.checked_headers
  ADD COLUMN pit_file_stability_fee_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.pit_file_ilk;
DROP TABLE maker.pit_file_stability_fee;
DROP TABLE maker.pit_file_debt_ceiling;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_debt_ceiling_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_ilk_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_stability_fee_checked;
