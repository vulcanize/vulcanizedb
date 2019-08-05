-- +goose Up
CREATE TABLE IF NOT EXISTS public.blocks (
  key TEXT UNIQUE NOT NULL,
  data BYTEA NOT NULL
);

-- +goose Down
DROP TABLE public.blocks;
