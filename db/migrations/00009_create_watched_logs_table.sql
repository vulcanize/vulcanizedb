-- +goose Up
CREATE TABLE public.watched_logs
(
    id               SERIAL PRIMARY KEY,
    contract_address VARCHAR(42),
    topic_zero       VARCHAR(66)
);

-- +goose Down
DROP TABLE public.watched_logs;
