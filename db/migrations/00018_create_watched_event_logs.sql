-- +goose Up
CREATE VIEW block_stats AS
SELECT max(block_number) AS max_block,
       min(block_number) AS min_block
FROM full_sync_logs;

CREATE VIEW watched_event_logs AS
SELECT log_filters.name,
       full_sync_logs.id,
       block_number,
       full_sync_logs.address,
       tx_hash,
       index,
       full_sync_logs.topic0,
       full_sync_logs.topic1,
       full_sync_logs.topic2,
       full_sync_logs.topic3,
       data,
       receipt_id
FROM log_filters
         CROSS JOIN block_stats
         JOIN full_sync_logs ON full_sync_logs.address = log_filters.address
    AND full_sync_logs.block_number >= coalesce(log_filters.from_block, block_stats.min_block)
    AND full_sync_logs.block_number <= coalesce(log_filters.to_block, block_stats.max_block)
WHERE (log_filters.topic0 = full_sync_logs.topic0 OR log_filters.topic0 ISNULL)
  AND (log_filters.topic1 = full_sync_logs.topic1 OR log_filters.topic1 ISNULL)
  AND (log_filters.topic2 = full_sync_logs.topic2 OR log_filters.topic2 ISNULL)
  AND (log_filters.topic3 = full_sync_logs.topic3 OR log_filters.topic3 ISNULL);

-- +goose Down
DROP VIEW watched_event_logs;
DROP VIEW block_stats;
