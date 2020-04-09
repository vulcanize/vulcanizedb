-- +goose Up
CREATE TABLE eth.queue_data (
  id SERIAL PRIMARY KEY,
  data BYTEA NOT NULL,
  height BIGINT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE eth.queue_data;