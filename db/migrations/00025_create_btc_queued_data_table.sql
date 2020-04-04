-- +goose Up
CREATE TABLE btc.queue_data (
  id SERIAL PRIMARY KEY,
  data BYTEA NOT NULL,
  height BIGINT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE btc.queue_data;