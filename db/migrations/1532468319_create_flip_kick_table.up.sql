CREATE TABLE maker.flip_kick (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id        NUMERIC NOT NULL,
  lot           NUMERIC,
  bid           NUMERIC,
  gal           VARCHAR,
	"end"         TIMESTAMP WITH TIME ZONE,
  urn           VARCHAR,
  tab           NUMERIC,
  tx_idx        INTEGER NOT NUll,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx)
);
