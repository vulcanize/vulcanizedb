-- +goose Up
CREATE TABLE eth.queued_nodes (
  id                    SERIAL PRIMARY KEY,
  header_id             INTEGER NOT NULL REFERENCES eth.header_cids (id) ON DELETE CASCADE,
  rlp                   BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.queued_nodes;