CREATE TABLE maker.price_feeds (
  id                  SERIAL PRIMARY KEY,
  block_number        BIGINT  NOT NULL,
  header_id           INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  medianizer_address  bytea,
  usd_value           NUMERIC,
  tx_idx              INTEGER NOT NULL,
  raw_log             JSONB,
  UNIQUE (header_id, medianizer_address, tx_idx)
);