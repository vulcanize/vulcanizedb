CREATE VIEW block_stats AS
  SELECT
    max(block_number) AS max_block,
    min(block_number) AS min_block
  FROM logs;

CREATE VIEW watched_event_logs AS
  SELECT
    log_filters.name,
    logs.id,
    block_number,
    logs.address,
    tx_hash,
    index,
    logs.topic0,
    logs.topic1,
    logs.topic2,
    logs.topic3,
    data,
    receipt_id
  FROM log_filters
    CROSS JOIN block_stats
    JOIN logs ON logs.address = log_filters.address
                 AND logs.block_number >= coalesce(log_filters.from_block, block_stats.min_block)
                 AND logs.block_number <= coalesce(log_filters.to_block, block_stats.max_block)
  WHERE (log_filters.topic0 = logs.topic0 OR log_filters.topic0 ISNULL)
        AND (log_filters.topic1 = logs.topic1 OR log_filters.topic1 ISNULL)
        AND (log_filters.topic2 = logs.topic2 OR log_filters.topic2 ISNULL)
        AND (log_filters.topic3 = logs.topic3 OR log_filters.topic3 ISNULL);