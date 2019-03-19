-- +goose Up
CREATE TABLE light_sync_transactions (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers(id) ON DELETE CASCADE,
  hash      TEXT,
  raw       JSONB,
  tx_index  INTEGER,
  tx_from   TEXT,
  tx_to     TEXT,
  UNIQUE (header_id, hash)
);

-- +goose Down
DROP TABLE light_sync_transactions;
