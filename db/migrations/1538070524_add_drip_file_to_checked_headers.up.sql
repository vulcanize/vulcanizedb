ALTER TABLE public.checked_headers
    ADD COLUMN drip_file_ilk_checked BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.checked_headers
    ADD COLUMN drip_file_repo_checked BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.checked_headers
    ADD COLUMN drip_file_vow_checked BOOLEAN NOT NULL DEFAULT FALSE;