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

package test_helpers

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

type TokenSupplyDBRow struct {
	ID           int64  `db:"id"`
	Supply       int64  `db:"supply"`
	BlockID      int64  `db:"block_id"`
	TokenAddress string `db:"token_address"`
}

type TransferDBRow struct {
	ID             int64 `db:"id"`
	VulcanizeLogID int64 `db:"vulcanize_log_id"`
}

func CreateNewDatabase() *postgres.DB {
	var node core.Node
	node = core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
	db := test_config.NewTestDB(node)

	_, err := db.Exec(`DELETE FROM logs`)
	Expect(err).NotTo(HaveOccurred())

	return db
}

func CreateBlock(blockNumber int64, repository repositories.BlockRepository) (blockId int64) {
	blockId, err := repository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
	Expect(err).NotTo(HaveOccurred())

	return blockId
}
