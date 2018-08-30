CREATE TABLE maker.bite (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk       bytea,
  lad       bytea,
  ink       VARCHAR,
  art       VARCHAR,
  iArt      VARCHAR,
  tab       NUMERIC,
  flip      VARCHAR,
  tx_idx    INTEGER NOT NUll,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx)
)