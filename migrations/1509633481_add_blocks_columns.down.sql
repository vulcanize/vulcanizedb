ALTER TABLE blocks
  Drop COLUMN block_difficulty,
  Drop COLUMN block_hash,
  drop COLUMN block_nonce,
  drop COLUMN block_parenthash,
  drop COLUMN block_size,
  drop COLUMN uncle_hash