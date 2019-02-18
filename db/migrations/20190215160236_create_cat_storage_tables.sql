-- +goose Up
CREATE TABLE maker.cat_nflip (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  nflip NUMERIC NOT NULL
);

CREATE TABLE maker.cat_live (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  live NUMERIC NOT NULL
);

CREATE TABLE maker.cat_vat (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  vat TEXT
);

CREATE TABLE maker.cat_pit (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  pit TEXT
);

CREATE TABLE maker.cat_vow (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  vow TEXT
);

CREATE TABLE maker.cat_ilk_flip (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  ilk TEXT,
  flip TEXT
);

CREATE TABLE maker.cat_ilk_chop (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  ilk TEXT,
  chop NUMERIC NOT NULL
);

CREATE TABLE maker.cat_ilk_lump (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  ilk TEXT,
  lump NUMERIC NOT NULL
);

CREATE TABLE maker.cat_flip_ilk (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  nflip NUMERIC NOT NULL,
  ilk TEXT
);

CREATE TABLE maker.cat_flip_urn (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  nflip NUMERIC NOT NULL,
  urn TEXT
);

CREATE TABLE maker.cat_flip_ink (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  nflip NUMERIC NOT NULL,
  ink NUMERIC NOT NULL
);

CREATE TABLE maker.cat_flip_tab (
  id SERIAL PRIMARY KEY,
  block_number BIGINT,
  block_hash TEXT,
  nflip NUMERIC NOT NULL,
  tab NUMERIC NOT NULL
);


-- +goose Down
DROP TABLE maker.cat_nflip;
DROP TABLE maker.cat_live;
DROP TABLE maker.cat_vat;
DROP TABLE maker.cat_pit;
DROP TABLE maker.cat_vow;
DROP TABLE maker.cat_ilk_flip;
DROP TABLE maker.cat_ilk_chop;
DROP TABLE maker.cat_ilk_lump;
DROP TABLE maker.cat_flip_ilk;
DROP TABLE maker.cat_flip_urn;
DROP TABLE maker.cat_flip_ink;
DROP TABLE maker.cat_flip_tab;