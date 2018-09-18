CREATE TABLE maker.tend (
  id            SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id           NUMERIC NOT NULL UNIQUE,
  lot              NUMERIC,
  bid              NUMERIC,
  guy              VARCHAR,
  tic              NUMERIC,
	tx_idx           INTEGER NOT NUll,
  raw_log          JSONB,
  UNIQUE (header_id, tx_idx)
);
