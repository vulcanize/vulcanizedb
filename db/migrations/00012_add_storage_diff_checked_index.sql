-- +goose Up
CREATE INDEX storage_diff_checked_index ON public.storage_diff (checked) WHERE checked = false;


-- +goose Down
DROP INDEX storage_diff_checked_index;