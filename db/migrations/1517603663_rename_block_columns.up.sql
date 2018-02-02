BEGIN;
ALTER TABLE blocks
  RENAME COLUMN block_number TO number;

ALTER TABLE blocks
  RENAME COLUMN block_gaslimit TO gaslimit;

ALTER TABLE blocks
  RENAME COLUMN block_gasused TO gasused;

ALTER TABLE blocks
  RENAME COLUMN block_time TO time;

ALTER TABLE blocks
  RENAME COLUMN block_difficulty TO difficulty;

ALTER TABLE blocks
  RENAME COLUMN block_hash TO hash;

ALTER TABLE blocks
  RENAME COLUMN block_nonce TO nonce;

ALTER TABLE blocks
  RENAME COLUMN block_parenthash TO parenthash;

ALTER TABLE blocks
  RENAME COLUMN block_size TO size;

ALTER TABLE blocks
  RENAME COLUMN block_miner TO miner;

ALTER TABLE blocks
  RENAME COLUMN block_extra_data TO extra_data;

ALTER TABLE blocks
  RENAME COLUMN block_reward TO reward;

ALTER TABLE blocks
  RENAME COLUMN block_uncles_reward TO uncles_reward;
COMMIT;










