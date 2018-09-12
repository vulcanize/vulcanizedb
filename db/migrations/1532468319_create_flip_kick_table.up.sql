CREATE TABLE maker.flip_kick (
  db_id            SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  id               NUMERIC NOT NULL UNIQUE,
  lot              NUMERIC,
  bid              NUMERIC,
  gal              VARCHAR,
	"end"            TIMESTAMP WITH TIME ZONE,
  urn              VARCHAR,
  tab              NUMERIC,
  raw_log          JSONB
);
