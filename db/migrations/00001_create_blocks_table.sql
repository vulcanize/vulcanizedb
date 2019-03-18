-- +goose Up
CREATE TABLE public.blocks (
  id            SERIAL PRIMARY KEY,
  difficulty    BIGINT,
  extra_data    VARCHAR,
  gaslimit      BIGINT,
  gasused       BIGINT,
  hash          VARCHAR(66),
  miner         VARCHAR(42),
  nonce         VARCHAR(20),
  "number"      BIGINT,
  parenthash    VARCHAR(66),
  reward        DOUBLE PRECISION,
  uncles_reward DOUBLE PRECISION,
  "size"        VARCHAR,
  "time"        BIGINT,
  is_final      BOOLEAN,
  uncle_hash    VARCHAR(66)
);


-- +goose Down
DROP TABLE public.blocks;