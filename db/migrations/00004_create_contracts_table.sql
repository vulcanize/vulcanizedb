-- +goose Up
CREATE TABLE watched_contracts
(
  contract_id   SERIAL PRIMARY KEY,
  contract_abi  json,
  contract_hash VARCHAR(66) UNIQUE
);

-- +goose Down
DROP TABLE watched_contracts;
