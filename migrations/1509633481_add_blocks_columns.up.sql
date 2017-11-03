ALTER TABLE blocks
  ADD COLUMN block_difficulty BIGINT,
  ADD COLUMN block_hash VARCHAR(66),
  ADD COLUMN block_nonce VARCHAR(20),
  ADD COLUMN block_parenthash VARCHAR(66),
  ADD COLUMN block_size BIGINT,
  ADD COLUMN uncle_hash VARCHAR(66)
