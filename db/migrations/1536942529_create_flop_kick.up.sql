CREATE TABLE maker.flop_kick (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id        NUMERIC NOT NULL,
  lot           NUMERIC NOT NULL,
  bid           NUMERIC NOT NULL,
  gal           VARCHAR,
  "end"         TIMESTAMP WITH TIME ZONE,
  tx_idx        INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN flop_kick_checked BOOLEAN NOT NULL DEFAULT FALSE;