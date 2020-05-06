-- +goose Up
CREATE TABLE btc.header_cids (
	id           SERIAL  PRIMARY KEY,
	block_number BIGINT NOT NULL,
	block_hash   VARCHAR(66) NOT NULL,
	parent_hash  VARCHAR(66) NOT NULL,
	cid          TEXT NOT NULL,
	timestamp    NUMERIC NOT NULL,
	bits         BIGINT NOT NULL,
	node_id      INTEGER NOT NULL REFERENCES nodes (id) ON DELETE CASCADE,
	UNIQUE (block_number, block_hash)
);

-- +goose Down
DROP TABLE btc.header_cids;