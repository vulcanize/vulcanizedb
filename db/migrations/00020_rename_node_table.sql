-- +goose Up
ALTER TABLE public.nodes RENAME TO eth_nodes;

ALTER TABLE public.eth_nodes RENAME COLUMN node_id TO eth_node_id;

ALTER TABLE public.eth_nodes DROP CONSTRAINT node_uc;
ALTER TABLE public.eth_nodes
  ADD CONSTRAINT eth_node_uc UNIQUE (genesis_block, network_id, eth_node_id);

ALTER TABLE public.eth_blocks RENAME COLUMN node_id TO eth_node_id;

ALTER TABLE public.eth_blocks DROP CONSTRAINT node_id_block_number_uc;
ALTER TABLE public.eth_blocks
  ADD CONSTRAINT eth_node_id_block_number_uc UNIQUE (number, eth_node_id);

ALTER TABLE public.eth_blocks DROP CONSTRAINT node_fk;
ALTER TABLE public.eth_blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (eth_node_id) REFERENCES eth_nodes (id) ON DELETE CASCADE;


-- +goose Down
ALTER TABLE public.eth_nodes
  RENAME TO nodes;

ALTER TABLE public.nodes
  RENAME COLUMN eth_node_id TO node_id;

ALTER TABLE public.nodes
  DROP CONSTRAINT eth_node_uc;
ALTER TABLE public.nodes
  ADD CONSTRAINT node_uc UNIQUE (genesis_block, network_id, node_id);

ALTER TABLE public.eth_blocks RENAME COLUMN eth_node_id TO node_id;

ALTER TABLE public.eth_blocks DROP CONSTRAINT eth_node_id_block_number_uc;
ALTER TABLE public.eth_blocks
  ADD CONSTRAINT node_id_block_number_uc UNIQUE (number, node_id);

ALTER TABLE public.eth_blocks DROP CONSTRAINT node_fk;
ALTER TABLE public.eth_blocks
  ADD CONSTRAINT node_fk
FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE;
