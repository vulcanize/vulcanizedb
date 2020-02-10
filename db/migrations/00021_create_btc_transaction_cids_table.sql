-- +goose Up
CREATE TABLE btc.transaction_cids (
  id           SERIAL PRIMARY KEY,
	header_id    INTEGER NOT NULL REFERENCES btc.header_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	index        INTEGER NOT NULL,
	tx_hash      VARCHAR(66) NOT NULL UNIQUE,
	cid          TEXT NOT NULL,
	segwit       BOOL NOT NULL,
	witness_hash VARCHAR(66)
);

-- +goose Down
DROP TABLE btc.transaction_cids;