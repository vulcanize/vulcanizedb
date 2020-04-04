-- +goose Up
CREATE TABLE eth.state_accounts (
  id                    SERIAL PRIMARY KEY,
  state_id              INTEGER NOT NULL REFERENCES eth.state_cids (id) ON DELETE CASCADE,
  balance               NUMERIC NOT NULL,
  nonce                 INTEGER NOT NULL,
  code_hash             BYTEA NOT NULL,
  storage_root          VARCHAR(66) NOT NULL,
  UNIQUE (state_id)
);

-- +goose Down
DROP TABLE eth.state_accounts;