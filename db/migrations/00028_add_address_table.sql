-- +goose Up
CREATE TABLE public.addresses
(
    id      SERIAL PRIMARY KEY,
    address character varying(42),
    UNIQUE (address)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE public.addresses;
