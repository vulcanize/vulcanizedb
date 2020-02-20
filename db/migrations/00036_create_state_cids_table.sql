-- +goose Up
CREATE TABLE public.state_cids (
  id                    SERIAL PRIMARY KEY,
  header_id             INTEGER NOT NULL REFERENCES header_cids (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  state_key             VARCHAR(66) NOT NULL,
  leaf                  BOOLEAN NOT NULL,
  cid                   TEXT NOT NULL,
  UNIQUE (header_id, state_key)
);

-- +goose Down
DROP TABLE public.state_cids;