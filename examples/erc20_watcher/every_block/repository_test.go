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

package every_block_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("ERC20 Token Supply Repository", func() {
	var db *postgres.DB
	var blockId int64
	var blockNumber int64
	var repository every_block.ERC20TokenRepository
	var blockRepository repositories.BlockRepository
	testAddress := "abc"

	BeforeEach(func() {
		db = test_helpers.CreateNewDatabase()
		repository = every_block.ERC20TokenRepository{DB: db}
		_, err := db.Query(`DELETE FROM token_supply`)
		Expect(err).NotTo(HaveOccurred())

		blockRepository = *repositories.NewBlockRepository(db)
		blockNumber = rand.Int63()
		blockId = test_helpers.CreateBlock(blockNumber, blockRepository)
	})

	Describe("Create", func() {
		It("creates a token supply record", func() {
			supply := supplyModel(blockNumber, testAddress, "100")
			err := repository.CreateSupply(supply)
			Expect(err).NotTo(HaveOccurred())

			dbResult := test_helpers.TokenSupplyDBRow{}
			expectedTokenSupply := test_helpers.TokenSupplyDBRow{
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
			err := repository.CreateSupply(errorSupply)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql"))
			Expect(err.Error()).To(ContainSubstring("block number -1"))
		})

		It("returns an error if inserting the token_supply fails", func() {
			errorSupply := supplyModel(blockNumber, "", "")
			err := repository.CreateSupply(errorSupply)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq"))
			Expect(err.Error()).To(ContainSubstring("token_supply for block number"))
		})
	})

	Describe("When there are multiple nodes", func() {
		var node2DB *postgres.DB
		var node2BlockRepo *repositories.BlockRepository
		var node2BlockId int64
		var node2TokenSupplyRepo every_block.ERC20TokenRepository
		var tokenSupply every_block.TokenSupply

		BeforeEach(func() {
			node2DB = createDbForAnotherNode()

			//create another block with the same number on node2
			node2BlockRepo = repositories.NewBlockRepository(node2DB)
			node2BlockId = test_helpers.CreateBlock(blockNumber, *node2BlockRepo)

			tokenSupply = supplyModel(blockNumber, "abc", "100")
			node2TokenSupplyRepo = every_block.ERC20TokenRepository{DB: node2DB}
		})

		It("only creates token_supply records for the current node (node2)", func() {
			err := node2TokenSupplyRepo.CreateSupply(tokenSupply)
			Expect(err).NotTo(HaveOccurred())

			var tokenSupplies []test_helpers.TokenSupplyDBRow
			err = node2TokenSupplyRepo.DB.Select(&tokenSupplies, `SELECT * FROM token_supply`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(tokenSupplies)).To(Equal(1))
			Expect(tokenSupplies[0].BlockID).To(Equal(node2BlockId))
		})

		It("only includes missing block numbers for the current node", func() {
			//create token_supply on original node
			err := repository.CreateSupply(tokenSupply)
			Expect(err).NotTo(HaveOccurred())

			originalNodeMissingBlocks, err := repository.MissingSupplyBlocks(blockNumber, blockNumber, testAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(originalNodeMissingBlocks)).To(Equal(0))

			node2MissingBlocks, err := node2TokenSupplyRepo.MissingSupplyBlocks(blockNumber, blockNumber, testAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingBlocks)).To(Equal(1))
		})
	})

	Describe("MissingBlocks", func() {
		It("returns the block numbers for which an associated TokenSupply record hasn't been created", func() {
			createTokenSupplyFor(repository, blockNumber, testAddress)

			newBlockNumber := blockNumber + 1
			test_helpers.CreateBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingSupplyBlocks(blockNumber, newBlockNumber, testAddress)

			Expect(blocks).To(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("only returns blocks within the given range", func() {
			newBlockNumber := blockNumber + 1
			test_helpers.CreateBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingSupplyBlocks(blockNumber, blockNumber, testAddress)

			Expect(blocks).NotTo(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return numbers that already have an associated TokenSupply record", func() {
			createTokenSupplyFor(repository, blockNumber, testAddress)
			blocks, err := repository.MissingSupplyBlocks(blockNumber, blockNumber, testAddress)

			Expect(blocks).To(BeEmpty())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("deletes the token supply record when the associated block is deleted", func() {
		err := repository.CreateSupply(every_block.TokenSupply{BlockNumber: blockNumber, TokenAddress: testAddress, Value: "0"})
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

var _ = Describe("ERC20 Token Balance Repository", func() {
	var db *postgres.DB
	var blockId int64
	var blockNumber int64
	var repository every_block.ERC20TokenRepository
	var blockRepository repositories.BlockRepository
	testTokenAddress := "abc"
	testHolderAddress := "def"

	BeforeEach(func() {
		db = test_helpers.CreateNewDatabase()
		repository = every_block.ERC20TokenRepository{DB: db}
		_, err := db.Query(`DELETE FROM token_balance`)
		Expect(err).NotTo(HaveOccurred())

		blockRepository = *repositories.NewBlockRepository(db)
		blockNumber = rand.Int63()
		blockId = test_helpers.CreateBlock(blockNumber, blockRepository)
	})

	Describe("Create", func() {
		It("creates a token balance record", func() {
			balance := balanceOfModel(blockNumber, testTokenAddress, testHolderAddress, "100")
			err := repository.CreateBalance(balance)
			Expect(err).NotTo(HaveOccurred())

			dbResult := test_helpers.TokenBalanceDBRow{}
			expectedTokenBalance := test_helpers.TokenBalanceDBRow{
				Balance:            int64(100),
				BlockID:            blockId,
				TokenAddress:       testTokenAddress,
				TokenHolderAddress: testHolderAddress,
			}

			var count int
			err = repository.DB.QueryRowx(`SELECT count(*) FROM token_balance`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			err = repository.DB.QueryRowx(`SELECT * FROM token_balance`).StructScan(&dbResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Balance).To(Equal(expectedTokenBalance.Balance))
			Expect(dbResult.BlockID).To(Equal(expectedTokenBalance.BlockID))
			Expect(dbResult.TokenAddress).To(Equal(expectedTokenBalance.TokenAddress))
			Expect(dbResult.TokenHolderAddress).To(Equal(expectedTokenBalance.TokenHolderAddress))
		})

		It("returns an error if fetching the block's id from the database fails", func() {
			errorBalance := balanceOfModel(-1, "", "", "")
			err := repository.CreateBalance(errorBalance)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql"))
			Expect(err.Error()).To(ContainSubstring("block number -1"))
		})

		It("returns an error if inserting the token_balance fails", func() {
			errorBalance := balanceOfModel(blockNumber, "", "", "")
			err := repository.CreateBalance(errorBalance)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq"))
			Expect(err.Error()).To(ContainSubstring("token_balance for block number"))
		})
	})

	Describe("When there are multiple nodes", func() {
		var node2DB *postgres.DB
		var node2BlockRepo *repositories.BlockRepository
		var node2BlockId int64
		var node2TokenSupplyRepo every_block.ERC20TokenRepository
		var tokenBalance every_block.TokenBalance

		BeforeEach(func() {
			node2DB = createDbForAnotherNode()

			//create another block with the same number on node2
			node2BlockRepo = repositories.NewBlockRepository(node2DB)
			node2BlockId = test_helpers.CreateBlock(blockNumber, *node2BlockRepo)

			tokenBalance = balanceOfModel(blockNumber, "abc", "def", "100")
			node2TokenSupplyRepo = every_block.ERC20TokenRepository{DB: node2DB}
		})

		It("only creates token_balance records for the current node (node2)", func() {
			err := node2TokenSupplyRepo.CreateBalance(tokenBalance)
			Expect(err).NotTo(HaveOccurred())

			var tokenBalances []test_helpers.TokenBalanceDBRow
			err = node2TokenSupplyRepo.DB.Select(&tokenBalances, `SELECT * FROM token_balance`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(tokenBalances)).To(Equal(1))
			Expect(tokenBalances[0].BlockID).To(Equal(node2BlockId))
		})

		It("only includes missing block numbers for the current node", func() {
			//create token_balance on original node
			err := repository.CreateBalance(tokenBalance)
			Expect(err).NotTo(HaveOccurred())

			originalNodeMissingBlocks, err := repository.MissingBalanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(originalNodeMissingBlocks)).To(Equal(0))

			node2MissingBlocks, err := node2TokenSupplyRepo.MissingBalanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingBlocks)).To(Equal(1))
		})
	})

	Describe("MissingBlocks", func() {
		It("returns the block numbers for which an associated TokenBalance record hasn't been created", func() {
			createTokenBalanceFor(repository, blockNumber, testTokenAddress, testHolderAddress)

			newBlockNumber := blockNumber + 1
			test_helpers.CreateBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingBalanceBlocks(blockNumber, newBlockNumber, testTokenAddress, testHolderAddress)

			Expect(blocks).To(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("only returns blocks within the given range", func() {
			newBlockNumber := blockNumber + 1
			test_helpers.CreateBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingBalanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress)

			Expect(blocks).NotTo(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return numbers that already have an associated TokenBalance record", func() {
			createTokenBalanceFor(repository, blockNumber, testTokenAddress, testHolderAddress)
			blocks, err := repository.MissingBalanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress)

			Expect(blocks).To(BeEmpty())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("deletes the token balance record when the associated block is deleted", func() {
		err := repository.CreateBalance(every_block.TokenBalance{
			BlockNumber:        blockNumber,
			TokenAddress:       testTokenAddress,
			TokenHolderAddress: testHolderAddress,
			Value:              "0",
		})
		Expect(err).NotTo(HaveOccurred())

		var count int
		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_balance`).Scan(&count)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))

		_, err = db.Query(`DELETE FROM blocks`)
		Expect(err).NotTo(HaveOccurred())

		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_balance`).Scan(&count)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(0))
	})
})

var _ = Describe("ERC20 Token Allowance Repository", func() {
	var db *postgres.DB
	var blockId int64
	var blockNumber int64
	var repository every_block.ERC20TokenRepository
	var blockRepository repositories.BlockRepository
	testTokenAddress := "abc"
	testHolderAddress := "def"
	testSpenderAddress := "ghi"

	BeforeEach(func() {
		db = test_helpers.CreateNewDatabase()
		repository = every_block.ERC20TokenRepository{DB: db}
		_, err := db.Query(`DELETE FROM token_allowance`)
		Expect(err).NotTo(HaveOccurred())

		blockRepository = *repositories.NewBlockRepository(db)
		blockNumber = rand.Int63()
		blockId = test_helpers.CreateBlock(blockNumber, blockRepository)
	})

	Describe("Create", func() {
		It("creates a token balance record", func() {
			allowance := allowanceModel(blockNumber, testTokenAddress, testHolderAddress, testSpenderAddress, "100")
			err := repository.CreateAllowance(allowance)
			Expect(err).NotTo(HaveOccurred())

			dbResult := test_helpers.TokenAllowanceDBRow{}
			expectedTokenAllowance := test_helpers.TokenAllowanceDBRow{
				Allowance:           int64(100),
				BlockID:             blockId,
				TokenAddress:        testTokenAddress,
				TokenHolderAddress:  testHolderAddress,
				TokenSpenderAddress: testSpenderAddress,
			}

			var count int
			err = repository.DB.QueryRowx(`SELECT count(*) FROM token_allowance`).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			err = repository.DB.QueryRowx(`SELECT * FROM token_allowance`).StructScan(&dbResult)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbResult.Allowance).To(Equal(expectedTokenAllowance.Allowance))
			Expect(dbResult.BlockID).To(Equal(expectedTokenAllowance.BlockID))
			Expect(dbResult.TokenAddress).To(Equal(expectedTokenAllowance.TokenAddress))
			Expect(dbResult.TokenHolderAddress).To(Equal(expectedTokenAllowance.TokenHolderAddress))
		})

		It("returns an error if fetching the block's id from the database fails", func() {
			errorAllowance := allowanceModel(-1, "", "", "", "")
			err := repository.CreateAllowance(errorAllowance)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql"))
			Expect(err.Error()).To(ContainSubstring("block number -1"))
		})

		It("returns an error if inserting the token_allowance fails", func() {
			errorAllowance := allowanceModel(blockNumber, "", "", "", "")
			err := repository.CreateAllowance(errorAllowance)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq"))
			Expect(err.Error()).To(ContainSubstring("token_allowance for block number"))
		})
	})

	Describe("When there are multiple nodes", func() {
		var node2DB *postgres.DB
		var node2BlockRepo *repositories.BlockRepository
		var node2BlockId int64
		var node2TokenSupplyRepo every_block.ERC20TokenRepository
		var tokenAllowance every_block.TokenAllowance

		BeforeEach(func() {
			node2DB = createDbForAnotherNode()

			//create another block with the same number on node2
			node2BlockRepo = repositories.NewBlockRepository(node2DB)
			node2BlockId = test_helpers.CreateBlock(blockNumber, *node2BlockRepo)

			tokenAllowance = allowanceModel(blockNumber, "abc", "def", "ghi", "100")
			node2TokenSupplyRepo = every_block.ERC20TokenRepository{DB: node2DB}
		})

		It("only creates token_allowance records for the current node (node2)", func() {
			err := node2TokenSupplyRepo.CreateAllowance(tokenAllowance)
			Expect(err).NotTo(HaveOccurred())

			var tokenAllowances []test_helpers.TokenAllowanceDBRow
			err = node2TokenSupplyRepo.DB.Select(&tokenAllowances, `SELECT * FROM token_allowance`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(tokenAllowances)).To(Equal(1))
			Expect(tokenAllowances[0].BlockID).To(Equal(node2BlockId))
		})

		It("only includes missing block numbers for the current node", func() {
			//create token_allowance on original node
			err := repository.CreateAllowance(tokenAllowance)
			Expect(err).NotTo(HaveOccurred())

			originalNodeMissingBlocks, err := repository.MissingAllowanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress, testSpenderAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(originalNodeMissingBlocks)).To(Equal(0))

			node2MissingBlocks, err := node2TokenSupplyRepo.MissingAllowanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress, testSpenderAddress)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(node2MissingBlocks)).To(Equal(1))
		})
	})

	Describe("MissingBlocks", func() {
		It("returns the block numbers for which an associated TokenAllowance record hasn't been created", func() {
			createTokenAllowanceFor(repository, blockNumber, testTokenAddress, testHolderAddress, testSpenderAddress)

			newBlockNumber := blockNumber + 1
			test_helpers.CreateBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingAllowanceBlocks(blockNumber, newBlockNumber, testTokenAddress, testHolderAddress, testSpenderAddress)

			Expect(blocks).To(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("only returns blocks within the given range", func() {
			newBlockNumber := blockNumber + 1
			test_helpers.CreateBlock(newBlockNumber, blockRepository)
			blocks, err := repository.MissingAllowanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress, testSpenderAddress)

			Expect(blocks).NotTo(ConsistOf(newBlockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return numbers that already have an associated TokenAllowance record", func() {
			createTokenAllowanceFor(repository, blockNumber, testTokenAddress, testHolderAddress, testSpenderAddress)
			blocks, err := repository.MissingAllowanceBlocks(blockNumber, blockNumber, testTokenAddress, testHolderAddress, testSpenderAddress)

			Expect(blocks).To(BeEmpty())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("deletes the token balance record when the associated block is deleted", func() {
		err := repository.CreateAllowance(every_block.TokenAllowance{
			BlockNumber:         blockNumber,
			TokenAddress:        testTokenAddress,
			TokenHolderAddress:  testHolderAddress,
			TokenSpenderAddress: testSpenderAddress,
			Value:               "0",
		})
		Expect(err).NotTo(HaveOccurred())

		var count int
		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_allowance`).Scan(&count)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))

		_, err = db.Query(`DELETE FROM blocks`)
		Expect(err).NotTo(HaveOccurred())

		err = repository.DB.QueryRowx(`SELECT count(*) FROM token_allowance`).Scan(&count)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(0))
	})
})

func supplyModel(blockNumber int64, tokenAddress, supplyValue string) every_block.TokenSupply {
	return every_block.TokenSupply{
		Value:        supplyValue,
		TokenAddress: tokenAddress,
		BlockNumber:  blockNumber,
	}
}

func balanceOfModel(blockNumber int64, tokenAddress, holderAddress, supplyValue string) every_block.TokenBalance {
	return every_block.TokenBalance{
		Value:              supplyValue,
		TokenAddress:       tokenAddress,
		TokenHolderAddress: holderAddress,
		BlockNumber:        blockNumber,
	}
}

func allowanceModel(blockNumber int64, tokenAddress, holderAddress, spenderAddress, supplyValue string) every_block.TokenAllowance {
	return every_block.TokenAllowance{
		Value:               supplyValue,
		TokenAddress:        tokenAddress,
		TokenHolderAddress:  holderAddress,
		TokenSpenderAddress: spenderAddress,
		BlockNumber:         blockNumber,
	}
}

func createTokenSupplyFor(repository every_block.ERC20TokenRepository, blockNumber int64, tokenAddress string) {
	err := repository.CreateSupply(every_block.TokenSupply{
		BlockNumber:  blockNumber,
		TokenAddress: tokenAddress,
		Value:        "0",
	})
	Expect(err).NotTo(HaveOccurred())
}

func createTokenBalanceFor(repository every_block.ERC20TokenRepository, blockNumber int64, tokenAddress, holderAddress string) {
	err := repository.CreateBalance(every_block.TokenBalance{
		BlockNumber:        blockNumber,
		TokenAddress:       tokenAddress,
		TokenHolderAddress: holderAddress,
		Value:              "0",
	})
	Expect(err).NotTo(HaveOccurred())
}

func createTokenAllowanceFor(repository every_block.ERC20TokenRepository, blockNumber int64, tokenAddress, holderAddress, spenderAddress string) {
	err := repository.CreateAllowance(every_block.TokenAllowance{
		BlockNumber:         blockNumber,
		TokenAddress:        tokenAddress,
		TokenHolderAddress:  holderAddress,
		TokenSpenderAddress: spenderAddress,
		Value:               "0",
	})
	Expect(err).NotTo(HaveOccurred())
}

func createDbForAnotherNode() *postgres.DB {
	anotherNode := core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "testNodeId",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}

	return test_config.NewTestDBWithoutDeletingRecords(anotherNode)
}
