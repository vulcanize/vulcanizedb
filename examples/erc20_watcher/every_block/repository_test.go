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

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
	"math/rand"
	"github.com/vulcanize/vulcanizedb/examples/erc20_test_helpers"
)

var _ = Describe("ERC20 Token Repository", func() {
	var db *postgres.DB
	var blockId int64
	var blockNumber int64
	var repository every_block.TokenSupplyRepository
	var blockRepository repositories.BlockRepository
	testAddress := "abc"

	BeforeEach(func() {
		node := test_config.NewTestNode()
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		repository = every_block.TokenSupplyRepository{DB: db}
		blockRepository = *repositories.NewBlockRepository(db)
		blockNumber = rand.Int63()
		blockId = test_config.NewTestBlock(blockNumber, blockRepository)
	})

	Describe("Create", func() {
		It("creates a token supply record", func() {
			supply := supplyModel(blockNumber, testAddress, "100")
			err := repository.Create(supply)
			Expect(err).NotTo(HaveOccurred())

			dbResult := erc20_test_helpers.TokenSupplyDBRow{}
			expectedTokenSupply := erc20_test_helpers.TokenSupplyDBRow{
				Supply:       int64(100),
				BlockID:      blockId,
				TokenAddress: testAddress,
			}

			var count int
			err = repository.DB.QueryRowx(`SELECT count(*) FROM token_supply`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			err = repository.DB.QueryRowx(`SELECT * FROM token_supply`).StructScan(&dbResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Supply).To(Equal(expectedTokenSupply.Supply))
			Expect(dbResult.BlockID).To(Equal(expectedTokenSupply.BlockID))
			Expect(dbResult.TokenAddress).To(Equal(expectedTokenSupply.TokenAddress))
		})

		It("returns an error if fetching the block's id from the database fails", func() {
			errorSupply := supplyModel(-1, "", "")
			err := repository.Create(errorSupply)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql"))
			Expect(err.Error()).To(ContainSubstring("block number -1"))
		})

		It("returns an error if inserting the token_supply fails", func() {
			errorSupply := supplyModel(blockNumber, "", "")
			err := repository.Create(errorSupply)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq"))
			Expect(err.Error()).To(ContainSubstring("token_supply for block number"))
		})
	})

	Describe("When there are multiple nodes", func() {
		var node2DB *postgres.DB
		var node2BlockRepo *repositories.BlockRepository
		var node2BlockId int64
		var node2TokenSupplyRepo every_block.TokenSupplyRepository
		var tokenSupply every_block.TokenSupply

		BeforeEach(func() {
			node2DB = createDbForAnotherNode()

			//create another block with the same number on node2
			node2BlockRepo = repositories.NewBlockRepository(node2DB)
			node2BlockId = test_config.NewTestBlock(blockNumber, *node2BlockRepo)

			tokenSupply = supplyModel(blockNumber, "abc", "100")
			node2TokenSupplyRepo = every_block.TokenSupplyRepository{DB: node2DB}
		})

		It("only creates token_supply records for the current node (node2)", func() {
			err := node2TokenSupplyRepo.Create(tokenSupply)
			Expect(err).NotTo(HaveOccurred())

			var tokenSupplies []erc20_test_helpers.TokenSupplyDBRow
			err = node2TokenSupplyRepo.DB.Select(&tokenSupplies, `SELECT * FROM token_supply`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(tokenSupplies)).To(Equal(1))
			Expect(tokenSupplies[0].BlockID).To(Equal(node2BlockId))
		})

		It("only includes missing block numbers for the current node", func() {
			//create token_supply on original node
			err := repository.Create(tokenSupply)
			Expect(err).NotTo(HaveOccurred())

			originalNodeMissingBlocks, err := repository.MissingBlocks(blockNumber, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(originalNodeMissingBlocks)).To(Equal(0))

			node2MissingBlocks, err := node2TokenSupplyRepo.MissingBlocks(blockNumber, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingBlocks)).To(Equal(1))
		})
	})

	Describe("MissingBlocks", func() {
		It("returns the block numbers for which an associated TokenSupply record hasn't been created", func() {
			createTokenSupplyFor(repository, blockNumber)

			newBlockNumber := blockNumber + 1
			test_config.NewTestBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingBlocks(blockNumber, newBlockNumber)

			Expect(blocks).To(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("only returns blocks within the given range", func() {
			newBlockNumber := blockNumber + 1
			test_config.NewTestBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingBlocks(blockNumber, blockNumber)

			Expect(blocks).NotTo(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return numbers that already have an associated TokenSupply record", func() {
			createTokenSupplyFor(repository, blockNumber)
			blocks, err := repository.MissingBlocks(blockNumber, blockNumber)

			Expect(blocks).To(BeEmpty())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("deletes the token supply record when the associated block is deleted", func() {
		err := repository.Create(every_block.TokenSupply{BlockNumber: blockNumber, Value: "0"})
		Expect(err).NotTo(HaveOccurred())

		var count int
		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_supply`).Scan(&count)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))

		_, err = db.Query(`DELETE FROM blocks`)
		Expect(err).NotTo(HaveOccurred())

		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_supply`).Scan(&count)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(0))
	})
})

func supplyModel(blockNumber int64, tokenAddress string, supplyValue string) every_block.TokenSupply {
	return every_block.TokenSupply{
		Value:        supplyValue,
		TokenAddress: tokenAddress,
		BlockNumber:  int64(blockNumber),
	}
}

func createTokenSupplyFor(repository every_block.TokenSupplyRepository, blockNumber int64) {
	err := repository.Create(every_block.TokenSupply{BlockNumber: blockNumber, Value: "0"})
	Expect(err).NotTo(HaveOccurred())
}

func createDbForAnotherNode() *postgres.DB {
	anotherNode := core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "testNodeId",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}

	return test_config.NewTestDB(anotherNode)
}
