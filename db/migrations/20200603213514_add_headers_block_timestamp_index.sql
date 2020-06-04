-- +goose Up
CREATE INDEX headers_block_timestamp_index ON public.headers (block_timestamp);


-- +goose Down
DROP INDEX headers_block_timestamp_index;