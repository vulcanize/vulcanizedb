-- +goose Up
CREATE TABLE eth_nodes
(
    id            SERIAL PRIMARY KEY,
    client_name   VARCHAR,
    genesis_block VARCHAR(66),
    network_id    NUMERIC,
    eth_node_id   VARCHAR(128),
    UNIQUE (genesis_block, network_id, eth_node_id)
);

-- +goose Down
DROP TABLE eth_nodes;
