-- +goose Up
CREATE TABLE eth.state_cache (
  id               SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES eth.header_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  state_path       BYTEA,
  mh_key           TEXT NOT NULL REFERENCES  public.blocks (key) ON DELETE CASCADE,
  UNIQUE (header_id, state_path)
);

-- +goose Down
DROP TABLE eth.state_cache;