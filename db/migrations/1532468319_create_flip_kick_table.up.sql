CREATE TABLE maker.flip_kick (
  db_id            SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  id               NUMERIC NOT NULL UNIQUE,
  mom              VARCHAR,
  vat              VARCHAR,
  ilk              VARCHAR,
  lot              NUMERIC,
  bid              NUMERIC,
  guy              VARCHAR,
  gal              VARCHAR,
	"end"            TIMESTAMP WITH TIME ZONE,
	era              TIMESTAMP WITH TIME ZONE,
  lad              VARCHAR,
  tab              NUMERIC
);
