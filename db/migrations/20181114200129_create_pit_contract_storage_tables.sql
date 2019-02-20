-- +goose Up
CREATE TABLE maker.pit_drip (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  drip          TEXT
);

CREATE TABLE maker.pit_ilk_spot (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id),
  spot          NUMERIC NOT NULL
);

CREATE TABLE maker.pit_ilk_line (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           INTEGER NOT NULL REFERENCES maker.ilks (id),
  line          NUMERIC NOT NULL
);

CREATE TABLE maker.pit_line (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  line          NUMERIC NOT NULL
);

CREATE TABLE maker.pit_live (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  live          NUMERIC NOT NULL
);

CREATE TABLE maker.pit_vat (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vat           TEXT
);

-- +goose Down
DROP TABLE maker.pit_drip;
DROP TABLE maker.pit_ilk_spot;
DROP TABLE maker.pit_ilk_line;
DROP TABLE maker.pit_line;
DROP TABLE maker.pit_live;
DROP TABLE maker.pit_vat;