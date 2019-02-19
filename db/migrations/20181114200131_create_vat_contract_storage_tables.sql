-- +goose Up
CREATE TABLE maker.vat_debt (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  debt          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_vice (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  vice          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_art (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           TEXT,
  art           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_ink (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           TEXT,
  ink           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_rate (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           TEXT,
  rate          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_ilk_take (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           TEXT,
  take          NUMERIC NOT NULL
);

CREATE TABLE maker.vat_urn_art (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           TEXT,
  urn           TEXT,
  art           TEXT
);

CREATE TABLE maker.vat_urn_ink (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           TEXT,
  urn           TEXT,
  ink           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_gem (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  ilk           TEXT,
  guy           TEXT,
  gem           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_dai (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  guy           TEXT,
  dai           NUMERIC NOT NULL
);

CREATE TABLE maker.vat_sin (
  id            SERIAL PRIMARY KEY,
  block_number  BIGINT,
  block_hash    TEXT,
  guy           TEXT,
  sin           NUMERIC NOT NULL
);

-- +goose Down
DROP TABLE maker.vat_debt;
DROP TABLE maker.vat_vice;
DROP TABLE maker.vat_ilk_art;
DROP TABLE maker.vat_ilk_ink;
DROP TABLE maker.vat_ilk_rate;
DROP TABLE maker.vat_ilk_take;
DROP TABLE maker.vat_urn_art;
DROP TABLE maker.vat_urn_ink;
DROP TABLE maker.vat_gem;
DROP TABLE maker.vat_dai;
DROP TABLE maker.vat_sin;
