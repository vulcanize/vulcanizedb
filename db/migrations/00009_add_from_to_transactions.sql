-- +goose Up
ALTER TABLE transactions
  ADD COLUMN tx_from VARCHAR(66);


-- +goose Down
ALTER TABLE transactions
  DROP COLUMN tx_from;
