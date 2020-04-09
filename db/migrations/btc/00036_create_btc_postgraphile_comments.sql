-- +goose Up
COMMENT ON TABLE btc.header_cids IS E'@name BtcHeaderCids';
COMMENT ON TABLE btc.transaction_cids IS E'@name BtcTransactionCids';
COMMENT ON TABLE btc.queue_data IS E'@name BtcQueueData';
COMMENT ON COLUMN btc.header_cids.node_id IS E'@name BtcNodeID';