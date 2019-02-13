-- +goose Up
CREATE TABLE maker.vow_vat (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vat           TEXT
);

CREATE TABLE maker.vow_cow (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  cow           TEXT
);

CREATE TABLE maker.vow_row (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  row           TEXT
);

CREATE TABLE maker.vow_sin (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  sin           numeric
);

CREATE TABLE maker.vow_woe (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  woe           numeric
);

CREATE TABLE maker.vow_ash (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ash           numeric
);

CREATE TABLE maker.vow_wait (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  wait          numeric
);

CREATE TABLE maker.vow_sump (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  sump          numeric
);

CREATE TABLE maker.vow_bump (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  bump          numeric
);

CREATE TABLE maker.vow_hump (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  hump          numeric
);

-- +goose Down
DROP TABLE maker.vow_vat;
DROP TABLE maker.vow_cow;
DROP TABLE maker.vow_row;
DROP TABLE maker.vow_sin;
DROP TABLE maker.vow_woe;
DROP TABLE maker.vow_ash;
DROP TABLE maker.vow_wait;
DROP TABLE maker.vow_sump;
DROP TABLE maker.vow_bump;
DROP TABLE maker.vow_hump;
