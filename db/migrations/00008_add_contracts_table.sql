-- +goose Up
CREATE TABLE watched_contracts
(
  contract_id SERIAL PRIMARY KEY,
  contract_hash VARCHAR(66)
)

-- +goose Down
DROP TABLE watched_contracts