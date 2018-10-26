DROP TABLE maker.drip_file_ilk;
DROP TABLE maker.drip_file_repo;
DROP TABLE maker.drip_file_vow;

ALTER TABLE public.checked_headers
  DROP COLUMN drip_file_ilk_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN drip_file_repo_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN drip_file_vow_checked;