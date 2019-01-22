-- +goose Up
BEGIN;
ALTER TABLE blocks
  ALTER COLUMN block_gaslimit TYPE BIGINT USING block_gaslimit :: BIGINT;

ALTER TABLE blocks
  ALTER COLUMN block_gasused TYPE BIGINT USING block_gasused :: BIGINT;

ALTER TABLE blocks
  ALTER COLUMN block_time TYPE BIGINT USING block_time :: BIGINT;

ALTER TABLE blocks
  ALTER COLUMN block_reward TYPE DOUBLE PRECISION USING block_time :: DOUBLE PRECISION;

ALTER TABLE blocks
  ALTER COLUMN block_uncles_reward TYPE DOUBLE PRECISION USING block_time :: DOUBLE PRECISION;

COMMIT;


-- +goose Down
-- +goose Up
BEGIN;
ALTER TABLE blocks
  ALTER COLUMN block_gaslimit TYPE DOUBLE PRECISION USING block_gaslimit :: DOUBLE PRECISION;

ALTER TABLE blocks
  ALTER COLUMN block_gasused TYPE DOUBLE PRECISION USING block_gasused :: DOUBLE PRECISION;

ALTER TABLE blocks
  ALTER COLUMN block_time TYPE DOUBLE PRECISION USING block_time :: DOUBLE PRECISION;

ALTER TABLE blocks
  ALTER COLUMN block_reward TYPE NUMERIC USING block_time :: NUMERIC;

ALTER TABLE blocks
  ALTER COLUMN block_uncles_reward TYPE NUMERIC USING block_time :: NUMERIC;

COMMIT;