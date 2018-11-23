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
	"math/big"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
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
		blockNumber = generic.DaiConfig.FirstBlock
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
		transformer, err := every_block.NewERC20TokenTransformer(db, blockChain, generic.DaiConfig)
		Expect(err).ToNot(HaveOccurred())
		transformer.Execute()

		var tokenSupplyCount int
		err = db.QueryRow(`SELECT COUNT(*) FROM token_supply where block_id = $1`, blockId).Scan(&tokenSupplyCount)
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
