-- +goose Up
CREATE TABLE maker.ilks (
  id        SERIAL PRIMARY KEY,
  ilk       TEXT UNIQUE
);

-- +goose Down
DROP TABLE maker.ilks;
