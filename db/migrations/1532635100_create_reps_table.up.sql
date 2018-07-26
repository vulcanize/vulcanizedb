CREATE TABLE maker.reps (
  id           SERIAL PRIMARY KEY,
  block_number BIGINT  NOT NULL,
  header_id    INTEGER NOT NULL,
  usd_value    NUMERIC,
  CONSTRAINT headers_fk FOREIGN KEY (header_id)
  REFERENCES headers (id)
  ON DELETE CASCADE
);