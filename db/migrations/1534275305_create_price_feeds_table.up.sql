CREATE TABLE maker.price_feeds (
  id                  SERIAL PRIMARY KEY,
  block_number        BIGINT  NOT NULL,
  header_id           INTEGER NOT NULL,
  medianizer_address  bytea,
  tx_idx              INTEGER NOT NULL,
  usd_value           NUMERIC,
  UNIQUE (header_id, medianizer_address, tx_idx),
  CONSTRAINT headers_fk FOREIGN KEY (header_id)
  REFERENCES headers (id)
  ON DELETE CASCADE
);