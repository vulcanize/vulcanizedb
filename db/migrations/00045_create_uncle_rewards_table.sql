-- +goose Up
CREATE TABLE public.uncle_rewards (
  id            SERIAL PRIMARY KEY,
  block_id      INTEGER,
  block_hash    VARCHAR(66) NOT NULL,
  uncle_hash    VARCHAR(66) NOT NULL,
  uncle_reward  NUMERIC NOT NULL,
  miner_address VARCHAR(66) NOT NULL,
  CONSTRAINT block_id_fk FOREIGN KEY (block_id)
  REFERENCES blocks (id)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE public.uncle_rewards;
