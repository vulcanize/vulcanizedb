-- +goose Up
ALTER TABLE blocks
  ADD COLUMN block_difficulty BIGINT,
  ADD COLUMN block_hash VARCHAR(66),
  ADD COLUMN block_nonce VARCHAR(20),
  ADD COLUMN block_parenthash VARCHAR(66),
  ADD COLUMN block_size BIGINT,
  ADD COLUMN uncle_hash VARCHAR(66);


-- +goose Down
ALTER TABLE blocks
  Drop COLUMN block_difficulty,
  Drop COLUMN block_hash,
  drop COLUMN block_nonce,
  drop COLUMN block_parenthash,
  drop COLUMN block_size,
  drop COLUMN uncle_hash;
