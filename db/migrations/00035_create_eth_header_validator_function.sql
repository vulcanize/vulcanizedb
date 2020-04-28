-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION eth_valid_section(hash VARCHAR(66), depth INT)
        RETURNS TABLE (
            id INT,
            block_number BIGINT
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
  SELECT id, block_number FROM validator ORDER BY block_number DESC;
$$ LANGUAGE SQL;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION eth_valid_section;