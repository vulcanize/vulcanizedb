-- +goose Up
CREATE TABLE public.eth_blocks (
  id            SERIAL PRIMARY KEY,
  difficulty    BIGINT,
  extra_data    VARCHAR,
  gas_limit      BIGINT,
  gas_used       BIGINT,
  hash          VARCHAR(66),
  miner         VARCHAR(42),
  nonce         VARCHAR(20),
  "number"      BIGINT,
  parent_hash    VARCHAR(66),
  reward        NUMERIC,
  uncles_reward NUMERIC,
  "size"        VARCHAR,
  "time"        BIGINT,
  is_final      BOOLEAN,
  uncle_hash    VARCHAR(66)
);


-- +goose Down
DROP TABLE public.eth_blocks;