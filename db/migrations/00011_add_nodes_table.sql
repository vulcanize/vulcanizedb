-- +goose Up
CREATE TABLE nodes (
  id            SERIAL PRIMARY KEY,
  genesis_block VARCHAR(66),
  network_id NUMERIC,
  CONSTRAINT node_uc UNIQUE (genesis_block, network_id)
);

-- +goose Down
DROP TABLE nodes;
