-- +goose Up
-- the header_id is the id for the header the state_cache row belongs to, not necessarily the header_id of the referenced state_cid
CREATE TABLE eth.state_cache (
  id               SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES eth.header_cids (id) ON DELETE CASCADE,
  state_id         INTEGER NOT NULL REFERENCES eth.state_cids (id) ON DELETE CASCADE,
  state_path       BYTEA,
  mh_key           TEXT NOT NULL REFERENCES public.blocks (key) ON DELETE CASCADE,
  UNIQUE (header_id, state_path)
);

-- +goose Down
DROP TABLE eth.state_cache;