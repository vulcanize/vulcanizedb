// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package every_block

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type ERC20RepositoryInterface interface {
	Create(supply TokenSupply) error
	MissingBlocks(startingBlock int64, highestBlock int64) ([]int64, error)
}

type TokenSupplyRepository struct {
	*postgres.DB
}

type repositoryError struct {
	err         string
	msg         string
	blockNumber int64
}

func (re *repositoryError) Error() string {
	return fmt.Sprintf(re.msg, re.blockNumber, re.err)
}

func newRepositoryError(err error, msg string, blockNumber int64) error {
	e := repositoryError{err.Error(), msg, blockNumber}
	log.Println(e.Error())
	return &e
}

const (
	GetBlockError          = "Error fetching block number %d: %s"
	InsertTokenSupplyError = "Error inserting token_supply for block number %d: %s"
	MissingBlockError      = "Error finding missing token_supply records starting at block %d: %s"
)

func (tsp *TokenSupplyRepository) Create(supply TokenSupply) error {
	var blockId int
	err := tsp.DB.Get(&blockId, `SELECT id FROM blocks WHERE number = $1 AND eth_node_id = $2`, supply.BlockNumber, tsp.NodeID)
	if err != nil {
		return newRepositoryError(err, GetBlockError, supply.BlockNumber)
	}

	_, err = tsp.DB.Exec(
		`INSERT INTO token_supply (supply, token_address, block_id)
                VALUES($1, $2, $3)`,
		supply.Value, supply.TokenAddress, blockId)
	if err != nil {
		return newRepositoryError(err, InsertTokenSupplyError, supply.BlockNumber)
	}
	return nil
}

func (tsp *TokenSupplyRepository) MissingBlocks(startingBlock int64, highestBlock int64) ([]int64, error) {
	blockNumbers := make([]int64, 0)

	err := tsp.DB.Select(
		&blockNumbers,
		`SELECT number FROM BLOCKS
               LEFT JOIN token_supply ON blocks.id = block_id
               WHERE block_id ISNULL
               AND eth_node_id = $1
               AND number >= $2
               AND number <= $3
               LIMIT 20`,
		tsp.NodeID,
		startingBlock,
		highestBlock,
	)
	if err != nil {
		return []int64{}, newRepositoryError(err, MissingBlockError, startingBlock)
	}
	return blockNumbers, err
}
