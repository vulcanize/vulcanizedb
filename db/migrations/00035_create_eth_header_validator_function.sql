-- +goose Up
-- +goose StatementBegin
-- Returns all of the header ids, in descending order, that are recursively validated from the provided hash
-- includes the header id for the provided hash and the id found at height = depth (so it returns depth+1 ids)
CREATE OR REPLACE FUNCTION eth_valid_section(hash VARCHAR(66), depth INT)
        RETURNS TABLE (
            id INT
) AS $$
  WITH RECURSIVE validator AS (
          SELECT id, parent_hash, block_number
          FROM eth.header_cids
          WHERE block_hash = hash
      UNION
          SELECT eth.header_cids.id, eth.header_cids.parent_hash, eth.header_cids.block_number
          FROM eth.header_cids
          INNER JOIN validator
             ON eth.header_cids.block_hash = validator.parent_hash
             AND eth.header_cids.block_number = validator.block_number - 1
             AND eth.header_cids.block_number >= depth
  )
  SELECT id FROM validator ORDER BY block_number DESC;
$$ LANGUAGE SQL;
-- +goose StatementEnd

-- +goose StatementBegin
-- Returns all of the latest state diffs (ids) for every path from height = n+1 to the height defined at the provided hash
-- Use this to return all of the diffs to add on top of the state_cache at `n` to build the state_cache at `hash`
CREATE OR REPLACE FUNCTION eth_valid_state_diffs(hash VARCHAR(66), n INT)
        RETURNS TABLE (
          state_id INT,
          state_path BYTEA,
          cid TEXT
) AS $$
  SELECT DISTINCT ON (state_path) state_cids.id, state_path, state_cids.cid
  FROM eth.state_cids INNER JOIN eth.header_cids ON state_cids.header_id = header_cids.id
  WHERE state_cids.header_id = ANY(SELECT * FROM eth_valid_section(hash, n+1))
  ORDER BY state_path, block_number DESC;
$$ LANGUAGE SQL;
-- +goose StatementEnd


-- +goose Down
DROP FUNCTION eth_valid_section;
DROP FUNCTION eth_valid_state_diffs;