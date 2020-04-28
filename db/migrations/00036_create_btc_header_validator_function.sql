-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION btc_valid_section(hash VARCHAR(66), depth INT)
        RETURNS TABLE (
            id INT,
            block_number BIGINT
) AS $$
  WITH RECURSIVE validator AS (
          SELECT id, parent_hash, block_number
          FROM btc.header_cids
          WHERE block_hash = hash
      UNION
          SELECT btc.header_cids.id, btc.header_cids.parent_hash, btc.header_cids.block_number
          FROM btc.header_cids
          INNER JOIN validator
             ON btc.header_cids.block_hash = validator.parent_hash
             AND btc.header_cids.block_number = validator.block_number - 1
             AND btc.header_cids.block_number >= depth
  )
  SELECT id, block_number FROM validator ORDER BY block_number DESC;
$$ LANGUAGE SQL;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION btc_valid_section;