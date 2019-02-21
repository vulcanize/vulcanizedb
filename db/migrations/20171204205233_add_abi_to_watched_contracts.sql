-- +goose Up
ALTER TABLE watched_contracts
  ADD COLUMN contract_abi json;

-- +goose Down
ALTER TABLE watched_contracts
    DROP COLUMN contract_abi;
