-- +goose Up
ALTER TABLE transactions
    ADD COLUMN tx_input_data VARCHAR;

-- +goose Down
ALTER TABLE transactions
    DROP COLUMN tx_input_data;
