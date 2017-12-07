ALTER TABLE watched_contracts
  ADD CONSTRAINT contract_hash_uc UNIQUE (contract_hash);
