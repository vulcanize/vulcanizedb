-- +goose Up
COMMENT ON TABLE public.nodes IS E'@name NodeInfo';
COMMENT ON TABLE public.headers IS E'@name EthHeaders';
COMMENT ON COLUMN public.headers.node_id IS E'@name EthNodeID';
COMMENT ON COLUMN public.nodes.node_id IS E'@name ChainNodeID';