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

type TokenBalanceDBRow struct {
	ID                 int64  `db:"id"`
	Balance            int64  `db:"balance"`
	BlockID            int64  `db:"block_id"`
	TokenAddress       string `db:"token_address"`
	TokenHolderAddress string `db:"token_holder_address"`
}

type TokenAllowanceDBRow struct {
	ID                  int64  `db:"id"`
	Allowance           int64  `db:"allowance"`
	BlockID             int64  `db:"block_id"`
	TokenAddress        string `db:"token_address"`
	TokenHolderAddress  string `db:"token_holder_address"`
	TokenSpenderAddress string `db:"token_spender_address"`
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
		ID:           "2ea672a45c4c7b96e3c4b130b21a22af390a552fd0b3cff96420b4bda26568d470dc56e05e453823f64f2556a6e4460ad1d4d00eb2d8b8fc16fcb1be73e86522",
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
