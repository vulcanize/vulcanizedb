CREATE TABLE maker.pit_file_ilk (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk           TEXT,
  what          TEXT,
  data          NUMERIC,
	tx_idx        INTEGER NOT NUll,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx)
);

CREATE TABLE maker.pit_file_stability_fee (
  id SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  what      TEXT,
  data      TEXT,
  tx_idx    INTEGER NOT NULL,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx)
)