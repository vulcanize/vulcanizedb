CREATE TABLE maker.drip_file_ilk (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk       TEXT,
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