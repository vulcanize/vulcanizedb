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

package every_block_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
	"strconv"
)

func setLastBlockOnChain(blockChain *fakes.MockBlockChain, blockNumber int64) {
	blockNumberString := strconv.FormatInt(blockNumber, 10)
	lastBlockOnChain := big.Int{}
	lastBlockOnChain.SetString(blockNumberString, 10)
	blockChain.SetLastBlock(&lastBlockOnChain)
}

var _ = Describe("Everyblock transformers", func() {
	var db *postgres.DB
	var blockChain *fakes.MockBlockChain
	var blockNumber int64
	var blockId int64
	var err error

	BeforeEach(func() {
		blockChain = fakes.NewMockBlockChain()
		blockNumber = erc20_watcher.DaiConfig.FirstBlock
		lastBlockNumber := blockNumber + 1
		db = test_helpers.CreateNewDatabase()
		setLastBlockOnChain(blockChain, lastBlockNumber)

		blockRepository := repositories.NewBlockRepository(db)

		blockId, err = blockRepository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
		Expect(err).NotTo(HaveOccurred())
		_, err = blockRepository.CreateOrUpdateBlock(core.Block{Number: lastBlockNumber})
		Expect(err).NotTo(HaveOccurred())
	})

	It("creates a token_supply record for each block in the given range", func() {
		initializer := every_block.ERC20TokenTransformerInitializer{Config: erc20_watcher.DaiConfig}
		transformer := initializer.NewERC20TokenTransformer(db, blockChain)
		transformer.Execute()

		var tokenSupplyCount int
		err := db.QueryRow(`SELECT COUNT(*) FROM token_supply where block_id = $1`, blockId).Scan(&tokenSupplyCount)
		Expect(err).ToNot(HaveOccurred())
		Expect(tokenSupplyCount).To(Equal(1))

		var tokenSupply test_helpers.TokenSupplyDBRow
		err = db.Get(&tokenSupply, `SELECT * from token_supply where block_id = $1`, blockId)
		Expect(err).ToNot(HaveOccurred())
		Expect(tokenSupply.BlockID).To(Equal(blockId))
		Expect(tokenSupply.TokenAddress).To(Equal(constants.DaiContractAddress))
		Expect(tokenSupply.Supply).To(Equal(int64(0)))
	})
})
