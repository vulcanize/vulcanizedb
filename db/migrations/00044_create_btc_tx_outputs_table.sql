-- +goose Up
CREATE TABLE btc.tx_outputs (
  id            SERIAL PRIMARY KEY,
	tx_id         INTEGER NOT NULL REFERENCES btc.transaction_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	index         INTEGER NOT NULL,
	value         INTEGER NOT NULL,
	pk_script     BYTEA NOT NULL,
	script_class  INTEGER NOT NULL,
	addresses     VARCHAR(66)[],
	required_sigs INTEGER NOT NULL,
	UNIQUE (tx_id, index)
);

-- +goose Down
DROP TABLE btc.tx_outputs;