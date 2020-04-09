-- +goose Up
COMMENT ON TABLE eth.transaction_cids IS E'@name EthTransactionCids';
COMMENT ON TABLE eth.header_cids IS E'@name EthHeaderCids';
COMMENT ON TABLE eth.queue_data IS E'@name EthQueueData';
COMMENT ON COLUMN eth.header_cids.node_id IS E'@name EthNodeID';
