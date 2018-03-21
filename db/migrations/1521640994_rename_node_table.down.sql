BEGIN;
ALTER TABLE public.eth_nodes
  RENAME TO nodes;

ALTER TABLE public.nodes
  RENAME COLUMN eth_node_id TO node_id;

ALTER TABLE public.nodes
  DROP CONSTRAINT eth_node_uc;
ALTER TABLE public.nodes
  ADD CONSTRAINT node_uc UNIQUE (genesis_block, network_id, node_id);

ALTER TABLE public.blocks RENAME COLUMN eth_node_id TO node_id;

ALTER TABLE public.blocks DROP CONSTRAINT eth_node_id_block_number_uc;
ALTER TABLE public.blocks
  ADD CONSTRAINT node_id_block_number_uc UNIQUE (number, node_id);

ALTER TABLE public.blocks DROP CONSTRAINT node_fk;
ALTER TABLE public.blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE;
COMMIT;