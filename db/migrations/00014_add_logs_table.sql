-- +goose Up
CREATE TABLE logs (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT,
  address      VARCHAR(66),
  tx_hash      VARCHAR(66),
  index        BIGINT,
  topic0       VARCHAR(66),
  topic1       VARCHAR(66),
  topic2       VARCHAR(66),
  topic3       VARCHAR(66),
  data         TEXT,
  CONSTRAINT log_uc UNIQUE (block_number, index)
);



-- +goose Down
DROP TABLE logs;
