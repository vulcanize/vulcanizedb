-- +goose Up
CREATE TABLE btc.transaction_cids (
  id           SERIAL PRIMARY KEY,
	header_id    INTEGER NOT NULL REFERENCES btc.header_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	index        INTEGER NOT NULL,
	tx_hash      VARCHAR(66) NOT NULL,
	cid          TEXT NOT NULL,
	has_witness  BOOL NOT NULL,
	witness_hash VARCHAR(66),
	UNIQUE (header_id, tx_hash)
);

-- +goose Down
DROP TABLE btc.transaction_cids;