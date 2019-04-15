-- +goose Up
CREATE TABLE public.storage_cids (
  id                    SERIAL PRIMARY KEY,
  state_id              INTEGER NOT NULL REFERENCES state_cids (id) ON DELETE CASCADE,
  storage_key           VARCHAR(66) NOT NULL,
  cid                   TEXT NOT NULL,
  UNIQUE (state_id, storage_key)
);

-- +goose Down
DROP TABLE public.storage_cids;