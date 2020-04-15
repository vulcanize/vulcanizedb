-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION public.create_back_filled_diff(block_height BIGINT, block_hash BYTEA, hashed_address BYTEA,
                                                          storage_key BYTEA, storage_value BYTEA) RETURNS VOID AS
$$
DECLARE
    last_storage_value  BYTEA := (
        SELECT storage_diff.storage_value
        FROM public.storage_diff
        WHERE storage_diff.block_height <= create_back_filled_diff.block_height
          AND storage_diff.hashed_address = create_back_filled_diff.hashed_address
          AND storage_diff.storage_key = create_back_filled_diff.storage_key
        ORDER BY storage_diff.block_height DESC
        LIMIT 1
    );
    empty_storage_value BYTEA := (
        SELECT '\x0000000000000000000000000000000000000000000000000000000000000000'::BYTEA
    );
BEGIN
    IF last_storage_value = create_back_filled_diff.storage_value THEN
        RETURN;
    END IF;

    IF last_storage_value is null and create_back_filled_diff.storage_value = empty_storage_value THEN
        RETURN;
    END IF;

    INSERT INTO public.storage_diff (block_height, block_hash, hashed_address, storage_key, storage_value,
                                     from_backfill)
    VALUES (create_back_filled_diff.block_height, create_back_filled_diff.block_hash,
            create_back_filled_diff.hashed_address, create_back_filled_diff.storage_key,
            create_back_filled_diff.storage_value, true)
    ON CONFLICT DO NOTHING;

    RETURN;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION public.create_back_filled_diff(block_height BIGINT, block_hash BYTEA, hashed_address BYTEA,
                                                          storage_key BYTEA, storage_value BYTEA) RETURNS VOID AS
$$
DECLARE
    last_storage_value  BYTEA := (
        SELECT storage_diff.storage_value
        FROM public.storage_diff
        WHERE storage_diff.hashed_address = create_back_filled_diff.hashed_address
          AND storage_diff.storage_key = create_back_filled_diff.storage_key
        ORDER BY storage_diff.block_height DESC
        LIMIT 1
    );
    empty_storage_value BYTEA := (
        SELECT '\x0000000000000000000000000000000000000000000000000000000000000000'::BYTEA
    );
BEGIN
    IF last_storage_value = create_back_filled_diff.storage_value THEN
        RETURN;
    END IF;

    IF last_storage_value is null and create_back_filled_diff.storage_value = empty_storage_value THEN
        RETURN;
    END IF;

    INSERT INTO public.storage_diff (block_height, block_hash, hashed_address, storage_key, storage_value,
                                     from_backfill)
    VALUES (create_back_filled_diff.block_height, create_back_filled_diff.block_hash,
            create_back_filled_diff.hashed_address, create_back_filled_diff.storage_key,
            create_back_filled_diff.storage_value, true)
    ON CONFLICT DO NOTHING;

    RETURN;
END
$$
    LANGUAGE plpgsql;
-- +goose StatementEnd