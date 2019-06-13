-- +goose Up
CREATE TABLE public.blocks (
  key TEXT UNIQUE NOT NULL,
  data BYTEA NOT NULL
);

-- +goose Down
DROP TABLE public.blocks;
