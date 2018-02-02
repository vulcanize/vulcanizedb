package repositories

import "github.com/vulcanize/vulcanizedb/pkg/core"

type NodeRepository interface {
	CreateNode(node *core.Node) error
}

func (pg *Postgres) CreateNode(node *core.Node) error {
	var nodeId int64
	err := pg.Db.QueryRow(
		`INSERT INTO nodes (genesis_block, network_id, node_id, client_name)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (genesis_block, network_id, node_id)
                  DO UPDATE
                    SET genesis_block = $1,
                        network_id = $2,
                        node_id = $3,
                        client_name = $4
                RETURNING id`,
		node.GenesisBlock, node.NetworkId, node.Id, node.ClientName).Scan(&nodeId)
	if err != nil {
		return ErrUnableToSetNode
	}
	pg.nodeId = nodeId
	return nil
}
