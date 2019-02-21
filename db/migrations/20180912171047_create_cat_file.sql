-- +goose Up
CREATE TABLE maker.cat_file_chop_lump (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id),
  what          TEXT,
  data          NUMERIC,
  tx_idx        INTEGER NOT NUll,
  log_idx       INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

CREATE TABLE maker.cat_file_flip (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           TEXT,
  what          TEXT,
  flip          TEXT,
  tx_idx        INTEGER NOT NUll,
  log_idx       INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

CREATE TABLE maker.cat_file_pit_vow (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  what          TEXT,
  data          TEXT,
  tx_idx        INTEGER NOT NUll,
  log_idx       INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN cat_file_chop_lump_checked BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.checked_headers
  ADD COLUMN cat_file_flip_checked BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.checked_headers
  ADD COLUMN cat_file_pit_vow_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.cat_file_chop_lump;
DROP TABLE maker.cat_file_flip;
DROP TABLE maker.cat_file_pit_vow;

ALTER TABLE public.checked_headers
  DROP COLUMN cat_file_chop_lump_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN cat_file_flip_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN cat_file_pit_vow_checked;
