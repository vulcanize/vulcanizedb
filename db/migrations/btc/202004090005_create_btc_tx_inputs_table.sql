-- +goose Up
CREATE TABLE btc.tx_inputs (
  id               SERIAL PRIMARY KEY,
	tx_id            INTEGER NOT NULL REFERENCES btc.transaction_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	index            INTEGER NOT NULL,
	witness          VARCHAR[],
	sig_script       BYTEA NOT NULL,
	outpoint_tx_hash VARCHAR(66) NOT NULL,
	outpoint_index   NUMERIC NOT NULL,
	UNIQUE (tx_id, index)
);

-- +goose Down
DROP TABLE btc.tx_inputs;