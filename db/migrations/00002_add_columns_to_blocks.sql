-- +goose Up
ALTER TABLE blocks
  ADD COLUMN block_gaslimit DOUBLE PRECISION,
  ADD COLUMN block_gasused DOUBLE PRECISION,
  ADD COLUMN block_time DOUBLE PRECISION;


-- +goose Down
ALTER TABLE blocks
  DROP COLUMN block_gaslimit,
  DROP COLUMN block_gasused,
  DROP COLUMN block_time;

