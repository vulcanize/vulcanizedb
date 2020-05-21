-- +goose Up
-- eth.storage_cids holds nodes emitted from diffing whereas this table holds nodes from iterating over the entire trie at a block
-- need to keep these separate as the diffed nodes represent the change in state between two blocks
-- whereas these nodes represent the entire state at a given block
CREATE TABLE eth.storage_trie_cids (
  id                    SERIAL PRIMARY KEY,
  state_id              INTEGER NOT NULL REFERENCES eth.state_trie_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  storage_path          BYTEA,
  storage_leaf_key      VARCHAR(66),
  node_type             INTEGER,
  cid                   TEXT NOT NULL,
  UNIQUE (state_id, storage_path)
);

-- +goose Down
DROP TABLE eth.storage_trie_cids;