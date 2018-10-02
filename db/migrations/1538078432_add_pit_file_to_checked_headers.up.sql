ALTER TABLE public.checked_headers
    ADD COLUMN pit_file_debt_ceiling_checked BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.checked_headers
    ADD COLUMN pit_file_ilk_checked BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE public.checked_headers
    ADD COLUMN pit_file_stability_fee_checked BOOLEAN NOT NULL DEFAULT FALSE;