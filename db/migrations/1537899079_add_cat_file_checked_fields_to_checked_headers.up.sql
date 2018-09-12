ALTER TABLE public.checked_headers
    ADD COLUMN cat_file_chop_lump_checked BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.checked_headers
    ADD COLUMN cat_file_flip_checked BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.checked_headers
    ADD COLUMN cat_file_pit_vow_checked BOOLEAN NOT NULL DEFAULT FALSE;
