CREATE TABLE maker.frob (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  tx_idx    INTEGER,
  ilk       bytea,
  lad       bytea,
  gem       NUMERIC,
  ink       NUMERIC,
  art       NUMERIC,
  era       NUMERIC,
  UNIQUE (header_id, tx_idx)
);