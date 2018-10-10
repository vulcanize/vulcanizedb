DROP TABLE maker.vat_heal;
ALTER TABLE public.checked_headers
    DROP COLUMN vat_heal_checked BOOLEAN NOT NULL DEFAULT FALSE;
