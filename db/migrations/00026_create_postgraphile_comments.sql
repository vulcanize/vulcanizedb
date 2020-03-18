-- +goose Up
COMMENT ON TABLE public.nodes IS E'@name NodeInfo';
COMMENT ON TABLE btc.header_cids IS E'@name BtcHeaderCids';
COMMENT ON TABLE btc.transaction_cids IS E'@name BtcTransactionCids';
COMMENT ON TABLE btc.queue_data IS E'@name BtcQueueData';
COMMENT ON TABLE eth.transaction_cids IS E'@name EthTransactionCids';
COMMENT ON TABLE eth.header_cids IS E'@name EthHeaderCids';
COMMENT ON TABLE eth.queue_data IS E'@name EthQueueData';
COMMENT ON TABLE public.headers IS E'@name EthHeaders';
COMMENT ON COLUMN public.headers.node_id IS E'@name EthNodeID';
COMMENT ON COLUMN public.nodes.node_id IS E'@name ChainNodeID';
COMMENT ON COLUMN eth.header_cids.node_id IS E'@name EthNodeID';
COMMENT ON COLUMN btc.header_cids.node_id IS E'@name BtcNodeID';