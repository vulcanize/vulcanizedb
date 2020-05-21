-- +goose Up
-- eth.state_cids holds nodes emitted from diffing whereas this table holds nodes from iterating over the entire trie at a block
-- need to keep these separate as the diffed nodes represent the change in state between two blocks
-- whereas these nodes represent the entire state at a given block
CREATE TABLE eth.state_trie_cids (
  id                    SERIAL PRIMARY KEY,
  header_id             INTEGER NOT NULL REFERENCES eth.header_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  state_path            BYTEA,
  state_leaf_key        VARCHAR(66),
  node_type             INTEGER,
  cid                   TEXT NOT NULL,
  UNIQUE (header_id, state_path)
);

-- +goose Down
DROP TABLE eth.state_trie_cids;