CREATE TABLE maker.pit_file (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           TEXT,
  what          TEXT,
  risk          NUMERIC,
	tx_idx        INTEGER NOT NUll,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx)
);