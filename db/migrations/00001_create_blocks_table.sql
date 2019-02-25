-- +goose Up
CREATE TABLE public.blocks
(
  block_number BIGINT
);


-- +goose Down
DROP TABLE public.blocks;
