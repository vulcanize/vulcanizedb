DROP TABLE maker.cat_file_chop_lump;
DROP TABLE maker.cat_file_flip;
DROP TABLE maker.cat_file_pit_vow;

ALTER TABLE public.checked_headers
  DROP COLUMN cat_file_chop_lump_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN cat_file_flip_checked;

ALTER TABLE public.checked_headers
  DROP COLUMN cat_file_pit_vow_checked;