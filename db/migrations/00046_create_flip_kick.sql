-- +goose Up
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

ALTER TABLE public.checked_headers
  ADD COLUMN flip_kick_checked BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP TABLE maker.flip_kick;

ALTER TABLE public.checked_headers
  DROP COLUMN flip_kick_checked;