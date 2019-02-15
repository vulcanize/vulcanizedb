// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package postgres

import (
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres driver
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type DB struct {
	*sqlx.DB
	Node   core.Node
	NodeID int64
}

var (
	ErrDBInsertFailed     = errors.New("postgres: insert failed")
	ErrDBDeleteFailed     = errors.New("postgres: delete failed")
	ErrDBConnectionFailed = errors.New("postgres: db connection failed")
	ErrUnableToSetNode    = errors.New("postgres: unable to set node")
)

func NewDB(databaseConfig config.Database, node core.Node) (*DB, error) {
	connectString := config.DbConnectionString(databaseConfig)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return &DB{}, ErrDBConnectionFailed
	}
	pg := DB{DB: db, Node: node}
	err = pg.CreateNode(&node)
	if err != nil {
		return &DB{}, ErrUnableToSetNode
	}
	return &pg, nil
}

func (db *DB) CreateNode(node *core.Node) error {
	var nodeId int64
	err := db.QueryRow(
		`INSERT INTO eth_nodes (genesis_block, network_id, eth_node_id, client_name)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (genesis_block, network_id, eth_node_id)
                  DO UPDATE
                    SET genesis_block = $1,
                        network_id = $2,
                        eth_node_id = $3,
                        client_name = $4
                RETURNING id`,
		node.GenesisBlock, node.NetworkID, node.ID, node.ClientName).Scan(&nodeId)
	if err != nil {
		return ErrUnableToSetNode
	}
	db.NodeID = nodeId
	return nil
}
