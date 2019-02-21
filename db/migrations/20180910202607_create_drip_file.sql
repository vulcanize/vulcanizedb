-- +goose Up
CREATE TABLE maker.drip_file_ilk (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk       INTEGER NOT NULL REFERENCES maker.ilks (id),
  vow       TEXT,
  tax       NUMERIC,
  log_idx   INTEGER NOT NUll,
  tx_idx    INTEGER NOT NUll,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

CREATE TABLE maker.drip_file_repo (
  id SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  what      TEXT,
  data      NUMERIC,
  log_idx   INTEGER NOT NULL,
  tx_idx    INTEGER NOT NULL,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

CREATE TABLE maker.drip_file_vow (
  id SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  what      TEXT,
  data      TEXT,
  log_idx   INTEGER NOT NULL,
  tx_idx    INTEGER NOT NULL,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN drip_file_ilk_checked BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.checked_headers
  ADD COLUMN drip_file_repo_checked BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.checked_headers
  ADD COLUMN drip_file_vow_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.drip_file_ilk;
DROP TABLE maker.drip_file_repo;
DROP TABLE maker.drip_file_vow;

ALTER TABLE public.checked_headers
  DROP COLUMN drip_file_ilk_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN drip_file_repo_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN drip_file_vow_checked;
