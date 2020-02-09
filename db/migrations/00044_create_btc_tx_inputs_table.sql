-- +goose Up
CREATE TABLE btc.tx_inputs (
  id             SERIAL PRIMARY KEY,
	tx_id          INTEGER NOT NULL REFERENCES btc.transaction_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	index          INTEGER NOT NULL,
	witness        BYTEA[],
	sig_script     BYTEA NOT NULL,
	outpoint_id    INTEGER REFERENCES btc.tx_outputs (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	UNIQUE (tx_id, index)
);

-- +goose Down
DROP TABLE btc.tx_inputs;