-- +goose Up
ALTER TABLE watched_contracts
  ADD CONSTRAINT contract_hash_uc UNIQUE (contract_hash);


-- +goose Down
ALTER TABLE watched_contracts
  DROP CONSTRAINT contract_hash_uc;