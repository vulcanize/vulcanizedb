-- +goose Up
CREATE TABLE public.checked_headers (
  id                  SERIAL PRIMARY KEY,
  header_id           INTEGER UNIQUE NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  price_feeds_checked BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE public.checked_headers;
