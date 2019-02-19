-- +goose Up
DROP TABLE maker.pit_file_stability_fee;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_stability_fee_checked;


-- +goose Down
CREATE TABLE maker.pit_file_stability_fee (
  id SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  what      TEXT,
  data      TEXT,
  log_idx   INTEGER NOT NULL,
  tx_idx    INTEGER NOT NULL,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN pit_file_stability_fee_checked BOOLEAN NOT NULL DEFAULT FALSE;
