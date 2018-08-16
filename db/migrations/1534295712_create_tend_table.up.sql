CREATE TABLE maker.tend (
  db_id            SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  id               NUMERIC NOT NULL UNIQUE,
  lot              NUMERIC,
  bid              NUMERIC,
  guy              BYTEA,
  tic              NUMERIC,
	era              TIMESTAMP WITH TIME ZONE,
	tx_idx           INTEGER NOT NUll,
  raw_log          JSONB
);
