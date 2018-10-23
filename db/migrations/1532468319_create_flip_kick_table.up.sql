CREATE TABLE maker.flip_kick (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id        NUMERIC NOT NULL,
  lot           NUMERIC,
  bid           NUMERIC,
  gal           TEXT,
	"end"         TIMESTAMP WITH TIME ZONE,
  urn           TEXT,
  tab           NUMERIC,
  tx_idx        INTEGER NOT NUll,
  log_idx       INTEGER NOT NUll,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);
