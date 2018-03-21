BEGIN;
ALTER TABLE public.nodes RENAME TO eth_nodes;

ALTER TABLE public.eth_nodes RENAME COLUMN node_id TO eth_node_id;

ALTER TABLE public.eth_nodes DROP CONSTRAINT node_uc;
ALTER TABLE public.eth_nodes
  ADD CONSTRAINT eth_node_uc UNIQUE (genesis_block, network_id, eth_node_id);

ALTER TABLE public.blocks RENAME COLUMN node_id TO eth_node_id;

ALTER TABLE public.blocks DROP CONSTRAINT node_id_block_number_uc;
ALTER TABLE public.blocks
  ADD CONSTRAINT eth_node_id_block_number_uc UNIQUE (number, eth_node_id);

ALTER TABLE public.blocks DROP CONSTRAINT node_fk;
ALTER TABLE public.blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (eth_node_id) REFERENCES eth_nodes (id) ON DELETE CASCADE;

COMMIT;