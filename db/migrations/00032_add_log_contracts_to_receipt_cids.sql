-- +goose Up
ALTER TABLE eth.receipt_cids
ADD COLUMN log_contracts VARCHAR(66)[];

-- +goose Down
ALTER TABLE eth.receipt_cids
DROP COLUMN log_contracts;