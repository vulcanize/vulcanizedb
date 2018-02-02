BEGIN;
ALTER TABLE blocks
  RENAME COLUMN number TO block_number;

ALTER TABLE blocks
  RENAME COLUMN gaslimit TO block_gaslimit;

ALTER TABLE blocks
  RENAME COLUMN gasused TO block_gasused;

ALTER TABLE blocks
  RENAME COLUMN TIME TO block_time;

ALTER TABLE blocks
  RENAME COLUMN difficulty TO block_difficulty;

ALTER TABLE blocks
  RENAME COLUMN HASH TO block_hash;

ALTER TABLE blocks
  RENAME COLUMN nonce TO block_nonce;

ALTER TABLE blocks
  RENAME COLUMN parenthash TO block_parenthash;

ALTER TABLE blocks
  RENAME COLUMN size TO block_size;

ALTER TABLE blocks
  RENAME COLUMN miner TO block_miner;

ALTER TABLE blocks
  RENAME COLUMN extra_data TO block_extra_data;

ALTER TABLE blocks
  RENAME COLUMN reward TO block_reward;

ALTER TABLE blocks
  RENAME COLUMN uncles_reward TO block_uncles_reward;
COMMIT;