-- +goose Up
ALTER TABLE blocks
  ADD COLUMN block_reward NUMERIC,
  ADD COLUMN block_uncles_reward NUMERIC;


-- +goose Down
ALTER TABLE blocks
  DROP COLUMN block_reward,
  DROP COLUMN block_uncles_reward;
