DROP TABLE maker.pit_file_stability_fee;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_stability_fee_checked;
