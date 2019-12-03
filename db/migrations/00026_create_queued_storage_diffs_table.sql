-- +goose Up
CREATE TABLE public.queued_storage
(
    id      SERIAL PRIMARY KEY,
    diff_id BIGINT UNIQUE NOT NULL REFERENCES public.storage_diff (id)
);

-- +goose Down
DROP TABLE public.queued_storage;
