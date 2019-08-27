-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE public.checked_logs
(
    id               SERIAL PRIMARY KEY,
    contract_address VARCHAR(42),
    topic_zero       VARCHAR(66)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE public.checked_logs;
