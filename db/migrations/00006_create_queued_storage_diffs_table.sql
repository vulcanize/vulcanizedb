-- +goose Up
CREATE TABLE public.queued_storage
(
    id      SERIAL PRIMARY KEY,
    diff_id BIGINT UNIQUE NOT NULL REFERENCES public.storage_diff (id)
);

COMMENT ON TABLE public.queued_storage
    IS E'@omit';

-- +goose Down
DROP TABLE public.queued_storage;
