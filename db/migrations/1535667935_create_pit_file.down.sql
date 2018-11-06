DROP TABLE maker.pit_file_ilk;
DROP TABLE maker.pit_file_stability_fee;
DROP TABLE maker.pit_file_debt_ceiling;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_debt_ceiling_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_ilk_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN pit_file_stability_fee_checked;