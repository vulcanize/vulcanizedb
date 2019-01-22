-- +goose Up
CREATE TABLE maker.flap_kick (
  id            SERIAL PRIMARY KEY,
  header_id     INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  bid_id        NUMERIC NOT NULL,
  lot           NUMERIC NOT NULL,
  bid           NUMERIC NOT NULL,
  gal           TEXT,
  "end"         TIMESTAMP WITH TIME ZONE,
  tx_idx        INTEGER NOT NULL,
  log_idx       INTEGER NOT NULL,
  raw_log       JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN flap_kick_checked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE maker.flap_kick;
ALTER TABLE public.checked_headers
  DROP COLUMN flap_kick_checked;