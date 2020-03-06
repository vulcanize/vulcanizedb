-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION public.get_or_create_header(block_number BIGINT, hash VARCHAR(66), raw JSONB, block_timestamp NUMERIC, eth_node_id INTEGER) RETURNS INTEGER AS
$$
DECLARE
    matching_header_id INTEGER := (
        SELECT id
        FROM public.headers
        WHERE headers.block_number = get_or_create_header.block_number
        AND headers.hash = get_or_create_header.hash
        );
    nonmatching_header_id INTEGER := (
        SELECT id
        FROM public.headers
        WHERE headers.block_number = get_or_create_header.block_number
        AND headers.hash != get_or_create_header.hash
        );
    inserted_header_id INTEGER;
BEGIN
    IF matching_header_id != 0 THEN
        RETURN matching_header_id;
    END IF;

    IF nonmatching_header_id != 0 THEN
        DELETE FROM public.headers WHERE id = nonmatching_header_id;
    end if;

    INSERT INTO public.headers (hash, block_number, raw, block_timestamp, eth_node_id)
        VALUES (get_or_create_header.hash, get_or_create_header.block_number, get_or_create_header.raw, get_or_create_header.block_timestamp, get_or_create_header.eth_node_id)
        RETURNING id INTO inserted_header_id;

    RETURN inserted_header_id;
END
    $$
    LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION public.get_or_create_header(block_number BIGINT, hash VARCHAR, raw JSONB, block_timestamp NUMERIC, eth_node_id INTEGER);
